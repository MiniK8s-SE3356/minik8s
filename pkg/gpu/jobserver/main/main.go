package main

import (
	"fmt"
	"os"

	"github.com/MiniK8s-SE3356/minik8s/pkg/gpu/jobserver/server"
	jobmanager_server "github.com/MiniK8s-SE3356/minik8s/pkg/gpu/server"
)

func main() {
	forever := make(chan bool)

	// Get JobManagerIP, JobManagerPort, JobName from env
	// TODO: We should remove package jobmanager_server because it will bring a lot of dependencies when building the binary
	JobManagerIP := os.Getenv(jobmanager_server.JobManagerIPEnv)
	JobManagerPort := os.Getenv(jobmanager_server.JobManagerPortEnv)
	JobName := os.Getenv(jobmanager_server.JobNameEnv)
	fmt.Println("JobManagerIP: ", JobManagerIP)
	fmt.Println("JobManagerPort: ", JobManagerPort)
	fmt.Println("JobName: ", JobName)

	server.JobManagerUrl = fmt.Sprintf("http://%s:%s", JobManagerIP, JobManagerPort)

	jobServer, err := server.NewJobServer(JobName)
	if err != nil {
		fmt.Println("failed to create new job server")
		// return
		<-forever
	}

	jobServer.Run()

	<-forever
}
