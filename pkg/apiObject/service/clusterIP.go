package service

import "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/selector"

const (
	CLUSTERIP_PREFIX = "MINIK8S-CLUSTERIP-"
	// CLUSTERIP_ID_ALLOCATED        = "ID Allocated"        /*状态1*/
	// CLUSTERIP_IP_ALLOCATED        = "IP Allocated"        /*状态2*/
	// CLUSTERIP_ENDPOINTS_ALLOCATED = "Endponits Allocated" /*状态3*/
	// CLUSTERIP_SUCCESS             = "Success"             /*状态4*/
	// CLUSTERIP_ERROR               = "Error"               /*状态5*/
	CLUSTERIP_READY    = "READY"
	CLUSTERIP_NOTREADY = "NOTREADY"
)

type ClusterIP struct {
	ApiVersion string            `json:"apiVersion" yaml:"apiVersion"`
	Kind       string            `json:"kind" yaml:"kind"`
	Metadata   ClusterIPMetadata `json:"metadata" yaml:"metadata"`
	Spec       ClusterIPSpec     `json:"spec" yaml:"spec"`
	Status     ClusterIPStatus   `json:"status" yaml:"status"`
}

type ClusterIPMetadata struct {
	Name      string            `json:"name" yaml:"name"`
	Namespace string            `json:"namespace" yaml:"namespace"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
	Ip        string            `json:"ip" yaml:"ip"`
	Id        string            `json:"id" yaml:"id"`
}

type ClusterIPSpec struct {
	Type     string              `json:"type" yaml:"type"`
	Selector selector.Selector   `json:"selector" yaml:"selector"`
	Ports    []ClusterIPPortInfo `json:"ports" yaml:"ports"`
}

type ClusterIPPortInfo struct {
	Protocal   string `json:"protocal" yaml:"protocal"`
	Port       uint16 `json:"port" yaml:"port"`
	TargetPort uint16 `json:"targetPort" yaml:"targetPort"`
}

type ClusterIPStatus struct {
	Phase          string              `json:"phase" yaml:"phase"` /*READY or NOTREADY */
	Version        int                 /* 版本号 */
	ServicesStatus map[uint16][]string `json:"servicesStatus" yaml:"servicesStatus"`
}
