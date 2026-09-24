package handlers

import (
	"back/internal/models/messages"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (h *MessageHandlers) RetrieveMessages(w http.ResponseWriter, r *http.Request) {

	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

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

	roomID, roomErr := strconv.Atoi(r.URL.Query().Get("roomID"))
	if roomErr != nil {
		log.Println("[ERROR] while retrieving roomID: ", roomErr)
		http.Error(w, "error parsing roomID parameter", http.StatusBadRequest)
		return
	}

	messages, err := h.repo.GetLast(r.Context(), roomID, limit, offset)
	if err != nil {
		log.Println("[ERROR] Couldn't retrieve messages: ", err)
		http.Error(w, "couldn't retrieve messages", http.StatusInternalServerError)
		return
	}

	// Convert each persisted message into its API send format (resolves sender).
	retrievedMessages := make([]message.MessageResponse, 0, len(messages))
	for i := range messages {
		sendFormat, err := messages[i].ToSendFormat()
		if err != nil {
			log.Println("[ERROR] Couldn't convert message to send format: ", err)
			http.Error(w, "couldn't format messages", http.StatusInternalServerError)
			return
		}
		retrievedMessages = append(retrievedMessages, *sendFormat)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(retrievedMessages)
	if err != nil {
		log.Println("[ERROR] Couldn't encode last messages: ", err)
		return
	}

	//log.Println("RETRIEVED MESSAGES: ", retrievedMessages)

}
