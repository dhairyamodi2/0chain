package sharder

import (
	"context"
	"math"
	"sort"
	"strconv"
	"time"

	"0chain.net/node"
	"0chain.net/round"
	"0chain.net/transaction"
	"0chain.net/util"
	metrics "github.com/rcrowley/go-metrics"

	"0chain.net/blockstore"
	"0chain.net/config"

	"0chain.net/block"
	"0chain.net/datastore"
	. "0chain.net/logging"
	"go.uber.org/zap"
)

var blockSaveTimer metrics.Timer

func init() {
	blockSaveTimer = metrics.GetOrRegisterTimer("block_save_time", nil)
}

/*UpdatePendingBlock - update the pending block */
func (sc *Chain) UpdatePendingBlock(ctx context.Context, b *block.Block, txns []datastore.Entity) {

}

/*UpdateFinalizedBlock - updates the finalized block */
func (sc *Chain) UpdateFinalizedBlock(ctx context.Context, b *block.Block) {
	fr := sc.GetRound(b.Round)
	Logger.Info("update finalized block", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.Any("lf_round", sc.LatestFinalizedBlock.Round), zap.Any("current_round", sc.CurrentRound))
	if config.Development() {
		for _, t := range b.Txns {
			if !t.DebugTxn() {
				continue
			}
			Logger.Info("update finalized block (debug transaction)", zap.String("txn", t.Hash), zap.String("block", b.Hash))
		}
	}
	sc.BlockCache.Add(b.Hash, b)
	sc.cacheBlockTxns(b.Hash, b.Txns)
	sc.StoreTransactions(ctx, b)
	err := sc.StoreBlockSummary(ctx, b)
	if err != nil {
		Logger.Error("db error (save block)", zap.Any("round", b.Round), zap.String("block", b.Hash), zap.Error(err))
	}
	self := node.GetSelfNode(ctx)
	if sc.IsBlockSharder(b, self.Node) {
		sc.storeBlock(b)
	}
	if fr != nil {
		sc.storeRound(ctx, fr, b)
	}
	sc.DeleteRoundsBelow(ctx, b.Round)
}

func (sc *Chain) cacheBlockTxns(hash string, txns []*transaction.Transaction) {
	for _, txn := range txns {
		txnSummary := txn.GetSummary()
		txnSummary.BlockHash = hash
		sc.BlockTxnCache.Add(txn.Hash, txnSummary)
	}
}

func (sc *Chain) processBlock(ctx context.Context, b *block.Block) {
	if err := sc.VerifyNotarization(ctx, b.Hash, b.VerificationTickets); err != nil {
		Logger.Error("notarization verification failed", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.Error(err))
		return
	}
	if err := b.Validate(ctx); err != nil {
		Logger.Error("block validation", zap.Any("round", b.Round), zap.Any("hash", b.Hash), zap.Error(err))
		return
	}
	er := sc.GetRound(b.Round)
	if er == nil {
		var r = round.NewRound(b.Round)
		er, _ = sc.AddRound(r).(*round.Round)
		sc.SetRandomSeed(er, b.RoundRandomSeed)
	}
	if sc.AddRoundBlock(er, b) != b {
		return
	}
	sc.SetRoundRank(er, b)
	Logger.Info("received block", zap.Int64("round", b.Round), zap.String("block", b.Hash), zap.String("client_state", util.ToHex(b.ClientStateHash)))
	sc.AddNotarizedBlock(ctx, er, b)
}

func (sc *Chain) CheckForMissingRounds(ctx context.Context, currRound int64) *round.Round {
	latestRounds := make([]*round.Round, 1)

	latestRoundHandler := func(ctx context.Context, entity datastore.Entity) (interface{}, error) {
		r, ok := entity.(*round.Round)
		if !ok {
			return nil, nil
		}
		latestRounds = append(latestRounds, r)
		return nil, nil
	}

	sc.Sharders.RequestEntityFromAll(ctx, LatestRoundRequestor, nil, latestRoundHandler)
	sort.Slice(latestRounds, func(i int, j int) bool { return latestRounds[i].Number >= latestRounds[j].Number })

	if len(latestRounds) > 0 {
		return latestRounds[0]
	}

	return nil
}

func (sc *Chain) GetMissingRounds(ctx context.Context, r int64, currRound int64) {
	currRound += 1
	blocks := make([]*block.Block, 1)
	var params map[string]string

	missingRoundBlocksHandler := func(ctx context.Context, entity datastore.Entity) (interface{}, error) {
		rb, ok := entity.(*block.Block)
		if !ok {
			return nil, nil
		}
		blocks = append(blocks, rb)
		return nil, nil
	}

	for currRound != r {
		params["round"] = strconv.FormatInt(currRound, 10)
		r := round.NewRound(currRound)
		sc.Sharders.RequestEntityFromAll(ctx, BlockRequestor, params, missingRoundBlocksHandler)
		//TODO any consensus to be done?? currently taking the first block that we received
		if len(blocks) > 0 {
			b := blocks[0]
			sc.storeRound(ctx, r, b)
			self := node.GetSelfNode(ctx)
			Logger.Info("bc-27 missed round stored in db", zap.Int64("round", r.GetRoundNumber()))
			if sc.IsBlockSharder(b, self.Node) {
				sc.storeBlock(b)
				Logger.Info("bc-27 missed block stored in db", zap.String("block-hash", b.Hash))
			}
			blocks = make([]*block.Block, 1)
		}
		currRound += 1
	}
}

func (sc *Chain) storeRound(ctx context.Context, r round.RoundI, b *block.Block) {
	r.Finalize(b)
	rImpl, _ := r.(*round.Round)
	err := sc.StoreRound(ctx, rImpl)
	Logger.Info("**!round stored in db", zap.Int64("round", r.GetRoundNumber()))
	if err != nil {
		Logger.Error("db error (save round)", zap.Int64("round", r.GetRoundNumber()), zap.Error(err))
	}
}

func (sc *Chain) storeBlock(b *block.Block) {
	sc.SharderStats.ShardedBlocksCount++
	ts := time.Now()
	err := blockstore.GetStore().Write(b)
	duration := time.Since(ts)
	blockSaveTimer.UpdateSince(ts)
	p95 := blockSaveTimer.Percentile(.95)
	if blockSaveTimer.Count() > 100 && 2*p95 < float64(duration) {
		Logger.Error("block save - slow", zap.Any("round", b.Round), zap.String("block", b.Hash), zap.Duration("duration", duration), zap.Duration("p95", time.Duration(math.Round(p95/1000000))*time.Millisecond))
	}
	if err != nil {
		Logger.Error("block save", zap.Any("round", b.Round), zap.Any("hash", b.Hash), zap.Error(err))
	}
}
