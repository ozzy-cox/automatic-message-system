package db_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	conn, err := sql.Open("sqlite3", ":memory:")

	_, err = conn.Exec(`
		CREATE TABLE messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content VARCHAR(1023) CHECK(length(content) < 1023),
			to_ VARCHAR(1023),
			is_sent  BOOLEAN DEFAULT false,
			sent_at TIMESTAMP NULL DEFAULT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	    `)
	require.NoError(t, err)

	return conn
}

func TestGetMessages(t *testing.T) {
	conn := setupTestDB(t)
	defer conn.Close()

	repo := db.NewMessageRepository(conn)

	now := time.Now()
	_, err := conn.Exec(`
		INSERT INTO messages (content, to_, is_sent, created_at) VALUES
		($1, $2, false, $3),
		($4, $5, false, $6),
		($7, $8, true, $9)
	`, "Test1", "1", now, "Test2", "2", now, "Test3", "3", now)
	require.NoError(t, err)

	messages := make([]*db.Message, 0)
	for msg, err := range repo.GetMessages(10, 0) {
		require.NoError(t, err)
		messages = append(messages, msg)
	}

	assert.Len(t, messages, 3)
	assert.False(t, messages[0].Sent)
	assert.False(t, messages[1].Sent)
}

func TestGetSentMessages(t *testing.T) {
	conn := setupTestDB(t)
	defer conn.Close()

	repo := db.NewMessageRepository(conn)

	// Insert test data
	now := time.Now()
	_, err := conn.Exec(`
		INSERT INTO messages (content, to_, is_sent, sent_at, created_at) VALUES
		($1, $2, true, $3, $4),
		($5, $6, true, $7, $8),
		($9, $10, false, NULL, $11)
	`, "Test1", "1", now, now, "Test2", "2", now, now, "Test3", "3", now)
	require.NoError(t, err)

	// Test getting sent messages
	messages := make([]*db.Message, 0)
	for msg, err := range repo.GetSentMessages(10, 0) {
		require.NoError(t, err)
		messages = append(messages, msg)
	}

	assert.Len(t, messages, 2)
	assert.True(t, messages[0].Sent)
	assert.True(t, messages[1].Sent)
}

func TestMarkMessageAsSent(t *testing.T) {
	conn := setupTestDB(t)
	defer conn.Close()

	repo := db.NewMessageRepository(conn)

	var messageID int
	err := conn.QueryRow(`
		INSERT INTO messages (content, to_, is_sent)
		VALUES ($1, $2, false)
		RETURNING id
	`, "Test", "123").Scan(&messageID)
	require.NoError(t, err)

	err = repo.MarkMessageAsSent(messageID)
	assert.NoError(t, err)

	var sent bool
	err = conn.QueryRow("SELECT is_sent FROM messages WHERE id = $1", messageID).Scan(&sent)
	require.NoError(t, err)
	assert.True(t, sent)
}

func TestMarkMessageAsSentNonExistent(t *testing.T) {
	conn := setupTestDB(t)
	defer conn.Close()

	repo := db.NewMessageRepository(conn)

	// Try to mark non-existent message as sent
	err := repo.MarkMessageAsSent(99999)
	assert.Error(t, err)
}

func TestMessageMaxLength(t *testing.T) {
	conn := setupTestDB(t)
	defer conn.Close()

	longContent := strings.Repeat("a", 1024) // One character over limit
	_, err := conn.ExecContext(context.Background(),
		"INSERT INTO messages (content, to_) VALUES (?, ?)",
		longContent, "test@example.com",
	)

	assert.Error(t, err, "Expected error when content exceeds 1023 characters")
}
