package pod

import (
	minik8s_apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"
	minik8s_container "github.com/MiniK8s-SE3356/minik8s/pkg/types/container"
)

type PodPhase string

// These are the valid statuses of pods.
const (
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	// Deprecated: It isn't being set since 2015 (74da3b14b0c0f658b3bb8d2def5094686d0e9095)
	PodUnknown PodPhase = "Unknown"
)

type ContainerStatus struct {
	Name string `yaml:"name" json:"name"`

	State minik8s_container.ContainerState `yaml:"state" json:"state"`

	Ready bool `yaml:"ready" json:"ready"`

	Started string `yaml:"started" json:"started"`

	Finished string `yaml:"finished" json:"finished"`
}

type PodStatus struct {
	Phase PodPhase `yaml:"phase" json:"phase"`

	PodIP string `yaml:"podIP" json:"podIP"`

	ContainerStatuses []ContainerStatus `yaml:"containerStatuses" json:"containerStatuses"`

	CPUUsage float64 `yaml:"cpuUsage" json:"cpuUsage"`

	MemoryUsage float64 `yaml:"memoryUsage" json:"memoryUsage"`
}

type HostPath struct {
	Path string `json:"path" yaml:"path"`

	Type string `json:"type" yaml:"type"`
}

// NOTE: modify for PV/PVC BEGIN

type Volume struct {
	Name     string   `json:"name" yaml:"name"`
	HostPath HostPath `json:"hostPath" yaml:"hostPath"`
	// 此处，PV的优先级高于PVC,如果两个都写，以PV的规则为准
	// 这俩的优先级均低于HostPath
	PersistentVolume      PodMountPV  `json:"persistentVolume" yaml:"persistentVolume"`
	PersistentVolumeClaim PodMountPVC `json:"persistentVolumeClaim" yaml:"persistentVolumeClaim"`
}

type PodMountPV struct {
	PvName string `json:"pvName" yaml:"pvName"`
}

type PodMountPVC struct {
	ClaimName string `json:"claimName" yaml:"claimName"`
}

// NOTE: modify for PV/PVC END

type PodSpec struct {
	NodeName string `json:"nodeName" yaml:"nodeName"`

	Containers []minik8s_container.Container `json:"containers" yaml:"containers"`

	Volumes []Volume `json:"volumes" yaml:"volumes"`

	// Some labels for the pod
	NodeSelector map[string]string `json:"nodeSelector" yaml:"nodeSelector"`
}

type Pod struct {
	minik8s_apiobject.Basic `yaml:",inline"`

	Spec PodSpec `json:"spec" yaml:"spec"`

	Status PodStatus `json:"status" yaml:"status"`
}
