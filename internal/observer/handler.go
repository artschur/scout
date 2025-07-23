package observer

import (
	"fmt"
	"go-observability-tool/internal/websocket"
	"log"
	"net/http"

	gw "github.com/gorilla/websocket"
)

type MetricsHandler struct {
	wsListener *websocket.Listener
}

func NewMetricsHandler(listener *websocket.Listener) *MetricsHandler {
	return &MetricsHandler{
		wsListener: listener,
	}
}

type Publisher struct {
	conn          *gw.Conn
	publisherName string
}

// endpoint for metric publishers
func (h *MetricsHandler) NewPublisher(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgradeToWS(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newPublisher := &Publisher{
		conn:          conn,
		publisherName: r.URL.Query().Get("name"),
	}
	if newPublisher.publisherName == "" {
		http.Error(w, "Publisher name is required", http.StatusBadRequest)
		return
	}

	h.wsListener.AddConnection(conn)
}

// endpoint for a client (dashboard) that will get metrics from other sys
func (h *MetricsHandler) NewMetricListener(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgradeToWS(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for metric := range h.wsListener.MetricsChan {
		err := conn.WriteJSON(metric)
		if err != nil {
			log.Printf("error sending metric to observer")
		}
	}

}

func (h *MetricsHandler) upgradeToWS(w http.ResponseWriter, r *http.Request) (*gw.Conn, error) {
	upgrader := gw.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("Error upgrading conn to websocket")
	}
	return conn, nil
}
