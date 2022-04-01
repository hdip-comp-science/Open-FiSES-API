package http

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/Open-FiSE/go-rest-api/internal/document"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// open app to localhost:4000 origin. Solves CORS issue with client app
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
}

// PostDocument - adds a new document
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// 1. parse input , type multipart/form-data
	r.ParseMultipartForm(3 << 30) // Maximum upload of 3 MB files

	// 2. retrieve data from file posted form-date
	//    get fileHeader for filename, size and headers
	file, fileHeader, err := r.FormFile("file") // where is the
	if err != nil {
		log.Warning(err)
		http.Error(w, "Error Retrieving file from form-data", http.StatusInternalServerError)
		return
	}

	defer file.Close()
	// generate sha256 hash value of file in memory.
	// source: https://yourbasic.org/golang/hash-md5-sha256-string-file/
	// create a new hash.Hash form crypto pkg "crypto/sha256"
	hash := sha256.New()
	// implments the io.Writer function, copy from src "hash" to dst "file"
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	// extract the checksum by calling its Sum function
	sum := hash.Sum(nil)
	// convert from hex to string using "encoding/hex" pkg
	log.Infof("%s", hex.EncodeToString(sum[:]))

	// print file data to console
	log.Infof("Uploading File: %+v\n", fileHeader.Filename)
	log.Infof("File Size: %+v\n", fileHeader.Size)
	log.Infof("MIME Header: %+v\n", fileHeader.Header)

	// func ReadAll(r io.Reader) ([]byte, error)
	// ReadAll is defined to read from src until EOF and returns the data it read
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
	}
	path := "/app/docs/"
	// Create the uploads folder if it doesn't already exist
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var document document.Document

	// assign document values to document struct fileds
	document.Path = path + fileHeader.Filename
	document.Title = fileHeader.Filename
	document.Version = "1.0"
	document.Author = ""                       //TBD when auth is implemented
	document.Hash = hex.EncodeToString(sum[:]) // `:` is needed because byte arrays cannot be directly turned to a string while slices can

	// Post to document service
	document, err = h.Service.PostDocument(document)
	if err != nil {
		fmt.Fprintf(w, "Failed to post new document")
	}
	// upload to final destination path
	//  print sha here
	// write data to named file. If file does not exist WriteFile creates it.
	err = os.WriteFile(path+fileHeader.Filename, fileBytes, 0644)
	if err != nil {
		log.Error(err)
	}

	// 4. return whether or not this has been successful
	log.Infof("Successfully uploaded file: %s\n", document.Title)
	w.WriteHeader(http.StatusOK)

}

// GetDocument - retrieve a single document by ID
func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
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

	filename := document.Path
	file, err := os.Open(filename)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}
	defer file.Close()

	//get file info for content-length header
	info, err := os.Stat(filename)
	if err != nil {
		http.Error(w, "File stat error: "+err.Error(), http.StatusBadRequest)
		return
	}

	//Set headers
	w.Header().Set("Content-type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))

	//Stream to response
	if _, err := io.Copy(w, file); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}
	// return
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
		log.Warning(err)
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

	documentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	document, err = h.Service.UpdateDocument(uint(documentID), document)
	if err != nil {
		fmt.Fprintf(w, "Failed to update document")
	}

	// Return the newly update document as json
	if err := json.NewEncoder(w).Encode(document); err != nil {
		log.Warning(err)
	}
}

// DeleteDocument - delete a document by ID
func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	id := vars["id"]

	documentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	err = h.Service.DeleteDocument(uint(documentID))
	if err != nil {
		fmt.Fprintf(w, "Failed to delete document")
	}

	if err := json.NewEncoder(w).Encode(Response{Message: "Successfully deleted document"}); err != nil {
		log.Warning(err)
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
		log.Warning(err)
	}
}
