package database

import (
	"github.com/Open-FiSE/go-rest-api/internal/booking"
	"github.com/Open-FiSE/go-rest-api/internal/document"
	"github.com/jinzhu/gorm"
)

// MigrateDB - will create and migrate the tables, and then creates relationships if necessary
func MigrateDB(db *gorm.DB) error {
	// AutoMigrate - takes in document model (struct) &
	// define DB columns Path | Body | Author as well as predefined gorm (ID, update time etc).
	if result := db.AutoMigrate(&document.Document{}, &booking.Booking{}); result.Error != nil {
		return result.Error
	}
	return nil
}
