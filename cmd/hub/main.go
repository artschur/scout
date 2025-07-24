package main

import (
	"context"
	observer "go-observability-tool/internal/hub"
	"go-observability-tool/internal/metrics"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.DefaultServeMux

	metricsChan := make(chan metrics.MetricsToDisplay)
	observer.AddRoutes(mux, metricsChan)

	metricsDisplayer := metrics.NewMetricsDisplay(metricsChan)

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	go metricsDisplayer.LogMetrics()

	go func() {
		log.Println("Server starting on :8082")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server gracefully stopped")
}
