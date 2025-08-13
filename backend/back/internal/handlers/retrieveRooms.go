package handlers

import (
	"back/internal/models/rooms"
	"encoding/json"
	"log"
	"net/http"
)

func RetrieveRooms(w http.ResponseWriter, r *http.Request) {
	retrievedRooms, err := rooms.GetRooms();
	if err != nil {
		log.Println("[ERROR] Retrieving rooms: ", err)
		http.Error(w, "Couldn't retrieve rooms please try later", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(retrievedRooms)
	if err != nil {
		log.Println("[ERROR] Ecoding retrived rooms: ", err)
		return
	}
}
