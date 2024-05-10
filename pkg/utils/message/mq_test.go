package message_test

import (
	"encoding/json"
	"testing"

	minik8s_mq "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/streadway/amqp"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Callback(delivery amqp.Delivery) {
	println("Received message")
	var person Person
	err := json.Unmarshal(delivery.Body, &person)
	if err != nil {
		println("Error unmarshalling json")
		panic(err)
	}
	println(person.Name)
}

func TestMain(m *testing.M) {
	test_mq, err := minik8s_mq.NewMQConnection(
		minik8s_mq.DefaultMQConfig,
	)

	if err != nil {
		println("Error creating mq connection")
		panic(err)
	}

	person := Person{
		Name: "Alice",
		Age:  20,
	}
	body, err := json.Marshal(person)

	if err != nil {
		println("Error marshalling json")
		panic(err)
	}

	// test publish
	err = test_mq.Publish(
		"minik8s_test",
		"kubelet",
		"application/json",
		body,
	)

	if err != nil {
		println("Error publishing message")
		panic(err)
	}

	done := make(chan bool)

	err = test_mq.Subscribe(
		"kubelet",
		Callback,
		done,
	)

	if err != nil {
		println("Error subscribing")
		panic(err)
	}
}
