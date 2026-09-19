package space

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

type CreateSpaceRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

type CreateSpaceResponse struct {
	Message string `json:"message"`
	SpaceID string `json:"space_id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateSpaceRequest 

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	spaceID, err := h.service.CreateSpace(r.Context(), req.Name, req.Address, req.Description)

	if err != nil {
		if err.Error() == "name and address are required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := CreateSpaceResponse {
		Message: "space created successfully",
		SpaceID: spaceID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

type GetAllsSpaceResponce struct {
	Spaces []Space `json:"spaces"`
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	spaces, err := h.service.GetAllSpace(r.Context())

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := GetAllsSpaceResponce{
		Spaces: spaces,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode responce", http.StatusInternalServerError)
	}
}