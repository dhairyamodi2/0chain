package stakepool

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *DelegatePool) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	// string "Balance"
	o = append(o, 0x86, 0xa7, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65)
	o, err = z.Balance.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Balance")
		return
	}
	// string "Reward"
	o = append(o, 0xa6, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64)
	o, err = z.Reward.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Reward")
		return
	}
	// string "Status"
	o = append(o, 0xa6, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73)
	o, err = z.Status.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Status")
		return
	}
	// string "RoundCreated"
	o = append(o, 0xac, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64)
	o = msgp.AppendInt64(o, z.RoundCreated)
	// string "DelegateID"
	o = append(o, 0xaa, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x49, 0x44)
	o = msgp.AppendString(o, z.DelegateID)
	// string "StakedAt"
	o = append(o, 0xa8, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x64, 0x41, 0x74)
	o, err = z.StakedAt.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "StakedAt")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DelegatePool) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Balance":
			bts, err = z.Balance.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Balance")
				return
			}
		case "Reward":
			bts, err = z.Reward.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Reward")
				return
			}
		case "Status":
			bts, err = z.Status.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Status")
				return
			}
		case "RoundCreated":
			z.RoundCreated, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "RoundCreated")
				return
			}
		case "DelegateID":
			z.DelegateID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "DelegateID")
				return
			}
		case "StakedAt":
			bts, err = z.StakedAt.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "StakedAt")
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
func (z *DelegatePool) Msgsize() (s int) {
	s = 1 + 8 + z.Balance.Msgsize() + 7 + z.Reward.Msgsize() + 7 + z.Status.Msgsize() + 13 + msgp.Int64Size + 11 + msgp.StringPrefixSize + len(z.DelegateID) + 9 + z.StakedAt.Msgsize()
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *DelegatePoolStat) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 12
	// string "ID"
	o = append(o, 0x8c, 0xa2, 0x49, 0x44)
	o = msgp.AppendString(o, z.ID)
	// string "Balance"
	o = append(o, 0xa7, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65)
	o, err = z.Balance.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Balance")
		return
	}
	// string "DelegateID"
	o = append(o, 0xaa, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x49, 0x44)
	o = msgp.AppendString(o, z.DelegateID)
	// string "Rewards"
	o = append(o, 0xa7, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73)
	o, err = z.Rewards.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Rewards")
		return
	}
	// string "UnStake"
	o = append(o, 0xa7, 0x55, 0x6e, 0x53, 0x74, 0x61, 0x6b, 0x65)
	o = msgp.AppendBool(o, z.UnStake)
	// string "ProviderId"
	o = append(o, 0xaa, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x49, 0x64)
	o = msgp.AppendString(o, z.ProviderId)
	// string "ProviderType"
	o = append(o, 0xac, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65)
	o, err = z.ProviderType.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "ProviderType")
		return
	}
	// string "TotalReward"
	o = append(o, 0xab, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64)
	o, err = z.TotalReward.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "TotalReward")
		return
	}
	// string "TotalPenalty"
	o = append(o, 0xac, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79)
	o, err = z.TotalPenalty.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "TotalPenalty")
		return
	}
	// string "Status"
	o = append(o, 0xa6, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73)
	o = msgp.AppendString(o, z.Status)
	// string "RoundCreated"
	o = append(o, 0xac, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64)
	o = msgp.AppendInt64(o, z.RoundCreated)
	// string "StakedAt"
	o = append(o, 0xa8, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x64, 0x41, 0x74)
	o, err = z.StakedAt.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "StakedAt")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DelegatePoolStat) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ID":
			z.ID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "Balance":
			bts, err = z.Balance.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Balance")
				return
			}
		case "DelegateID":
			z.DelegateID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "DelegateID")
				return
			}
		case "Rewards":
			bts, err = z.Rewards.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Rewards")
				return
			}
		case "UnStake":
			z.UnStake, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "UnStake")
				return
			}
		case "ProviderId":
			z.ProviderId, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ProviderId")
				return
			}
		case "ProviderType":
			bts, err = z.ProviderType.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "ProviderType")
				return
			}
		case "TotalReward":
			bts, err = z.TotalReward.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "TotalReward")
				return
			}
		case "TotalPenalty":
			bts, err = z.TotalPenalty.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "TotalPenalty")
				return
			}
		case "Status":
			z.Status, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Status")
				return
			}
		case "RoundCreated":
			z.RoundCreated, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "RoundCreated")
				return
			}
		case "StakedAt":
			bts, err = z.StakedAt.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "StakedAt")
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
func (z *DelegatePoolStat) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID) + 8 + z.Balance.Msgsize() + 11 + msgp.StringPrefixSize + len(z.DelegateID) + 8 + z.Rewards.Msgsize() + 8 + msgp.BoolSize + 11 + msgp.StringPrefixSize + len(z.ProviderId) + 13 + z.ProviderType.Msgsize() + 12 + z.TotalReward.Msgsize() + 13 + z.TotalPenalty.Msgsize() + 7 + msgp.StringPrefixSize + len(z.Status) + 13 + msgp.Int64Size + 9 + z.StakedAt.Msgsize()
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Settings) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "DelegateWallet"
	o = append(o, 0x85, 0xae, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74)
	o = msgp.AppendString(o, z.DelegateWallet)
	// string "MinStake"
	o = append(o, 0xa8, 0x4d, 0x69, 0x6e, 0x53, 0x74, 0x61, 0x6b, 0x65)
	o, err = z.MinStake.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "MinStake")
		return
	}
	// string "MaxStake"
	o = append(o, 0xa8, 0x4d, 0x61, 0x78, 0x53, 0x74, 0x61, 0x6b, 0x65)
	o, err = z.MaxStake.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "MaxStake")
		return
	}
	// string "MaxNumDelegates"
	o = append(o, 0xaf, 0x4d, 0x61, 0x78, 0x4e, 0x75, 0x6d, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x73)
	o = msgp.AppendInt(o, z.MaxNumDelegates)
	// string "ServiceChargeRatio"
	o = append(o, 0xb2, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x68, 0x61, 0x72, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6f)
	o = msgp.AppendFloat64(o, z.ServiceChargeRatio)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Settings) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "DelegateWallet":
			z.DelegateWallet, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "DelegateWallet")
				return
			}
		case "MinStake":
			bts, err = z.MinStake.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "MinStake")
				return
			}
		case "MaxStake":
			bts, err = z.MaxStake.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxStake")
				return
			}
		case "MaxNumDelegates":
			z.MaxNumDelegates, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxNumDelegates")
				return
			}
		case "ServiceChargeRatio":
			z.ServiceChargeRatio, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ServiceChargeRatio")
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
func (z *Settings) Msgsize() (s int) {
	s = 1 + 15 + msgp.StringPrefixSize + len(z.DelegateWallet) + 9 + z.MinStake.Msgsize() + 9 + z.MaxStake.Msgsize() + 16 + msgp.IntSize + 19 + msgp.Float64Size
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *StakePool) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Pools"
	o = append(o, 0x84, 0xa5, 0x50, 0x6f, 0x6f, 0x6c, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Pools)))
	keys_za0001 := make([]string, 0, len(z.Pools))
	for k := range z.Pools {
		keys_za0001 = append(keys_za0001, k)
	}
	msgp.Sort(keys_za0001)
	for _, k := range keys_za0001 {
		za0002 := z.Pools[k]
		o = msgp.AppendString(o, k)
		if za0002 == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = za0002.MarshalMsg(o)
			if err != nil {
				err = msgp.WrapError(err, "Pools", k)
				return
			}
		}
	}
	// string "Reward"
	o = append(o, 0xa6, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64)
	o, err = z.Reward.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Reward")
		return
	}
	// string "Settings"
	o = append(o, 0xa8, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73)
	o, err = z.Settings.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Settings")
		return
	}
	// string "Minter"
	o = append(o, 0xa6, 0x4d, 0x69, 0x6e, 0x74, 0x65, 0x72)
	o, err = z.Minter.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Minter")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StakePool) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Pools":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Pools")
				return
			}
			if z.Pools == nil {
				z.Pools = make(map[string]*DelegatePool, zb0002)
			} else if len(z.Pools) > 0 {
				for key := range z.Pools {
					delete(z.Pools, key)
				}
			}
			for zb0002 > 0 {
				var za0001 string
				var za0002 *DelegatePool
				zb0002--
				za0001, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Pools")
					return
				}
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					za0002 = nil
				} else {
					if za0002 == nil {
						za0002 = new(DelegatePool)
					}
					bts, err = za0002.UnmarshalMsg(bts)
					if err != nil {
						err = msgp.WrapError(err, "Pools", za0001)
						return
					}
				}
				z.Pools[za0001] = za0002
			}
		case "Reward":
			bts, err = z.Reward.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Reward")
				return
			}
		case "Settings":
			bts, err = z.Settings.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Settings")
				return
			}
		case "Minter":
			bts, err = z.Minter.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Minter")
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
func (z *StakePool) Msgsize() (s int) {
	s = 1 + 6 + msgp.MapHeaderSize
	if z.Pools != nil {
		for za0001, za0002 := range z.Pools {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001)
			if za0002 == nil {
				s += msgp.NilSize
			} else {
				s += za0002.Msgsize()
			}
		}
	}
	s += 7 + z.Reward.Msgsize() + 9 + z.Settings.Msgsize() + 7 + z.Minter.Msgsize()
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *StakePoolRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ProviderType"
	o = append(o, 0x82, 0xac, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65)
	o, err = z.ProviderType.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "ProviderType")
		return
	}
	// string "ProviderID"
	o = append(o, 0xaa, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x49, 0x44)
	o = msgp.AppendString(o, z.ProviderID)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StakePoolRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ProviderType":
			bts, err = z.ProviderType.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "ProviderType")
				return
			}
		case "ProviderID":
			z.ProviderID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ProviderID")
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
func (z *StakePoolRequest) Msgsize() (s int) {
	s = 1 + 13 + z.ProviderType.Msgsize() + 11 + msgp.StringPrefixSize + len(z.ProviderID)
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *StakePoolStat) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 8
	// string "ID"
	o = append(o, 0x88, 0xa2, 0x49, 0x44)
	o = msgp.AppendString(o, z.ID)
	// string "Balance"
	o = append(o, 0xa7, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65)
	o, err = z.Balance.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Balance")
		return
	}
	// string "StakeTotal"
	o = append(o, 0xaa, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x54, 0x6f, 0x74, 0x61, 0x6c)
	o, err = z.StakeTotal.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "StakeTotal")
		return
	}
	// string "UnstakeTotal"
	o = append(o, 0xac, 0x55, 0x6e, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x54, 0x6f, 0x74, 0x61, 0x6c)
	o, err = z.UnstakeTotal.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "UnstakeTotal")
		return
	}
	// string "Delegate"
	o = append(o, 0xa8, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Delegate)))
	for za0001 := range z.Delegate {
		o, err = z.Delegate[za0001].MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "Delegate", za0001)
			return
		}
	}
	// string "Penalty"
	o = append(o, 0xa7, 0x50, 0x65, 0x6e, 0x61, 0x6c, 0x74, 0x79)
	o, err = z.Penalty.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Penalty")
		return
	}
	// string "Rewards"
	o = append(o, 0xa7, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73)
	o, err = z.Rewards.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Rewards")
		return
	}
	// string "Settings"
	o = append(o, 0xa8, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73)
	o, err = z.Settings.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Settings")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StakePoolStat) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ID":
			z.ID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "Balance":
			bts, err = z.Balance.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Balance")
				return
			}
		case "StakeTotal":
			bts, err = z.StakeTotal.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "StakeTotal")
				return
			}
		case "UnstakeTotal":
			bts, err = z.UnstakeTotal.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "UnstakeTotal")
				return
			}
		case "Delegate":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Delegate")
				return
			}
			if cap(z.Delegate) >= int(zb0002) {
				z.Delegate = (z.Delegate)[:zb0002]
			} else {
				z.Delegate = make([]DelegatePoolStat, zb0002)
			}
			for za0001 := range z.Delegate {
				bts, err = z.Delegate[za0001].UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "Delegate", za0001)
					return
				}
			}
		case "Penalty":
			bts, err = z.Penalty.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Penalty")
				return
			}
		case "Rewards":
			bts, err = z.Rewards.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Rewards")
				return
			}
		case "Settings":
			bts, err = z.Settings.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Settings")
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
func (z *StakePoolStat) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID) + 8 + z.Balance.Msgsize() + 11 + z.StakeTotal.Msgsize() + 13 + z.UnstakeTotal.Msgsize() + 9 + msgp.ArrayHeaderSize
	for za0001 := range z.Delegate {
		s += z.Delegate[za0001].Msgsize()
	}
	s += 8 + z.Penalty.Msgsize() + 8 + z.Rewards.Msgsize() + 9 + z.Settings.Msgsize()
	return
}

// MarshalMsg implements msgp.Marshaler
func (z UserPoolStat) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 0
	o = append(o, 0x80)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *UserPoolStat) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
func (z UserPoolStat) Msgsize() (s int) {
	s = 1
	return
}
