package metrics

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func BroadcastSystemMetrics(metricsChan chan MetricsReceived) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		metrics, err := getSystemMetrics()
		if err != nil {
			return fmt.Errorf("error broadcasting metrics: %v", err)
		}
		metricsChan <- metrics
	}
	return nil
}

func getSystemMetrics() (MetricsReceived, error) {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return MetricsReceived{}, fmt.Errorf("error getting cpu percentage: %v", err)
	}
	cpuPercent := cpuPercents[0]

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return MetricsReceived{}, fmt.Errorf("error getting memory: %v", err)
	}
	memUsedMB := float64(vmStat.Used) / 1024 / 1024

	diskStat, err := disk.Usage("/")
	if err != nil {
		return MetricsReceived{}, fmt.Errorf("error getting disk: %v", err)
	}
	diskUsedGB := float64(diskStat.Used) / 1024 / 1024 / 1024

	netIOs, err := net.IOCounters(false)
	if err != nil {
		return MetricsReceived{}, fmt.Errorf("error net io: %v", err)
	}
	netMBps := float64(netIOs[0].BytesSent+netIOs[0].BytesRecv) / 1024 / 1024

	return MetricsReceived{
		CPUUsage:    cpuPercent,
		MemoryUsage: memUsedMB,
		DiskUsage:   diskUsedGB,
		NetworkIO:   netMBps,
	}, nil
}
