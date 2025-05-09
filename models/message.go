package models

import "time"

type Message struct {
	ID          int
	SenderID    int
	RecipientID int
	Content     string
	CreatedAt   time.Time
}
