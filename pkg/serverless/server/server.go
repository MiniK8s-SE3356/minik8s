package server

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/etcdclient"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/registry"
	"github.com/gin-gonic/gin"
)

const (
	version = "v1"
	// 这里暂时是localhost
	Port    = "8081"
	RootURL = "http://localhost:" + Port
	prefix  = "/api/" + version

	TriggerServerlessFunction = prefix + "TriggerServerlessFunction"
	TriggerServerlessWorkflow = prefix + "TriggerServerlessWorkflow"
	CreateFunction            = prefix + "CreateFunction"
)

type Server struct {
	R *gin.Engine
}

func bind(r *gin.Engine) {
	r.POST(TriggerServerlessFunction, triggerServerlessFunction)
	r.POST(TriggerServerlessWorkflow, triggerServerlessWorkflow)
	r.POST(CreateFunction, createFunction)
}

func (s *Server) Init() {
	var err error
	_, err = registry.RegistryInit()
	if err != nil {
		fmt.Println("failed to init registry ", err.Error())
		return
	}

	for i := 0; i < 100; i++ {
		EtcdCli, err = etcdclient.Connect([]string{etcdclient.EtcdURL}, etcdclient.Timeout)
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Println("failed to connect to etcd ", err.Error())
		return
	}

	s.R = gin.Default()
	bind(s.R)
}

func (s *Server) Run() {
	http.ListenAndServe(":8081", s.R)
}
