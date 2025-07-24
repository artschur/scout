package publisher

import (
	"fmt"
	"go-observability-tool/internal/metrics"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

func metricsLoop(metricsChan chan metrics.MetricsReceived) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		metrics, err := getMetrics()
		if err != nil {
			fmt.Printf("Error getting metrics: %v\n", err)
			continue
		}
		metricsChan <- metrics
	}
}

func getMetrics() (metrics.MetricsReceived, error) {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return metrics.MetricsReceived{}, fmt.Errorf("error getting cpu percentage: %v", err)
	}
	cpuPercent := cpuPercents[0]

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return metrics.MetricsReceived{}, fmt.Errorf("error getting memory: %v", err)
	}
	memUsedMB := float64(vmStat.Used) / 1024 / 1024

	diskStat, err := disk.Usage("/")
	if err != nil {
		return metrics.MetricsReceived{}, fmt.Errorf("error getting disk: %v", err)
	}
	diskUsedGB := float64(diskStat.Used) / 1024 / 1024 / 1024

	netIOs, err := net.IOCounters(false)
	if err != nil {
		return metrics.MetricsReceived{}, fmt.Errorf("error net io: %v", err)
	}
	netMBps := float64(netIOs[0].BytesSent+netIOs[0].BytesRecv) / 1024 / 1024

	return metrics.MetricsReceived{
		CPUUsage:    cpuPercent,
		MemoryUsage: memUsedMB,
		DiskUsage:   diskUsedGB,
		NetworkIO:   netMBps,
	}, nil
}
