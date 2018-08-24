package util

import (
	"bytes"
	"context"

	"0chain.net/config"
	. "0chain.net/logging"
	"github.com/tecbot/gorocksdb"
	"go.uber.org/zap"
)

/*PNodeDB - a node db that is persisted */
type PNodeDB struct {
	dataDir string
	db      *gorocksdb.DB
	ro      *gorocksdb.ReadOptions
	wo      *gorocksdb.WriteOptions
	to      *gorocksdb.TransactionOptions
	fo      *gorocksdb.FlushOptions
}

/*NewPNodeDB - create a new PNodeDB */
func NewPNodeDB(dataDir string) (*PNodeDB, error) {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	bbto.SetFilterPolicy(gorocksdb.NewBloomFilter(10))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	opts.OptimizeForPointLookup(64)
	db, err := gorocksdb.OpenDb(opts, dataDir)
	if err != nil {
		return nil, err
	}
	pnodedb := &PNodeDB{db: db}
	pnodedb.dataDir = dataDir
	pnodedb.ro = gorocksdb.NewDefaultReadOptions()
	pnodedb.wo = gorocksdb.NewDefaultWriteOptions()
	pnodedb.wo.SetSync(true)
	pnodedb.to = gorocksdb.NewDefaultTransactionOptions()
	pnodedb.fo = gorocksdb.NewDefaultFlushOptions()
	return pnodedb, nil
}

/*GetNode - implement interface */
func (pndb *PNodeDB) GetNode(key Key) (Node, error) {
	data, err := pndb.db.Get(pndb.ro, key)
	if err != nil {
		return nil, err
	}
	defer data.Free()
	buf := data.Data()
	if buf == nil || len(buf) == 0 {
		return nil, ErrNodeNotFound
	}
	return CreateNode(bytes.NewReader(buf))
}

/*PutNode - implement interface */
func (pndb *PNodeDB) PutNode(key Key, node Node) error {
	data := node.Encode()
	err := pndb.db.Put(pndb.wo, key, data)
	return err
}

/*DeleteNode - implement interface */
func (pndb *PNodeDB) DeleteNode(key Key) error {
	err := pndb.db.Delete(pndb.wo, key)
	return err
}

/*MultiPutNode - implement interface */
func (pndb *PNodeDB) MultiPutNode(keys []Key, nodes []Node) error {
	wb := gorocksdb.NewWriteBatch()
	defer wb.Destroy()
	for idx, key := range keys {
		wb.Put(key, nodes[idx].Encode())
	}
	return pndb.db.Write(pndb.wo, wb)
}

/*MultiDeleteNode - implement interface */
func (pndb *PNodeDB) MultiDeleteNode(keys []Key) error {
	wb := gorocksdb.NewWriteBatch()
	defer wb.Destroy()
	for _, key := range keys {
		wb.Delete(key)
	}
	return pndb.db.Write(pndb.wo, wb)
}

/*Iterate - implement interface */
func (pndb *PNodeDB) Iterate(ctx context.Context, handler NodeDBIteratorHandler) error {
	ro := gorocksdb.NewDefaultReadOptions()
	ro.SetFillCache(false)
	it := pndb.db.NewIterator(ro)
	defer it.Close()
	for it.SeekToFirst(); it.Valid(); it.Next() {
		key := it.Key()
		value := it.Value()
		kdata := key.Data()
		vdata := value.Data()
		node, err := CreateNode(bytes.NewReader(vdata))
		if err != nil {
			key.Free()
			value.Free()
			Logger.Error("iterate - create node", zap.String("key", ToHex(kdata)), zap.Error(err))
			continue
		}
		err = handler(ctx, kdata, node)
		if err != nil {
			key.Free()
			value.Free()
			Logger.Error("iterate - create node handler error", zap.String("key", ToHex(kdata)), zap.Any("data", vdata), zap.Error(err))
			break
		}
		key.Free()
		value.Free()
	}
	return nil
}

/*Flush - flush the db */
func (pndb *PNodeDB) Flush() {
	pndb.db.Flush(pndb.fo)
}

/*PruneBelowOrigin - prune the state below the given origin */
func (pndb *PNodeDB) PruneBelowOrigin(ctx context.Context, origin Origin) error {
	BatchSize := 64
	ps := GetPruneStats(ctx)
	var total int64
	var count int64
	batch := make([]Key, BatchSize)
	handler := func(ctx context.Context, key Key, node Node) error {
		total++
		if node.GetOrigin() >= origin {
			return nil
		}
		count++
		if config.DevConfiguration.State {
			Logger.Debug("prune below origin - deleting node", zap.String("key", ToHex(key)), zap.Any("old_origin", node.GetOrigin()), zap.Any("new_origin", origin))
		}
		tkey := make([]byte, len(key))
		copy(tkey, key)
		batch = append(batch, tkey)
		if len(batch) == BatchSize {
			var err error
			if config.DevConfiguration.State {
				err = pndb.MultiDeleteNode(batch)
				Logger.Debug("prune below origin - deleting nodes", zap.String("key", ToHex(key)), zap.Any("old_origin", node.GetOrigin()), zap.Any("new_origin", origin))
			}
			batch = batch[:0]
			if err != nil {
				Logger.Error("prune below origin - error deleting node", zap.String("key", ToHex(key)), zap.Any("old_origin", node.GetOrigin()), zap.Any("new_origin", origin), zap.Error(err))
				return err
			}
		}
		return nil
	}
	err := pndb.Iterate(ctx, handler)
	if err != nil {
		return err
	}
	if len(batch) > 0 {
		if config.DevConfiguration.State {
			err := pndb.MultiDeleteNode(batch)
			if err != nil {
				Logger.Error("prune below origin - error deleting node", zap.Any("new_origin", origin), zap.Error(err))
				return err
			}
		}
	}
	pndb.Flush()
	if ps != nil {
		ps.Total = total
		ps.Deleted = count
	}
	if config.DevConfiguration.State {
		size := pndb.Size(ctx)
		Logger.Info("prune below origin", zap.Int64("size", size), zap.Any("stats", ps))
	}
	return err
}

/*Size - count number of keys in the db */
func (pndb *PNodeDB) Size(ctx context.Context) int64 {
	var count int64
	handler := func(ctx context.Context, key Key, node Node) error {
		count++
		return nil
	}
	err := pndb.Iterate(ctx, handler)
	if err != nil {
		Logger.Error("count", zap.Error(err))
		return -1
	}
	return count
}
