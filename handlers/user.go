package handlers

import (
	"encoding/json"
	"net/http"
	"realtime-chat/db"
	"realtime-chat/models"
)

type UserHandler struct {
	pg *db.PostgresDB
	rd *db.RedisDB
}

func NewUserHandler(pg *db.PostgresDB, rd *db.RedisDB) *UserHandler {
	return &UserHandler{pg: pg, rd: rd}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO users (username) VALUES ($1) RETURNING id, created_at`
	err := h.pg.DB.QueryRow(query, user.Username).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Set user presence in Redis
	h.rd.Client.Set(r.Context(), "presence:"+user.Username, "online", 0)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetPresence(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	status, err := h.rd.Client.Get(r.Context(), "presence:"+username).Result()
	if err != nil {
		status = "offline"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"username": username, "status": status})
}
