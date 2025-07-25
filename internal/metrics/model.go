package metrics

type MetricsReceived struct {
	CPUUsage         float64 `json:"cpu_usage"`
	MemoryUsageMB    float64 `json:"mem_usage"`
	CPUTemperature   float64 `json:"cpu_temperature,omitempty"`
	MemoryPercentage float64 `json:"mem_percentage"`
}

type MetricsToDisplay struct {
	MetricsReceived MetricsReceived
	Name            string `json:"name"`
	IP              string `json:"ip"`
}
