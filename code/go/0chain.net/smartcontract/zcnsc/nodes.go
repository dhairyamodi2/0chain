package zcnsc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"0chain.net/chaincore/tokenpool"
	"0chain.net/smartcontract/dbs/event"
	"gorm.io/gorm"

	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/state"
	"0chain.net/core/common"
	"0chain.net/core/datastore"
	"0chain.net/core/encryption"
	"0chain.net/core/util"
	"0chain.net/smartcontract"
)

//msgp:ignore AuthorizerNode
//go:generate msgp -v -io=false -tests=false -unexported

// ------------- GlobalNode ------------------------

type GlobalNode struct {
	ID                 string         `json:"id"`
	MinMintAmount      state.Balance  `json:"min_mint_amount"`
	MinBurnAmount      state.Balance  `json:"min_burn_amount"`
	MinStakeAmount     state.Balance  `json:"min_stake_amount"`
	MaxFee             state.Balance  `json:"max_fee"`
	PercentAuthorizers float64        `json:"percent_authorizers"`
	MinAuthorizers     int64          `json:"min_authorizers"`
	BurnAddress        string         `json:"burn_address"`
	OwnerId            string         `json:"owner_id"`
	Cost               map[string]int `json:"cost"`
}

func (gn *GlobalNode) UpdateConfig(cfg *smartcontract.StringMap) error {
	var (
		err error
	)

	for key, value := range cfg.Fields {
		switch key {
		case MinMintAmount:
			amount, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("key %s, unable to convert %v to state.Balance", key, value)
			}
			gn.MinMintAmount = state.Balance(amount * 1e10)
		case MinBurnAmount:
			amount, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("key %s, unable to convert %v to state.Balance", key, value)
			}
			gn.MinBurnAmount = state.Balance(amount * 1e10)
		case BurnAddress:
			if value == "" {
				return fmt.Errorf("key %s is empty", key)
			}
			gn.BurnAddress = value
		case PercentAuthorizers:
			gn.PercentAuthorizers, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("key %s, unable to convert %v to float64", key, value)
			}
		case MinAuthorizers:
			gn.MinAuthorizers, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("key %s, unable to convert %v to int64", key, value)
			}
		case MinStakeAmount:
			amount, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("key %s, unable to convert %v to state.Balance", key, value)
			}
			gn.MinStakeAmount = state.Balance(amount * 1e10)
		case MaxFee:
			amount, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("key %s, unable to convert %v to state.Balance", key, value)
			}
			gn.MaxFee = state.Balance(amount * 1e10)
		case OwnerID:
			gn.OwnerId = value
		default:
			return gn.setCostValue(key, value)
		}
	}

	return nil
}

func (gn *GlobalNode) setCostValue(key, value string) error {
	if !strings.HasPrefix(key, fmt.Sprintf("%s.", Cost)) {
		return fmt.Errorf("key %s not recognised as setting", key)
	}

	costKey := strings.ToLower(strings.TrimPrefix(key, fmt.Sprintf("%s.", Cost)))
	for _, costFunction := range CostFunctions {
		if costKey != strings.ToLower(costFunction) {
			continue
		}
		costValue, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("key %s, unable to convert %v to integer", key, value)
		}

		if costValue < 0 {
			return fmt.Errorf("cost.%s contains invalid value %s", key, value)
		}

		gn.Cost[costKey] = costValue

		return nil
	}

	return fmt.Errorf("cost config setting %s not found", costKey)
}

func (gn *GlobalNode) Validate() error {
	const (
		Code = "failed to validate global node"
	)

	switch {
	case gn.MinStakeAmount < 1:
		return common.NewError(Code, fmt.Sprintf("min stake amount (%v) is less than 1", gn.MinStakeAmount))
	case gn.MinMintAmount < 1:
		return common.NewError(Code, fmt.Sprintf("min mint amount (%v) is less than 1", gn.MinMintAmount))
	case gn.MaxFee < 1:
		return common.NewError(Code, fmt.Sprintf("max fee (%v) is less than 1", gn.MaxFee))
	case gn.MinAuthorizers < 20:
		return common.NewError(Code, fmt.Sprintf("min quantity of authorizers (%v) is less than 20", gn.MinAuthorizers))
	case gn.MinBurnAmount < 1:
		return common.NewError(Code, fmt.Sprintf("min burn amount (%v) is less than 1", gn.MinBurnAmount))
	case gn.PercentAuthorizers < 70:
		return common.NewError(Code, fmt.Sprintf("min percentage of authorizers (%v) is less than 70", gn.PercentAuthorizers))
	case gn.BurnAddress == "":
		return common.NewError(Code, fmt.Sprintf("burn address (%v) is not valid", gn.BurnAddress))
	case gn.OwnerId == "":
		return common.NewError(Code, fmt.Sprintf("owner id (%v) is not valid", gn.OwnerId))
	}
	return nil
}

func (gn *GlobalNode) GetKey() datastore.Key {
	return fmt.Sprintf("%s:%s:%s", ADDRESS, GlobalNodeType, gn.ID)
}

func (gn *GlobalNode) GetHash() string {
	return util.ToHex(gn.GetHashBytes())
}

func (gn *GlobalNode) GetHashBytes() []byte {
	return encryption.RawHash(gn.Encode())
}

func (gn *GlobalNode) Encode() []byte {
	buff, _ := json.Marshal(gn)
	return buff
}

func (gn *GlobalNode) Decode(input []byte) error {
	err := json.Unmarshal(input, gn)
	return err
}

func (gn *GlobalNode) Save(balances cstate.StateContextI) (err error) {
	_, err = balances.InsertTrieNode(gn.GetKey(), gn)
	return
}

// ----- AuthorizerConfig --------------------

type AuthorizerConfig struct {
	Fee state.Balance `json:"fee"`
}

func (c *AuthorizerConfig) Decode(input []byte) (err error) {
	err = json.Unmarshal(input, c)
	return
}

// ----- AuthorizerNode --------------------

type AuthorizerNode struct {
	ID        string                    `json:"id"`
	PublicKey string                    `json:"public_key"`
	Staking   *tokenpool.ZcnLockingPool `json:"staking"`
	URL       string                    `json:"url"`
	Config    *AuthorizerConfig         `json:"config"`
}

// NewAuthorizer To review: tokenLock init values
// PK = authorizer node public key
// ID = authorizer node public id = Client ID
func NewAuthorizer(ID string, PK string, URL string) *AuthorizerNode {
	return &AuthorizerNode{
		ID:        ID,
		PublicKey: PK,
		URL:       URL,
		Staking: &tokenpool.ZcnLockingPool{
			ZcnPool: tokenpool.ZcnPool{
				TokenPool: tokenpool.TokenPool{
					ID:      "", // must be filled when DigPool is invoked. Usually this is a trx.Hash
					Balance: 0,  // filled when we dig pool
				},
			},
			TokenLockInterface: &TokenLock{
				StartTime: 0,
				Duration:  0,
				Owner:     ID,
			},
		},
		Config: &AuthorizerConfig{
			Fee: 0,
		},
	}
}

func (an *AuthorizerNode) UpdateConfig(cfg *AuthorizerConfig) error {
	if cfg == nil {
		return errors.New("config not initialized")
	}
	an.Config = cfg

	return nil
}

func (an *AuthorizerNode) GetKey() string {
	return fmt.Sprintf("%s:%s:%s", ADDRESS, AuthorizerNodeType, an.ID)
}

func (an *AuthorizerNode) Encode() []byte {
	bytes, _ := json.Marshal(an)
	return bytes
}

func (an *AuthorizerNode) Decode(input []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(input, &objMap)
	if err != nil {
		return err
	}

	id, ok := objMap["id"]
	if ok {
		var idStr *string
		err = json.Unmarshal(*id, &idStr)
		if err != nil {
			return err
		}
		an.ID = *idStr
	}

	pk, ok := objMap["public_key"]
	if ok {
		var pkStr *string
		err = json.Unmarshal(*pk, &pkStr)
		if err != nil {
			return err
		}
		an.PublicKey = *pkStr
	}

	url, ok := objMap["url"]
	if ok {
		var urlStr *string
		err = json.Unmarshal(*url, &urlStr)
		if err != nil {
			return err
		}
		an.URL = *urlStr
	}

	if an.Staking == nil {
		an.Staking = &tokenpool.ZcnLockingPool{
			ZcnPool: tokenpool.ZcnPool{
				TokenPool: tokenpool.TokenPool{},
			},
		}
	}

	staking, ok := objMap["staking"]
	if ok && staking != nil {
		tokenlock := &TokenLock{}
		err = an.Staking.Decode(*staking, tokenlock)
		if err != nil {
			return err
		}
	}

	rawCfg, ok := objMap["config"]
	if ok {
		var cfg = &AuthorizerConfig{}
		err = cfg.Decode(*rawCfg)
		if err != nil {
			return err
		}

		an.Config = cfg
	}

	return nil
}

type authorizerNodeDecode AuthorizerNode

func (an *AuthorizerNode) MarshalMsg(o []byte) ([]byte, error) {
	d := authorizerNodeDecode(*an)
	return d.MarshalMsg(o)
}

func (an *AuthorizerNode) UnmarshalMsg(data []byte) ([]byte, error) {
	d := authorizerNodeDecode{Staking: &tokenpool.ZcnLockingPool{TokenLockInterface: &TokenLock{}}}
	o, err := d.UnmarshalMsg(data)
	if err != nil {
		return nil, err
	}

	*an = AuthorizerNode(d)
	return o, nil
}

func (an *AuthorizerNode) Save(ctx cstate.StateContextI) (err error) {
	_, err = ctx.InsertTrieNode(an.GetKey(), an)
	if err != nil {
		return common.NewError("save_auth_node_failed", "saving authorizer node: "+err.Error())
	}
	return nil
}

func (an *AuthorizerNode) ToEvent() ([]byte, error) {
	if an.Config == nil {
		an.Config = new(AuthorizerConfig)
	}
	data, err := json.Marshal(&event.Authorizer{
		Model:        gorm.Model{},
		Fee:          an.Config.Fee,
		AuthorizerID: an.ID,
		URL:          an.URL,
	})
	if err != nil {
		return nil, fmt.Errorf("marshalling authorizer event: %v", err)
	}

	return data, nil
}

func AuthorizerFromEvent(buf []byte) (*AuthorizerNode, error) {
	ev := &event.Authorizer{}
	err := json.Unmarshal(buf, ev)
	if err != nil {
		return nil, err
	}

	return &AuthorizerNode{
		ID:        ev.AuthorizerID,
		URL:       ev.URL,
		PublicKey: "",  // fetch this from MPT
		Staking:   nil, // fetch this from MPT
	}, nil
}

// ----- UserNode ------------------

type UserNode struct {
	ID    string `json:"id"`
	Nonce int64  `json:"nonce"`
}

func NewUserNode(id string, nonce int64) *UserNode {
	return &UserNode{
		ID:    id,
		Nonce: nonce,
	}
}

func (un *UserNode) GetKey() datastore.Key {
	return fmt.Sprintf("%s:%s:%s", ADDRESS, UserNodeType, un.ID)
}

func (un *UserNode) GetHash() string {
	return util.ToHex(un.GetHashBytes())
}

func (un *UserNode) GetHashBytes() []byte {
	return encryption.RawHash(un.Encode())
}

func (un *UserNode) Encode() []byte {
	buff, _ := json.Marshal(un)
	return buff
}

func (un *UserNode) Decode(input []byte) error {
	err := json.Unmarshal(input, un)
	return err
}

func (un *UserNode) Save(balances cstate.StateContextI) (err error) {
	_, err = balances.InsertTrieNode(un.GetKey(), un)
	return
}