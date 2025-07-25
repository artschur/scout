package metrics

import "fmt"

type MetricsDisplay struct {
	metricsChan chan MetricsToDisplay
	metricsMap  map[string]MetricsToDisplay
}

func NewMetricsDisplay(metricsChan chan MetricsToDisplay) *MetricsDisplay {
	return &MetricsDisplay{
		metricsChan: metricsChan,
		metricsMap:  make(map[string]MetricsToDisplay),
	}
}

func (md *MetricsDisplay) LogMetrics() {
	printHeader := true
	var prevRows int

	for metric := range md.metricsChan {
		md.metricsMap[metric.Name] = metric

		// Print header once
		if printHeader {
			fmt.Println("Name         | CPU    | CPU Temp | Memory (MB) | Memory (%)")
			fmt.Println("-------------------------------------------------------------")
			printHeader = false
		} else if prevRows > 0 {
			// Move cursor up to start of data rows
			fmt.Printf("\033[%dA", prevRows)
		}

		// Print all metrics and count rows
		rows := 0
		for _, metric := range md.metricsMap {
			printMetric(metric)
			rows++
		}

		// If number of hosts decreased, clear extra lines
		for i := rows; i < prevRows; i++ {
			fmt.Print("\033[2K\n") // Clear line and move down
		}
		if rows < prevRows {
			// Move cursor up to stay at the end of the table
			fmt.Printf("\033[%dA", prevRows-rows)
		}

		prevRows = rows
	}
}

func (md *MetricsDisplay) printAllMetrics() int {
	fmt.Println("Name         | CPU    | CPU Temp | Memory (MB) | Memory (%)")
	fmt.Println("-------------------------------------------------------------")
	lines := 2
	for _, metric := range md.metricsMap {
		printMetric(metric)
		lines++
	}
	return lines
}

func printMetric(metric MetricsToDisplay) {
	cpuTemp := "N/A"
	if metric.MetricsReceived.CPUTemperature > 0 {
		cpuTemp = fmt.Sprintf("%.2f°C", metric.MetricsReceived.CPUTemperature)
	}
	fmt.Printf(
		"%-12s | %6.2f%% | %8s | %10.2f | %9.2f%%\n",
		metric.Name,
		metric.MetricsReceived.CPUUsage,
		cpuTemp,
		metric.MetricsReceived.MemoryUsageMB,
		metric.MetricsReceived.MemoryPercentage,
	)
}
