package room

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type CreateRoomRequest struct {
	SpaceID      string  `json:"space_id"`
	Name         string  `json:"name"`
	Capacity     int     `json:"capacity"`
	PricePerHour float64 `json:"price_per_hour"`
}

type CreateRoomResponse struct {
	Message string `json:"message"`
	RoomID  string `json:"room_id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	roomID, err := h.service.CreateRoom(r.Context(), req.SpaceID, req.Name, req.Capacity, req.PricePerHour)

	if err != nil {
		if err.Error() == "name, capacity and pricePerHour are required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
	}

	http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := CreateRoomResponse {
		Message: "room created successfully",
		RoomID: roomID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

type GetAllRoomsResponse struct {
	Rooms []Room `json:"rooms"`
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	spaceID := r.PathValue("id")

	rooms, err := h.service.GetAllRooms(r.Context(), spaceID)
	if err != nil {
		if err.Error() == "space_id is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := GetAllRoomsResponse{
		Rooms: rooms,
	}

	// Отправляем JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}