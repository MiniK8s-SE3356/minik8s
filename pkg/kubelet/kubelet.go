package kubelet

import (
	msgproxy "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet/msg_proxy"
	kubelet_worker "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet/worker"
)

type Kubelet struct {
	kubeletConfig *KubeletConfig
	msgProxy      *msgproxy.MsgProxy
	podManager    kubelet_worker.PodManager
}

func NewKubelet(config *KubeletConfig) *Kubelet {
	kubelet := &Kubelet{
		kubeletConfig: config,
		msgProxy:      msgproxy.NewMsgProxy(&config.MQConfig),
		podManager:    kubelet_worker.NewPodManager(),
	}
	return kubelet
}

func (k *Kubelet) Run() {
	go k.msgProxy.Run()
}
