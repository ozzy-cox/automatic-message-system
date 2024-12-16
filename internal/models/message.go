package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"name"`
	To        string    `json:"to"`
	Sent      bool      `json:"sent"`
	SentAt    time.Time `json:"sent_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
