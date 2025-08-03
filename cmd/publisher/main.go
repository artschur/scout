package main

import (
	"context"
	"go-observability-tool/internal/publisher"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := publisher.Config{
		HostName:       "my-machine",
		HubAddress:     "localhost:8082",
		MetricInterval: 500,
	}

	pub, err := publisher.NewPublisher(config)
	if err != nil {
		panic(err)
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		cancel()
	}()

	if err := pub.Run(ctx); err != nil {
		log.Printf("publisher stopped with error: %v", err)
	} else {
		log.Println("publisher stopped gracefully")
	}
}
