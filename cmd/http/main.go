package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ozzy-cox/automatic-message-system/config"
	_ "github.com/ozzy-cox/automatic-message-system/docs"
	"github.com/ozzy-cox/automatic-message-system/internal/db"
	"github.com/ozzy-cox/automatic-message-system/internal/handlers"

	"github.com/go-chi/chi"
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
	cfg, err := config.GetAPIConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	_, err = db.GetConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Could not load database: %v", err)
	}

	addr := ":" + cfg.HTTP.Port
	swaggerAddr := cfg.HTTP.Host + addr
	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", swaggerAddr)),
	))

	r.Get("/sent-messages", handlers.HandleGetSentMessages)
	r.Post("/toggle-worker", handlers.HandleToggleWorker)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
