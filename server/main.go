package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Open-FiSE/go-rest-api/internal/database"
	"github.com/Open-FiSE/go-rest-api/internal/document"
	transportHTTP "github.com/Open-FiSE/go-rest-api/internal/transport/http"
	"github.com/golang/glog"
)

type App struct{}

// Run - sets up the application
func (app *App) Run() error {
	glog.Info("Setting up App")

	db, err := database.NewDatabase()
	if err != nil {
		glog.Error("Error: Failed to setup database connection")
	}
	//
	err = database.MigrateDB(db)
	if err != nil {
		glog.Error("Error: Failed to migrate database")
	}

	documentService := document.NewService(db)

	handler := transportHTTP.NewHandler(documentService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":3000", handler.Router); err != nil {
		return fmt.Errorf("failed to setup web server, %v", err)
	}

	return nil
}

func main() {

	// This is needed to make `glog` believe that the flags have already been parsed, otherwise
	// every log messages is prefixed by an error message stating the the flags haven't been parsed.
	_ = flag.CommandLine.Parse([]string{})

	// Always log to stderr by default
	if err := flag.Set("logtostderr", "true"); err != nil {
		glog.Infof("Unable to set logtostderr to true")
	}

	glog.Info("Welcome to the beginning of this project")

	app := App{}
	if err := app.Run(); err != nil {
		glog.Warningf("Failed to start App, %v", err)
	}

}
