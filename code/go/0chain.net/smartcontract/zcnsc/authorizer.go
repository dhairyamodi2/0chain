package zcnsc

import (
	"fmt"

	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/state"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	. "0chain.net/core/logging"
	"go.uber.org/zap"
)

// AddAuthorizer sc API function
// Transaction must include ClientID, ToClientID, PublicKey, Hash, Value
// inputData is a publicKey in case public key in Tx is missing
// Either PK or inputData must be present
// balances have `GetTriedNode` implemented to get nodes
// ContractMap contains all the SC addresses
// ToClient is an SC address
func (zcn *ZCNSmartContract) AddAuthorizer(tran *transaction.Transaction, inputData []byte, ctx cstate.StateContextI) (string, error) {
	var (
		//publicKey    = tran.PublicKey  // authorizer public key
		authorizerID = tran.ClientID   // sender address
		recipientID  = tran.ToClientID // smart contract address
		node         *AuthorizerNode
		err          error
		errorCode    = "failed to add authorizer"
	)

	if inputData == nil {
		msg := "input data is nil"
		err = common.NewError(errorCode, msg)
		Logger.Error(msg, zap.Error(err))
		return "", err
	}

	node, err = GetAuthorizerNode(authorizerID, ctx)
	if err == nil {
		msg := fmt.Sprintf("authorizer(authorizerID: %v) already exists: %v", authorizerID, err)
		err = common.NewError(errorCode, msg)
		Logger.Warn("get authorizer node", zap.Error(err))
		return "", err
	}

	gn, err := GetGlobalNode(ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to get global node, authorizer(authorizerID: %v), err: %v", authorizerID, err)
		err := common.NewError(errorCode, msg)
		Logger.Error("get global node", zap.Error(err))
		return "", err
	}

	// compare the global min of authorizerNode Authorizer to that of the transaction amount
	if gn.MinStakeAmount > tran.Value {
		msg := fmt.Sprintf("min stake amount '(%d)' > transaction value '(%d)'", gn.MinStakeAmount, tran.Value)
		err = common.NewError(errorCode, msg)
		Logger.Error("min stake amount > transaction value", zap.Error(err))
		return "", err
	}

	authParam := AuthorizerParameter{}
	err = authParam.Decode(inputData)
	if err != nil {
		err = common.NewError(errorCode, "failed to decode AuthorizerParameter")
		Logger.Error("public key error", zap.Error(err))
		return "", err
	}

	if authParam.PublicKey == "" {
		err = common.NewError(errorCode, "public key was not included with transaction")
		Logger.Error("public key error", zap.Error(err))
		return "", err
	}

	// Save authorizer
	// authorizerID = tran.ClientID
	node = CreateAuthorizer(authorizerID, authParam.PublicKey, authParam.URL)

	err = node.Save(ctx)
	if err != nil {
		msg := fmt.Sprintf("error saving authorizer(authorizerID: %v), err: %v", authorizerID, err)
		err = common.NewError(errorCode, msg)
		Logger.Error("saving authorizer node", zap.Error(err))
		return "", err
	}

	//Dig pool for authorizer

	transfer, response, err := node.Staking.DigPool(tran.Hash, tran)
	if err != nil {
		err = common.NewError(errorCode, fmt.Sprintf("error digging pool, err: (%v)", err))
		return "", err
	}

	err = ctx.AddTransfer(transfer)
	if err != nil {
		msg := "Error: '%v', transaction.ClientId: '%s', transaction.ToClientId: '%s', transfer.ClientID: '%s', transfer.ToClientID: '%s'"
		err = common.NewError(errorCode, fmt.Sprintf(msg, err, authorizerID, recipientID, transfer.ClientID, transfer.ToClientID))
		return "", err
	}

	return response, err
}

func (zcn *ZCNSmartContract) DeleteAuthorizer(tran *transaction.Transaction, _ []byte, ctx cstate.StateContextI) (string, error) {
	var (
		authorizerID = tran.ClientID
		node         *AuthorizerNode
		transfer     *state.Transfer
		response     string
		err          error
		errorCode    = "failed to delete authorizer"
	)

	node, err = GetAuthorizerNode(authorizerID, ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to get authorizer (authorizerID: %v), err: %v", authorizerID, err)
		err = common.NewError(errorCode, msg)
		Logger.Error("get authorizer node", zap.Error(err))
		return "", err
	}

	if node == nil {
		msg := fmt.Sprintf("authorizer (authorizerID: %v) not found, err: %v", authorizerID, err)
		err = common.NewError(errorCode, msg)
		Logger.Error("authorizer node not found", zap.Error(err))
		return "", err
	}

	gn, err := GetGlobalNode(ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to get global node (authorizerID: %v), err: %v", authorizerID, err)
		err = common.NewError(errorCode, msg)
		Logger.Error("get global node", zap.Error(err))
		return "", err
	}

	// empty the authorizer's pool
	pool := node.Staking
	if pool == nil {
		msg := "pool is not created"
		err := common.NewError(errorCode, msg)
		Logger.Error("node staking pool", zap.Error(err))
		return "", err
	}

	transfer, response, err = pool.EmptyPool(gn.ID, tran.ClientID, tran)
	if err != nil {
		msg := fmt.Sprintf("error emptying pool, err: (%v)", err)
		err = common.NewError(errorCode, msg)
		Logger.Error("empty pool", zap.Error(err))
		return response, err
	}

	// transfer tokens back to authorizer account
	err = ctx.AddTransfer(transfer)
	if err != nil {
		msg := fmt.Sprintf("error adding transfer: (%v)", err)
		err = common.NewError(errorCode, msg)
		Logger.Error("add transfer", zap.Error(err))
		return response, err
	}

	// delete authorizer node
	_, err = ctx.DeleteTrieNode(node.GetKey())
	if err != nil {
		msg := fmt.Sprintf("failed to delete authorizerID: (%v), node key: (%v), err: %v", authorizerID, node.GetKey(), err)
		err = common.NewError(errorCode, msg)
		Logger.Error("delete trie node", zap.Error(err))
		return "", err
	}

	Logger.Info(
		"deleted authorizer",
		zap.String("hash", tran.Hash),
		zap.String("authorizerID", authorizerID),
	)

	return response, err
}