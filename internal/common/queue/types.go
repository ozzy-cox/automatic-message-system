package queue

import "time"

type MessagePayload struct {
	ID        int        `json:"id"`
	Content   string     `json:"content"`
	To        string     `json:"to"`
	CreatedAt *time.Time `json:"created_at"`
}
