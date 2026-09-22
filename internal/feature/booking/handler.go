package booking

import (
	"encoding/json"
	"net/http"

	"github.com/artemida000/go-coworking-booking/internal/core/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type CreateBookingRequest struct {
	RoomID    string `json:"room_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	bookingID, totalPrice, err := h.service.CreateBooking(r.Context(), userID, req.RoomID, req.StartTime, req.EndTime)
	if err != nil {
		if err.Error() == "room is already booked for this time" || 
		   err.Error() == "cannot book in the past" || 
		   err.Error() == "end_time must be after start_time" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "booking created successfully",
		"booking_id":  bookingID,
		"total_price": totalPrice,
	})
}

func (h *Handler) GetMyBookings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	bookings, err := h.service.GetMyBookings(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"bookings": bookings,
	})
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	bookingID := r.PathValue("id")

	err := h.service.CancelBooking(r.Context(), bookingID, userID)
	if err != nil {
		if err.Error() == "booking not found or cannot be cancelled" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "booking cancelled successfully",
	})
}