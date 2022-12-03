package event

type (
	EventType int
	EventTag  int
)

const (
	TypeNone EventType = iota
	TypeError
	TypeChain
	TypeStats
	NumberOfTypes
)

func (t EventType) String() string {
	if int(t) < len(TagString) && int(t) >= 0 {
		return TagString[t]
	}

	return "unknown tag"
}

func (t EventType) Int() int {
	return int(t)
}

const (
	TagNone EventTag = iota
	TagAddBlobber
	TagUpdateBlobber
	TagUpdateBlobberAllocatedHealth
	TagUpdateBlobberTotalStake
	TagUpdateBlobberTotalOffers
	TagDeleteBlobber
	TagAddAuthorizer
	TagUpdateAuthorizer
	TagDeleteAuthorizer
	TagAddTransactions
	TagAddOrOverwriteUser
	TagAddWriteMarker
	TagAddBlock
	TagAddOrOverwiteValidator
	TagUpdateValidator
	TagAddReadMarker
	TagAddOrOverwriteMiner
	TagUpdateMiner
	TagDeleteMiner
	TagAddOrOverwriteSharder
	TagUpdateSharder
	TagDeleteSharder
	TagAddOrOverwriteCurator
	TagRemoveCurator
	TagStakePoolReward
	TagAddOrOverwriteDelegatePool
	TagUpdateDelegatePool
	TagAddAllocation
	TagUpdateAllocationStakes
	TagUpdateAllocation
	TagMintReward
	TagAddChallenge
	TagUpdateChallenge
	TagUpdateBlobberChallenge
	TagUpdateAllocationChallenge
	TagAddChallengeToAllocation
	TagAddOrOverwriteAllocationBlobberTerm
	TagUpdateAllocationBlobberTerm
	TagDeleteAllocationBlobberTerm
	TagAddOrUpdateChallengePool
	TagUpdateAllocationStat
	TagUpdateBlobberStat
	TagCollectProviderReward
	TagSendTransfer
	TagReceiveTransfer
	TagLockStakePool
	TagUnlockStakePool
	TagLockWritePool
	TagUnlockWritePool
	TagLockReadPool
	TagUnlockReadPool
	TagToChallengePool
	TagFromChallengePool
	TagAddMint
	TagBurn
	TagAllocValueChange
	TagAllocBlobberValueChange
	TagUpdateBlobberOpenChallenges
	TagUpdateValidatorStakeTotal
	NumberOfTags
)

var (
	TypeString []string
	TagString  []string
)

func init() {
	initTypeString()
	initTagString()
}

func initTypeString() {
	TypeString = make([]string, NumberOfTypes+1)
	TypeString[TypeNone] = "none"
	TypeString[TypeError] = "error"
	TypeString[TypeChain] = "chain"
	TypeString[TypeStats] = "stats"
	TypeString[NumberOfTypes] = "invalid"
}

func initTagString() {
	TagString = make([]string, NumberOfTags+1)
	TagString[TagNone] = "none"
	TagString[TagAddBlobber] = "TagAddBlobber"
	TagString[TagUpdateBlobber] = "TagUpdateBlobber"
	TagString[TagUpdateBlobberAllocatedHealth] = "TagUpdateBlobberAllocatedHealth"
	TagString[TagUpdateBlobberTotalStake] = "TagUpdateBlobberTotalStake"
	TagString[TagUpdateBlobberTotalOffers] = "TagUpdateBlobberTotalOffers"
	TagString[TagDeleteBlobber] = "TagDeleteBlobber"
	TagString[TagAddAuthorizer] = "TagAddAuthorizer"
	TagString[TagUpdateAuthorizer] = "TagUpdateAuthorizer"
	TagString[TagDeleteAuthorizer] = "TagDeleteAuthorizer"
	TagString[TagAddTransactions] = "TagAddTransactions"
	TagString[TagAddOrOverwriteUser] = "TagAddOrOverwriteUser"
	TagString[TagAddWriteMarker] = "TagAddWriteMarker"
	TagString[TagAddBlock] = "TagAddBlock"
	TagString[TagAddOrOverwiteValidator] = "TagAddOrOverwiteValidator"
	TagString[TagUpdateValidator] = "TagUpdateValidator"
	TagString[TagAddReadMarker] = "TagAddReadMarker"
	TagString[TagAddOrOverwriteMiner] = "TagAddOrOverwriteMiner"
	TagString[TagUpdateMiner] = "TagUpdateMiner"
	TagString[TagDeleteMiner] = "TagDeleteMiner"
	TagString[TagAddOrOverwriteSharder] = "TagAddOrOverwriteSharder"
	TagString[TagUpdateSharder] = "TagUpdateSharder"
	TagString[TagDeleteSharder] = "TagDeleteSharder"
	TagString[TagAddOrOverwriteCurator] = "TagAddOrOverwriteCurator"
	TagString[TagRemoveCurator] = "TagRemoveCurator"
	TagString[TagStakePoolReward] = "TagStakePoolReward"
	TagString[TagAddOrOverwriteDelegatePool] = "TagAddOrOverwriteDelegatePool"
	TagString[TagUpdateDelegatePool] = "TagUpdateDelegatePool"
	TagString[TagAddAllocation] = "TagAddAllocation"
	TagString[TagUpdateAllocationStakes] = "TagUpdateAllocationStakes"
	TagString[TagUpdateAllocation] = "TagUpdateAllocation"
	TagString[TagMintReward] = "TagMintReward"
	TagString[TagAddChallenge] = "TagAddChallenge"
	TagString[TagUpdateChallenge] = "TagUpdateChallenge"
	TagString[TagUpdateBlobberChallenge] = "TagUpdateBlobberChallenge"
	TagString[TagUpdateAllocationChallenge] = "TagUpdateAllocationChallenge"
	TagString[TagAddChallengeToAllocation] = "TagAddChallengeToAllocation"
	TagString[TagAddOrOverwriteAllocationBlobberTerm] = "TagAddOrOverwriteAllocationBlobberTerm"
	TagString[TagUpdateAllocationBlobberTerm] = "TagUpdateAllocationBlobberTerm"
	TagString[TagDeleteAllocationBlobberTerm] = "TagDeleteAllocationBlobberTerm"
	TagString[TagAddOrUpdateChallengePool] = "TagAddOrUpdateChallengePool"
	TagString[TagUpdateAllocationStat] = "TagUpdateAllocationStat"
	TagString[TagUpdateBlobberStat] = "TagUpdateBlobberStat"
	TagString[TagCollectProviderReward] = "TagCollectProviderReward"
	TagString[TagSendTransfer] = "TagSendTransfer"
	TagString[TagReceiveTransfer] = "TagReceiveTransfer"
	TagString[TagLockStakePool] = "TagLockStakePool"
	TagString[TagUnlockStakePool] = "TagUnlockStakePool"
	TagString[TagLockWritePool] = "TagLockWritePool"
	TagString[TagUnlockWritePool] = "TagUnlockWritePool"
	TagString[TagLockReadPool] = "TagLockReadPool"
	TagString[TagUnlockReadPool] = "TagUnlockReadPool"
	TagString[TagToChallengePool] = "TagToChallengePool"
	TagString[TagFromChallengePool] = "TagFromChallengePool"
	TagString[TagAddMint] = "TagAddMint"
	TagString[TagBurn] = "TagBurn"
	TagString[TagAllocValueChange] = "TagAllocValueChange"
	TagString[TagAllocBlobberValueChange] = "TagAllocBlobberValueChange"
	TagString[TagUpdateBlobberOpenChallenges] = "TagUpdateBlobberOpenChallenges"
	TagString[TagUpdateValidatorStakeTotal] = "TagUpdateValidatorStakeTotal"
	TagString[NumberOfTags] = "invalid"
}

func (tag EventTag) String() string {
	if int(tag) < len(TagString) && int(tag) >= 0 {
		return TagString[tag]
	}
	return "unknown tag"
}

func (tag EventTag) Int() int {
	return int(tag)
}