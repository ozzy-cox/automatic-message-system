package api

import (
	"database/sql"
	"log"

	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
)

type APIDeps struct {
	DBConnection *sql.DB
	Logger       *logger.Logger
}

func NewAPIDeps(cfg APIConfig) *APIDeps {
	loggerInst, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	dbConnection, err := db.NewConnection(cfg.Database)
	if err != nil {
		loggerInst.Fatalf("Could not load database: %v", err)
	}

	return &APIDeps{
		DBConnection: dbConnection,
		Logger:       loggerInst,
	}
}

func (d *APIDeps) Cleanup() {
	d.Logger.Println("Cleaning up")
}
