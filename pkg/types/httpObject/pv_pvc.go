package httpobject

import "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/persistVolume"

type HTTPRequest_UpdatePersistVolume struct {
	/* pv name -> pv */
	Pv map[string]persistVolume.PersistVolume `json:"pv" yaml:"pv"`
	/* pvc name -> pvc */
	Pvc map[string]persistVolume.PersistVolumeClaim `json:"pvc" yaml:"pvc"`
}

type HTTPResponse_GetAllPersistVolume struct {
	/* pv name -> pv */
	Pv map[string]persistVolume.PersistVolume `json:"pv" yaml:"pv"`
	/* pvc name -> pvc */
	Pvc map[string]persistVolume.PersistVolumeClaim `json:"pvc" yaml:"pvc"`
}

type HTTPReuqest_AddPVImmediately struct {
	PvName string `json:"pvName" yaml:"pvName"`
	PvType string `json:"pvType" yaml:"pvType"`
}

type HTTPRequest_AddPV struct {
	Pv persistVolume.PersistVolume `json:"pv" yaml:"pv"`
}

type HTTPRequest_AddPVC struct {
	Pvc persistVolume.PersistVolumeClaim `json:"pvc" yaml:"pvc"`
}
