package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/amrittb/choto-link-api/internal/allocator"
	"github.com/amrittb/choto-link-api/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func getIdAllocator(etcdClient *clientv3.Client) (*allocator.IdAllocator, error) {
	etcdKey := os.Getenv("ETCD_ID_ALLOCATION_KEY")
	if etcdKey == "" {
		return nil, fmt.Errorf("ETCD_ID_ALLOCATION_KEY environment variable is empty")
	}
	rangeSize, err := strconv.ParseUint(os.Getenv("ID_ALLOCATION_RANGE_SIZE"), 10, 64)
	if err != nil {
		return nil, err
	}

	return allocator.NewIdAllocator(etcdClient, etcdKey, rangeSize)
}

func getUrlRepository(pgxPool *pgxpool.Pool) (*repository.UrlRepository, error) {
	return repository.NewUrlRepository(pgxPool), nil
}
