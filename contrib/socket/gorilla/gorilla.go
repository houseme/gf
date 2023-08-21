package gorilla

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	w http.ResponseWriter
	r *http.Request
)

func init() {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	g.Log().Debug(context.Background(), "websocket connect success")
}

type DefaultWebSocket struct {
	// websocket internal state
	conn *websocket.Conn
}

func (ws *DefaultWebSocket) Connect() error {
	// connect to websocket
	conn, _, err := websocket.DefaultDialer.Dial("ws://example.com/ws", nil)
	if err != nil {
		return err
	}

	ws.conn = conn
	return nil
}

func (ws *DefaultWebSocket) ReadMessage() (int, []byte, error) {
	// read message from websocket
	_, message, err := ws.conn.ReadMessage()
	if err != nil {
		return 0, nil, err
	}

	return websocket.TextMessage, message, nil
}

func (ws *DefaultWebSocket) SendMessage(messageType int, data []byte) error {
	// send message via websocket
	err := ws.conn.WriteMessage(messageType, data)
	if err != nil {
		return err
	}

	return nil
}

func (ws *DefaultWebSocket) Close() error {
	// close websocket
	err := ws.conn.Close()
	return err
}
