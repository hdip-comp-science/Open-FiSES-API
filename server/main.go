package main

import (
	"net/http"
	"os"

	"github.com/Open-FiSE/go-rest-api/internal/database"
	"github.com/Open-FiSE/go-rest-api/internal/document"
	transportHTTP "github.com/Open-FiSE/go-rest-api/internal/transport/http"
	log "github.com/sirupsen/logrus"
)

// App - defines the application properties.
type App struct {
	Name    string
	Version string
}

const port string = ":4000"

// Run - sets up and starts the application
func (app *App) Run() error {
	// change the o/p of log to json format
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(
		log.Fields{
			"AppName":    app.Name,
			"AppVersion": app.Version,
		}).Info("Setting up Application")

	os.Setenv("DB_USERNAME", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_TABLE", "postgres")
	os.Setenv("DB_PORT", "5433")

	// connection to the DB will be used concurrently across all incoming API-calls
	db, err := database.NewDatabase()
	if err != nil {
		log.Error("Error: Failed to setup database connection")
	}
	//
	err = database.MigrateDB(db)
	if err != nil {
		log.Error("Error: Failed to migrate database")
	}

	documentService := document.NewService(db)

	handler := transportHTTP.NewHandler(documentService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(port, handler.Router); err != nil {
		log.Error("failed to setup web server")
	}

	return nil
}

func main() {
	// instantiate the application
	app := App{
		Name:    "FiSES API Service",
		Version: "1.0.0",
	}
	if err := app.Run(); err != nil {
		log.Error("Failed to start App")
		log.Fatal(err)
	}

}
