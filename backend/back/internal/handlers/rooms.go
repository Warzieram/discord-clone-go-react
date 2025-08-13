package handlers

import (
	"back/internal/models/rooms"
	"encoding/json"
	"log"
	"net/http"
)

type CreateRoomRequestBody struct {
	Name string `json:"name"`
}

type DeleteRoomRequestBody struct {
	ID int `json:"id"`
}

func RetrieveRooms(w http.ResponseWriter, r *http.Request) {
	retrievedRooms, err := rooms.GetRooms()
	if err != nil {
		log.Println("[ERROR] Retrieving rooms: ", err)
		http.Error(w, "Couldn't retrieve rooms please try later", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(retrievedRooms)
	if err != nil {
		log.Println("[ERROR] Ecoding retrieved rooms: ", err)
		return
	}
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	body := &CreateRoomRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("[ERROR] Couldn't decode create room request body: ", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)
	room, err := rooms.CreateRoom(body.Name, userID)
	if err != nil {
		log.Println("[ERROR] Couldn't create room: ", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = room.Save()
	if err != nil {
		log.Println("[ERROR] Couldn't save room: ", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(room)
	if err != nil {
		log.Println("[ERROR] Ecoding retrived rooms: ", err)
		return
	}

}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	body := &DeleteRoomRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("[ERROR] Couldn't decode delete room request body: ", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	room, err := rooms.GetRoomByID(body.ID)
	if err != nil {
		log.Println("[ERROR] Couldn't get room: ", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value("user_id").(int)
	if room.CreatorID != userID {
		http.Error(w, "Only room creator can delete room", http.StatusUnauthorized)
		return
	}

	err = rooms.MarkAsDeleted(room.ID)
	if err != nil {
		log.Println("[ERROR] Couldn't mark room as deleted: ", err)
		http.Error(w, "Something went wrong while deleting the room", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
