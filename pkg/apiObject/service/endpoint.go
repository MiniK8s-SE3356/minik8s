package service

type EndPoint struct {
	Id      string `json:"id" yaml:"id"`
	PodID   string `json:"podID" yaml:"podID"`
	PodIP   string `json:"podIP" yaml:"podIP"`
	PodPort int16  `json:"podPort" yaml:"podPort"`
}
