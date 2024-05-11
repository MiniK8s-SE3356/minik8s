package node

// NotReady				表示节点不健康而且不能接收 Pod
// Ready	 			表示节点是健康的并已经准备好接收 Pod
// Unknown 				表示节点控制器在最近 node-monitor-grace-period 期间（默认 40 秒）没有收到节点的消息
// DiskPressure			表示节点存在磁盘空间压力，即磁盘可用量低
// MemoryPressure		表示节点存在内存压力，即节点内存可用量低
// PIDPressure	True 	表示节点存在进程压力，即节点上进程过多
// NetworkUnavailable	表示节点网络配置不正确

const (
	NODE_PREFIX             = "MINIK8S-NODE-"
	NODE_NOTREADY           = "NotReady"
	NODE_Ready              = "Ready"
	NODE_Unknown            = "Unknown"
	NODE_DiskPressure       = "DiskPressure"
	NODE_MemoryPressure     = "MemoryPressure"
	NODE_PIDPressure        = "PIDPressure"
	NODE_NetworkUnavailable = "NetworkUnavailable"
)

type Node struct {
	Metadata NodeMetadata
	Status   NodeStatus
}

type NodeStatus struct {
	Hostname   string   `json:"hostname" yaml:"hostname"`
	Ip         string   `json:"ip" yaml:"ip"`
	Condition  []string `json:"condition" yaml:"condition"` /*对应上述的NODE状态*/
	CpuPercent float64  `json:"cpuPercent" yaml:"cpuPercent"`
	MemPercent float64  `json:"memPercent" yaml:"memPercent"`
	NumPods    int      `json:"numPods" yaml:"numPods"`
	UpdateTime string   `json:"updateTime" yaml:"updateTime"`
}

type NodeMetadata struct {
	Id     string            `json:"id" yaml:"id"`
	Name   string            `json:"name" yaml:"name"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}
