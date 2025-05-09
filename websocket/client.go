package websocket

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	Send     chan []byte
	username string
}

func NewClient(hub *Hub, conn *websocket.Conn, username string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		Send:     make(chan []byte),
		username: username,
	}
}

func (c *Client) ReadPump() {
	if c.username != "" {
		c.hub.rd.Client.Set(context.Background(), "presence:"+c.username, "online", 0)
	}

	defer func() {
		if c.username != "" {
			c.hub.rd.Client.Set(context.Background(), "presence:"+c.username, "offline", 0)
		}
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// Gửi tin nhắn hiện tại
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
			// Kiểm tra và gửi tất cả tin nhắn còn lại trong kênh Send
			for len(c.Send) > 0 {
				message, ok = <-c.Send
				if !ok {
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}
			}
		}
	}
}
