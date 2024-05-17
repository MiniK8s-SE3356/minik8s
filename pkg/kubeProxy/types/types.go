package types

type KpClusterIP struct {
	Version int      `json:"version" yaml:"version"`
	Vports  []uint16 `json:"vports" yaml:"vports"`
}

type KpNodePort struct {
	Version int      `json:"version" yaml:"version"`
	Nports  []uint16 `json:"nports" yaml:"nports"`
}

type KpServicesStatus struct {
	ClusterIP map[string]KpClusterIP `json:"clusterIP" yaml:"clusterIP"`
	NodePort  map[string]KpNodePort  `json:"nodePort" yaml:"nodePort"`
}
