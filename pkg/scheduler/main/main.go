package main

import (
	"flag"

	minik8s_scheduler "github.com/MiniK8s-SE3356/minik8s/pkg/scheduler"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
)

func main() {
	apiServerIP := flag.String("apiserverip", "127.0.0.1", "APIServer IP address")
	apiServerPort := flag.String("apiserverport", "8080", "APIServer port")

	scheduler, err := minik8s_scheduler.NewScheduler(
		minik8s_message.DefaultMQConfig,
		minik8s_scheduler.RoundRobin,
		*apiServerIP,
		*apiServerPort,
	)
	if err != nil {
		panic(err)
	}

	scheduler.Run()
}
