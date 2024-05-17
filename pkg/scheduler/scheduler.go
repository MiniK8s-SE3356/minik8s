package scheduler

import (
	"encoding/json"
	"math/rand"
	"sync"
	"time"

	minik8s_apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"
	minik8s_pod "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	minik8s_yaml "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	minik8s_container "github.com/MiniK8s-SE3356/minik8s/pkg/types/container"
	minik8s_node "github.com/MiniK8s-SE3356/minik8s/pkg/types/node"
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
	nodes := []*minik8s_node.Node{
		{
			Metadata: minik8s_node.NodeMetadata{
				Id:   "1",
				Name: "node1",
				Labels: map[string]string{
					"zone": "zone1",
				},
			},
			Status: minik8s_node.NodeStatus{
				Hostname:   "node1",
				Ip:         "127.0.0.1",
				Condition:  []string{minik8s_node.NODE_Ready},
				CpuPercent: 0.1,
				MemPercent: 0.2,
				NumPods:    1,
				UpdateTime: time.Now().String(),
			},
		},
	}

	return nodes, nil
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
	var parsed_msg minik8s_message.Message
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

	// Unmarshal pod_desc from msg body
	pod_desc := minik8s_yaml.PodDesc{}
	err = json.Unmarshal([]byte(parsed_msg.Body), &pod_desc)
	if err != nil {
		println("Error unmarshalling pod_desc")
		return
	}

	selected_node := s.SelectNode(available_nodes)
	if selected_node == nil {
		println("No available nodes")
		return
	}

	pod := minik8s_pod.Pod{
		Basic: minik8s_apiobject.Basic{
			APIVersion: pod_desc.ApiVersion,
			Kind:       pod_desc.Kind,
			Metadata: minik8s_apiobject.Metadata{
				// TODO: Generate UUID
				UUID: "uuid",
				Name: pod_desc.Metadata.Name,
				// TODO: Add namespace
				Namespace: "default",
				Labels:    pod_desc.Metadata.Labels,
				// TODO: Add annotations
				Annotations: map[string]string{},
			},
		},
		Spec: minik8s_pod.PodSpec{
			NodeName:   selected_node.Metadata.Name,
			Containers: []minik8s_container.Container{},
			// TODO: Add volumes
			Volumes: []minik8s_pod.Volume{},
			// TODO: Add nodeSelector
			NodeSelector: map[string]string{},
		},
	}

	// Add containers to pod
	for _, container_desc := range pod_desc.Spec {
		pod.Spec.Containers = append(
			pod.Spec.Containers,
			minik8s_container.Container{
				Name:  container_desc.Name,
				Image: container_desc.Image,
				// TODO: combine some fields
			},
		)
	}

	// Publish pod to apiserver

	// Publish pod update to mq
	pod_json, err := json.Marshal(pod)
	if err != nil {
		println("Error marshalling pod")
		return
	}

	s.mqConn.Publish(
		"minik8s",
		// TODO: Node's ID or Name?
		selected_node.Metadata.Id,
		"application/json",
		pod_json,
	)
}

func (s *Scheduler) Run() {
	done := make(chan bool)
	s.mqConn.Subscribe("scheduler", s.ScheduleHandler, done)
}
