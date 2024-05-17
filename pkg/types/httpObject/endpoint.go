package httpobject

import (
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
)

type HTTPResponse_GetAllEndpoint map[string]service.EndPoint

type HTTPRequest_AddorDeleteEndpoint struct {
	Delete []string           `json:"delete" yaml:"delete"`
	Add    []service.EndPoint `json:"add" yaml:"add"`
}
