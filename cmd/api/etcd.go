package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func getEtcdClient() (*clientv3.Client, error) {
	etcdEndpointRaw := os.Getenv("ETCD_ENDPOINT")
	if etcdEndpointRaw == "" {
		return nil, fmt.Errorf("ETCD_ENDPOINT environment variable is empty")
	}

	etcdEndpoints := strings.Split(etcdEndpointRaw, ",")
	if len(etcdEndpoints) == 0 {
		return nil, fmt.Errorf("ETCD_ENDPOINT environment variable is empty")
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
