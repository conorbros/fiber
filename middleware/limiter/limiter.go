// inspired by golang/time/blob/master/rate/rate.go

package limiter

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	// Storage ErrNotExist
	errNotExist = "key does not exist"

	// X-RateLimit-* headers
	xRateLimitLimit     = "X-RateLimit-Limit"
	xRateLimitRemaining = "X-RateLimit-Remaining"
	xRateLimitReset     = "X-RateLimit-Reset"
)

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	var (
		// Limiter variables
		mux        = &sync.RWMutex{}
		max        = strconv.Itoa(cfg.Max)
		expiration = uint64(cfg.Expiration.Seconds())
	)

	// Create manager to simplify storage operations ( see manager.go )
	manager := newManager(cfg.Storage)

	durationFromTokens := func(tokens float64) time.Duration {
		seconds := tokens / float64(cfg.Max)
		return time.Nanosecond * time.Duration(1e9*seconds)
	}

	tokensFromDuration := func(d time.Duration) float64 {
		sec := float64(d/time.Second) * float64(cfg.Max)
		nsec := float64(d%time.Second) * float64(cfg.Max)
		return sec + nsec/1e9
	}

	setHeaders := func(tokens float64, c *fiber.Ctx) *fiber.Ctx {
		c.Set(xRateLimitLimit, max)
		c.Set(xRateLimitRemaining, strconv.Itoa(int(tokens)))
		c.Set(xRateLimitReset, strconv.FormatUint(uint64(expiration), 10))
		return c
	}

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Get key from request
		key := cfg.KeyGenerator(c)

		// Lock entry
		mux.Lock()
		defer mux.Unlock()

		// Get entry from pool and release when finished
		e := manager.get(key)

		// If first time seen key or last entry had expired and was garbage collected
		if e.last.Equal(time.Time{}) && e.tokens == 0 {
			e.tokens = float64(cfg.Max) - 1
			e.last = time.Now()
			manager.set(key, e, cfg.Expiration)
			return setHeaders(e.tokens, c).Next()
		}

		now := time.Now()
		last := e.last
		if now.Before(last) {
			last = now
		}

		maxElapsed := durationFromTokens(float64(cfg.Max) - e.tokens)
		elapsed := now.Sub(last)
		if elapsed > maxElapsed {
			elapsed = maxElapsed
		}

		// new bucket amount = gained tokens - token used for this request
		tokens := e.tokens + tokensFromDuration(elapsed) - 1
		if tokens > float64(cfg.Max) {
			tokens = float64(cfg.Max)
		}

		// Calculate when it resets in seconds
		expire := e.exp - ts

		elapsed := ts - (e.exp - expiration)
		revoked := int(float64(cfg.Max) / float64(expiration) * float64(elapsed))

		if e.hits -= revoked; e.hits < 0 {
			e.hits = 0
		}

		// Set how many hits we have left
		remaining := cfg.Max - e.hits

		fmt.Printf("%v %v %v\n\n", e.hits, revoked, ts)

		// Update storage
		manager.set(key, e, cfg.Expiration)

		// Check if no tokens remaining
		if tokens < 0 {
			retryAfter := durationFromTokens(-tokens)
			// Return response with Retry-After header
			// https://tools.ietf.org/html/rfc6584
			c.Set(fiber.HeaderRetryAfter, strconv.FormatUint(uint64(retryAfter), 10))
			c.Set(xRateLimitLimit, max)
			// Call LimitReached handler
			return cfg.LimitReached(c)
		}

		// We can continue, update RateLimit headers and continue stack
		return setHeaders(tokens, c).Next()
	}
}
