package servicesController

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/servicesController/ipAllocater"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
)

type ServicesController struct {
	Ipallocater ipAllocater.IpAllocater
}

func NewServicesController() *(ServicesController) {
	fmt.Printf("New ServicesController...\n")
	return &ServicesController{
		Ipallocater: ipAllocater.IpAllocater{},
	}
}

func (sc *ServicesController) Init() {
	fmt.Printf("Init ServicesController ...\n")
	// 获取子网段和子网掩码位数
	sc.Ipallocater.Init()
}

func (sc *ServicesController) Run() {
	fmt.Printf("Run ServicesController ...\n")
	poller.PollerStaticPeriod(3*time.Second, sc.routine, true)
	// can not reach here
}

func (sc *ServicesController) routine() {
	fmt.Printf("ServicesController routine\n")

	// NOTE: ServiceController只分配ClusterIP的虚拟IP,对于NodePort没有影响

	var response_object httpobject.HTTPResponse_GetAllServices

	status, err := httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllService", nil, &response_object)

	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}

	fmt.Println(response_object)

	update_map := httpobject.HTTPRequest_UpdateServices{}

	// 遍历nodeport,为nodeport分配对应的ClusterIP
	for _,value:=range response_object.NodePort{
		if (value.Status.ClusterIPID==""){
			// 分配一个新的clusterIP
			new_clu:=service.NodePort2ClusterIP(&value)
			value.Status.ClusterIPID=new_clu.Metadata.Id
			
			response_object.ClusterIP=append(response_object.ClusterIP, new_clu)

			//nodeport加入更新队列
			json_byte, _ := json.Marshal(value)
			update_map[value.Metadata.Name] = string(json_byte)
		}
	}
	

	// 遍历clusterIP，分配虚拟IP
	// 但是注意clusterIP还没有READY,其需要绑定endponits后才会READY
	// 也要注意，由于clusterIP没有READY,其（可能的）绑定对应的NodePort也不可能READY
	for _, value := range response_object.ClusterIP {
		if value.Metadata.Ip == "" {
			value.Metadata.Ip = sc.Ipallocater.AllocateIP()
			// // value.Status.Phase = service.CLUSTERIP_IP_ALLOCATED/*只允许此文件上述的和状态有关的const常量，可参考飞书《Service设计方案》*/
			json_byte, _ := json.Marshal(value)
			update_map[value.Metadata.Name] = string(json_byte)
		}
	}

	status, err = httpRequest.PostRequestByObject("http://192.168.1.6:8080/api/v1/UpdateService", update_map, nil)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}

}
