package kubelet

import (
	minik8s_types "github.com/MiniK8s-SE3356/minik8s/pkg/types"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
)

type Kubelet struct {
	kubeletConfiguration *minik8s_types.KubeletConfiguration
	mqConn               *minik8s_message.MQConnection
}
