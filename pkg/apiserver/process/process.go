package process

import (
	"sync"

	"github.com/MiniK8s-SE3356/minik8s/pkg/etcdclient"
)

const (
	prefix           = "/minik8s"
	namespacePrefix  = prefix + "/namespace/"
	nodePrefix       = prefix + "/node/"
	podPrefix        = prefix + "/pod/"
	servicePrefix    = prefix + "/service/"
	endpointPrefix   = prefix + "/endpoint/"
	replicasetPrefix = prefix + "/replicaset"
)

var EtcdCli *etcdclient.EtcdClient
var mu sync.RWMutex
