package hub

import (
	"context"
	"go-observability-tool/internal/metrics"
	"go-observability-tool/internal/websocket"
)

type Hub struct {
	conns          map[*websocket.Connection]bool
	MetricsChan    chan metrics.MetricsToDisplay
	registerChan   chan *websocket.Connection
	unregisterChan chan *websocket.Connection
}

func NewHub() *Hub {
	return &Hub{
		conns:          make(map[*websocket.Connection]bool),
		registerChan:   make(chan *websocket.Connection),
		unregisterChan: make(chan *websocket.Connection),
		MetricsChan:    make(chan metrics.MetricsToDisplay),
	}
}

func (l *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for conn := range l.conns {
				conn.Conn.Close()
			}
			return
		case metric := <-l.MetricsChan:
			for conn := range l.conns {
				if conn.Role == "subscriber" {
					err := conn.Conn.WriteJSON(metric)
					if err != nil {
						conn.Conn.Close()
						delete(l.conns, conn)
					}
				}
			}

		case conn := <-l.registerChan:
			l.addConnection(conn)
		case conn := <-l.unregisterChan:
			l.removeConnection(conn)
		}
	}
}

func (l *Hub) addConnection(conn *websocket.Connection) {
	l.conns[conn] = true
}

func (l *Hub) removeConnection(conn *websocket.Connection) {
	delete(l.conns, conn)
	conn.Conn.Close()
	if conn.Role == "publisher" {
		conn.Conn.Close()
	}
}
