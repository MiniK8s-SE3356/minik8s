package controller_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// pre-test code
	// controller.StartController()
	// test func
	exitCode := m.Run()

	// post-test code
	os.Exit(exitCode)
}

func TestController(t *testing.T) {

}
