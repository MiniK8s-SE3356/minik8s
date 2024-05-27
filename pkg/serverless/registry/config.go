package registry

import (
	"github.com/docker/docker/api/types/registry"
)

const (
	REGISTRY_IMAGE = "registry:2"
	REGISTRY_NAME  = "minik8s-registry"
	REGISTRY_PORT  = "5000"
)

var REGISTRY_IP = "localhost"

var REGISTRY_AUTH_CONFIG = registry.AuthConfig{
	Username: "minik8s",
	Password: "minik8s",
}
