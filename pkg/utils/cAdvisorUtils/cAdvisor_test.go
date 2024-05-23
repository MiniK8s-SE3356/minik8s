package cadvisorutils_test

import (
	"testing"

	cadvisorutils "github.com/MiniK8s-SE3356/minik8s/pkg/utils/cAdvisorUtils"
)

func TestMain(m *testing.M) {
	cadvisorutils.GetMachineInfo()
	// cadvisorutils.GetContainerInfo(
	// 	"localhost",
	// 	"8080",
	// 	"some-rabbit",
	// )
	cadvisorutils.GetContainerCPUandMemory(
		"localhost",
		"8080",
		"some-rabbit",
	)
}
