package document

import (
	"os"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

//  Service - the struct for the document service
type Service struct {
	DB *gorm.DB
}

// Document - Defines the Document Model Structure
type Document struct {
	gorm.Model
	Path    string  `json:"path"`
	Title   string  `json:"title"`
	Version float32 `json:"version"`
	Author  string  `json:"author"`
	Body    string  `json:"body"`
	Hash    string  `json:"hash"`
}

// https://www.baeldung.com/linux/sha-256-from-command-line
// Hash    hash.Hash `json:"hash"`
// DocumentService - Defines the contract in against which you have to
// implement the document service
type DocumentService interface {
	GetDocument(ID uint) (Document, error)
	GetDocumentByPath(path string) ([]Document, error)
	PostDocument(document Document) (Document, error)
	UpdateDocument(ID uint, newDocument Document) (Document error)
	DeleteDocument(ID uint) error
	GetAllDocuments() ([]Document, error)
}

// NewService - takes in a pointer to the DB & returns a pointer to a new document service
func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

// GetDocument - retrieves documents by their ID from the database
func (s *Service) GetDocument(ID uint) (Document, error) {
	var document Document
	// retrieve the first document in the DB with the passed in ID
	if result := s.DB.First(&document, ID); result.Error != nil {
		return Document{}, result.Error
	}
	// read the filename and return the contents (bytes).
	body, err := os.ReadFile(document.Path)
	if err != nil {
		log.Error("unable to read file")
		log.Error(err)
	}
	document.Body = string(body)

	return document, nil
}

// PostDocument - adds a new document to the database
func (s *Service) PostDocument(document Document) (Document, error) {
	if result := s.DB.Save(&document); result.Error != nil {
		return Document{}, result.Error
	}
	return document, nil
}

// UpdateDocument - updates a document by ID with new document info
func (s *Service) UpdateDocument(ID uint, newDocument Document) (Document, error) {
	document, err := s.GetDocument(ID)
	if err != nil {
		return Document{}, err
	}

	if result := s.DB.Model(&document).Updates(newDocument); result.Error != nil {
		return Document{}, result.Error
	}

	return document, nil
}

// DeleteDocument - deletes a document from the database by ID
func (s *Service) DeleteDocument(ID uint) error {
	if result := s.DB.Delete(&Document{}, ID); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAllDocuments() - retrieves all documents from the database
func (s *Service) GetAllDocuments() ([]Document, error) {
	var documents []Document
	if result := s.DB.Find(&documents); result.Error != nil {
		return documents, result.Error
	}
	return documents, nil
}
