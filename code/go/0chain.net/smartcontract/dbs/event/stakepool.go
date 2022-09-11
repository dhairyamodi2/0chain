package event

import (
	"fmt"
	"time"

	"0chain.net/chaincore/currency"
	"0chain.net/core/logging"
	"0chain.net/smartcontract/stakepool/spenum"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"0chain.net/smartcontract/dbs"
)

type providerAggregateStats struct {
	Rewards     currency.Coin `json:"rewards"`
	TotalReward currency.Coin `json:"total_reward"`
}

type providerRewardsDelegates struct {
	rewards []ProviderRewards

	delegateRewards   []DelegatePool
	delegatePenalties []DelegatePool
	desc              [][]string
}

func aggregateProviderRewards(spus []dbs.StakePoolReward) (*providerRewardsDelegates, error) {
	var (
		rewards           = make([]ProviderRewards, 0, len(spus))
		delegateRewards   = make([]DelegatePool, 0, len(spus))
		delegatePenalties = make([]DelegatePool, 0, len(spus))
		descs             = make([][]string, 0, len(spus))
	)

	for i, sp := range spus {
		descs = append(descs, sp.Desc)
		if sp.Reward != 0 {
			rewards = append(rewards, ProviderRewards{
				ProviderID:   sp.ProviderId,
				Rewards:      sp.Reward,
				TotalRewards: sp.Reward,
			})
		}

		for k, v := range spus[i].DelegateRewards {
			delegateRewards = append(delegateRewards, DelegatePool{
				ProviderID:   sp.ProviderId,
				ProviderType: sp.ProviderType,
				PoolID:       k,
				Reward:       currency.Coin(v),
				TotalReward:  currency.Coin(v),
			})
		}

		for k, v := range spus[i].DelegatePenalties {
			delegatePenalties = append(delegatePenalties, DelegatePool{
				ProviderID:   sp.ProviderId,
				ProviderType: sp.ProviderType,
				PoolID:       k,
				TotalPenalty: currency.Coin(v),
			})
		}
	}

	return &providerRewardsDelegates{
		rewards:           rewards,
		delegateRewards:   delegateRewards,
		delegatePenalties: delegatePenalties,
		desc:              descs,
	}, nil
}

func (edb *EventDb) rewardUpdate(spus []dbs.StakePoolReward) error {
	if len(spus) == 0 {
		return nil
	}

	ts := time.Now()
	rewards, err := aggregateProviderRewards(spus)
	if err != nil {
		return err
	}

	defer func() {
		du := time.Since(ts)
		n := len(rewards.rewards) + len(rewards.delegateRewards) + len(rewards.delegatePenalties)
		if du > 50*time.Millisecond {
			logging.Logger.Debug("event db - update reward slow",
				zap.Any("duration", du),
				zap.Int("update items", n),
				zap.Any("desc", rewards.desc))
		}
	}()

	if len(rewards.rewards) > 0 {
		if err := edb.rewardProviders(rewards.rewards); err != nil {
			return fmt.Errorf("could not rewards providers: %v", err)
		}
	}

	rpdu := time.Since(ts)
	if rpdu.Milliseconds() > 50 {
		logging.Logger.Debug("event db - reward providers slow", zap.Any("duration", rpdu))
	}

	if len(rewards.delegateRewards) > 0 {
		if err := rewardProviderDelegates(edb, rewards.delegateRewards); err != nil {
			return fmt.Errorf("could not rewards delegate pool: %v", err)
		}
	}

	if len(rewards.delegatePenalties) > 0 {
		if err := penaltyProviderDelegates(edb, rewards.delegatePenalties); err != nil {
			return fmt.Errorf("could not penalty delegate pool: %v", err)
		}
	}

	return nil
}

func rewardProvider[T any](edb *EventDb, tableName, index string, providers []T) error {
	vs := map[string]interface{}{
		"rewards":      gorm.Expr(fmt.Sprintf("%s.rewards + excluded.rewards", tableName)),
		"total_reward": gorm.Expr(fmt.Sprintf("%s.total_reward + excluded.total_reward", tableName)),
	}

	return edb.Store.Get().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: index}},
		DoUpdates: clause.Assignments(vs),
	}).Create(&providers).Error
}

func (edb *EventDb) rewardProviders(rewards []ProviderRewards) error {
	vs := map[string]interface{}{
		"rewards":       gorm.Expr(fmt.Sprintf("provider_rewards.rewards + excluded.rewards")),
		"total_rewards": gorm.Expr(fmt.Sprintf("provider_rewards.total_rewards + excluded.total_rewards")),
	}

	return edb.Store.Get().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider_id"}},
		DoUpdates: clause.Assignments(vs),
	}).Create(rewards).Error
}

func rewardProviderDelegates(edb *EventDb, rewards []DelegatePool) error {
	vs := map[string]interface{}{
		"reward":       gorm.Expr("delegate_pools.reward + excluded.reward"),
		"total_reward": gorm.Expr("delegate_pools.total_reward + excluded.total_reward"),
	}

	return edb.Store.Get().Clauses(clause.OnConflict{
		Where: clause.Where{
			Exprs: []clause.Expression{gorm.Expr("delegate_pools.status != ?", spenum.Deleted)},
		},
		Columns: []clause.Column{
			{Name: "provider_type"},
			{Name: "provider_id"},
			{Name: "pool_id"},
		},
		DoUpdates: clause.Assignments(vs),
	}).Create(&rewards).Error
}

func penaltyProviderDelegates(edb *EventDb, penalties []DelegatePool) error {
	vs := map[string]interface{}{
		"total_penalty": gorm.Expr("delegate_pools.total_penalty + excluded.total_penalty"),
	}

	return edb.Store.Get().Clauses(clause.OnConflict{
		Where: clause.Where{
			Exprs: []clause.Expression{gorm.Expr("delegate_pools.status != ?", spenum.Deleted)},
		},
		Columns: []clause.Column{
			{Name: "provider_type"},
			{Name: "provider_id"},
			{Name: "pool_id"},
		},
		DoUpdates: clause.Assignments(vs),
	}).Create(&penalties).Error
}

type rewardInfo struct {
	pool  string
	value int64
}

func (edb *EventDb) rewardProvider(spu dbs.StakePoolReward) error {
	if spu.Reward == 0 {
		return nil
	}

	var provider interface{}
	switch spenum.Provider(spu.ProviderType) {
	case spenum.Blobber:
		provider = &Blobber{BlobberID: spu.ProviderId}
	case spenum.Validator:
		provider = &Validator{ValidatorID: spu.ProviderId}
	case spenum.Miner:
		provider = &Miner{MinerID: spu.ProviderId}
	case spenum.Sharder:
		provider = &Sharder{SharderID: spu.ProviderId}
	default:
		return fmt.Errorf("not implented provider type %v", spu.ProviderType)
	}

	vs := map[string]interface{}{
		"rewards":      gorm.Expr("rewards + ?", spu.Reward),
		"total_reward": gorm.Expr("total_reward + ?", spu.Reward),
	}

	return edb.Store.Get().Model(provider).Where(provider).Updates(vs).Error
}
