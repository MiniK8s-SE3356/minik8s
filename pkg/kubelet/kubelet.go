package kubelet

import (
	kubelet_worker "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet/worker"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
)

type Kubelet struct {
	kubeletConfig *KubeletConfig
	mqConn        *minik8s_message.MQConnection
	podManager    *kubelet_worker.PodManager
}

func NewKubelet(config *KubeletConfig) *Kubelet {
	kubelet := &Kubelet{
		kubeletConfig: config,
	}
	return kubelet
}
