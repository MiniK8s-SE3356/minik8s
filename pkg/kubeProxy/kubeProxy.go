package kubeProxy

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/app"
)

func StartKubeProxy() {

	fmt.Printf("Hello KubeProxy!\n")
	kube_proxy := app.NewKubeProxy()
	kube_proxy.Init()
	kube_proxy.Run()
}
