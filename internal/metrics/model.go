package metrics

type MetricsReceived struct {
	CPUUsage         float64 `json:"cpu_usage"`
	MemoryUsageMB    float64 `json:"mem_usage"`
	CPUTemperature   float64 `json:"cpu_temperature,omitempty"`
	MemoryPercentage float64 `json:"mem_percentage"`
}

type MetricsToDisplay struct {
	Name            string
	MetricsReceived MetricsReceived
	CPUHistory      []float64 // Added: Store historical CPU usage
	MemoryHistory   []float64 // Added: Store historical memory usage
	Ip              string
}

const (
	HistoryLength = 20 // Number of data points to store for graphs
)
