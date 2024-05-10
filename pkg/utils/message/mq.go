package message

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type MQConfig struct {
	// RabbitMQ connection url config
	// example: amqp://user:password@localhost:5672/vhost
	User     string
	Password string
	Host     string
	Port     string
	Vhost    string

	// Some configurations for the reconnect
	MaxRetry   int
	RetryDelay time.Duration
}

var DefaultMQConfig = &MQConfig{
	User:       "guest",
	Password:   "guest",
	Host:       "localhost",
	Port:       "5672",
	Vhost:      "/",
	MaxRetry:   5,
	RetryDelay: 5 * time.Second,
}

type MQConnection struct {
	Conn   *amqp.Connection
	Config *MQConfig
}

func NewMQConnection(config *MQConfig) (*MQConnection, error) {
	conn, err := amqp.Dial("amqp://" + config.User + ":" + config.Password + "@" + config.Host + ":" + config.Port + config.Vhost)
	fmt.Println(config.User + ":" + config.Password + "@" + config.Host + ":" + config.Port + config.Vhost)

	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ")
		return nil, err
	}

	mqConn := &MQConnection{
		Conn:   conn,
		Config: config,
	}

	// Create a goroutine to retry connecting when connection is closed abnormally
	go mqConn.handleReconnect()

	return mqConn, nil
}

func connectWithRetry(config *MQConfig) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < config.MaxRetry; i++ {
		conn, err = amqp.Dial("amqp://" + config.User + ":" + config.Password + "@" + config.Host + ":" + config.Port + config.Vhost)
		if err == nil {
			return conn, nil
		}
		fmt.Println("Failed to connect to RabbitMQ, error message: ", err, "retrying...")
		time.Sleep(config.RetryDelay)
	}

	return nil, err
}

func (mq *MQConnection) handleReconnect() {
	// Use 'make' to create a channel to receive type '*amqp.Error' data
	// Due to this channel has no buffer the operation sending message to this channel
	// will be blocked until a receive goroutine is ready to receive.
	// 'NotifyClose' will register a channel. When AMQP is closed abnormally,
	// the error message will be sent to this channel. You can listen to
	// 'notify' to get error message
	notify := mq.Conn.NotifyClose(
		make(chan *amqp.Error),
	)

	for err := range notify {
		fmt.Println("RabbitMQ connection closed, error message: ", err)

		// Try to reconnect
		for {
			time.Sleep(mq.Config.RetryDelay)

			conn, err := connectWithRetry(mq.Config)
			if err == nil {
				mq.Conn = conn
				fmt.Println("RabbitMQ reconnected successfullt")
				// A new goroutine to handle next reconnection
				go mq.handleReconnect()
				return
			}
		}
	}
}

func (mq *MQConnection) Publish(exchange string, routingKey string, contentType string, body []byte) error {
	ch, err := mq.Conn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel, error message: ", err)
		return err
	}
	defer ch.Close()

	// TODO: Check queue exists or create and bind queue

	err = ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		},
	)

	if err != nil {
		fmt.Println("Failed to publish a message, error message: ", err)
		return err
	}

	return nil
}

// Subscribe starts listening on a queue and calls the callback function when a message is received
func (mq *MQConnection) Subscribe(queue string, callback func(amqp.Delivery), done <-chan bool) error {
	ch, err := mq.Conn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel, error message: ", err)
		return err
	}
	defer ch.Close()

	msgChannel, err := ch.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("Failed to consume message, error message: ", err)
		return err
	}

	// Start a goroutine to handle messages
	go func() {
		// for {
		// 	select {
		// 	case msg := <-msgChannel:
		// 		callback(msg)
		// 	case <-done:
		// 		fmt.Println("Unsubscribing from queue: ", queue)
		// 		return
		// 	}
		// }
		for msg := range msgChannel {
			callback(msg)
		}
	}()
	fmt.Println("Subscribed to queue: ", queue)

	return nil
}

func (mq *MQConnection) GetChannel() (*amqp.Channel, error) {
	return mq.Conn.Channel()
}

func (mq *MQConnection) Close() {
	mq.Conn.Close()
}
