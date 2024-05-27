package controller

import (
	"flag"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/app"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/config"
)

// NOTE: 由于DNS/反向代理需要，nginx必须部署在和controller相同的
func StartController() {
	apiServerIP := flag.String("apiserverip", "127.0.0.1", "APIServer IP address")
	apiServerPort := flag.String("apiserverport", "8080", "APIServer port")

	config.HTTPURL_AddorDeleteEndpoint=fmt.Sprintf(config.HTTPURL_AddorDeleteEndpoint_Template,*apiServerIP,*apiServerPort)
	config.HTTPURL_GetAllDNS=fmt.Sprintf(config.HTTPURL_GetAllDNS_Template,*apiServerIP,*apiServerPort)
	config.HTTPURL_GetAllPod=fmt.Sprintf(config.HTTPURL_GetAllPod_Template,*apiServerIP,*apiServerPort)
	config.HTTPURL_GetAllService=fmt.Sprintf(config.HTTPURL_GetAllService_Template,*apiServerIP,*apiServerPort)
	config.HTTPURL_UpdateDNS=fmt.Sprintf(config.HTTPURL_UpdateDNS_Template,*apiServerIP,*apiServerPort)
	config.HTTPURL_UpdateService=fmt.Sprintf(config.HTTPURL_UpdateService_Template,*apiServerIP,*apiServerPort)

	fmt.Printf("Hello Controller!\n")
	controller := app.NewController()
	controller.Init()
	controller.Run()
}
