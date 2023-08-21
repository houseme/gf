package gws

import (
	"net/http"

	"github.com/lxzan/gws"

	"github.com/gogf/gf/v2/net/gsocket"
)

func New() gsocket.WebSocketServer {
	return &BuiltinEventHandler{}
}

// BuiltinEventHandler is the default event handler for websocket server.
type BuiltinEventHandler struct {
	conn *gws.Conn
}

// Connect is the callback function for new connection creating.
func (h *BuiltinEventHandler) Connect(r *http.Request, w http.ResponseWriter) (err error) {
	upgrader := gws.NewUpgrader(new(gws.BuiltinEventHandler), &gws.ServerOption{})
	if h.conn, err = upgrader.Upgrade(w, r); err != nil {
		return err
	}
	return nil
}

// ReadMessage reads and returns the message from the client.
func (h *BuiltinEventHandler) ReadMessage() (messageType int, p []byte, err error) {
	h.conn.SessionStorage.Store("username", "lxzan")
	go h.conn.ReadLoop()

	return
}

// SendMessage writes the message to the client.
func (h *BuiltinEventHandler) SendMessage(messageType int, data []byte) error {
	return nil
}

// Close closes the connection.
func (h *BuiltinEventHandler) Close() error {
	return nil
}
