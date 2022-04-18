package http

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
	// (*w).Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
}

// PostDocument - adds a new document
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	// 1. parse input , type multipart/form-data
	//    defines how bug the chunk size of the data that will be received
	err := r.ParseMultipartForm(32 << 20) // chuck size in bytes - max upload of 32 MB files
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. retrieve data from file posted form-data
	//    get a reference to the fileHeaders, they are accessible only after ParseMultipartForm is called
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Error(err)
		http.Error(w, "Error Retrieving file from form-data", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	// generate sha256 hash value of file in memory.
	// create a new hash.Hash form crypto pkg "crypto/sha256"
	sha256 := sha256.New()
	// Copy the uploaded file to the filesystem at the specified destination
	if _, err := io.Copy(sha256, file); err != nil {
		log.Fatal(err)
	}

	// extract the checksum by calling its Sum function
	sum := sha256.Sum(nil)
	// convert from hex to string using "encoding/hex" pkg
	log.Infof("%s", hex.EncodeToString(sum[:]))

	shaStr := hex.EncodeToString(sum)

	// Seek sets the offset for the next Read or Write to offsetback to the start of the file again
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Errorf("Seek start of file error: %v", err)
	}

	// print file data to console
	log.Infof("Uploading File: %+v\n", fileHeader.Filename)
	log.Infof("File Size: %+v\n", fileHeader.Size)
	log.Infof("MIME Header: %+v\n", fileHeader.Header)

	path := "/app/docs/"

	// Create the uploads folder if it doesn't already exist
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var document document.Document

	// assign document version an initial value of 1.0
	var docVer float32 = document.Version + 1.0

	// get first matched record of hash value in documents table
	if dbErr := h.Service.DB.Where("hash = ?", shaStr).First(&document).Error; dbErr != nil {
		// no hash match in db
		log.Warnf(dbErr.Error())
	}
	// if no has match check if filename exists. If so, update the version no. else post new document
	if err := h.Service.DB.Where("title = ?", document.Title).First(&document).Error; err != nil {
		document.Path = path + fileHeader.Filename
		document.Title = fileHeader.Filename
		document.Version = docVer
		document.Author = ""                       //TBD when auth is implemented
		document.Hash = hex.EncodeToString(sum[:]) // `:` is needed because byte arrays cannot be directly turned to a string while slices can

		document, err = h.Service.PostDocument(document)
		if err != nil {
			fmt.Fprintf(w, "Failed to update document version no.")
		}
	} else {
		document.Version = docVer + 1.0
		log.Infof("updating document version to:", document.Version)

		document, err = h.Service.UpdateDocument(document.ID, document)
		if err != nil {
			fmt.Fprintf(w, "Failed to update document version no.")
		}

	}

	// Create a new file in the '/app/docs' directory
	dst, err := os.Create(path + fileHeader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem.
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. return whether or not this has been successful
	log.Infof("Successfully uploaded file: %s\n", document.Title)
}

// GetDocument - retrieve a single document by ID
func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	id := vars["id"]

	// GetDocument is expecting a uint, parse string and set to base 10, size 64
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
