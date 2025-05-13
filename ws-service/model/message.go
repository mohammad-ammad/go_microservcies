package model

import "time"

type Message struct {
	UserID    uint      `json:"userId"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
