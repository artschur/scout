package publisher

import (
	"fmt"
	"go-observability-tool/internal/metrics"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
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
	memUsedMB := float64(vmStat.Used) / (1024 * 1024) // Convert bytes to MB

	temps, err := sensors.SensorsTemperatures()
	if err != nil {
		return metrics.MetricsReceived{}, fmt.Errorf("error getting CPU temperature: %v", err)
	}
	var cpuTemp float64
	if len(temps) == 0 {
		cpuTemp = 0
	} else {
		for _, t := range temps {
			if strings.Contains(t.SensorKey, "tdie") && t.Temperature > cpuTemp {
				cpuTemp = t.Temperature
			}
		}
	}

	return metrics.MetricsReceived{
		CPUUsage:         cpuPercent,
		MemoryUsageMB:    memUsedMB,
		CPUTemperature:   cpuTemp,
		MemoryPercentage: vmStat.UsedPercent,
	}, nil
}
