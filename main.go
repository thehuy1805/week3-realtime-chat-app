package main

import (
	"log"
	"net/http"
	"realtime-chat/config"
	"realtime-chat/db"
	"realtime-chat/handlers"
	"realtime-chat/websocket"
	"time"

	"github.com/didip/tollbooth"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()

	// Khởi tạo cơ sở dữ liệu
	pg, err := db.NewPostgresDB(cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	rd, err := db.NewRedisDB(cfg.RedisAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer rd.Close()

	// Khởi tạo WebSocket hub với Redis
	hub := websocket.NewHub(rd)
	go hub.Run()

	// Khởi tạo handlers
	userHandler := handlers.NewUserHandler(pg, rd)
	chatHandler := handlers.NewChatHandler(hub, pg, rd)

	// Thiết lập router
	r := chi.NewRouter()

	// Áp dụng giới hạn tỷ lệ (100 yêu cầu mỗi phút cho mỗi IP)
	limiter := tollbooth.NewLimiter(100/60, nil)
	limiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	limiter.SetTokenBucketExpirationTTL(time.Minute)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := tollbooth.LimitByRequest(limiter, w, r); err != nil {
				http.Error(w, err.Message, err.StatusCode)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Định nghĩa các tuyến đường (routes)
	r.Post("/register", userHandler.Register)
	r.Get("/ws", chatHandler.HandleWebSocket)
	r.Get("/messages", chatHandler.GetMessageHistory)
	r.Get("/presence", userHandler.GetPresence)

	// Khởi động server
	log.Printf("Server is starting on port %s", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
