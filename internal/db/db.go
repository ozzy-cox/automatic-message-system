package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ozzy-cox/automatic-message-system/config"
)

var DbConnection *sql.DB

func GetConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	if DbConnection != nil {
		return DbConnection, nil
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	DbConnection = db

	return db, nil
}
