package servicesController

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	// TODO: 只做了ClusterIP ，NodePort相关还没添加

	var response_object httpobject.HTTPResponse_GetAllServices

	status, err := httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllService", nil, &response_object)

	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}

	fmt.Println(response_object)

	update_map := make(httpobject.HTTPRequest_UpdateServices)
	for _, value := range response_object.ClusterIP {
		if value.Metadata.Ip == "" {

			value.Metadata.Ip = sc.Ipallocater.AllocateIP()
			// value.Status.Phase = service.CLUSTERIP_IP_ALLOCATED/*只允许此文件上述的和状态有关的const常量，可参考飞书《Service设计方案》*/

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
