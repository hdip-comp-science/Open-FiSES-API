package http

import (
	"encoding/json"
	"net/http"

	"github.com/Open-FiSE/go-rest-api/internal/document"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handler - stores pointer to our document service
type Handler struct {
	Router  *mux.Router
	Service *document.Service
}

// Response - an object to store repsonses from the API
type Response struct {
	Message string
}

// NewHandler - returns a pointer to a Handler
func NewHandler(service *document.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// Logger - adds middleware around endpoints
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//  improve visibility by checking endpoints and paths being consumed.
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			}).Info("handled request")
		next.ServeHTTP(w, r)
	})
}

// SetupRoutes - sets up all routes for the application
func (h *Handler) SetupRoutes() {
	log.Info("Setting Up Routes")
	h.Router = mux.NewRouter()

	// tell router to use middleware function
	// improves the way logging is handled.
	h.Router.Use(Logger)

	h.Router.HandleFunc("/api/v1/document", h.GetAllDocuments).Methods("GET")
	h.Router.HandleFunc("/api/v1/document", h.PostDocument).Methods("POST")
	h.Router.HandleFunc("/api/v1/document/{id}", h.UpdateDocument).Methods("PUT")
	h.Router.HandleFunc("/api/v1/document/{id}", h.GetDocument).Methods("GET")
	h.Router.HandleFunc("/api/v1/document/{id}", h.DeleteDocument).Methods("DELETE")

	h.Router.HandleFunc("/api/v1/upload", h.Upload)

	h.Router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{Message: "HTTP Status: 200 OK"}); err != nil {
			log.Warning(err)
		}
	})
	log.Info("App Setup Complete...")
}
