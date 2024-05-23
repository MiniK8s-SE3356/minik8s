package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	minik8s_kubelet "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet"
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
	apiServerIP := flag.String("apiserverip", "127.0.0.1", "APIServer IP address")
	apiServerPort := flag.String("apiserverport", "8080", "APIServer port")
	nodeIP := flag.String("nodeip", defaultNodeIP, "Node IP address")
	hostName := flag.String("hostname", defaultHostname, "Node hostname")

	// Get command line arguments about labels
	var labels []string
	pflag.StringArrayVar(&labels, "label", []string{}, "Node labels")
	pflag.Parse()

	// TODO : Run cAdvisor container

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

	kubelet.Run()
}
