package endpointsController

import (
	"fmt"
	"net/http"
	"time"

	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
)

type EndpointsController struct {
	// Endpoints
	eid_list []string
}

func NewEndpointsController() *(EndpointsController) {
	fmt.Printf("New EndpointsController...\n")
	return &EndpointsController{}
}

func (ec *EndpointsController) Init() {
	fmt.Printf("Init EndpointsController ...\n")

}

func (ec *EndpointsController) Run() {
	fmt.Printf("Run EndpointsController ...\n")
	poller.PollerStaticPeriod(1*time.Second, ec.routine, true)
}

func (ec *EndpointsController) routine() {
	fmt.Printf("EndpointsController routine ...\n")

	// 获得所有service
	var response_object httpobject.HTTPResponse_GetAllServices = httpobject.HTTPResponse_GetAllServices{}
	status, err := httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllService", nil, &response_object)

	if status != http.StatusOK || err != nil {
		fmt.Printf("EndpointsController routine error get, status %d, return\n", status)
		return
	}

	// 获得所有pod

	// 遍历每个service,通过selector筛选出对应的pod，并为其创建endponits端口，并修改状态

	// service更新

	// endpoints更新

}
