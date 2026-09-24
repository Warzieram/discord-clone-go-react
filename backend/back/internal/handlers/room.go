package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"back/internal/models/room"
	roomrepository "back/internal/repositories/room_repository"
)

type RoomHandlers struct {
	repo roomrepository.RoomRepository
}

func NewRoomHandlers(repo roomrepository.RoomRepository) *RoomHandlers {
	return &RoomHandlers{repo: repo}
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

// CreateRoom handles POST /api/rooms with a {"name": "..."} body
func (h *RoomHandlers) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	rm, err := room.CreateRoom(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repo.Save(r.Context(), rm)
	if err != nil {
		log.Println("[ERROR] Couldn't save room: ", err)
		http.Error(w, "couldn't create room", http.StatusInternalServerError)
		return
	}
	rm.Id = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(rm); err != nil {
		log.Println("[ERROR] Couldn't encode room: ", err)
	}
}

// ListRooms handles GET /api/rooms
func (h *RoomHandlers) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Println("[ERROR] Couldn't retrieve rooms: ", err)
		http.Error(w, "couldn't retrieve rooms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		log.Println("[ERROR] Couldn't encode rooms: ", err)
	}
}
