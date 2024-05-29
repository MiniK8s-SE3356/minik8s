package persistVolumeController

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/persistVolume"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/config"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/selectorUtils"
)

const (
	MINIK8S_PV_PATH = "/var/lib/minik8s/volumes"
)

type PersistVolumeController struct {
	// // nfspv name -> int(1)
	// // we only care if nfspv exits
	// nfs_pv_list map[string]int
}

func NewersistVolumeController() *PersistVolumeController {
	fmt.Printf("New PersistVolumeController...\n")
	return &PersistVolumeController{}
}

func (pvc *PersistVolumeController) Init() {
	fmt.Printf("Init PersistVolumeController...\n")
	// pvc.nfs_pv_list = make(map[string]int)
}

func (pvc *PersistVolumeController) Run() {
	fmt.Printf("Run PersistVolumeController...\n")
	poller.PollerStaticPeriod(5*time.Second, pvc.routine, true)
}

func (pvc *PersistVolumeController) routine() {
	// 请求所有的pod
	var pod_list httpobject.HTTPResponse_GetAllPod
	status, err := httpRequest.GetRequestByObject(config.HTTPURL_GetAllPod, nil, &pod_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error GetAllPod, status %d, return\n", status)
		return
	}
	// 请求所有的pv/pvc
	var pv_pvc_list httpobject.HTTPResponse_GetAllPersistVolume
	status, err = httpRequest.GetRequestByObject(config.HTTPURL_GetAllPersistVolume, nil, &pv_pvc_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}

	// 分为两部分，一部分将当前请求下的pv中没有创建持久化卷的创建，另一部分负责pod pv pvc三者绑定

	// 持久化卷创建
	new_available_pv := pvc.syncPV(&(pv_pvc_list.Pv))

	// pod pv pvc 三者绑定
	pv_pvc_renew_list := httpobject.HTTPRequest_UpdatePersistVolume{
		Pv:  make(map[string]persistVolume.PersistVolume),
		Pvc: make(map[string]persistVolume.PersistVolumeClaim),
	}

	// 首先将pv绑定pod
	pv2pod_list := map[string][]string{}
	// 遍历pod,取出对应的映射填写pv2pod_list
	for podname, podvalue := range pod_list {
		for _, volumevalue := range podvalue.Spec.Volumes {
			if volumevalue.PersistentVolumeClaim.ClaimName != "" || volumevalue.PersistentVolume.PvName != "" {
				selectorUtils.AddToMayNilMap(volumevalue.Name, podname, &pv2pod_list)
			}
		}
	}

	for pvname, pvvalue := range pv_pvc_list.Pv {
		need_renew := false
		if selectorUtils.IsStrInArray(pvname, &new_available_pv) {
			pvvalue.Status.Phase = persistVolume.PV_PHASE_AVAILABLE
			need_renew = true
		}

		if new_pod_list, exist := pv2pod_list[pvname]; exist {
			if !selectorUtils.IsSameSet(&new_pod_list, &pvvalue.Status.MountPod) {
				// 如果本轮podlist存在，但不一致，要重新赋值
				pvvalue.Status.MountPod = new_pod_list
				need_renew = true
			}
		} else {
			// 如果本轮podlist不存在，但是之前存在或nil,则改为空数组
			if pvvalue.Status.MountPod == nil || len(pvvalue.Status.MountPod) != 0 {
				pvvalue.Status.MountPod = []string{}
				need_renew = true
			}
		}

		if need_renew {
			pv_pvc_renew_list.Pv[pvname] = pvvalue
		}

	}

	// 其次pvc绑定pv
	for pvcname, pvcvalue := range pv_pvc_list.Pvc {
		need_renew := false
		new_pv_list := selectorUtils.SelectPVNameList(&pvcvalue.Spec.Selector, &pv_pvc_list.Pv)
		if pvcvalue.Status.BoundPV == nil || !selectorUtils.IsSameSet(&new_pv_list, &pvcvalue.Status.BoundPV) {
			pvcvalue.Status.BoundPV = new_pv_list
			need_renew = true
		}
		if need_renew {
			pv_pvc_renew_list.Pvc[pvcname] = pvcvalue
		}
	}

	// 更新pv和pvc
	status, err = httpRequest.PostRequestByObject(config.HTTPURL_UpdatePersistVolume, pv_pvc_renew_list, nil)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error UpdatePersistVolume, status %d, return\n", status)
		return
	}

	return
}

// syncPV 	会根据请求到的etcd中pv数据更新pv对应的底层存储对象
//
//					目前仅支持对nfs PV的管理
//					目前仅涉及添加不存在的PV,删除以后再做
//		 @receiver pvc
//		 @param pv_list
//		 @param wg
//	  @return 返回的卷将由CREATED变为AVAILABLE
func (pvc *PersistVolumeController) syncPV(pv_list *(map[string]persistVolume.PersistVolume)) []string {
	result := []string{}
	for key, pv := range *pv_list {
		// 非PV_PHASE_CREATED的都已经有卷了，不必新建
		if pv.Status.Phase != persistVolume.PV_PHASE_CREATED {
			continue
		}
		switch pv.Spec.Type {
		case persistVolume.PV_TYPE_NFS:
			{
				// if _, exist := pvc.nfs_pv_list[key]; !exist {
				err := pvc.createNFSVolume(key)
				if err == nil {
					// 本地记录
					// pvc.nfs_pv_list[key] = 1
					result = append(result, key)
				}
				// }
				break
			}
		default:
			{
				fmt.Printf("Unknown PV type %s\n", pv.Spec.Type)
				break
			}
		}
	}
	return result
}

func (pvc *PersistVolumeController) CreatePVImmediately(pvName string, pvType string) error {
	fmt.Printf("Create PV Immediately...\n")
	switch pvType {
	case persistVolume.PV_TYPE_NFS:
		{
			return pvc.createNFSVolume(pvName)
		}
	}
	return errors.New("unknown pv type")
}

func (pvc *PersistVolumeController) createNFSVolume(pvName string) error {
	cmd := exec.Command("mkdir", MINIK8S_PV_PATH+"/"+pvName+"/_data", "-p")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command in createNFSVolume:", err)
		return err
	}
	fmt.Printf("createVolume succeed! Output: %s\n", output)
	return nil
}
