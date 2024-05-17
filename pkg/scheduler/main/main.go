package main

import (
	minik8s_scheduler "github.com/MiniK8s-SE3356/minik8s/pkg/scheduler"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
)

func main() {
	scheduler, err := minik8s_scheduler.NewScheduler(
		minik8s_message.DefaultMQConfig,
		minik8s_scheduler.RoundRobin,
		"127.0.0.1",
		// TODO: API server port should be configurable
		"8080",
	)
	if err != nil {
		panic(err)
	}

	scheduler.Run()
}
