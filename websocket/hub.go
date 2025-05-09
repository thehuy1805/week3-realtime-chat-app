package websocket

import (
	"context"
	"log"
	"realtime-chat/db"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	rd         *db.RedisDB
}

func NewHub(rd *db.RedisDB) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		rd:         rd,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered. Total clients: %d", len(h.clients))
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				close(client.Send)
				delete(h.clients, client)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered. Total clients: %d", len(h.clients))
		case message := <-h.broadcast:
			// Lưu tin nhắn vào Redis trước khi phát sóng
			if err := h.rd.Client.LPush(context.Background(), "messages", message).Err(); err != nil {
				log.Printf("Lỗi khi lưu tin nhắn vào Redis: %v", err)
			}
			// Phát sóng tin nhắn đến tất cả client
			h.mu.RLock()
			for client := range h.clients {
				client.Send <- message
			}
			h.mu.RUnlock()
		}
	}
}
