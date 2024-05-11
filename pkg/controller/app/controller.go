package app

import (
	"fmt"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/endpointsController"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/nodesController"
)

type Controller struct {
	endpointsController *(endpointsController.EndpointsController)
	nodesController     *(nodesController.NodesController)
}

func NewController() *(Controller) {
	fmt.Printf("New Controller...\n")
	endpoints_controller := endpointsController.NewEndpointsController()
	nodes_controller := nodesController.NewNodesController()
	return &Controller{
		endpointsController: endpoints_controller,
		nodesController:     nodes_controller,
	}
}

func (co *Controller) Init() {
	fmt.Printf("Init Controller ...\n")
	co.endpointsController.Init()
	co.nodesController.Init()
}

func (co *Controller) Run() {
	fmt.Printf("Run Controller ...\n")
	go co.endpointsController.Run()
	go co.nodesController.Run()

	time.Sleep(2 * time.Second)
}
