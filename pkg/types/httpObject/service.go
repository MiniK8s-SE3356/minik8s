package httpobject

import (
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
)

type HTTPResponse_GetAllServices struct {
	ClusterIP []service.ClusterIP `json:"clusterIP" yaml:"clusterIP"`
	NodePort  []service.NodePort  `json:"nodePort" yaml:"nodePort"`
}

type HTTPRequest_UpdateServices map[string]string
