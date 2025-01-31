package limiter

import (
	"strconv"
	"sync"
	"sync/atomic"
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
		timestamp  = uint64(time.Now().Unix())
		expiration = uint64(cfg.Expiration.Seconds())
	)

	// Create manager to simplify storage operations ( see manager.go )
	manager := newManager(cfg.Storage)

	// Update timestamp every second
	go func() {
		for {
			atomic.StoreUint64(&timestamp, uint64(time.Now().Unix()))
			time.Sleep(1 * time.Second)
		}
	}()

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
		// Get entry from pool and release when finished
		e := manager.get(key)

		// Get timestamp
		ts := atomic.LoadUint64(&timestamp)
		// Set expiration if entry does not exist
		if e.exp == 0 {
			e.exp = ts + expiration

		} else if ts >= e.exp {
			// Check if entry is expired
			e.prevHits = e.currHits
			e.currHits = 0

			// checks how into the current window it is and sets the
			// expiry based on that, otherwise this would only reset on
			// the next request and not show the correct expiry
			elapsed := ts - e.exp
			if elapsed >= expiration {
				e.exp = ts + expiration
			} else {
				e.exp = ts + expiration - elapsed
			}
		}

		// Increment hits
		e.currHits++

		// Calculate when it resets in seconds
		expire := e.exp - ts

		// weight = time until current window reset / total window length
		weight := float64(expire) / float64(expiration)

		// rate = request count in previous window - weight + request count in current window
		rate := int(float64(e.prevHits)*weight) + e.currHits

		// Calculate how many hits can be made based on the current rate
		remaining := cfg.Max - rate

		// Update storage. Garbage collect when the next window ends.
		// |-------------------------|-------------------------|
		//               ^           ^              ^          ^
		//              ts        e.exp  End sample window   End next window
		//               <----------->
		// 				    expire
		// expire = e.exp - ts - time until end of current window.
		// duration + expiration = end of next window.
		// Because we don't want to garbage collect in the middle of a window
		// we add the expiration to the duration.
		// Otherwise after the end of "sample window", attackers could launch
		// a new request with the full window length.
		manager.set(key, e, time.Duration(expire+expiration)*time.Second)

		// Unlock entry
		mux.Unlock()

		// Check if hits exceed the cfg.Max
		if remaining < 0 {
			// Return response with Retry-After header
			// https://tools.ietf.org/html/rfc6584
			c.Set(fiber.HeaderRetryAfter, strconv.FormatUint(expire, 10))

			// Call LimitReached handler
			return cfg.LimitReached(c)
		}

		// We can continue, update RateLimit headers
		c.Set(xRateLimitLimit, max)
		c.Set(xRateLimitRemaining, strconv.Itoa(int(remaining)))
		c.Set(xRateLimitReset, strconv.FormatUint(expire, 10))

		// Continue stack
		return c.Next()
	}
}
