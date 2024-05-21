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
	Key   string `json:"key"`
	Value string `json:"value"`
}

const (
	EtcdURL = "http://127.0.0.1:2379"
	Timeout = 5 * 1000 * 1000
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

	timeoutContext, cancel := context.WithTimeout(context.Background(), Timeout)
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
		return []byte{}, err
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

func (c *EtcdClient) Ls(dir string) ([]KVPair, error) {
	// gpt说是可以这么写，但是get出来的dir和file分别是什么格式还不清楚
	response, err := c.cl.Get(context.TODO(), dir)
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

func (c *EtcdClient) Mkdir(dirPath string) error {
	// 虽然etcd提供了ls mkdir这样的命令，但是实际上这只是对prefix的一种等效处理
	// 并不存在真正的创建
	// 所以put的时候不需要先mkdir再put，put进去就直接“创建完成”了
	// 不过这里还是提供了一个mkdir方法
	if _, err := c.cl.Put(context.TODO(), dirPath, ""); err != nil {
		fmt.Println("failed to put to etcd ", err.Error())
		return err
	}

	return nil
}

func (c *EtcdClient) Exist(key string) (bool, error) {
	resp, err := c.cl.Get(context.TODO(), key)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return false, err
	}

	if len(resp.Kvs) == 0 {
		return false, nil
	}

	return true, nil

}
