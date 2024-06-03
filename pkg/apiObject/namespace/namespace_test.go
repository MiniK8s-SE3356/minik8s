package namespace_test

import (
	"os"
	"testing"
)

/*
*	NOTE: This package only includes struct definition
*	and we don not need to test it
 */

func TestMain(m *testing.M) {
	// pre-test code

	// test func
	exitCode := m.Run()

	// post-test code

	os.Exit(exitCode)
}

func TestNameSpace(t *testing.T) {

}
