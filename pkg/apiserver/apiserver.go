package main

import (
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/handler"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
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
	r.POST(url.AddPodURL, handler.AddPod)
	r.GET(url.GetPodURL, handler.GetPod)
	r.POST(url.RemovePodURL, handler.RemovePod)
	r.GET(url.DescribePodURL, handler.DescribePod)

	r.POST(url.AddNamespaceURL, handler.AddNamespace)
	r.GET(url.GetNamespaceURL, handler.GetNamespace)
	r.POST(url.RemoveNamespaceURL, handler.RemoveNamespace)
	r.GET(url.DescribeNamespaceURL, handler.DescribeNamespace)

	r.GET(url.GetNodesURL, example)

	r.POST(url.AddNamespaceURL, handler.AddNamespace)
	r.GET(url.GetNamespaceURL, handler.GetNamespace)
	r.POST(url.RemoveNamespaceURL, handler.RemoveNamespace)
	r.GET(url.DescribeNamespaceURL, handler.DescribeNamespace)

	r.POST(url.AddServiceURL, handler.AddService)
	r.GET(url.GetServiceURL, handler.GetService)
	r.POST(url.RemoveServiceURL, handler.RemoveService)
	r.GET(url.DescribeServiceURL, handler.DescribeService)
}

func main() {
	r := gin.Default()
	bind(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
