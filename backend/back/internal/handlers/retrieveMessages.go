package handlers

import (
	"back/internal/models/messages"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func RetrieveMessages(w http.ResponseWriter, r *http.Request) {

	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	roomIDParam := r.URL.Query().Get("room")


	roomID, roomIDErr := strconv.Atoi(roomIDParam)
	if roomIDErr != nil {
		log.Println("[ERROR] while retrieving room_ID: ", roomIDErr)
		http.Error(w, "error parsing room parameter", http.StatusBadRequest)
		return
	}
	log.Printf("Retriving messages in room %v ", roomID)

	limit, limitErr := strconv.Atoi(limitParam)
	if limitErr != nil {
		log.Println("[ERROR] while retrieving limit: ", limitErr)
		http.Error(w, "error parsing limit parameter", http.StatusBadRequest)
		return
	}

	offset, offsetErr := strconv.Atoi(offsetParam)
	if offsetErr != nil {
		log.Println("[ERROR] while retrieving offset: ", offsetErr)
		http.Error(w, "error parsing offset parameter", http.StatusBadRequest)
		return
	}

	retrievedMessages, err := message.GetLastMessages(roomID, limit, offset )
	if err != nil {
		log.Println("[ERROR] Couldn't retrieve messages: ", err)
		http.Error(w, "error while retrieving messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(retrievedMessages)
	if err != nil {
		log.Println("[ERROR] Couldn't encode last messages: ", err)
		return
	}

	//log.Println("RETRIEVED MESSAGES: ", retrievedMessages)

}
