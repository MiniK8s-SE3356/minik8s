package cadvisorutils

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	nodestatusutils "github.com/MiniK8s-SE3356/minik8s/pkg/utils/nodeStatusUtils"
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

func GetContainerInfo(cAdvisorIP string, cAdvisorPort string, containerName string) {
	request_url := fmt.Sprintf("http://%s:%s/api/v1.3/docker/%s", cAdvisorIP, cAdvisorPort, containerName)
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

	// TODO: Calculate the CPU usage

	for _, container := range containerInfo {
		fmt.Println(container.Spec.CreationTime)
		for _, stat := range container.Stats {
			fmt.Println(stat.Timestamp)
			fmt.Println(stat.Cpu.Usage.Total)
			fmt.Println(stat.Memory.Usage)
		}
	}
}

func GetContainerCPUandMemory(cAdvisorIP string, cAdvisorPort string, containerName string) (float64, float64) {
	request_url := fmt.Sprintf("http://%s:%s/api/v1.3/docker/%s", cAdvisorIP, cAdvisorPort, containerName)
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

	// calculate cpu and memory usage

	for _, container := range containerInfo {
		var totalCPUUsage float64
		var cpuCount float64
		var averageCPUUsage float64

		totalCPUUsage = 0.0
		cpuCount = 0.0
		averageCPUUsage = 0.0

		for i := 1; i < len(container.Stats); i++ {
			t1 := container.Stats[i-1].Timestamp
			t2 := container.Stats[i].Timestamp

			duration := t2.Sub(t1).Seconds()
			cpuDelta := float64(container.Stats[i].Cpu.Usage.Total-container.Stats[i-1].Cpu.Usage.Total) / 1e9

			if duration > 0 {
				cpuUsage := (cpuDelta / duration) // * 100
				totalCPUUsage += cpuUsage
				cpuCount++
			}
		}
		if cpuCount > 0 {
			averageCPUUsage = totalCPUUsage / cpuCount
		} else {
			fmt.Println("No valid data to calculate average CPU usage.")
		}

		var totalMemoryUsage uint64
		var averageMemoryUsage uint64
		for _, stat := range container.Stats {
			totalMemoryUsage += stat.Memory.Usage
		}
		averageMemoryUsage = totalMemoryUsage / uint64(len(container.Stats))

		totalMemCapacity := nodestatusutils.GetTotalMem()

		memUsePercentage := float64(averageMemoryUsage) / float64(totalMemCapacity)

		fmt.Println("Average CPU Usage: ", averageCPUUsage)
		fmt.Println("Memory Usage Percentage: ", memUsePercentage)

		return averageCPUUsage, memUsePercentage
	}
	return 0.0, 0.0
}
