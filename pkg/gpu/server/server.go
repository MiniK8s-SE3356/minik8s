package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MiniK8s-SE3356/minik8s/pkg/etcdclient"
	"github.com/gin-gonic/gin"
)

const (
	version = "v1"
	// TODO: the port should be configurable
	Port    = "8083"
	RootURL = "http://localhost:" + Port
	prefix  = "/api/" + version

	SubmitGPUJob = prefix + "/SubmitGPUJob"
	GetGPUJob    = prefix + "/GetGPUJob"
	// Job pod will request this url to get the job details
	//! I think the job server in the pod should be stateless.
	//! Otherwise these job servers will consume a lot of resources to listening.
	RequireGPUJob = prefix + "/RequireGPUJob"
)

type Server struct {
	R *gin.Engine
}

func NewServer() *Server {
	return &Server{}
}

func bind(r *gin.Engine) {
	r.POST(SubmitGPUJob, SubmitGPUJobHandler)
	r.GET(GetGPUJob, GetGPUJobHandler)
	r.GET(RequireGPUJob, RequireGPUJobHandler)
}

var EtcdCli *etcdclient.EtcdClient
var GPUJobZipDir string
var ControlPanelIP string
var APIServerPort string
var JobManagerPort string

func (s *Server) Init() {

	// Connect to etcd
	var err error
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

	// We should have a directory for persistent user uploaded .zip files
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("failed to get user home directory ", err.Error())
		return
	}
	GPUJobZipDir = homeDir + "/minik8s_gpujob_zip"
	// Create the directory if it does not exist
	if _, err := os.Stat(GPUJobZipDir); os.IsNotExist(err) {
		err = os.Mkdir(GPUJobZipDir, 0777)
		if err != nil {
			fmt.Println("failed to create GPUJobZipDir ", err.Error())
			return
		}
	}

	s.R = gin.Default()
	bind(s.R)
}

func (s *Server) Run() {
	http.ListenAndServe(":"+Port, s.R)
}
