package handlers

import (
	"back/internal/models/messages"
	"log"
	"net/http"
	"strconv"
)

func RetrieveMessages(w http.ResponseWriter, r *http.Request)  {

	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	limit, limitErr := strconv.Atoi(limitParam)
	if limitErr != nil {
		log.Println("[ERROR] while retrieving limit", limitErr)
		return
	}

	offset, offsetErr := strconv.Atoi(offsetParam)
	if offsetErr != nil {
		log.Println("[ERROR] while retrieving offset: ", offsetErr)
		return
	}



	retrievedMessages, err := message.GetLastMessages(limit, offset)
	if err != nil {
		log.Println("[ERROR] Couldn't retrieve messages: ", err)
	}



	
}
