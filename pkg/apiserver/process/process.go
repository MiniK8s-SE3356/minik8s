package process

import "github.com/MiniK8s-SE3356/minik8s/pkg/etcdclient"

const (
	prefix          = "/minik8s"
	namespacePrefix = prefix + "/namespace/"
)

var EtcdCli *etcdclient.EtcdClient
