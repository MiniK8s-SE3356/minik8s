package dns_test

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

	// post-test codek

	os.Exit(exitCode)
}

func TestDns(t *testing.T) {

}
