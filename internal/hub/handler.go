package hub

import (
	"fmt"
	"go-observability-tool/internal/websocket"
	"log"
	"net/http"

	gw "github.com/gorilla/websocket"
)

type MetricsHandler struct {
	wsHub *Hub
}

func NewMetricsHandler(Hub *Hub) *MetricsHandler {
	return &MetricsHandler{
		wsHub: Hub,
	}
}

// endpoint for metric publishers
func (h *MetricsHandler) NewPublisher(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgradeToWS(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	publisherName := r.URL.Query().Get("name")

	if publisherName == "" {
		http.Error(w, "Publisher name is required", http.StatusBadRequest)
		return
	}

	newPublisher := &websocket.Connection{
		Conn: conn,
		IP:   r.RemoteAddr,
		Name: publisherName,
		Role: "publisher",
	}

	h.wsHub.registerChan <- newPublisher
}

// endpoint for a client (dashboard) that will get metrics from other sys
func (h *MetricsHandler) NewMetricSubscriber(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgradeToWS(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for metric := range h.wsHub.MetricsChan {
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
