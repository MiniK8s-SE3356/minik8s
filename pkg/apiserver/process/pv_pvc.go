package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/persistVolume"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

func AddPV(namespace string, pv *persistVolume.PersistVolume) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Printf("failed to generate uuid for PV\n")
		return "failed to generate uuid", err
	}

	pv.Metadata.Id = id
	if pv.Metadata.Labels == nil {
		pv.Metadata.Labels = make(map[string]string)
	}
	pv.Status.Phase = persistVolume.PV_PHASE_CREATED
	pv.Status.MountPod = make([]string, 0)

	value, err := json.Marshal((*pv))
	if err != nil {
		fmt.Println("failed to translate pv into json ", err.Error())
		return "failed to translate pv into json ", err
	}
	// 先查看一下key是否已经存在
	tmp, err := EtcdCli.Exist(pvPrefix + namespace + "/" + pv.Metadata.Name)
	if err != nil {
		fmt.Println("failed to check pv existence in etcd")
		return "failed to check pv existence in etcd", err
	}
	if tmp {
		fmt.Println("pv has existed")
		return "pv has existed", nil
	}
	// 然后存入etcd
	err = EtcdCli.Put(pvPrefix+namespace+"/"+pv.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write pv to etcd ", err.Error())
		return "failed to write pv to etcd", err
	}

	return "add PV to minik8s", nil
}

func AddPVC(namespace string, pvc *persistVolume.PersistVolumeClaim) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate pvc uuid")
		return "failed to generate pvc uuid", err
	}

	pvc.Metadata.Id=id
	if(pvc.Spec.Selector.MatchLabels==nil){
		pvc.Spec.Selector.MatchLabels=make(map[string]string)
	}
	pvc.Status.Phase=persistVolume.PVC_PHASE_AVAILABLE
	if(pvc.Status.BoundPV==nil){
		pvc.Status.BoundPV=make([]string, 0)
	}

	value, err := json.Marshal(*pvc)
	if err != nil {
		fmt.Println("failed to translate pvc into json ", err.Error())
		return "failed to translate pvc into json ", err
	}

	// 先查看一下key是否已经存在
	tmp, err := EtcdCli.Exist(pvcPrefix + namespace + "/" + pvc.Metadata.Name)
	if err != nil {
		fmt.Println("failed to check pvc existence in etcd")
		return "failed to check pvc existence in etcd", err
	}
	if tmp {
		fmt.Println("pvc has existed")
		return "pvc has existed", nil
	}
	// 然后存入etcd
	err = EtcdCli.Put(pvcPrefix+namespace+"/"+pvc.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write pvc to etcd ", err.Error())
		return "failed to write pvc to etcd", err
	}

	return "add pvc to minik8s", nil
}
