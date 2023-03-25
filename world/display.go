package world

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
)

type DisplayClient struct {
	conn *websocket.Conn
}

func (c *DisplayClient) Send(data []byte) {
	PrintErr(c.conn.WriteMessage(websocket.TextMessage, data))
}

func NewDisplayClient(name string) *DisplayClient {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/world", RawQuery: fmt.Sprintf("name=%s", name)}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil
	}

	return &DisplayClient{conn: c}
}
