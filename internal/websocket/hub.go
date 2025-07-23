package websocket

import (
	"go-observability-tool/internal/metrics"
)

type Hub struct {
	connections []*Connection

	MetricsChan chan metrics.MetricsReceived
}

func NewHub(metricsChan chan metrics.MetricsReceived) *Hub {
	return &Hub{
		connections: make([]*Connection, 0),
		MetricsChan: metricsChan,
	}
}

func (l *Hub) AddConnection(conn *Connection) {
	l.connections = append(l.connections, conn)
	go func(c *Connection) {
		for {
			var metrics metrics.MetricsReceived
			err := c.Conn.ReadJSON(&metrics)
			if err != nil {
				break
			}
			l.MetricsChan <- metrics
		}
	}(conn)
}

func (l *Hub) RemoveConnection(conn *Connection) {
	for i, c := range l.connections {
		if c == conn {
			l.connections = append(l.connections[:i], l.connections[i+1:]...)
			break
		}
	}
}
