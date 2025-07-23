package websocket

import gw "github.com/gorilla/websocket"

type Connection struct {
	Conn *gw.Conn
	IP   string
	Name string
	Role string // "publisher" or listener
}
