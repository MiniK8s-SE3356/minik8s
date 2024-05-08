package etcdclient

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	cl *clientv3.Client
}

type KVPair struct {
	Key   string
	Value string
}

const (
	etcdURL = ""
	timeout = 5 * 1000 * 1000
)

func Connect(endpoints []string, dialTimeout time.Duration) (*EtcdClient, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout, // 连接超时时间
	})

	if err != nil {
		fmt.Println("failed to connect to etcd ", err.Error())
		return nil, err
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = cli.Status(timeoutContext, endpoints[0])
	if err != nil {
		return nil, err
	}

	return &EtcdClient{cl: cli}, nil
}

func (c *EtcdClient) Get(key string) ([]byte, error) {
	resp, err := c.cl.Get(context.TODO(), key)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return []byte{}, nil
	}
	// 理论上Get到是唯一的

	if len(resp.Kvs) == 0 {
		// key not found的情况怎么办，还需要再考虑
		fmt.Println("key not found ")
		return []byte{}, nil
	}

	if len(resp.Kvs) > 1 {
		// 这里还要再处理
		return []byte{}, nil
	}
	kv := resp.Kvs[0]

	return kv.Value, nil
}

func (c *EtcdClient) Put(key string, value string) error {
	if _, err := c.cl.Put(context.TODO(), key, value); err != nil {
		fmt.Println("failed to put to etcd ", err.Error())
		return err
	}

	return nil
}

func (c *EtcdClient) Del(key string) error {
	if _, err := c.cl.Delete(context.TODO(), key); err != nil {
		fmt.Println("failed to del in etcd ", err.Error())
		return err
	}

	return nil
}

func (c *EtcdClient) DelAll() error {
	if _, err := c.cl.Delete(context.TODO(), "", clientv3.WithPrefix()); err != nil {
		fmt.Println("failed to delall in etcd ", err.Error())
		return err
	}

	return nil
}

func (c *EtcdClient) GetWithPrefix(keyPrefix string) ([]KVPair, error) {
	response, err := c.cl.Get(context.TODO(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("failed to get with prefix ", err.Error())
		return []KVPair{}, err
	}

	var result []KVPair
	for _, kvp := range response.Kvs {
		result = append(result, KVPair{Key: string(kvp.Key), Value: string(kvp.Value)})
	}

	return result, nil
}
