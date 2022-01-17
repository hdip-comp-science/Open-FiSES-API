package http

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

// Handler - stores pointer to our document service
type Handler struct {
	Router *mux.Router
}

// NewHandler - returns a pointer to a Handler
func NewHandler() *Handler {
	return &Handler{}
}

// SetupRoutes - sets up all routes for the application
func (h *Handler) SetupRoutes() {
	glog.Info("Setting Up Routes")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		glog.Info("health check success - 200 OK")
	})
}
