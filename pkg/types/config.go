package types

import (
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

// CreateContainerConfig is the configuration for creating a container
// When you create a container, you need 3 types of configurations:
// 1. container.Config
// 2. container.HostConfig
// 3. network.NetworkingConfig
// This struct is a wrapper for all 3 configurations
type CreateContainerConfig struct {
	// 1. Config:
	Image        string              // Name of the image as it was passed by the operator (e.g. could be symbolic)
	Volumes      map[string]struct{} // List of volumes (mounts) used for the container
	WorkingDir   string              // Current directory (PWD) in the command will be launched
	Env          []string            // List of environment variable to set in the container
	Entrypoint   strslice.StrSlice   // Entrypoint to run when starting the container
	Cmd          strslice.StrSlice   // Command to run when starting the container
	ExposedPorts nat.PortSet         `json:",omitempty"` // List of exposed ports
	Labels       map[string]string   // List of labels set to this container

	// 2. HostConfig:
	Binds        []string    // List of volume bindings for this container
	PortBindings nat.PortMap // Port mapping between the exposed port (container) and the host

	// shared workspace in a pod
	NetworkMode string   // Network mode to use for the container
	IpcMode     string   // IPC namespace to use for the container
	PidMode     string   // PID namespace to use for the container
	VolumesFrom []string // List of volumes to take from other containers

	// resource constraints
	// TODO: There are still a lot of resource could be limited. Now only CPU & Memory
	CPU    int64 // CPU shares (relative weight)
	Memory int64 // Memory limit for the container

	// 3.NetworkingConfig:
	// TODO: There is no need to modify NetworkingConfig for now, maybe add in the future
}

type KubeletConfiguration struct {
}
