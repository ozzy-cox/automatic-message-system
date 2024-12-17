package db

import (
	"database/sql"
)

type Message struct {
	ID        int          `json:"id" db:"id"`
	Content   string       `json:"content" db:"content"`
	To        string       `json:"to" db:"to_"`
	Sent      bool         `json:"sent" db:"sent"`
	SentAt    sql.NullTime `json:"sent_at,omitempty" db:"sent_at"`
	CreatedAt sql.NullTime `json:"created_at" db:"created_at"`
}
