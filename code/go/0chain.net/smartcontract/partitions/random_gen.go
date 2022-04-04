package partitions

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z ItemType) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendInt(o, int(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ItemType) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 int
		zb0001, bts, err = msgp.ReadIntBytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = ItemType(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ItemType) Msgsize() (s int) {
	s = msgp.IntSize
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *randomSelectorDecode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Name"
	o = append(o, 0x84, 0xa4, 0x4e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "PartitionSize"
	o = append(o, 0xad, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x69, 0x7a, 0x65)
	o = msgp.AppendInt(o, z.PartitionSize)
	// string "NumPartitions"
	o = append(o, 0xad, 0x4e, 0x75, 0x6d, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	o = msgp.AppendInt(o, z.NumPartitions)
	// string "ItemType"
	o = append(o, 0xa8, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, int(z.ItemType))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *randomSelectorDecode) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Name")
				return
			}
		case "PartitionSize":
			z.PartitionSize, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "PartitionSize")
				return
			}
		case "NumPartitions":
			z.NumPartitions, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "NumPartitions")
				return
			}
		case "ItemType":
			{
				var zb0002 int
				zb0002, bts, err = msgp.ReadIntBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "ItemType")
					return
				}
				z.ItemType = ItemType(zb0002)
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
func (z *randomSelectorDecode) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 14 + msgp.IntSize + 14 + msgp.IntSize + 9 + msgp.IntSize
	return
}