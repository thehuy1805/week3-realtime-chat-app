package models

import "time"

type User struct {
	ID        int
	Username  string
	CreatedAt time.Time
}
