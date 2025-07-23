package publisher

import (
	"fmt"
	"go-observability-tool/internal/metrics"

	"github.com/gorilla/websocket"
)

type SenderClient struct {
	listenerAddress string
	metrics         chan metrics.MetricsReceived
}

func NewSender(listenerAddress string) *SenderClient {
	return &SenderClient{
		listenerAddress: listenerAddress,
		metrics:         make(chan metrics.MetricsReceived),
	}
}

func (s *SenderClient) Run() {
	conn, err := s.connectToSender()
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	go metricsLoop(s.metrics)
	for metric := range s.metrics {
		if err := conn.WriteJSON(metric); err != nil {
			fmt.Println("Error sending metric:", err)
			break
		}
	}
}

func (s *SenderClient) connectToSender() (*websocket.Conn, error) {
	websocketEndpoint := fmt.Sprintf("ws://%v/send", s.listenerAddress)

	conn, _, err := websocket.DefaultDialer.Dial(websocketEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %w", err)
	}
	return conn, nil
}
