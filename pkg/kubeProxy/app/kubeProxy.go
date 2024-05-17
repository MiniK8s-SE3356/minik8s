package app

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/iptablesController"
	kptype "github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/types"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
)

type KubeProxy struct {
	iptables_controller *(iptablesController.IptablesController)
	services_status     kptype.KpServicesStatus
	mutex               sync.Mutex
}

var need_send bool = false

func NewKubeProxy() *KubeProxy {
	fmt.Printf("Create KubeProxy ...\n")
	iptables_controller := iptablesController.NewIptablesController()
	kube_proxy := &KubeProxy{
		iptables_controller: iptables_controller,
		services_status:     kptype.KpServicesStatus{},
		mutex:               sync.Mutex{},
	}
	return kube_proxy
}

func (kp *KubeProxy) Init() {
	fmt.Printf("Init KubeProxy ...\n")
	kp.iptables_controller.Init()
}

func (kp *KubeProxy) Run() {
	fmt.Printf("Run KubeProxy ...\n")

	poller.PollerStaticPeriod(10*time.Second, kp.routine, true)
	poller.PollerStaticPeriod(5*time.Second, kp.syncServicesToKubelet, true)
}

func (kp *KubeProxy) routine() {
	fmt.Printf("KubeProxy routine...\n")

	var service_list httpobject.HTTPResponse_GetAllServices
	status, err := httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllService", nil, &service_list)

	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}

	var endpoint_list httpobject.HTTPResponse_GetAllEndpoint
	status, err = httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllEndpoint", nil, &endpoint_list)

	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}

	new_service, err := kp.iptables_controller.SyncConfig(&service_list, &endpoint_list)
	if err != nil {
		return
	}

	err = kp.iptables_controller.SyncIptables()
	if err != nil {
		return
	} else {
		kp.mutex.Lock()
		kp.services_status = new_service
		need_send = true
		kp.mutex.Unlock()
	}

}

func (kp *KubeProxy) syncServicesToKubelet() {
	// 向kubelet更新本机上的service信息
	if need_send {
		kp.mutex.Lock()
		// TODO: 将node上的service数据发送至kubelet
		need_send = false
		kp.mutex.Unlock()
	}
}
