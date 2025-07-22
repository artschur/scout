package routes

import (
	"go-observability-tool/internal/handler"
	"go-observability-tool/internal/metrics"
	"go-observability-tool/internal/websocket"
	"net/http"
)

func AddRoutes(mux *http.ServeMux, metricsChan chan metrics.MetricsReceived) {
	webhookListener := websocket.NewListener(metricsChan)

	metricsHandler := handler.NewMetricsHandler(webhookListener)

	mux.HandleFunc("GET /listen", metricsHandler.NewMetricObserver)
	mux.HandleFunc("GET /send", metricsHandler.NewSender)
}
