package main

import (
	"net/http"
	"os"

	"github.com/Open-FiSE/go-rest-api/internal/database"
	"github.com/Open-FiSE/go-rest-api/internal/document"
	transportHTTP "github.com/Open-FiSE/go-rest-api/internal/transport/http"
	log "github.com/sirupsen/logrus"
)

type App struct {
	Name    string
	Version string
}

// Run - sets up the application
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
	os.Setenv("DB_PORT", "5432")

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

	if err := http.ListenAndServe(":4000", handler.Router); err != nil {
		log.Error("failed to setup web server")
	}

	return nil
}

func main() {

	// // This is needed to make `glog` believe that the flags have already been parsed, otherwise
	// // every log messages is prefixed by an error message stating the the flags haven't been parsed.
	// _ = flag.CommandLine.Parse([]string{})

	// // Always log to stderr by default
	// if err := flag.Set("logtostderr", "true"); err != nil {
	// 	glog.Infof("Unable to set logtostderr to true")
	// }
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
