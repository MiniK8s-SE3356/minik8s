package message

import "time"

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

type MQQueue string

type MQExchange string

const (
	DefaultExchange MQExchange = "minik8s"
)
