package app

import (
	"fmt"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/endpointsController"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/nodesController"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/servicesController"
)

type Controller struct {
	endpointsController *(endpointsController.EndpointsController)
	nodesController     *(nodesController.NodesController)
	servicesController  *(servicesController.ServicesController)
}

func NewController() *(Controller) {
	fmt.Printf("New Controller...\n")
	endpoints_controller := endpointsController.NewEndpointsController()
	nodes_controller := nodesController.NewNodesController()
	services_controller := servicesController.NewServicesController()
	return &Controller{
		endpointsController: endpoints_controller,
		nodesController:     nodes_controller,
		servicesController:  services_controller,
	}
}

func (co *Controller) Init() {
	fmt.Printf("Init Controller ...\n")
	co.endpointsController.Init()
	co.nodesController.Init()
	co.servicesController.Init()
}

func (co *Controller) Run() {
	fmt.Printf("Run Controller ...\n")
	go co.endpointsController.Run()
	go co.nodesController.Run()
	go co.servicesController.Run()

	time.Sleep(2 * time.Second)
}
