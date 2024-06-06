package server

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/handler"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/MiniK8s-SE3356/minik8s/pkg/etcdclient"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/gin-gonic/gin"
)

// handler、process相当于spring里面的controller、service，临时先用这个名字

// func example(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "pong",
// 	})
// }

func bind(r *gin.Engine) {
	// Pod
	r.POST(url.AddPod, handler.AddPod)
	r.GET(url.GetPod, handler.GetPod)
	r.POST(url.RemovePod, handler.RemovePod)
	r.POST(url.UpdatePod, handler.UpdatePod)
	r.GET(url.GetAllPod, handler.GetAllPod)
	r.POST(url.AddServerlessFuncPod, handler.AddServerlessFuncPod)
	r.GET(url.GetServerlessFuncPod, handler.GetServerlessFuncPod)

	r.POST(url.AddNamespace, handler.AddNamespace)
	r.GET(url.GetNamespace, handler.GetNamespace)
	r.POST(url.RemoveNamespace, handler.RemoveNamespace)

	r.POST(url.AddNode, handler.AddNode)
	r.GET(url.GetNode, handler.GetNode)
	r.POST(url.RemoveNode, handler.RemoveNode)
	r.POST(url.NodeHeartBeat, handler.NodeHeartBeat)

	r.POST(url.AddReplicaset, handler.AddReplicaSet)
	r.GET(url.GetReplicaset, handler.GetReplicaSet)
	r.POST(url.RemoveReplicaset, handler.RemoveReplicaSet)

	r.POST(url.AddHPA, handler.AddHPA)
	r.GET(url.GetHPA, handler.GetHPA)
	r.POST(url.RemoveHPA, handler.RemoveHPA)
	r.POST(url.UpdateHPA, handler.UpdateHPA)

	r.POST(url.AddDNS, handler.AddDNS)
	r.GET(url.GetDNS, handler.GetDNS)
	r.POST(url.RemoveDNS, handler.RemoveDNS)
	r.POST(url.UpdateDNS, handler.UpdateDNS)
	r.GET(url.GetAllDNS, handler.GetAllDNS)

	r.POST(url.AddService, handler.AddService)
	r.GET(url.GetAllService, handler.GetAllService)
	r.GET(url.GetService, handler.GetService)
	r.GET(url.GetFilteredService, handler.GetFilteredService)
	r.POST(url.UpdateService, handler.UpdateService)
	r.POST(url.RemoveService, handler.RemoveService)

	r.GET(url.GetAllEndpoint, handler.GetAllEndpoint)
	r.POST(url.UpdateEndpointBatch, handler.UpdateEndpointBatch)
	r.POST(url.AddorDeleteEndpoint, handler.AddorDeleteEndpoint)

	r.GET(url.GetAllServerlessFunction, handler.GetAllServerlessFunction)

	r.GET(url.GetFunction, handler.GetFunction)
}

func Start() {
	var err error
	for i := 0; i < 100; i++ {
		process.EtcdCli, err = etcdclient.Connect([]string{etcdclient.EtcdURL}, etcdclient.Timeout)
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Println("failed to connect to etcd ", err.Error())
		return
	}

	process.Mq, err = message.NewMQConnection(message.DefaultMQConfig)
	if err != nil {
		fmt.Println("failed to connect to rabbitmq", err.Error())
		return
	}

	go process.CheckNodeWrapper()
	r := gin.Default()
	bind(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
