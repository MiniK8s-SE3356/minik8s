package controller

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/app"
)

func StartController() {

	fmt.Printf("Hello Controller!\n")
	controller := app.NewController()
	controller.Init()
	controller.Run()
}
