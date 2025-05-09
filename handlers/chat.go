package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"realtime-chat/db"
	ws "realtime-chat/websocket"

	gorilla "github.com/gorilla/websocket"
)

type ChatHandler struct {
	hub *ws.Hub
	pg  *db.PostgresDB
	rd  *db.RedisDB
}

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Điều chỉnh cho môi trường production
	},
}

func NewChatHandler(hub *ws.Hub, pg *db.PostgresDB, rd *db.RedisDB) *ChatHandler {
	return &ChatHandler{hub: hub, pg: pg, rd: rd}
}

func (h *ChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Lỗi nâng cấp WebSocket:", err)
		return
	}

	client := ws.NewClient(h.hub, conn, username)
	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *ChatHandler) GetMessageHistory(w http.ResponseWriter, r *http.Request) {
	messages, err := h.rd.Client.LRange(r.Context(), "messages", 0, 100).Result()
	if err != nil {
		http.Error(w, "Không thể lấy lịch sử tin nhắn", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
