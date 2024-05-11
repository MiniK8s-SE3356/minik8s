package controller

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/app"
)

func StartController() {

	fmt.Printf("Hello Controller!\n")
	kube_proxy := app.NewController()
	kube_proxy.Init()
	kube_proxy.Run()
}
