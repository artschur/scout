package publisher

import (
	"context"
	"fmt"
	"go-observability-tool/internal/metrics"

	"github.com/gorilla/websocket"
)

type PublisherClient struct {
	Config

	metrics chan metrics.MetricsReceived
}

func NewPublisher(Config Config) (*PublisherClient, error) {
	if Config.HostName == "" {
		return nil, fmt.Errorf("publisher name must be provided")
	}
	if Config.HubAddress == "" {
		return nil, fmt.Errorf("hub address must be provided")
	}
	return &PublisherClient{
		Config:  Config,
		metrics: make(chan metrics.MetricsReceived),
	}, nil
}

func (s *PublisherClient) Run(ctx context.Context) {
	conn, err := s.connectToSender()
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	
	defer conn.Close()
	go metricsLoop(ctx, s.metrics)
	for metric := range s.metrics {
		fmt.Println("Sending metric:", metric)
		if err := conn.WriteJSON(metric); err != nil {
			fmt.Println("Error sending metric:", err)
			break
		}
	}
}

func (s *PublisherClient) connectToSender() (*websocket.Conn, error) {
	websocketEndpoint := fmt.Sprintf("ws://%v/send?name=%s", s.Config.HubAddress, s.Config.HostName)

	conn, _, err := websocket.DefaultDialer.Dial(websocketEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %w", err)
	}
	return conn, nil
}
