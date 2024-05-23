package kubelet

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	msgproxy "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet/msg_proxy"
	kubelet_worker "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet/worker"
	minik8s_node "github.com/MiniK8s-SE3356/minik8s/pkg/types/node"
	httpRequest "github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
)

type Kubelet struct {
	kubeletConfig *KubeletConfig
	msgProxy      *msgproxy.MsgProxy
	podManager    kubelet_worker.PodManager
	minik8s_node.Node
}

func NewKubelet(config *KubeletConfig) *Kubelet {
	// We shouldn't create connection until the node is registered into apiserver
	kubelet := &Kubelet{
		kubeletConfig: config,
		podManager: kubelet_worker.NewPodManager(
			kubelet_worker.APIServer{
				IP:   config.APIServerIP,
				Port: config.APIServerPort,
			},
		),
	}
	return kubelet
}

func (k *Kubelet) Proxy() {
	// update := <-k.msgProxy.PodUpdateChannel
	// k.podManager.AddPod(update.Pod)

	// go k.Proxy()
	for update := range k.msgProxy.PodUpdateChannel {
		switch update.Type {
		case kubelet_worker.Task_Add:
			fmt.Println("Kubelet Add pod")
			k.podManager.AddPod(update.Pod, nil)
		case kubelet_worker.Task_Update:
			k.podManager.UpdatePod(update.Pod, nil)
		case kubelet_worker.Task_Remove:
			fmt.Println("Kubelet Remove pod")
			k.podManager.RemovePod(update.Pod, nil)
		}
	}
}

func (k *Kubelet) RegisterNode() error {
	url := fmt.Sprintf("http://%s:%s/api/v1/AddNode", k.kubeletConfig.APIServerIP, k.kubeletConfig.APIServerPort)

	node := minik8s_node.Node{
		Metadata: minik8s_node.NodeMetadata{},
		Status: minik8s_node.NodeStatus{
			Hostname: k.kubeletConfig.NodeHostName,
			Ip:       k.kubeletConfig.NodeIP,
			Condition: []string{
				minik8s_node.NODE_Ready,
			},
		},
	}

	// TODO: response target should include error message?
	var responseNode minik8s_node.Node
	_, err := httpRequest.PostRequestByObject(
		url,
		node,
		&responseNode,
	)
	if err != nil {
		fmt.Println("Register node failed")
		return err
	}

	// We check the response is valid or not in a naive method:
	// only check the parsed responseNode's metadata is empty or not
	// if empty, we think the response is invalid
	if responseNode.Metadata.Id == "" {
		fmt.Println("Register node failed")
		return fmt.Errorf("register node failed")
	}

	json_responseNode, _ := json.Marshal(responseNode)
	fmt.Println("Register node success: \n", string(json_responseNode))

	k.Node = responseNode

	// We have got the node id, so we can create the connection to RabbitMQ queue now
	k.msgProxy = msgproxy.NewMsgProxy(
		&k.kubeletConfig.MQConfig,
		k.Node.Metadata.Name,
	)

	return nil
}

func (k *Kubelet) GetNodeStatus() (minik8s_node.NodeStatus, error) {
	nodeStatus, err := kubelet_worker.NodeRuntimeMangaer.GetNodeStatus()
	if err != nil {
		return minik8s_node.NodeStatus{}, err
	}

	nodeStatus.NumPods = k.podManager.GetPodNum()

	return nodeStatus, nil
}

func (k *Kubelet) HeartBeat() {
	nodeStatus, err := k.GetNodeStatus()
	if err != nil {
		return
	}
	k.Node.Status = nodeStatus

	pods, err := k.podManager.FetchLocalPods()
	if err != nil {
		return
	}

	request_url := fmt.Sprintf("http://%s:%s%s", k.kubeletConfig.APIServerIP, k.kubeletConfig.APIServerPort, url.NodeHeartBeat)
	request_body := make(map[string]interface{})
	request_body["nodeStatus"] = k.Node.Status
	request_body["pods"] = pods
	request_body_data, _ := json.Marshal(request_body)
	response, err := httpRequest.PostRequest(
		request_url,
		request_body_data,
	)
	if err != nil {
		fmt.Println("Error posting request: ", err)
		return
	}
	fmt.Println("\nHeartbeat response: ", response)
}

func (k *Kubelet) Run() {
	// TODO: run a cAdvisor container to monitor the node and container status
	forever := make(chan bool)

	err := k.RegisterNode()
	if err != nil {
		return
	}

	go k.Proxy()
	go k.msgProxy.Run()

	go poller.PollerStaticPeriod(
		time.Duration(10*time.Second),
		k.HeartBeat,
		true,
	)

	<-forever
}
