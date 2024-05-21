package scheduler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	minik8s_pod "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	minik8s_node "github.com/MiniK8s-SE3356/minik8s/pkg/types/node"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/streadway/amqp"
)

var (
	counter int
	lock    sync.Mutex
)

type Policy string

const (
	RoundRobin Policy = "RoundRobin"
	Random     Policy = "Random"
)

type Scheduler struct {
	mqConn        *minik8s_message.MQConnection
	policy        Policy
	APIServerIP   string
	APIServerPort string
}

func NewScheduler(mqConfig *minik8s_message.MQConfig, policy Policy, APIServerIP string, APIServerPort string) (*Scheduler, error) {
	newConn, err := minik8s_message.NewMQConnection(mqConfig)

	if err != nil {
		return nil, err
	}

	scheduler := &Scheduler{
		mqConn:        newConn,
		policy:        policy,
		APIServerIP:   APIServerIP,
		APIServerPort: APIServerPort,
	}
	return scheduler, nil
}

func (s *Scheduler) GetNodes() ([]*minik8s_node.Node, error) {
	// fake
	// nodes := []*minik8s_node.Node{
	// 	{
	// 		Metadata: minik8s_node.NodeMetadata{
	// 			Id:   "1",
	// 			Name: "node1",
	// 			Labels: map[string]string{
	// 				"zone": "zone1",
	// 			},
	// 		},
	// 		Status: minik8s_node.NodeStatus{
	// 			Hostname:   "node1",
	// 			Ip:         "127.0.0.1",
	// 			Condition:  []string{minik8s_node.NODE_Ready},
	// 			CpuPercent: 0.1,
	// 			MemPercent: 0.2,
	// 			NumPods:    1,
	// 			UpdateTime: time.Now().String(),
	// 		},
	// 	},
	// }

	// Implement get nodes from apiserver
	request_url := fmt.Sprintf("http://%s:%s/api/v1/GetNode", s.APIServerIP, s.APIServerPort)
	var responseNodes []minik8s_node.Node
	responseStatus, err := httpRequest.GetRequestByObject(
		request_url,
		nil,
		&responseNodes,
	)
	if err != nil {
		fmt.Println("Get nodes failed")
		return nil, err
	}
	fmt.Println("Get nodes response status: ", responseStatus)

	responseNodesJson, _ := json.Marshal(responseNodes)
	fmt.Println("Get nodes response: \n", string(responseNodesJson))

	// return nodes' pointer array, switch to pointer array
	returnNodes := make([]*minik8s_node.Node, 0)
	for _, node := range responseNodes {
		returnNodes = append(returnNodes, &node)
	}

	return returnNodes, nil
}

func RoundRobinSelect(node []*minik8s_node.Node) *minik8s_node.Node {
	lock.Lock()
	defer lock.Unlock()

	if len(node) == 0 {
		return nil
	}

	if counter >= len(node) {
		counter = 0
	}
	nodeToReturn := node[counter]
	counter++
	return nodeToReturn
}

func RandomSelect(node []*minik8s_node.Node) *minik8s_node.Node {
	lock.Lock()
	defer lock.Unlock()

	if len(node) == 0 {
		return nil
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	return node[r.Intn(len(node))]
}

func (s *Scheduler) SelectNode(node []*minik8s_node.Node) *minik8s_node.Node {
	switch s.policy {
	case RoundRobin:
		return RoundRobinSelect(node)
	case Random:
		return RandomSelect(node)
	default:
		return nil
	}
}

func (s *Scheduler) ScheduleHandler(delivery amqp.Delivery) {
	var parsed_msg map[string]interface{}
	err := json.Unmarshal(delivery.Body, &parsed_msg)

	if err != nil {
		println("Error unmarshalling json")
		return
	}

	// Get all nodes from apiserver
	nodes, err := s.GetNodes()
	if err != nil {
		println("Error getting nodes")
		return
	}

	// Select available nodes
	available_nodes := make([]*minik8s_node.Node, 0)
	for _, node := range nodes {
		for _, condition := range node.Status.Condition {
			if condition == minik8s_node.NODE_Ready {
				available_nodes = append(available_nodes, node)
			}
		}
	}

	var msgType minik8s_message.MsgType
	msgTypeData, _ := json.Marshal(parsed_msg["type"])
	err = json.Unmarshal(msgTypeData, &msgType)
	if err != nil {
		println("Error unmarshalling msgType")
		return
	}

	// Unmarshal pod_desc from msg body
	pod := minik8s_pod.Pod{}
	contentData, _ := json.Marshal(parsed_msg["content"])
	err = json.Unmarshal(contentData, &pod)
	if err != nil {
		println("Error unmarshalling pod_desc")
		return
	}

	selected_node := s.SelectNode(available_nodes)
	if selected_node == nil {
		println("No available nodes")
		return
	}

	pod.Spec.NodeName = selected_node.Metadata.Name

	// Publish pod to apiserver
	request_url := fmt.Sprintf("http://%s:%s/api/v1/UpdatePod", s.APIServerIP, s.APIServerPort)
	requestBody := make(map[string]interface{})
	requestBody["namespace"] = "default"
	requestBody["pod"] = pod
	requestBodyData, _ := json.Marshal(requestBody)
	response, err := httpRequest.PostRequest(
		request_url,
		requestBodyData,
	)
	if err != nil {
		println("Error posting request: ", err)
		return
	}
	fmt.Println("\nUpdate pod response: ", response)

	// Publish pod update to mq
	msgBody := make(map[string]interface{})
	msgBody["type"] = msgType
	msgBody["content"] = pod
	msgBodyJson, err := json.Marshal(msgBody)
	if err != nil {
		println("Error marshalling msgBody")
		return
	}

	s.mqConn.Publish(
		minik8s_message.DefaultExchangeName,
		// We use node id as routing key
		selected_node.Metadata.Name,
		"application/json",
		msgBodyJson,
	)
}

func (s *Scheduler) Run() {
	done := make(chan bool)
	s.mqConn.Subscribe(
		minik8s_message.DefaultSchedulerQueueName,
		s.ScheduleHandler,
		done,
	)
}
