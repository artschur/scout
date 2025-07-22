package metrics

import "fmt"

type MetricsDisplay struct {
	metricsChan chan MetricsReceived
}

func NewMetricsDisplay(metricsChan chan MetricsReceived) *MetricsDisplay {
	return &MetricsDisplay{
		metricsChan: metricsChan,
	}
}

func (md *MetricsDisplay) LogMetrics() {
	for metric := range md.metricsChan {
		printMetric(metric)
	}
}

func printMetric(metric MetricsReceived) {
	fmt.Printf(
		"CPU: %.2f%% | Memory: %.2fMB | Disk: %.2fGB | Network: %.2fMB/s\n",
		metric.CPUUsage,
		metric.MemoryUsage,
		metric.DiskUsage,
		metric.NetworkIO,
	)
}
