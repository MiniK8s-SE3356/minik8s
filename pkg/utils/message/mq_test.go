package message_test

import (
	"encoding/json"
	"testing"
	"time"

	minik8s_mq "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/streadway/amqp"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Callback(delivery amqp.Delivery) {
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

	// test publish again
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

	// err = test_mq.Subscribe(
	// 	"kubelet",
	// 	Callback,
	// 	done,
	// )
	go func() {
		err = test_mq.Subscribe(
			"kubelet",
			Callback,
			done,
		)
		if err != nil {
			println("Error subscribing")
			panic(err)
		}
	}()

	time.Sleep(5 * time.Second)
	done <- true
	time.Sleep(5 * time.Second)

	//.................................................................................................//

	// ch, err := test_mq.Conn.Channel()
	// if err != nil {
	// 	fmt.Println("Failed to open a channel, error message: ", err)
	// }
	// defer ch.Close()

	// q, err := ch.QueueDeclare(
	// 	"minik8s_test",
	// 	false,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	fmt.Println("Failed to declare a queue, error message: ", err)
	// }

	// msgs, err := ch.Consume(
	// 	q.Name, // queue
	// 	"",     // consumer
	// 	true,   // auto-ack
	// 	false,  // exclusive
	// 	false,  // no-local
	// 	false,  // no-wait
	// 	nil,    // args
	// )

	// if err != nil {
	// 	log.Fatalf("Failed to register a consumer: %v", err)
	// }

	// forever := make(chan bool)

	// go func() {
	// 	for d := range msgs {
	// 		log.Printf("Received a message: %s", d.Body)
	// 	}
	// }()

	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// <-forever

}
