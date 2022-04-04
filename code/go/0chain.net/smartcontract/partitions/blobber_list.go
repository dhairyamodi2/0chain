package partitions

import (
	"0chain.net/chaincore/chain/state"
	"0chain.net/core/datastore"
	"0chain.net/core/util"
	"encoding/json"
	"errors"
	"fmt"
)

//go:generate msgp -io=false -tests=false -unexported=true -v

//------------------------------------------------------------------------------

type BlobberNode struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

func (bn *BlobberNode) Encode() []byte {
	var b, err = json.Marshal(bn)
	if err != nil {
		panic(err)
	}
	return b
}

func (bn *BlobberNode) Decode(b []byte) error {
	return json.Unmarshal(b, bn)
}

func (bn *BlobberNode) Data() string {
	return bn.Url
}

func (bn *BlobberNode) Name() string {
	return bn.ID
}

//------------------------------------------------------------------------------

type blobberItemList struct {
	Key     string        `json:"-"`
	Items   []BlobberNode `json:"items"`
	Changed bool          `json:"-"`
}

func (il *blobberItemList) Encode() []byte {
	var b, err = json.Marshal(il)
	if err != nil {
		panic(err)
	}
	return b
}

func (il *blobberItemList) Decode(b []byte) error {
	return json.Unmarshal(b, il)
}

func (il *blobberItemList) save(balances state.StateContextI) error {
	_, err := balances.InsertTrieNode(il.Key, il)
	return err
}

func (il *blobberItemList) get(key datastore.Key, balances state.StateContextI) error {
	err := balances.GetTrieNode(key, il)
	if err != nil {
		if err != util.ErrValueNotPresent {
			return err
		}
		il.Key = key
	}
	return nil
}

func (il *blobberItemList) add(it PartitionItem) error {
	for _, bi := range il.Items {
		if bi.ID == it.Name() {
			return errors.New("blobber item already exists")
		}
	}
	il.Items = append(il.Items, BlobberNode{
		ID:  it.Name(),
		Url: it.Data(),
	})
	il.Changed = true
	return nil
}

func (il *blobberItemList) update(it PartitionItem) error {
	val, ok := it.(*BlobberNode)
	if !ok {
		return errors.New("invalid item")
	}

	for i := 0; i < il.length(); i++ {
		if il.Items[i].Name() == it.Name() {
			newItem := *val
			il.Items[i] = newItem
			il.Changed = true
			return nil
		}
	}
	return errors.New("item not found")
}

func (il *blobberItemList) remove(item PartitionItem) error {
	if len(il.Items) == 0 {
		return fmt.Errorf("searching empty partition")
	}
	index := il.find(item)
	if index == notFound {
		return fmt.Errorf("cannot find item %v in partition", item)
	}
	il.Items[index] = il.Items[len(il.Items)-1]
	il.Items = il.Items[:len(il.Items)-1]
	il.Changed = true
	return nil
}

func (il *blobberItemList) cutTail() PartitionItem {
	if len(il.Items) == 0 {
		return nil
	}

	tail := il.Items[len(il.Items)-1]
	il.Items = il.Items[:len(il.Items)-1]
	il.Changed = true
	return &tail
}

func (il *blobberItemList) length() int {
	return len(il.Items)
}

func (il *blobberItemList) changed() bool {
	return il.Changed
}

func (il *blobberItemList) itemRange(start, end int) []PartitionItem {
	if start > end || end > len(il.Items) {
		return nil
	}

	var rtv []PartitionItem
	for i := start; i < end; i++ {
		rtv = append(rtv, &il.Items[i])
	}
	return rtv
}

func (il *blobberItemList) find(searchItem PartitionItem) int {
	for i, item := range il.Items {
		if item.Name() == searchItem.Name() {
			return i
		}
	}
	return notFound
}

//------------------------------------------------------------------------------