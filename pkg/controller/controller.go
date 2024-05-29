package controller

import (
	"flag"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/app"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/config"
)

// NOTE: 由于DNS/反向代理需要，nginx必须部署在和controller相同的节点
// NOTE: 由于PV目前的需要，controller所在节点必须为nfs server
// NOTE: 由于PV/PVC和DNS/反向代理需要，controller必须以sudo能力运行
func StartController() {
	apiServerIP := flag.String("apiserverip", "127.0.0.1", "APIServer IP address")
	apiServerPort := flag.String("apiserverport", "8080", "APIServer port")

	config.HTTPURL = fmt.Sprintf(config.HTTPURL_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_AddorDeleteEndpoint = fmt.Sprintf(config.HTTPURL_AddorDeleteEndpoint_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_GetAllDNS = fmt.Sprintf(config.HTTPURL_GetAllDNS_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_GetAllPod = fmt.Sprintf(config.HTTPURL_GetAllPod_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_GetAllService = fmt.Sprintf(config.HTTPURL_GetAllService_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_UpdateDNS = fmt.Sprintf(config.HTTPURL_UpdateDNS_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_UpdateService = fmt.Sprintf(config.HTTPURL_UpdateService_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_GetAllPersistVolume = fmt.Sprintf(config.HTTPURL_GetAllPersistVolume_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_UpdatePersistVolume = fmt.Sprintf(config.HTTPURL_UpdatePersistVolume_Template, *apiServerIP, *apiServerPort)

	fmt.Printf("Hello Controller!\n")
	controller := app.NewController()
	controller.Init()
	controller.Run()
}
