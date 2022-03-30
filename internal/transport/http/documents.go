package http

import (
	"bytes"
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

// func getHash(filename string) (uint32, error) {
// 	// open the file
// 	f, err := os.Open(filename)
// 	if err != nil {
// 		return 0, err
// 	}
// 	// remember to always close opened files
// 	defer f.Close()

// 	// create a hasher

// 	h := crc32.NewIEEE()
// 	// copy the file into the hasher
// 	// - copy takes (dst, src) and returns (bytesWritten, error)
// 	_, err = io.Copy(h, f)
// 	// we don't care about how many bytes were written, but we do want to
// 	// handle the error
// 	if err != nil {
// 		return 0, err
// 	}
// 	return h.Sum32(), nil
// }

const chunkSize = 64000

func deepCompare(file1, file2 string) bool {
	// Check file size ...

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

// PostDocument - adds a new document
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	path := "/app/docs/"
	enableCors(&w)

	// 1. parse input , type multipart/form-data
	r.ParseMultipartForm(3 << 30) // set constraints on file upload size

	// 2. retrieve data from file posted form-date
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Warning(err)
		http.Error(w, "Error Retrieving file from form-data", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	// implement file check here...

	deepCompare(path+fileHeader.Filename, path+"cheat_sheet_bash.pdf")

	// Stat returns a FileInfo describing the named file.
	// if _, err := os.Stat(document.Path); err == nil {
	// 	log.Errorf("%s file already exists \n", document.Path)
	// } else {
	// 	log.Infof("%s is a new file \n", document.Path)
	// }

	// print headers to console
	log.Infof("Uploading File: %+v\n", fileHeader.Filename)
	log.Infof("File Size: %+v\n", fileHeader.Size)
	log.Infof("MIME Header: %+v\n", fileHeader.Header)

	// func ReadAll(r io.Reader) ([]byte, error)
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
	}

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
	document.Author = "" //TBD when auth is implemented

	// Post to document service
	document, err = h.Service.PostDocument(document)
	if err != nil {
		fmt.Fprintf(w, "Failed to post new document")
	}

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
