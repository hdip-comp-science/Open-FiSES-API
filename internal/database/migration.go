package database

import (
	"github.com/Open-FiSE/go-rest-api/internal/document"
	"github.com/jinzhu/gorm"
)

// MigrateDB - migrates the database and creates the doucment table
func MigrateDB(db *gorm.DB) error {
	// AutoMigrate - takes in document model (struct)
	if result := db.AutoMigrate(&document.Document{}); result.Error != nil {
		return result.Error
	}
	return nil
}
