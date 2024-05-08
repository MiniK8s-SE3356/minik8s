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
	r.GET(url.AddPodURL, handler.AddPod)
	r.GET(url.GetPodURL, example)
	r.GET(url.GetNodesURL, example)
}

func main() {
	r := gin.Default()
	bind(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
