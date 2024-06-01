package app

import (
	"flag"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/gpu/server"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/nettools"
)

type JobManager struct {
	s (*server.Server)
}

func NewJobManager() *JobManager {
	fmt.Printf("New Job Manager...\n")
	return &JobManager{
		s: server.NewServer(),
	}
}

func (jm *JobManager) Init() {
	fmt.Printf("Init Job Manager...\n")
	jm.s.Init()
}

func (jm *JobManager) Run() {
	fmt.Printf("Run Job Manager...\n")
	jm.s.Run()
}

func StartJobManager() {
	defaultControlPanelIP := nettools.KubeletDefaultIP()

	controlPanelIP := flag.String("controlpanelip", defaultControlPanelIP, "Control Panel IP address")
	apiServerPort := flag.String("apiserverport", "8080", "APIServer port")
	jobManagerPort := flag.String("jobmanagerport", "8083", "Job Manager port")

	server.ControlPanelIP = *controlPanelIP
	server.APIServerPort = *apiServerPort
	server.JobManagerPort = *jobManagerPort

	jm := NewJobManager()
	jm.Init()
	jm.Run()
}
