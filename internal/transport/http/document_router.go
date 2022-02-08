package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// SetupRoutes - sets up all routes for the application
func (h *Handler) SetupRoutes() {
	glog.Info("Setting Up Routes")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/v1/document", h.GetAllDocuments).Methods("GET")
	h.Router.HandleFunc("/api/v1/document", h.PostDocument).Methods("POST")
	h.Router.HandleFunc("/api/v1/document/{id}", h.UpdateDocument).Methods("PUT")
	h.Router.HandleFunc("/api/v1/document/{id}", h.GetDocument).Methods("GET")
	h.Router.HandleFunc("/api/v1/document/{id}", h.DeleteDocument).Methods("DELETE")

	h.Router.HandleFunc("/api/v1/upload", h.UploadFile)

	h.Router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{Message: "HTTP Status: 200 OK"}); err != nil {
			glog.Warning(err)
		}
	})
}

// open app to localhost:4000 origin. Solves CORS issue with client app
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	glog.Info(w, "Uploading File\n")

	// 1. parse the input, type multipart/for-data
	r.ParseMultipartForm(10 << 20)
	// 2. retriev file from posted form data
	file, handler, err := r.FormFile("tempFile")
	if err != nil {
		glog.Errorf("Error Retrieving file from form-data", err)
		return
	}
	defer file.Close()
	glog.Info("Uploaded File: %+v\n", handler.Filename)
	glog.Info("File Size: %+v\n", handler.Size)
	glog.Info("MIME Header: %+v\n", handler.Header)

	// 3. write temporary file on the server
	// create temporary file in our project directory
	tempFile, err := ioutil.TempFile("temp-file", "upload-*.pdf")
	if err != nil {
		glog.Errorf("Error creating temp file", err)
		return
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		glog.Warningf("Warning: no file to read", err)
		return
	}
	tempFile.Write(fileBytes)
	// 4. return wether or not this has been succesful
	fmt.Fprintf(w, "Successfully Uploaded File\n")

}

// GetDocument - retrieve a single document by ID
func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

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

	if err := json.NewEncoder(w).Encode(document); err != nil {
		glog.Warning(err)
	}
}

// GetAllDocuments - fetch all documents from the document service
func (h *Handler) GetAllDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	documents, err := h.Service.GetAllDocuments()
	if err != nil {
		fmt.Fprintf(w, "Failed to retrieve documents")
	}
	if err := json.NewEncoder(w).Encode(documents); err != nil {
		glog.Warning(err)
	}
}

// PostDocument - adds a new document
func (h *Handler) PostDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	var document document.Document
	// Parse the request body as document
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		fmt.Fprintf(w, "Failed to decode JSON Body")
	}
	// Post to document service
	document, err := h.Service.PostDocument(document)
	if err != nil {
		fmt.Fprintf(w, "Failed to post new document")
	}
	// return the document
	if err := json.NewEncoder(w).Encode(document); err != nil {
		glog.Warning(err)
	}
}

// UpdateDocument - update an exisiting document by ID
func (h *Handler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	var document document.Document
	// Parse the request body as document
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		fmt.Fprintf(w, "Failed to decode JSON Body")
	}

	vars := mux.Vars(r)
	id := vars["id"]

	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	document, err = h.Service.UpdateDocument(uint(commentID), document)
	if err != nil {
		fmt.Fprintf(w, "Failed to update document")
	}

	// Return the newly update document as json
	if err := json.NewEncoder(w).Encode(document); err != nil {
		glog.Warning(err)
	}
}

// DeleteDocument - delete a document by ID
func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

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

	if err := json.NewEncoder(w).Encode(Response{Message: "Successfully deleted document"}); err != nil {
		glog.Warning(err)
	}
}
