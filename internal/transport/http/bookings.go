package http

// Define endpoints and map them to the booking service.
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Open-FiSE/go-rest-api/internal/booking"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	// retrieve the ID of the booking you want to fetch
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}
	// GetBooking is expecting a uint, so parse string to uint.
	booking, err := h.BookService.GetBooking(uint(i))
	if err != nil {
		fmt.Fprintf(w, "Error retrieving Booking by ID")
	}

	// Return the newly update booking as json
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Warning(err)
	}

}

// GetAllBookings - fetch all bookings from the booking service
func (h *Handler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	bookings, err := h.BookService.GetAllBookings()
	if err != nil {
		fmt.Fprintf(w, "Failed to retrieve bookings")
	}
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		log.Warning(err)
	}
}

// UpdateBooking - update an exisiting booking by ID
func (h *Handler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	var booking booking.Booking
	// Parse the request body as booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		fmt.Fprintf(w, "Failed to decode JSON Body")
	}

	vars := mux.Vars(r)
	id := vars["id"]

	bookingID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	booking, err = h.BookService.UpdateBooking(uint(bookingID), booking)
	if err != nil {
		fmt.Fprintf(w, "Failed to update booking")
	}

	// Return the newly update booking as json
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Warning(err)
	}
}

// DeleteBooking - delete a booking by ID
func (h *Handler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	id := vars["id"]

	bookingID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	err = h.BookService.DeleteBooking(uint(bookingID))
	if err != nil {
		fmt.Fprintf(w, "Failed to delete booking")
	}

	if err := json.NewEncoder(w).Encode(Response{Message: "Successfully deleted booking"}); err != nil {
		log.Warning(err)
	}
}

// PostBooking - adds a new booking
func (h *Handler) PostBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charcet=UTF-8")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	var booking booking.Booking
	// Parse the request body as booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		fmt.Fprintf(w, "Failed to decode JSON Body")
	}
	// Post to booking service
	booking, err := h.BookService.PostBooking(booking)
	if err != nil {
		fmt.Fprintf(w, "Failed to post new booking")
	}
	// return the booking
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Warning(err)
	}
}
