package server

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/handler"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/MiniK8s-SE3356/minik8s/pkg/etcdclient"
	"github.com/gin-gonic/gin"
)

// handler、process相当于spring里面的controller、service，临时先用这个名字

func example(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func bind(r *gin.Engine) {
	// Pod
	// r.POST(url.AddPod, handler.AddPod)
	// r.GET(url.GetPod, handler.GetPod)
	// r.POST(url.RemovePod, handler.RemovePod)
	// r.GET(url.DescribePod, handler.DescribePod)

	r.POST(url.AddNamespace, handler.AddNamespace)
	r.GET(url.GetNamespace, handler.GetNamespace)
	r.POST(url.RemoveNamespace, handler.RemoveNamespace)
	// r.GET(url.DescribeNamespace, handler.DescribeNamespace)

	r.POST(url.AddNode, handler.AddNode)
	r.GET(url.GetNode, handler.GetNode)
	r.POST(url.RemoveNode, handler.RemoveNode)

	// r.POST(url.AddReplicaset, handler.AddReplicaSet)
	// r.GET(url.GetReplicaset, handler.GetReplicaSet)
	// r.POST(url.RemoveReplicaset, handler.RemoveReplicaSet)
	// r.GET(url.DescribeReplicaset, handler.DescribeReplicaSet)

	// r.POST(url.AddService, handler.AddService)
	// r.GET(url.GetService, handler.GetService)
	// r.POST(url.RemoveService, handler.RemoveService)
	// r.GET(url.DescribeService, handler.DescribeService)
}

func Start() {
	var err error
	process.EtcdCli, err = etcdclient.Connect([]string{etcdclient.EtcdURL}, etcdclient.Timeout)
	if err != nil {
		fmt.Println("failed to connect to etcd ", err.Error())
		return
	}

	r := gin.Default()
	bind(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
