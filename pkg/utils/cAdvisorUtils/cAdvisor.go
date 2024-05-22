package cadvisorutils

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
)

func GetMachineInfo() {
	request_url := fmt.Sprintf("http://%s:%s/api/v1.3/machine", "localhost", "8080")
	var machineInfo MachineInfo

	responseStatus, err := httpRequest.GetRequestByObject(
		request_url,
		nil,
		&machineInfo,
	)
	if err != nil {
		fmt.Println("Get machine info failed")
	}
	fmt.Println("Get machine info response status: ", responseStatus)

	fmt.Println(machineInfo.CPUVendorID)
	fmt.Println(machineInfo.MemoryCapacity)
	fmt.Println(machineInfo.Timestamp)

}

func GetContainerInfo() {
	request_url := fmt.Sprintf("http://%s:%s/api/v1.3/docker/some-rabbit", "localhost", "8080")
	var containerInfo map[string]ContainerInfo

	responseStatus, err := httpRequest.GetRequestByObject(
		request_url,
		nil,
		&containerInfo,
	)
	if err != nil {
		fmt.Println("Get container info failed")
	}

	fmt.Println("Get container info response status: ", responseStatus)

	for _, container := range containerInfo {
		fmt.Println(container.Name)
		for _, stats := range container.Stats {
			fmt.Println(stats.Timestamp)
			fmt.Println(stats.Cpu.Usage.Total)
			fmt.Println(stats.Memory.Usage)
		}
	}

}
