package http

import (
	"encoding/json"
	"net/http"

	"github.com/Open-FiSE/go-rest-api/internal/booking"
	"github.com/Open-FiSE/go-rest-api/internal/document"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handler - store a pointer to the router and document service that the app uses
type Handler struct {
	Router      *mux.Router
	Service     *document.Service
	BookService *booking.BookService
}

// Response - an object to store repsonses from the API
type Response struct {
	Message string
}

// NewHandler - returns a pointer to a Handler
func NewHandler(service *document.Service, bookservice *booking.BookService) *Handler {
	return &Handler{
		Service:     service,
		BookService: bookservice,
	}
}

// Logger - is a middleware handler available globally that wraps around all endpoints
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//  improves visibility by checking endpoints and paths being consumed.
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			}).Info("handled request")
		// Pass down the request to the next handler
		next.ServeHTTP(w, r)
	})
}

const apiPrefix string = "/api/v1/"

// SetupRoutes - sets up all routes for the application
func (h *Handler) SetupRoutes() {
	log.Info("Setting Up Routes")
	h.Router = mux.NewRouter()

	// Logger middleware function installed globally, which will be called for each request
	// improves the way logging is handled.
	h.Router.Use(Logger)

	// Document Service Routes
	h.Router.HandleFunc(apiPrefix+"document", h.GetAllDocuments).Methods("GET")
	h.Router.HandleFunc(apiPrefix+"document", h.PostDocument).Methods("POST")
	h.Router.HandleFunc(apiPrefix+"document/{id}", h.UpdateDocument).Methods("PUT")
	h.Router.HandleFunc(apiPrefix+"document/{id}", h.GetDocument).Methods("GET")
	h.Router.HandleFunc(apiPrefix+"document/{id}", h.DeleteDocument).Methods("DELETE")
	h.Router.HandleFunc(apiPrefix+"upload", h.Upload).Methods("POST")

	// Booking Service Routes
	h.Router.HandleFunc(apiPrefix+"booking", h.GetAllBookings).Methods("GET")
	h.Router.HandleFunc(apiPrefix+"booking", h.PostBooking).Methods("POST")
	h.Router.HandleFunc(apiPrefix+"booking/{id}", h.UpdateBooking).Methods("PUT")
	h.Router.HandleFunc(apiPrefix+"booking/{id}", h.GetBooking).Methods("GET")
	h.Router.HandleFunc(apiPrefix+"booking/{id}", h.DeleteBooking).Methods("DELETE")

	h.Router.HandleFunc(apiPrefix+"health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{Message: "HTTP Status: 200 OK"}); err != nil {
			log.Warning(err)
		}
	})
	log.Info("Routes Setup Complete...")
}
