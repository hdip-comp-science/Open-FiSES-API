package document

import (
	"os"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

//  Service - the struct for the document service
type Service struct {
	DB *gorm.DB
}

// Document - Defines the Document Model Structure
type Document struct {
	gorm.Model
	Path    string `json:"path"`
	Title   string `json:"title"`
	Version string `json:"version"`
	Author  string `json:"author"`
	Body    string `json:"body"`
}

// DocumentService - Defines the contract in against which you have to
// implement the document service
type DocumentService interface {
	GetDocument(ID uint) (Document, error)
	GetDocumentByPath(path string) ([]Document, error)
	PostDocument(document Document) (Document, error)
	UpdateDocument(ID uint, newDocument Document) (Comment error)
	DeleteDocument(ID uint) error
	GetAllDocuments() ([]Document, error)
}

// NewService - returns a pointer to a new document service
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
		glog.Errorf("unable to read file: %v", err)
	}
	document.Body = string(body)

	return document, nil
}

// GetDocumentByPath - retrieves all documents by path (path - /sop/name/ )
// func (s *Service) GetDocumentByPath(path string) ([]Document, error) {
// 	var documents []Document
// 	if result := s.DB.First(&documents).Where("path=?", path); result.Error != nil {
// 		return []Document{}, result.Error
// 	}
// 	return documents, nil
// }

// PostDocument - adds a new document to the database
func (s *Service) PostDocument(document Document) (Document, error) {
	if result := s.DB.Save(&document); result.Error != nil {
		return Document{}, result.Error
	}
	return document, nil
}

// UpdateDocument - updates a document by ID with new document info
func (s *Service) UpdateDocument(ID uint, newComment Document) (Document, error) {
	document, err := s.GetDocument(ID)
	if err != nil {
		return Document{}, err
	}

	if result := s.DB.Model(&document).Updates(newComment); result.Error != nil {
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
