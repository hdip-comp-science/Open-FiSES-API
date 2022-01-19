package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Open-FiSE/go-rest-api/internal/document"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

// Handler - stores pointer to our document service
type Handler struct {
	Router  *mux.Router
	Service *document.Service
}

// NewHandler - returns a pointer to a Handler
func NewHandler(service *document.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// SetupRoutes - sets up all routes for the application
func (h *Handler) SetupRoutes() {
	glog.Info("Setting Up Routes")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/v1/document", h.GetAllDocuments).Methods("GET")
	h.Router.HandleFunc("/api/v1/document", h.PostDocument).Methods("POST")
	h.Router.HandleFunc("/api/v1/document/{id}", h.UpdateDocument).Methods("PUT")
	h.Router.HandleFunc("/api/v1/document/{id}", h.GetDocument).Methods("GET")
	h.Router.HandleFunc("/api/v1/document/{id}", h.DeleteDocument).Methods("DELETE")

	h.Router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		glog.Info("health check success - 200 OK")
	})
}

// GetDocument - retrieve a single document by ID
func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	document, err := h.Service.GetDocument(uint(i))
	if err != nil {
		fmt.Fprintf(w, "Error Retrieving Document by ID")
	}

	fmt.Fprintf(w, "%+v", document)
}

// GetAllDocuments - fetch all documents from the document service
func (h *Handler) GetAllDocuments(w http.ResponseWriter, r *http.Request) {

	documents, err := h.Service.GetAllDocuments()
	if err != nil {
		fmt.Fprintf(w, "Failed to retrieve documents")
	}
	fmt.Fprintf(w, "%+v", documents)
}

// PostDocument - adds a new document
func (h *Handler) PostDocument(w http.ResponseWriter, r *http.Request) {

	document, err := h.Service.PostDocument(document.Document{
		Path: "/",
	})
	if err != nil {
		fmt.Fprintf(w, "Failed to post new document")
	}
	fmt.Fprintf(w, "%+v", document)
}

// UpdateDocument - update an exisiting document by ID
func (h *Handler) UpdateDocument(w http.ResponseWriter, r *http.Request) {

	document, err := h.Service.UpdateDocument(1, document.Document{
		Path: "/new",
	})
	if err != nil {
		fmt.Fprintf(w, "Failed to update document")
	}
	fmt.Fprintf(w, "%+v", document)
}

// DeleteDocument - delete a document by ID
func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	err = h.Service.DeleteDocument(uint(commentID))
	if err != nil {
		fmt.Fprintf(w, "Failed to delete document")
	}
	fmt.Fprintf(w, "Successfully deleted document")
}
