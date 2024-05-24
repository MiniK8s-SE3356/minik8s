package service

import "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/selector"

const (
	NODEPORT_PREFIX = "MINIK8S-NODEPORT-"
	// NODEPORT_ID_ALLOCATED     = "ID Allocated"   /*状态1*/
	// NODEPORT_CLUSTERIP_FINISH = "Cluster Finish" /*状态2*/
	// NODEPORT_SUCCESS          = "Success"        /*状态3*/
	// NODEPORT_ERROR            = "Error"          /*状态4*/
	NODEPORT_READY    = "READY"
	NODEPORT_NOTREADY = "NOTREADY"
)

type NodePort struct {
	ApiVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   NodePortMetadata `json:"metadata" yaml:"metadata"`
	Spec       NodePortSpec     `json:"spec" yaml:"spec"`
	Status     NodePortStatus   `json:"status" yaml:"status"`
}

type NodePortMetadata struct {
	Name      string            `json:"name" yaml:"name"`
	Namespace string            `json:"namespace" yaml:"namespace"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
	Id        string            `json:"id" yaml:"id"`
}

type NodePortSpec struct {
	Type     string             `json:"type" yaml:"type"`
	Selector selector.Selector  `json:"selector" yaml:"selector"`
	Ports    []NodePortPortInfo `json:"ports" yaml:"ports"`
}

type NodePortPortInfo struct {
	NodePort   uint16 `json:"nodePort" yaml:"nodePort"`
	Protocol   string `json:"protocol" yaml:"protocol"`
	Port       uint16 `json:"port" yaml:"port"`
	TargetPort uint16 `json:"targetPort" yaml:"targetPort"`
}

type NodePortStatus struct {
	Phase       string `json:"phase" yaml:"phase"` /*READY or NOTREADY，初始为NOTREADY*/
	Version     int    `json:"version" yaml:"version"`	/*初始为0*/
	ClusterIPID string `json:"clusterIPID" yaml:"clusterIPID"`
}
