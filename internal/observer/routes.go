package observer

import (
	"go-observability-tool/internal/metrics"
	"go-observability-tool/internal/websocket"
	"net/http"
)

func AddRoutes(mux *http.ServeMux, metricsChan chan metrics.MetricsReceived) {
	webhookListener := websocket.NewListener(metricsChan)

	metricsHandler := NewMetricsHandler(webhookListener)

	mux.HandleFunc("GET /client", metricsHandler.NewMetricListener)
	mux.HandleFunc("GET /send", metricsHandler.NewPublisher)
}
