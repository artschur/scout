package websocket

import (
	"go-observability-tool/internal/metrics"

	"github.com/gorilla/websocket"
)

type Listener struct {
	connections []*websocket.Conn
	MetricsChan chan metrics.MetricsReceived
}

func NewListener(metricsChan chan metrics.MetricsReceived) *Listener {
	return &Listener{
		connections: make([]*websocket.Conn, 0),
		MetricsChan: metricsChan,
	}
}

func (l *Listener) AddConnection(conn *websocket.Conn) {
	l.connections = append(l.connections, conn)
}

func (l *Listener) RemoveConnection(conn *websocket.Conn) {
	for i, c := range l.connections {
		if c == conn {
			l.connections = append(l.connections[:i], l.connections[i+1:]...)
			break
		}
	}
}

func (l *Listener) GetMetricsFromConnection(conn *websocket.Conn) metrics.MetricsReceived {
	var metrics metrics.MetricsReceived
	err := conn.ReadJSON(&metrics)
	if err != nil {
		return metrics
	}

	return metrics
}

func (l *Listener) Listen() {
	for _, conn := range l.connections {
		go func(c *websocket.Conn) {
			for {
				var metrics metrics.MetricsReceived
				err := c.ReadJSON(&metrics)
				if err != nil {
					break
				}
				l.MetricsChan <- metrics
			}
		}(conn)
	}
}
