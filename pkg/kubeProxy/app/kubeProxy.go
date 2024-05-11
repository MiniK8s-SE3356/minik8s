package app

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/iptablesController"
)

type KubeProxy struct {
	iptablesController *(iptablesController.IptablesController)
}

func NewKubeProxy() *KubeProxy {
	fmt.Printf("Create KubeProxy ...\n")
	iptables_controller := iptablesController.NewIptablesController()
	kube_proxy := &KubeProxy{
		iptablesController: iptables_controller,
	}
	return kube_proxy
}

func (kp *KubeProxy) Init() {
	fmt.Printf("Init KubeProxy ...\n")
	kp.iptablesController.Init()
}

func (kp *KubeProxy) Run() {
	fmt.Printf("Run KubeProxy ...\n")
	kp.iptablesController.SyncConfig()
	kp.iptablesController.SyncIptables()
}
