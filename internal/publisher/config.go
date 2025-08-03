package publisher

type Config struct {
	HostName       string `json:"pub_name"`
	HubAddress     string `json:"hub_address"`
	MetricInterval int    `json:"metric_interval"` // in milliseconds
}
