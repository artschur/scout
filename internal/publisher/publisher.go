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

func (s *PublisherClient) Run(ctx context.Context) error {
	conn, err := s.connectToSender()
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer conn.Close()

	go metricsLoop(ctx, s.metrics)

	for {
		select {
		case <-ctx.Done():
			return nil
		case metric, ok := <-s.metrics:
			if !ok {
				return nil
			}
			if err := conn.WriteJSON(metric); err != nil {
				return fmt.Errorf("error sending metric: %w", err)
			}
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
