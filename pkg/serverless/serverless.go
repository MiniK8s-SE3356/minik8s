package serverless

import (
	"flag"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/app"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/config"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/registry"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/nettools"
)

func StartServerless() {
	fmt.Printf("Hello Serverless!\n")

	apiServerIP := flag.String("apiserverip", "127.0.0.1", "APIServer IP address")
	apiServerPort := flag.String("apiserverport", "8080", "APIServer port")

	defaultRegistryIP := nettools.KubeletDefaultIP()
	registryIP := flag.String("registryip", defaultRegistryIP, "Registry IP address")
	registry.REGISTRY_IP = *registryIP
	fmt.Printf("Registry IP: %s\n", *registryIP)

	config.HTTPURL_AddServerlessFuncPod = fmt.Sprintf(config.HTTPURL_AddServerlessFuncPod_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_GetAllPod = fmt.Sprintf(config.HTTPURL_GetAllPod_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_GetAllServerlessFunction = fmt.Sprintf(config.HTTPURL_GetAllServerlessFunction_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_root = fmt.Sprintf(config.HTTPURL_root_Template, *apiServerIP, *apiServerPort)

	server := app.NewServerlessServer()
	server.Init()
	server.Run()
}
