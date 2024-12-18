package db

import (
	"database/sql"
	"fmt"
	"iter"
)

type MessageRepository struct {
	DB *sql.DB
}

func (r *MessageRepository) GetUnsentMessagesFromDb(limit, offset int) iter.Seq2[*Message, error] {
	rows, err := r.DB.Query("SELECT * FROM messages WHERE is_sent is false LIMIT $1 OFFSET $2", limit, offset)
	return func(yield func(*Message, error) bool) {
		if err != nil {
			if !yield(nil, fmt.Errorf("Error reading unsent messages from db: %s\n", err)) {
				return
			}
		}

		for rows.Next() {
			var msg Message
			err := rows.Scan(
				&msg.ID,
				&msg.Content,
				&msg.To,
				&msg.Sent,
				&msg.SentAt,
				&msg.CreatedAt,
			)

			if err != nil {
				if !yield(nil, fmt.Errorf("Failed to scan messages: %s\n", err)) {
					return
				}
			}
			if !yield(&msg, nil) {
				return
			}
		}
	}
}

func (r *MessageRepository) GetSentMessagesFromDb(limit, offset int) iter.Seq2[*Message, error] {
	rows, err := r.DB.Query("SELECT * FROM messages WHERE sending_status is true LIMIT $1 OFFSET $2", limit, offset)
	return func(yield func(*Message, error) bool) {
		if err != nil {
			if !yield(nil, fmt.Errorf("Error reading sent messages from db: %s\n", err)) {
				return
			}
		}

		for rows.Next() {
			var msg Message
			err := rows.Scan(
				&msg.ID,
				&msg.Content,
				&msg.To,
				&msg.Sent,
				&msg.SentAt,
				&msg.CreatedAt,
			)

			if err != nil {
				if !yield(nil, fmt.Errorf("Failed to scan messages: %s\n", err)) {
					return
				}
			}
			if !yield(&msg, nil) {
				return
			}
		}
	}
}

func (r *MessageRepository) SetMessageSent(messageId int) error { // TODO do batching
	result, err := r.DB.Exec("UPDATE messages SET is_sent = true WHERE id = $1", messageId)

	if err != nil {
		return fmt.Errorf("Error updating message sent state: %s\n", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Rows affected may not be supported: %s\n", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Row with id: %d not found: %s\n", messageId, err)
	}

	return nil
}
