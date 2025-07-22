package handler

import (
	"fmt"
	"go-observability-tool/internal/websocket"
	"log"
	"net/http"

	gw "github.com/gorilla/websocket"
)

type MetricsHandler struct {
	wsListener websocket.Listener
}

func NewMetricsHandler(listener *websocket.Listener) *MetricsHandler {
	go listener.Listen()

	return &MetricsHandler{
		wsListener: *listener,
	}
}

func (h *MetricsHandler) NewSender(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgradeToWS(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.wsListener.AddConnection(conn)
}

func (h *MetricsHandler) NewMetricObserver(w http.ResponseWriter, r *http.Request) {
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
