package endpointsController

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/selectorUtils"
	"github.com/google/uuid"
)

type EndpointsController struct {
	// 记录上一次service id绑定的pod **name**组,这对service version的更新有帮助
	last_service2pod map[string][]string
}

func NewEndpointsController() *(EndpointsController) {
	fmt.Printf("New EndpointsController...\n")
	return &EndpointsController{
		last_service2pod: map[string]([]string){},
	}
}

func (ec *EndpointsController) Init() {
	fmt.Printf("Init EndpointsController ...\n")

}

func (ec *EndpointsController) Run() {
	fmt.Printf("Run EndpointsController ...\n")
	poller.PollerStaticPeriod(1*time.Second, ec.routine, true)
}

func isSameNameList(old_nl *[]string,new_nl *[]string)bool{
	if(len(*old_nl)!=len(*new_nl))	{
		return false
	}

	for _,ov:=range(*old_nl){
		canfind:=false
		for _,nv:=range(*new_nl){
			if(ov==nv){
				canfind=true
			}
		}
		if(canfind==false){
			return false
		}
	}
	return true
}


func (ec *EndpointsController) routine() {
	fmt.Printf("EndpointsController routine ...\n")

	// 获得所有service
	var service_list httpobject.HTTPResponse_GetAllServices = httpobject.HTTPResponse_GetAllServices{}
	status, err := httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllService", nil, &service_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("EndpointsController routine error get, status %d, return\n", status)
		return
	}

	// 获得所有pod
	var pod_list httpobject.HTTPResponse_GetAllPod
	status, err = httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllPod", nil, &pod_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}

	// 待更新的service列表(请求状态)
	renew_service_request := httpobject.HTTPRequest_UpdateServices{}
	// 待更新的clusterip列表
	renew_clusterip_list:=[]service.ClusterIP{}
	// 待更新的nodeport列表

	// // 本轮service与pod的绑定状态
	// new_service2pod := map[string][]string{}
	// 准备调整的endpoints列表
	endponits_list := httpobject.HTTPRequest_AddorDeleteEndpoint{
		Add: []service.EndPoint{},
		Delete: []string{},
	}



	// 遍历每个clusterIP
	for _,clusterIP:=range(service_list.ClusterIP){
		// 如果未分配id,则不予更新
		if(clusterIP.Metadata.Ip==""){
			continue
		}

		new_c2p:=selectorUtils.SelectPodNameList(&clusterIP.Spec.Selector,&pod_list)
		if old_c2p,exist:=ec.last_service2pod[clusterIP.Metadata.Id];!exist{
			// 之前不存在这个clusterip,是本轮新加入的
			// 加入绑定集合
			ec.last_service2pod[clusterIP.Metadata.Id]=new_c2p
			// clusterip加入待更新集合
			renew_clusterip_list=append(renew_clusterip_list, clusterIP)
		}else{
			//之前已经存在这个clusterip,这轮需要查看是否更新
			if(isSameNameList(&old_c2p,&new_c2p)){
				//相等，不需要更新这个clusterip
			}else{
				//不相等，需要更新这个clusterip
				// 更新绑定集合
				ec.last_service2pod[clusterIP.Metadata.Id]=new_c2p
				// clusterip加入待更新集合
				renew_clusterip_list=append(renew_clusterip_list, clusterIP)
				// 原有的endponit全部加入待删除集合（粒度较粗）
				for _,elist:=range(clusterIP.Status.ServicesStatus){
					endponits_list.Delete=append(endponits_list.Delete, elist...)
				}	
			}
		}
	}

	// 遍历所有待更新的clusterip,更新版本号，设置READY,分配新的endpoints
	for _,clusterip:=range(renew_clusterip_list){
		clusterip.Status.Phase=service.CLUSTERIP_READY
		clusterip.Status.Version+=1

		// 清空ServicesStatus
		clusterip.Status.ServicesStatus=map[uint16][]string{}

		for _,port_info:=range(clusterip.Spec.Ports){
			ep_id_list:=[]string{}
			// 为每个选中的pod创建endpoint
			for _,podname:=range(ec.last_service2pod[clusterip.Metadata.Id]){
				// 创建新ep
				new_endpoint:=service.EndPoint{
					Id:uuid.NewString(),
					Protocol: port_info.Protocal,
					PodID: pod_list[podname].Metadata.UUID,
					PodIP: pod_list[podname].Status.PodIP,
					PodPort: port_info.TargetPort,
				}
				// ep加入http add
				endponits_list.Add=append(endponits_list.Add, new_endpoint)
				// ep加入cluster的service status对应的port下
				ep_id_list=append(ep_id_list, new_endpoint.Id)
			}
			//更新此port在ServicesStatus下的状态
			clusterip.Status.ServicesStatus[port_info.Port]=ep_id_list
		}

		// 更新完毕，cluster加入http更新object中
		json_byte, _ := json.Marshal(clusterip)
		renew_service_request[clusterip.Metadata.Name]=string(json_byte)
	}

	// 遍历每个nodeport
	for _,nodeport:=range(service_list.NodePort){
		if(nodeport.Status.ClusterIPID==""||nodeport.Status.Phase==service.NODEPORT_READY){
			// 尚未分配clusterip和已经ready的nodeport不需要更新
			continue
		}
		nodeport.Status.Phase=service.NODEPORT_READY
		// 更新完毕，nodeport加入http更新object中
		json_byte, _ := json.Marshal(nodeport)
		renew_service_request[nodeport.Metadata.Name]=string(json_byte)
	}

	// 请求service更新
	status, err = httpRequest.PostRequestByObject("http://192.168.1.6:8080/api/v1/UpdateService", renew_service_request, nil)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}
	// 请求endpoints更新
	status, err = httpRequest.PostRequestByObject("http://192.168.1.6:8080/api/v1/AddorDeleteEndpoint", endponits_list, nil)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}
}
