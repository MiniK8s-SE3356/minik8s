package test

import (
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/app"
)


func TestWhole(t *testing.T){
	kp:=app.NewKubeProxy()
	kp.Init()
}