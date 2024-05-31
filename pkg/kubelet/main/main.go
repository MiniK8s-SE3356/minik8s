package main

import (
	"fmt"
	"os"
	"time"

	minik8s_kubelet "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet"
	minik8s_runtime "github.com/MiniK8s-SE3356/minik8s/pkg/runtime"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/nettools"
	"github.com/spf13/pflag"
)

func main() {

	// Get default node ip and hostname
	defaultNodeIP := nettools.KubeletDefaultIP()
	defaultHostname, err := os.Hostname()
	if err != nil {
		defaultHostname = "node"
	}
	// Get command line arguments
	// Command line arguments must include API server address and port
	apiServerIP := pflag.String("apiserverip", "127.0.0.1", "APIServer IP address")
	apiServerPort := pflag.String("apiserverport", "8080", "APIServer port")
	nodeIP := pflag.String("nodeip", defaultNodeIP, "Node IP address")
	hostName := pflag.String("hostname", defaultHostname, "Node hostname")
	exporter := pflag.String("exporter", "", "port of exporter")

	// Get command line arguments about labels
	var labels []string
	pflag.StringArrayVar(&labels, "label", []string{}, "Node labels")
	pflag.Parse()

	// run a cAdvisor container to monitor the node and container status
	minik8s_runtime.NodeRuntimeMangaer.RuncAdvicorContainer()

	kubeletConfig := minik8s_kubelet.KubeletConfig{
		MQConfig: minik8s_message.MQConfig{
			User:       "guest",
			Password:   "guest",
			Host:       *apiServerIP,
			Port:       "5672",
			Vhost:      "/",
			MaxRetry:   5,
			RetryDelay: 5 * time.Second,
		},
		APIServerIP:   *apiServerIP,
		APIServerPort: *apiServerPort,
		NodeIP:        *nodeIP,
		NodeHostName:  *hostName,
	}

	fmt.Println(
		"Kubelet Config: \n",
		"API Server IP: ", kubeletConfig.APIServerIP,
		"\nAPI Server Port: ", kubeletConfig.APIServerPort,
		"\nNode IP: ", kubeletConfig.NodeIP,
		"\nNode Hostname: ", kubeletConfig.NodeHostName,
	)

	kubelet := minik8s_kubelet.NewKubelet(&kubeletConfig)
	kubelet.Metadata.Labels = make(map[string]string)
	fmt.Println(*exporter)
	if *exporter != "" {
		kubelet.Metadata.Labels["metric_port"] = *exporter
	}

	kubelet.Run()
}
