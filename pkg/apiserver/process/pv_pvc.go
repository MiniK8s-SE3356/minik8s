package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/persistVolume"
	persistVolumeController "github.com/MiniK8s-SE3356/minik8s/pkg/controller/PersistVolumeController"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

func DeletePV(namespace string, name string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(pvPrefix + namespace + "/" + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "service not found", nil
	}

	pvvalue, err := EtcdCli.Get(pvPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Printf("PV does not exist in DeletePV, namespace %s, name %s\n", namespace, name)
		return IsPVC_PV_exist_Return_PVMISS, err
	}
	var pv persistVolume.PersistVolume
	err = json.Unmarshal(pvvalue, &pv)
	if err != nil {
		fmt.Printf("can not unmarshal pvc in DeletePVC,error msg: %s\n", err.Error())
		return "PV umarsharshal error", err
	}
	if(pv.Status.MountPod!=nil&&len(pv.Status.MountPod)!=0){
		fmt.Println("Can't Delete PV: it still bind Pod!")
		return "Can't Delete PV: it still bind Pod!",errors.New("Can't Delete PV: it still bind Pod!")
	}

	cmd := exec.Command("rm","-r", persistVolumeController.MINIK8S_PV_PATH+"/"+name)
	
	_, err = cmd.Output()
	if err != nil {
		fmt.Println("Error executing command in delNFSVolume:", err)
		return "Error executing command in delNFSVolume:",err
	}

	err = EtcdCli.Del(pvPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}


	return "del successfully", nil
}

func DeletePVC(namespace string, name string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(pvcPrefix + namespace + "/" + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "service not found", nil
	}

	pvcvalue, err := EtcdCli.Get(pvcPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Printf("PVC does not exist in DeletePVC, namespace %s, name %s\n", namespace, name)
		return IsPVC_PV_exist_Return_PVCMISS, err
	}
	var pvc persistVolume.PersistVolumeClaim
	err = json.Unmarshal(pvcvalue, &pvc)
	if err != nil {
		fmt.Printf("can not unmarshal pvc in DeletePVC,error msg: %s\n", err.Error())
		return IsPVC_PV_exist_Return_PVCVALUEERROR, err
	}
	if(pvc.Status.BoundPV!=nil&&len(pvc.Status.BoundPV)!=0){
		fmt.Println("Can't Delete PVC: it still bind PV!")
		return "Can't Delete PVC: it still bind PV!",errors.New("Can't Delete PVC: it still bind PV!")
	}

	err = EtcdCli.Del(pvcPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del successfully", nil
}



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

	pvc.Metadata.Id = id
	if pvc.Spec.Selector.MatchLabels == nil {
		pvc.Spec.Selector.MatchLabels = make(map[string]string)
	}
	pvc.Status.Phase = persistVolume.PVC_PHASE_AVAILABLE
	if pvc.Status.BoundPV == nil {
		pvc.Status.BoundPV = make([]string, 0)
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

func GetAllPersistVolume() (map[string]persistVolume.PersistVolume, error) {
	mu.RLock()
	defer mu.RUnlock()
	result := make(map[string]persistVolume.PersistVolume, 0)

	pairs, err := EtcdCli.GetWithPrefix(pvPrefix)
	if err != nil {
		fmt.Println("failed to get pv from etcd")
		return result, err
	}

	for _, p := range pairs {
		var tmp persistVolume.PersistVolume
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to translate pv into json")
		} else {
			parts := strings.Split(p.Key, "/")
			result[parts[len(parts)-1]] = tmp
		}
	}

	return result, nil
}

func GetAllPersistVolumeClaim() (map[string]persistVolume.PersistVolumeClaim, error) {
	mu.RLock()
	defer mu.RUnlock()
	result := make(map[string]persistVolume.PersistVolumeClaim, 0)

	pairs, err := EtcdCli.GetWithPrefix(pvcPrefix)
	if err != nil {
		fmt.Println("failed to get pvc from etcd")
		return result, err
	}

	for _, p := range pairs {
		var tmp persistVolume.PersistVolumeClaim
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to translate pvc into json")
		} else {
			parts := strings.Split(p.Key, "/")
			result[parts[len(parts)-1]] = tmp
		}
	}

	return result, nil
}

func UpdatePersistVolume(namespace string, name string, value string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	err := EtcdCli.Put(pvPrefix+namespace+"/"+name, value)
	if err != nil {
		fmt.Println("failed to update pv in etcd")
		return "failed to update pv in etcd", err
	}

	return "update pv successfully", nil
}

func UpdatePersistVolumeClaim(namespace string, name string, value string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	err := EtcdCli.Put(pvcPrefix+namespace+"/"+name, value)
	if err != nil {
		fmt.Println("failed to update pvc in etcd")
		return "failed to update pvc in etcd", err
	}

	return "update pvc successfully", nil
}

const (
	IsPVAvailable_Return_MISS        = "PV_MISS"
	IsPVAvailable_Return_VALUEERROR  = "PV_VALUEERROR"
	IsPVAvailable_Return_UNAVAILABLE = "PV_UNAVAILABLE"
	IsPVAvailable_Return_OK          = "OK"
)

func IsPVAvailable(namespace string, name string) (string,string) {

	pvvalue, err := EtcdCli.Get(pvPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Printf("PV does not exist, namespace %s, name %s\n", namespace, name)
		return IsPVAvailable_Return_MISS,""
	}
	var pv persistVolume.PersistVolume
	err = json.Unmarshal(pvvalue, &pv)
	if err != nil {
		fmt.Printf("can not unmarshal pv in IsPVAvailable,error msg: %s\n", err.Error())
		return IsPVAvailable_Return_VALUEERROR,""
	}
	if pv.Status.Phase != persistVolume.PV_PHASE_AVAILABLE {
		fmt.Printf("pv not available in IsPVAvailable\n")
		return IsPVAvailable_Return_UNAVAILABLE,""
	}
	return IsPVAvailable_Return_OK,pv.Spec.Capacity
}

const (
	// 不正常的error，拒绝执行
	// PVC不存在
	IsPVC_PV_exist_Return_PVCMISS = "PVCMISS"
	// PVC无法Unmarshal
	IsPVC_PV_exist_Return_PVCVALUEERROR = "PVCVALUEERROR"
	// PVC PV均存在，但是PV不可用
	IsPVC_PV_exist_return_PVUNAVAILABLE = "PVUNAVAILABLE"
	// PVC和PV均存在，但是PVC下没有绑定PV,我们不允许重名PV，所以请删除老的PV,或者等待其绑定完成
	IsPVC_PV_exist_return_PVNOTBIND = "PV_NOTBIND"

	// 正常的fault，可以补救一下再执行
	// PVC正常，但是需要的PVC下PV不存在（handler会要求controller立刻创建）
	IsPVC_PV_exist_Return_PVMISS = "PVMISS"

	// 很正常
	IsPVC_PV_exist_Return_OK = "OK"
)

func IsPVC_PV_exist(namespace string, pvcname string, pvname string) (string, persistVolume.PersistVolumeClaim) {
	// 先检查这个pvc是否存在，且绑定了这个pv
	pvcvalue, err := EtcdCli.Get(pvcPrefix + namespace + "/" + pvcname)
	if err != nil {
		fmt.Printf("PVC does not exist in IsPVC_PV_exist, namespace %s, name %s\n", namespace, pvcname)
		return IsPVC_PV_exist_Return_PVCMISS, persistVolume.PersistVolumeClaim{}
	}
	var pvc persistVolume.PersistVolumeClaim
	err = json.Unmarshal(pvcvalue, &pvc)
	if err != nil {
		fmt.Printf("can not unmarshal pvc in IsPVC_PV_exist,error msg: %s\n", err.Error())
		return IsPVC_PV_exist_Return_PVCVALUEERROR, persistVolume.PersistVolumeClaim{}
	}
	pv_is_in_pvc_bind := false
	for _, pvbind_item := range pvc.Status.BoundPV {
		if pvname == pvbind_item {
			pv_is_in_pvc_bind = true
		}
	}

	// 检查这个pv是否存在
	result_str,_ := IsPVAvailable(namespace, pvname)
	switch result_str {
	case IsPVAvailable_Return_MISS:
		{
			// PV实际etcd缺失，则直接判断为PVMISS，外层需要立刻想conreoller请求
			return IsPVC_PV_exist_Return_PVMISS, pvc
		}
	case IsPVAvailable_Return_VALUEERROR:
		{
			// PV值不可UnMarshal,认为也是存在但UNAVAILABLE的
			return IsPVC_PV_exist_return_PVUNAVAILABLE, persistVolume.PersistVolumeClaim{}
		}
	case IsPVAvailable_Return_UNAVAILABLE:
		{
			// PV存在有效但UNAVAILABLE,要么才创建，要么Release,都不可用，反馈给用户，拒绝pod创建
			return IsPVC_PV_exist_return_PVUNAVAILABLE, persistVolume.PersistVolumeClaim{}
		}
	case IsPVAvailable_Return_OK:
		{
			// PV存在且十分正常
			if pv_is_in_pvc_bind {
				// 绑定好了，外层直接可用
				return IsPVC_PV_exist_Return_OK, persistVolume.PersistVolumeClaim{}
			} else {
				// 未绑定好，那也不可用（因为不允许重名）
				return IsPVC_PV_exist_return_PVNOTBIND, persistVolume.PersistVolumeClaim{}
			}
		}
	}

	return IsPVC_PV_exist_Return_OK, persistVolume.PersistVolumeClaim{}
}

func AddPVImmediately(pvname string, pvc persistVolume.PersistVolumeClaim) error {
	// 先立刻请求创建pv

	requestMsg := httpobject.HTTPReuqest_AddPVImmediately{
		PvName: pvname,
		PvType: pvc.Spec.Type,
	}
	status, err := httpRequest.PostRequestByObject("http://localhost:8082/api/v1/AddPVImmediately", requestMsg, nil)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error post AddPVImmediately, status %d, return\n", status)
		return err
	}
	// 成功后，pv进入etcd
	new_pv := persistVolume.PersistVolume{}
	new_pv.ApiVersion = pvc.ApiVersion
	new_pv.Kind = pvc.Kind
	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Printf("failed to generate uuid for PV in AddPVImmediately\n")
		return err
	}
	new_pv.Metadata.Id = id
	new_pv.Metadata.Name = pvname
	new_pv.Spec.Type = pvc.Spec.Type
	new_pv.Spec.Capacity = pvc.Spec.Capacity
	new_pv.Status.Phase = persistVolume.PV_PHASE_AVAILABLE
	new_pv.Status.MountPod = []string{}

	// 为newpv打标签
	new_pv.Metadata.Labels = make(map[string]string)
	for k, v := range pvc.Spec.Selector.MatchLabels {
		new_pv.Metadata.Labels[k] = v
	}

	// newpv序列化存入etcd
	mresult, err := json.Marshal(new_pv)
	if err != nil {
		fmt.Printf("can't marshal new pv in AddPVImmediately %s\n", err.Error())
		return err
	}
	err = EtcdCli.Put(pvPrefix+DefaultNamespace+"/"+pvname, string(mresult))
	if err != nil {
		fmt.Println("failed to add pv to etcd in AddPVImmediately")
		return err
	}

	return nil
}
