package limiter

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/gofiber/fiber/v2/internal/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *item) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxvk uint32
	zxvk, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxvk > 0 {
		zxvk--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "last":
			z.last, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "tokens":
			z.tokens, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z item) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "last"
	err = en.Append(0x82, 0xa4, 0x6c, 0x61, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteTime(z.last)
	if err != nil {
		return
	}
	// write "tokens"
	err = en.Append(0xa6, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.tokens)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z item) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "last"
	o = append(o, 0x82, 0xa4, 0x6c, 0x61, 0x73, 0x74)
	o = msgp.AppendTime(o, z.last)
	// string "tokens"
	o = append(o, 0xa6, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73)
	o = msgp.AppendFloat64(o, z.tokens)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *item) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "last":
			z.last, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				return
			}
		case "tokens":
			z.tokens, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z item) Msgsize() (s int) {
	s = 1 + 5 + msgp.TimeSize + 7 + msgp.Float64Size
	return
}
