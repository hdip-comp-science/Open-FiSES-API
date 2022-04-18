package booking

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BookService struct {
	DB *gorm.DB
}

// Job has 1-1 with booking, Job can consists of 0-* equipment
type Job struct {
	ID              uint
	SerialNo        string
	InstrumentModel string
	Manufacturer    string
	Customer        Customer `gorm:"foreignKey:ID"`
}

// Booking has one Job association and one Customer
type Booking struct {
	gorm.Model
	// bookingNo uint
	Summary       string
	Description   string
	Location      string
	StartDateTime time.Time
	EndDateTime   time.Time
	Customer      Customer `gorm:"foreignKey:ID"` //Embeded Type
	Job           Job      `gorm:"foreignKey:ID"` //Embeded Type
}

// Customer may have 0-* bookings
type Customer struct {
	ID      int
	Name    string
	Booking []Booking `gorm:"foreignKey:ID"`
}

// BookingService - the interface for our boooking service
type BookingService interface {
	GetBooking(ID uint) (Booking, error)
	PostBooking(booking Booking) (Booking, error)
	UpdateBooking(ID uint, newBooking Booking) (Booking error)
	DeleteBookings(ID uint) error
	GetAllBookings() ([]Booking, error)
}

// NewService - takes in a pointer to the DB & returns a pointer to a new booking service
func NewService(db *gorm.DB) *BookService {
	return &BookService{
		DB: db,
	}
}

// GetBooking - retrieves bookings by ID from the database
func (s *BookService) GetBooking(ID uint) (Booking, error) {
	var booking Booking // define a new booking variable
	// retireive the 1st booking from the DB with the passed in Id & populate the booking var with the result obj
	if result := s.DB.First(&booking, ID); result.Error != nil {
		return Booking{}, result.Error
	}
	return booking, nil
}

func (s *BookService) PostBooking(booking Booking) (Booking, error) {
	if result := s.DB.Save(&booking); result.Error != nil {
		return Booking{}, result.Error
	}
	return booking, nil
}

// UpdateDocument - updates a booking by ID with new document info
func (s *BookService) UpdateBooking(ID uint, newBooking Booking) (Booking, error) {
	booking, err := s.GetBooking(ID)
	if err != nil {
		return Booking{}, err
	}
	if result := s.DB.Model(&booking).Updates(newBooking); result.Error != nil {
		return Booking{}, result.Error
	}
	// return booking once it has been updated by gorm.
	return booking, nil
}

// DeleteBooking - deletes a booking from the database by ID
func (s *BookService) DeleteBooking(ID uint) error {
	// pass in empty comment obj and ID of booking to delete
	if result := s.DB.Delete(&Booking{}, ID); result.Error != nil {
		return result.Error
	}
	// if ID passed in is successfully deleted, return nil
	return nil
}

// GetAllBookings() - retrieves all bookings from the database
func (s *BookService) GetAllBookings() ([]Booking, error) {
	var bookings []Booking
	if result := s.DB.Find(&bookings); result.Error != nil {
		return bookings, result.Error
	}
	return bookings, nil
}
