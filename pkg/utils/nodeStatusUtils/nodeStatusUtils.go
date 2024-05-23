package nodestatusutils

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetNodeCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0] / 100
}

func GetNodeMemPercent() float64 {
	vMem, _ := mem.VirtualMemory()
	return vMem.UsedPercent / 100
}

func GetTotalMem() uint64 {
	vMem, _ := mem.VirtualMemory()
	return vMem.Total
}
