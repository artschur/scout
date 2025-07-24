package websocket

import (
	"go-observability-tool/internal/metrics"
)

type Hub struct {
	connections []*Connection

	MetricsChan chan metrics.MetricsToDisplay
}

func NewHub(metricsChan chan metrics.MetricsToDisplay) *Hub {
	return &Hub{
		connections: make([]*Connection, 0),
		MetricsChan: metricsChan,
	}
}

func (l *Hub) AddConnection(conn *Connection) {
	l.connections = append(l.connections, conn)
	go func(c *Connection) {
		for {
			var metricsReceived metrics.MetricsReceived
			err := c.Conn.ReadJSON(&metricsReceived)
			if err != nil {
				break
			}

			metricsToDisplay := metrics.MetricsToDisplay{
				MetricsReceived: metricsReceived,
				Name:            c.Name,
				IP:              c.IP,
			}

			l.MetricsChan <- metricsToDisplay
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
