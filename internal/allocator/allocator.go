package allocator

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.etcd.io/etcd/client/v3"
)

type IdAllocator struct {
	etcdClient *clientv3.Client
	etcdKey    string
	rangeSize  uint64
	current    atomic.Uint64
	end        uint64
	mu         sync.Mutex
}

func NewIdAllocator(client *clientv3.Client, etcdKey string, rangeSize uint64) (*IdAllocator, error) {
	allocator := &IdAllocator{etcdClient: client, etcdKey: etcdKey, rangeSize: rangeSize}

	err := allocator.allocateNewRange()
	if err != nil {
		return nil, err
	}
	return allocator, nil
}

func (ka *IdAllocator) allocateNewRange() error {
	// 1. Lock Key Allocator
	// 1. Start TXN
	// 2. Get current value
	// 3. Set next value using Compare-And-Swap
	// 4. Unlock Key Allocator
	ka.mu.Lock()
	defer ka.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := ka.etcdClient.Get(ctx, ka.etcdKey)
	if err != nil {
		return err
	}

	var current uint64
	var version int64

	if len(res.Kvs) > 0 {
		current, _ = strconv.ParseUint(string(res.Kvs[0].Value), 10, 64)
		version = res.Kvs[0].Version
	}

	newEnd := current + ka.rangeSize
	txnRes, err := ka.etcdClient.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(ka.etcdKey), "=", version)).
		Then(clientv3.OpPut(ka.etcdKey, fmt.Sprintf("%d", newEnd))).Commit()

	if err != nil {
		return err
	}

	if !txnRes.Succeeded {
		return fmt.Errorf("etcd id range modified by someone else")
	}

	ka.current.Store(current)
	ka.end = newEnd
	log.Printf("Allocated new ID range: [%d - %d]", current, newEnd)
	return nil
}

func (ka *IdAllocator) NextId() (uint64, error) {
	nextId := ka.current.Add(1) - 1
	if nextId < ka.end {
		// If within range, then return as it is.
		return nextId, nil
	}

	err := ka.allocateNewRange()
	if err != nil {
		return 0, err
	}

	// Return from new range
	return ka.current.Add(1) - 1, nil
}
