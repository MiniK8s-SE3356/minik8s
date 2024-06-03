package app

import (
	"fmt"
	"os"
	"time"

	minik8s_kubelet "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet"
	minik8s_runtime "github.com/MiniK8s-SE3356/minik8s/pkg/runtime"
	"github.com/MiniK8s-SE3356/minik8s/pkg/types/node"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/nettools"
	"gopkg.in/yaml.v3"
)

func KubeletExec() {

	// Get default node ip and hostname
	defaultNodeIP := nettools.KubeletDefaultIP()
	defaultHostname, err := os.Hostname()
	if err != nil {
		defaultHostname = "node"
	}

	args := os.Args

	//确保只有一个参数
	if len(args) != 2 {
		fmt.Println("Usage: ./kubelet [node.yaml]")
		return
	}

	// Get command line arguments
	nodeYamlFilePath := args[1]
	yamlFile, err := os.ReadFile(nodeYamlFilePath)
	if err != nil {
		fmt.Println("failed to read yaml file")
		return
	}

	//!debug//
	fmt.Println("nodeYamlFile:", string(yamlFile))
	//!debug//

	var nodeStartInfo node.NodeStartInfo
	err = yaml.Unmarshal(yamlFile, &nodeStartInfo)
	if err != nil {
		fmt.Println("failed to unmarshal yaml file, err:", err)
		return
	}

	apiServerIP := "127.0.0.1"
	apiServerPort := "8080"
	nodeIP := defaultNodeIP
	hostName := defaultHostname

	if nodeStartInfo.APIServerIP != "" {
		apiServerIP = nodeStartInfo.APIServerIP
	}
	if nodeStartInfo.APIServerPort != "" {
		apiServerPort = nodeStartInfo.APIServerPort
	}
	if nodeStartInfo.NodeIP != "" {
		nodeIP = nodeStartInfo.NodeIP
	}
	if nodeStartInfo.NodeHostName != "" {
		hostName = nodeStartInfo.NodeHostName
	}

	// // Get command line arguments
	// // Command line arguments must include API server address and port
	// apiServerIP := flag.String("apiserverip", "127.0.0.1", "APIServer IP address")
	// apiServerPort := flag.String("apiserverport", "8080", "APIServer port")
	// nodeIP := flag.String("nodeip", defaultNodeIP, "Node IP address")
	// hostName := flag.String("hostname", defaultHostname, "Node hostname")

	// // Get command line arguments about labels
	// var labels []string
	// pflag.StringArrayVar(&labels, "label", []string{}, "Node labels")
	// pflag.Parse()

	// run a cAdvisor container to monitor the node and container status
	minik8s_runtime.NodeRuntimeMangaer.RuncAdvicorContainer()

	kubeletConfig := minik8s_kubelet.KubeletConfig{
		MQConfig: minik8s_message.MQConfig{
			User:       "guest",
			Password:   "guest",
			Host:       apiServerIP,
			Port:       "5672",
			Vhost:      "/",
			MaxRetry:   5,
			RetryDelay: 5 * time.Second,
		},
		APIServerIP:   apiServerIP,
		APIServerPort: apiServerPort,
		NodeIP:        nodeIP,
		NodeHostName:  hostName,
		Labels:        nodeStartInfo.Labels,
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
