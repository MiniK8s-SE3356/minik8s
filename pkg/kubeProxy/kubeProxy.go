package kubeProxy

import (
	"flag"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/app"
	"github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/config"
)

func StartKubeProxy() {
	apiServerIP := flag.String("apiserverip", "127.0.0.1", "APIServer IP address")
	apiServerPort := flag.String("apiserverport", "8080", "APIServer port")
	nginxIP := flag.String("nginxip", "127.0.0.1", "DNS/Proxy Nginx IP address")
	flag.Parse()

	config.HTTPURL_GetAllDNS = fmt.Sprintf(config.HTTPURL_GetAllDNS_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_GetAllEndpoint = fmt.Sprintf(config.HTTPURL_GetAllEndpoint_Template, *apiServerIP, *apiServerPort)
	config.HTTPURL_GetAllService = fmt.Sprintf(config.HTTPURL_GetAllService_Template, *apiServerIP, *apiServerPort)
	config.NGINX_IP = *nginxIP

	fmt.Printf("Hello KubeProxy!\n")
	fmt.Println("GetAllDNS",config.HTTPURL_GetAllDNS)
	fmt.Println("GetAllEndpoint",config.HTTPURL_GetAllEndpoint)
	fmt.Println("GetAllService",config.HTTPURL_GetAllService)
	fmt.Println("NGINX_IP",config.NGINX_IP)
	kube_proxy := app.NewKubeProxy()
	kube_proxy.Init()
	kube_proxy.Run()
}
