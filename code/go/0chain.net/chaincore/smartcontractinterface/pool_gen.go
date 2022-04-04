package smartcontractinterface

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *PoolStats) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 7
	// string "DelegateID"
	o = append(o, 0x87, 0xaa, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x49, 0x44)
	o = msgp.AppendString(o, z.DelegateID)
	// string "High"
	o = append(o, 0xa4, 0x48, 0x69, 0x67, 0x68)
	o, err = z.High.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "High")
		return
	}
	// string "Low"
	o = append(o, 0xa3, 0x4c, 0x6f, 0x77)
	o, err = z.Low.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Low")
		return
	}
	// string "InterestPaid"
	o = append(o, 0xac, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x65, 0x73, 0x74, 0x50, 0x61, 0x69, 0x64)
	o, err = z.InterestPaid.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "InterestPaid")
		return
	}
	// string "RewardPaid"
	o = append(o, 0xaa, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x50, 0x61, 0x69, 0x64)
	o, err = z.RewardPaid.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "RewardPaid")
		return
	}
	// string "NumRounds"
	o = append(o, 0xa9, 0x4e, 0x75, 0x6d, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x73)
	o = msgp.AppendInt64(o, z.NumRounds)
	// string "Status"
	o = append(o, 0xa6, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73)
	o = msgp.AppendString(o, z.Status)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PoolStats) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "DelegateID":
			z.DelegateID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "DelegateID")
				return
			}
		case "High":
			bts, err = z.High.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "High")
				return
			}
		case "Low":
			bts, err = z.Low.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Low")
				return
			}
		case "InterestPaid":
			bts, err = z.InterestPaid.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "InterestPaid")
				return
			}
		case "RewardPaid":
			bts, err = z.RewardPaid.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "RewardPaid")
				return
			}
		case "NumRounds":
			z.NumRounds, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "NumRounds")
				return
			}
		case "Status":
			z.Status, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Status")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *PoolStats) Msgsize() (s int) {
	s = 1 + 11 + msgp.StringPrefixSize + len(z.DelegateID) + 5 + z.High.Msgsize() + 4 + z.Low.Msgsize() + 13 + z.InterestPaid.Msgsize() + 11 + z.RewardPaid.Msgsize() + 10 + msgp.Int64Size + 7 + msgp.StringPrefixSize + len(z.Status)
	return
}