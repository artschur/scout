package metrics

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

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

func (md *MetricsDisplay) LogMetrics(ctx context.Context) {
	printHeader := true
	var prevRows int

	for {
		select {
		case <-ctx.Done():
			return
		case metric := <-md.metricsChan:
			// Update history before storing
			if existing, ok := md.metricsMap[metric.Name]; ok {
				// Append new values and limit history length
				metric.CPUHistory = append(existing.CPUHistory, metric.MetricsReceived.CPUUsage)
				if len(metric.CPUHistory) > HistoryLength {
					metric.CPUHistory = metric.CPUHistory[len(metric.CPUHistory)-HistoryLength:]
				}

				metric.MemoryHistory = append(existing.MemoryHistory, metric.MetricsReceived.MemoryPercentage)
				if len(metric.MemoryHistory) > HistoryLength {
					metric.MemoryHistory = metric.MemoryHistory[len(metric.MemoryHistory)-HistoryLength:]
				}
			} else {
				// Initialize history for new metrics
				metric.CPUHistory = []float64{metric.MetricsReceived.CPUUsage}
				metric.MemoryHistory = []float64{metric.MetricsReceived.MemoryPercentage}
			}

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

			// Sort host names
			var hostNames []string
			for name := range md.metricsMap {
				hostNames = append(hostNames, name)
			}
			sort.Strings(hostNames)

			// Print all metrics in order and count rows
			rows := 0
			for _, name := range hostNames {
				// Each metric now takes 3 rows (stats + 2 graph lines)
				rows += printMetric(md.metricsMap[name])
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
}

func printMetric(metric MetricsToDisplay) int {
	cpuTemp := "N/A"
	if metric.MetricsReceived.CPUTemperature > 0 {
		cpuTemp = fmt.Sprintf("%.2f°C", metric.MetricsReceived.CPUTemperature)
	}

	// Print the main metrics line
	fmt.Printf(
		"%-12s | %6.2f%% | %8s | %10.2f | %9.2f%%\n",
		metric.Name,
		metric.MetricsReceived.CPUUsage,
		cpuTemp,
		metric.MetricsReceived.MemoryUsageMB,
		metric.MetricsReceived.MemoryPercentage,
	)

	// Print colored CPU graph
	fmt.Printf("  CPU:       %s\n", renderColoredSparkline(metric.CPUHistory, 100))

	// Print colored Memory graph
	fmt.Printf("  Memory:    %s\n", renderColoredSparkline(metric.MemoryHistory, 100))

	return 3 // Return number of lines printed
}

// renderColoredSparkline generates a colored text-based graph using Unicode block characters
func renderColoredSparkline(values []float64, maxValue float64) string {
	if len(values) == 0 {
		return ""
	}

	// Unicode block elements from lowest to highest
	blocks := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	var result strings.Builder
	for _, v := range values {
		// Scale value to 0-7 range for block selection
		normalizedValue := v / maxValue * float64(len(blocks)-1)
		if normalizedValue < 0 {
			normalizedValue = 0
		} else if normalizedValue >= float64(len(blocks)) {
			normalizedValue = float64(len(blocks) - 1)
		}

		blockIndex := int(normalizedValue)

		// Choose color based on value
		var colorCode string
		percentage := v / maxValue * 100
		switch {
		case percentage < 30:
			colorCode = colorGreen // Green for low values
		case percentage < 80:
			colorCode = colorYellow // Yellow for medium values
		default:
			colorCode = colorRed // Red for high values
		}

		result.WriteString(colorCode)
		result.WriteRune(blocks[blockIndex])
		result.WriteString(colorReset) // Reset color after each block
	}

	return result.String()
}
