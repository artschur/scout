package hub

import (
	"net/http"
)

func AddRoutes(mux *http.ServeMux, metricsHandler *MetricsHandler) {
	mux.HandleFunc("GET /client", metricsHandler.NewMetricSubscriber)
	mux.HandleFunc("GET /send", metricsHandler.NewPublisher)
}
