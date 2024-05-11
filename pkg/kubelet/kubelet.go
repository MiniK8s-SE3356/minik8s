package kubelet

import (
	minik8s_types "github.com/MiniK8s-SE3356/minik8s/pkg/types"
)

type Kubelet struct {
	kubeletConfiguration *minik8s_types.KubeletConfiguration
}
