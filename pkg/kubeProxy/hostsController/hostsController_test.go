package hostsController_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// pre-test code

	// test func
	exitCode := m.Run()

	// post-test code

	os.Exit(exitCode)
}

func TestXxx(t *testing.T) {

}
