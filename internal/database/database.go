package database

//  Use gorm effectively communicate with the DB.
// Perform actions such as migration and pinging the DB

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	// wrapper around DB driver to open up communications to the database
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewDatabase - returns a pointer to the postgres database object
func NewDatabase() (*gorm.DB, error) {
	log.Info("Setting up database connection")

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbTable := os.Getenv("DB_TABLE")
	dbPort := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUsername, dbTable, dbPassword)

	// open up DB connection
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return db, fmt.Errorf("unable to connect to database: %v", err)
	}
	// Ping the DB using the credentials obtained via env variables.
	if err := db.DB().Ping(); err != nil {
		return db, err
	}
	log.Info("Database connection established")
	return db, nil
}
