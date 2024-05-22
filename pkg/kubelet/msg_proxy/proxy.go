package msgproxy

import (
	"encoding/json"

	minik8s_pod "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	minik8s_worker "github.com/MiniK8s-SE3356/minik8s/pkg/kubelet/worker"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/streadway/amqp"
)

type MsgProxy struct {
	mqConn           *minik8s_message.MQConnection
	listenQueueName  string
	PodUpdateChannel chan minik8s_worker.Task
}

func NewMsgProxy(mqConfig *minik8s_message.MQConfig, queueName string) *MsgProxy {
	newConn, _ := minik8s_message.NewMQConnection(mqConfig)
	return &MsgProxy{
		mqConn:           newConn,
		listenQueueName:  queueName,
		PodUpdateChannel: make(chan minik8s_worker.Task),
	}
}

func (mp *MsgProxy) handleReceive(delivery amqp.Delivery) {
	var parsed_msg map[string]interface{}
	err := json.Unmarshal(delivery.Body, &parsed_msg)
	if err != nil {
		return
	}

	typeData, _ := json.Marshal(parsed_msg["type"])
	var msgType minik8s_message.MsgType
	err = json.Unmarshal(typeData, &msgType)
	if err != nil {
		return
	}

	//!debug//
	// fmt.Println("msgType: ", msgType)
	//!debug//

	contentData, _ := json.Marshal(parsed_msg["content"])
	var parsed_pod minik8s_pod.Pod
	err = json.Unmarshal(contentData, &parsed_pod)
	if err != nil {
		return
	}

	//!debug//
	// podJson, _ := json.Marshal(parsed_pod)
	// fmt.Println("pod: ", string(podJson))
	// fmt.Println("container numbs: ", len(parsed_pod.Spec.Containers))
	//!debug//

	switch msgType {
	case minik8s_message.CreatePod:
		mp.PodUpdateChannel <- minik8s_worker.Task{
			Type: minik8s_worker.Task_Add,
			Pod:  &parsed_pod,
		}
	case minik8s_message.RemovePod:
		mp.PodUpdateChannel <- minik8s_worker.Task{
			Type: minik8s_worker.Task_Remove,
			Pod:  &parsed_pod,
		}
	// TODO: more actions
	default:
		return
	}
}

func (mp *MsgProxy) Run() {
	done := make(chan bool)

	err := mp.mqConn.Subscribe(
		mp.listenQueueName,
		mp.handleReceive,
		done,
	)

	if err != nil {
		return
	}
}
