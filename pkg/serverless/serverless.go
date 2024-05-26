package serverless

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/app"
)

func StartServerless() {
	fmt.Printf("Hello Serverless!\n")
	server := app.NewServerlessServer()
	server.Init()
	server.Run()
}
