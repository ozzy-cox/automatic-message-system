package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/ozzy-cox/automatic-message-system/docs"
	"github.com/ozzy-cox/automatic-message-system/internal/api"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/swaggo/http-swagger/v2"
)

//	@title			Automatic Message System API
//	@version		1.0
//	@description	Automatic message sending service
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
func main() {
	cfg, err := api.GetAPIConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	loggerInst, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	dbConn, err := db.NewConnection(cfg.Database)
	if err != nil {
		loggerInst.Fatalf("Could not load database: %v", err)
	}

	service := api.Service{
		Config:            cfg,
		MessageRepository: db.NewMessageRepository(dbConn),
		Logger:            loggerInst,
	}

	addr := ":" + cfg.Port
	swaggerAddr := cfg.Host + addr
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", swaggerAddr)),
	))
	r.Get("/sent-messages", service.HandleGetSentMessages)
	r.Post("/toggle-worker", service.HandleToggleWorker)

	if err := http.ListenAndServe(addr, r); err != nil {
		loggerInst.Fatalf("Could not start server: %v", err)
	}
}
