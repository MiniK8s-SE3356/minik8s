package service

type EndPoint struct {
	Id       string `json:"id" yaml:"id"`
	Protocol string `json:"protocol" yaml:"protocol"`
	PodID    string `json:"podID" yaml:"podID"`
	PodIP    string `json:"podIP" yaml:"podIP"`
	PodPort  uint16 `json:"podPort" yaml:"podPort"`
}
