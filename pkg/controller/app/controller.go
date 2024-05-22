package app

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/dnsController"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/endpointsController"

	replicasetcontroller "github.com/MiniK8s-SE3356/minik8s/pkg/controller/replicasetController"

	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/servicesController"
)

type Controller struct {
	endpointsController  *(endpointsController.EndpointsController)
	servicesController   *(servicesController.ServicesController)
	replicasetController *(replicasetcontroller.ReplicasetController)
	dnsController        *(dnsController.DnsController)
}

func NewController() *(Controller) {
	fmt.Printf("New Controller...\n")
	endpoints_controller := endpointsController.NewEndpointsController()
	services_controller := servicesController.NewServicesController()
	replicaset_controller := replicasetcontroller.NewReplicasetController()
	dns_controller := dnsController.NewDnsController()
	return &Controller{
		endpointsController:  endpoints_controller,
		servicesController:   services_controller,
		replicasetController: replicaset_controller,
		dnsController:        dns_controller,
	}
}

func (co *Controller) Init() {
	fmt.Printf("Init Controller ...\n")
	co.endpointsController.Init()
	co.servicesController.Init()
	co.dnsController.Init()
}

func (co *Controller) Run() {
	fmt.Printf("Run Controller ...\n")
	go co.endpointsController.Run()
	go co.servicesController.Run()
	go co.replicasetController.Run()

	// TODO:主线程暂时没有要做的事情，先跑dnscontroller,后续需要补充
	co.dnsController.Run()
}
