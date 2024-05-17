package msgproxy

import (
	minik8s_worker "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet/worker"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
)

type MsgProxy struct {
	mqConn           *minik8s_message.MQConnection
	PodUpdateChannel chan<- *minik8s_worker.Task
}

func NewMsgProxy(mqConfig *minik8s_message.MQConfig) *MsgProxy {
	newConn, _ := minik8s_message.NewMQConnection(mqConfig)
	return &MsgProxy{
		mqConn:           newConn,
		PodUpdateChannel: make(chan *minik8s_worker.Task),
	}
}
