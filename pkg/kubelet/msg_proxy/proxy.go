package msgproxy

import (
	minik8s_worker "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet/worker"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
)

type MsgProxy struct {
	mqConn           *minik8s_message.MQConnection
	PodUpdateChannel chan<- *minik8s_worker.Task
}
