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

func (k *Kubelet) Proxy() {
	// update := <-k.msgProxy.PodUpdateChannel
	// k.podManager.AddPod(update.Pod)

	// go k.Proxy()
	for update := range k.msgProxy.PodUpdateChannel {
		switch update.Type {
		case kubelet_worker.Task_Add:
			k.podManager.AddPod(update.Pod)
		case kubelet_worker.Task_Update:
			k.podManager.UpdatePod(update.Pod)
		case kubelet_worker.Task_Remove:
			k.podManager.RemovePod(update.Pod)
		}
	}
}

func (k *Kubelet) Run() {
	go k.Proxy()
	go k.msgProxy.Run()
}
