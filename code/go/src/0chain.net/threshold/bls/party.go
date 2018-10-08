package bls

import (
	"fmt"
)

/* BLS implementation */

type SimpleBLS struct {
	t             int
	n             int
	msg           Message
	sigShare      Sign
	gpPubKey      GroupPublicKey
	verifications []VerificationKey

	partyKeyShare AfterDKGKeyShare
}

func MakeSimpleBLS(dkg *BLSSimpleDKG, mg Message) SimpleBLS {

	return SimpleBLS{
		t:             dkg.T,
		n:             dkg.N,
		msg:           mg,
		sigShare:      Sign{},
		gpPubKey:      GroupPublicKey,
		verifications: nil,
		partyKeyShare: AfterDKGKeyShare{},
	}

}

func (bs *SimpleBLS) SignMsg() Sign {

	priKeyShare := bs.partyKeyShare.m
	sigShare := *priKeyShare.Sign(bs.msg)
	return sigShare
}

/* VerifySign - Verifies the signature share with the public key share */
func (bs *SimpleBLS) VerifySign(sigShare Sign) bool {

	pubKeyShare := bs.partyKeyShare.v
	if !sigShare.Verify(&pubKeyShare, bs.msg) {
		fmt.Println("Signature does not verify")
		return false
	}
	return true
}

func (bs *SimpleBLS) RecoverGroupSig(from []PartyId, shares []Sign) interface{} {

	signVec := shares
	idVec := from

	bs.sigShare.Recover(signVec, idVec)

	return bs.sigShare
}

func (bs *SimpleBLS) VerifyGroupSig(GroupSig) bool {
	//TODO
	return true
}