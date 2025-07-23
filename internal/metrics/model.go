package metrics

type MetricsReceived struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"mem_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	NetworkIO   float64 `json:"network_io"`
}
