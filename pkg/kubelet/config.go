package kubelet

import (
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
)

type KubeletConfig struct {
	minik8s_message.MQConfig
	APIServerIP   string
	APIServerPort string
	NodeIP        string
	NodeHostName  string
	Labels        map[string]string
}
