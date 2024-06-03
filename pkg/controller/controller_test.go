package controller_test

import (
	"os"
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/controller"
)

func TestMain(m *testing.M) {
	// pre-test code
	controller.StartController()
	// test func
	exitCode := m.Run()

	// post-test code
	os.Exit(exitCode)
}

func TestController(t *testing.T) {
	
}
