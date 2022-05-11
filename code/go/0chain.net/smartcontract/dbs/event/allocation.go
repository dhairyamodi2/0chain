package event

import (
	"errors"
	"fmt"
	"time"

	"0chain.net/chaincore/state"
	"gorm.io/gorm"
)

type Allocation struct {
	gorm.Model
	AllocationID               string        `json:"allocation_id" gorm:"uniqueIndex"`
	AllocationName             string        `json:"allocation_name" gorm:"column:allocation_name;size:64;"`
	TransactionID              string        `json:"transaction_id"`
	DataShards                 int           `json:"data_shards"`
	ParityShards               int           `json:"parity_shards"`
	Size                       int64         `json:"size"`
	Expiration                 int64         `json:"expiration"`
	Terms                      string        `json:"terms"`
	Owner                      string        `json:"owner"`
	OwnerPublicKey             string        `json:"owner_public_key"`
	IsImmutable                bool          `json:"is_immutable"`
	ReadPriceMin               state.Balance `json:"read_price_min"`
	ReadPriceMax               state.Balance `json:"read_price_max"`
	WritePriceMin              state.Balance `json:"write_price_min"`
	WritePriceMax              state.Balance `json:"write_price_max"`
	MaxChallengeCompletionTime int64         `json:"max_challenge_completion_time"`
	ChallengeCompletionTime    int64         `json:"challenge_completion_time"`
	StartTime                  int64         `json:"start_time"`
	Finalized                  bool          `json:"finalized"`
	Cancelled                  bool          `json:"cancelled"`
	UsedSize                   int64         `json:"used_size"`
	MovedToChallenge           state.Balance `json:"moved_to_challenge"`
	MovedBack                  state.Balance `json:"moved_back"`
	MovedToValidators          state.Balance `json:"moved_to_validators"`
	TimeUnit                   int64         `json:"time_unit"`
	NumWrites                  int64         `json:"num_writes"`
	NumReads                   int64         `json:"num_reads"`
	TotalChallenges            int64         `json:"total_challenges"`
	OpenChallenges             int64         `json:"open_challenges"`
	SuccessfulChallenges       int64         `json:"successful_challenges"`
	FailedChallenges           int64         `json:"failed_challenges"`
	LatestClosedChallengeTxn   string        `json:"latest_closed_challenge_txn"`
}

type AllocationTerm struct {
	BlobberID    string `json:"blobber_id"`
	AllocationID string `json:"allocation_id"`
	// ReadPrice is price for reading. Token / GB (no time unit).
	ReadPrice state.Balance `json:"read_price"`
	// WritePrice is price for reading. Token / GB / time unit. Also,
	// it used to calculate min_lock_demand value.
	WritePrice state.Balance `json:"write_price"`
	// MinLockDemand in number in [0; 1] range. It represents part of
	// allocation should be locked for the blobber rewards even if
	// user never write something to the blobber.
	MinLockDemand float64 `json:"min_lock_demand"`
	// MaxOfferDuration with this prices and the demand.
	MaxOfferDuration time.Duration `json:"max_offer_duration"`
	// ChallengeCompletionTime is duration required to complete a challenge.
	ChallengeCompletionTime time.Duration `json:"challenge_completion_time"`
}

func (edb EventDb) GetAllocation(id string) (*Allocation, error) {
	var alloc Allocation
	err := edb.Store.Get().Model(&Allocation{}).Where("allocation_id = ?", id).First(&alloc).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving allocation: %v, error: %v", id, err)
	}

	return &alloc, nil
}

func (edb EventDb) GetClientsAllocation(clientID string) ([]Allocation, error) {
	allocs := make([]Allocation, 0)
	result := edb.Store.Get().Model(&Allocation{}).Where("owner = ?", clientID).Find(&allocs)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving allocation for client: %v, error: %v", clientID, result.Error)
	}

	return allocs, nil
}

func (edb EventDb) GetActiveAllocationsCount() (int64, error) {
	var count int64
	result := edb.Store.Get().Model(&Allocation{}).Where("finalized = ? AND cancelled = ?", false, false).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("error retrieving active allocations , error: %v", result.Error)
	}

	return count, nil
}

func (edb EventDb) GetActiveAllocsBlobberCount() (int64, error) {
	var count int64
	err := edb.Store.Get().
		Raw("SELECT SUM(parity_shards) + SUM(data_shards) FROM allocations WHERE finalized = ? AND cancelled = ?",
			false, false).
		Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("error retrieving blobber allocations count, error: %v", err)
	}

	return count, nil
}

func (edb *EventDb) overwriteAllocation(alloc *Allocation) error {
	return edb.Store.Get().Model(&Allocation{}).Where("allocation_id = ?", alloc.AllocationID).Updates(alloc).Error
}

func (edb *EventDb) addOrOverwriteAllocation(alloc *Allocation) error {
	exists, err := alloc.exists(edb)
	if err != nil {
		return err
	}

	if exists {
		return edb.overwriteAllocation(alloc)
	}

	return edb.Store.Get().Create(&alloc).Error
}

func (alloc *Allocation) exists(edb *EventDb) (bool, error) {
	var data Allocation
	err := edb.Store.Get().Model(&Allocation{}).Where("allocation_id = ?", alloc.AllocationID).Take(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("error searching for allocation %v, error %v", alloc.AllocationID, err)
	}

	return true, nil
}