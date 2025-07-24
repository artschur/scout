package metrics

import "fmt"

type MetricsDisplay struct {
	metricsChan chan MetricsToDisplay
}

func NewMetricsDisplay(metricsChan chan MetricsToDisplay) *MetricsDisplay {
	return &MetricsDisplay{
		metricsChan: metricsChan,
	}
}

func (md *MetricsDisplay) LogMetrics() {
	for metric := range md.metricsChan {
		printMetric(metric)
	}
}

func printMetric(metric MetricsToDisplay) {
	fmt.Printf(
		"Name: %s| CPU: %.2f%% | Memory: %.2fMB\n",
		metric.Name,
		metric.MetricsReceived.CPUUsage,
		metric.MetricsReceived.MemoryUsage,
	)
}
