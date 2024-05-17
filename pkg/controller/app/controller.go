package app

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/endpointsController"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/servicesController"
)

type Controller struct {
	endpointsController *(endpointsController.EndpointsController)
	servicesController  *(servicesController.ServicesController)
}

func NewController() *(Controller) {
	fmt.Printf("New Controller...\n")
	endpoints_controller := endpointsController.NewEndpointsController()
	services_controller := servicesController.NewServicesController()
	return &Controller{
		endpointsController: endpoints_controller,
		servicesController:  services_controller,
	}
}

func (co *Controller) Init() {
	fmt.Printf("Init Controller ...\n")
	co.endpointsController.Init()
	co.servicesController.Init()
}

func (co *Controller) Run() {
	fmt.Printf("Run Controller ...\n")
	go co.endpointsController.Run()
	go co.servicesController.Run()

	// TODO:主线程暂时没有要做的事情，先while1,后续需要补充
	for {

	}
}
