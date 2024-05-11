package kubelet

import (
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
)

type Kubelet struct {
	kubeletConfiguration *KubeletConfiguration
	mqConn               *minik8s_message.MQConnection
}
