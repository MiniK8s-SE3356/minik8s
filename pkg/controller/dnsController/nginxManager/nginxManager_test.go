package nginxManager_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M){
	// pre-test code
	
	// test func
	exitCode := m.Run()

	// post-test code

	// 返回测试运行的退出码
	os.Exit(exitCode)
}

func TestXxx(t *testing.T) {
	
}