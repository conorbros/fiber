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
		case "currHits":
			z.currHits, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "prevHits":
			z.prevHits, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "exp":
			z.exp, err = dc.ReadUint64()
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
	// map header, size 3
	// write "currHits"
	err = en.Append(0x83, 0xa8, 0x63, 0x75, 0x72, 0x72, 0x48, 0x69, 0x74, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.currHits)
	if err != nil {
		return
	}
	// write "prevHits"
	err = en.Append(0xa8, 0x70, 0x72, 0x65, 0x76, 0x48, 0x69, 0x74, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.prevHits)
	if err != nil {
		return
	}
	// write "exp"
	err = en.Append(0xa3, 0x65, 0x78, 0x70)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.exp)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z item) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "currHits"
	o = append(o, 0x83, 0xa8, 0x63, 0x75, 0x72, 0x72, 0x48, 0x69, 0x74, 0x73)
	o = msgp.AppendInt(o, z.currHits)
	// string "prevHits"
	o = append(o, 0xa8, 0x70, 0x72, 0x65, 0x76, 0x48, 0x69, 0x74, 0x73)
	o = msgp.AppendInt(o, z.prevHits)
	// string "exp"
	o = append(o, 0xa3, 0x65, 0x78, 0x70)
	o = msgp.AppendUint64(o, z.exp)
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
		case "currHits":
			z.currHits, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "prevHits":
			z.prevHits, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "exp":
			z.exp, bts, err = msgp.ReadUint64Bytes(bts)
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
	s = 1 + 9 + msgp.IntSize + 9 + msgp.IntSize + 4 + msgp.Uint64Size
	return
}
