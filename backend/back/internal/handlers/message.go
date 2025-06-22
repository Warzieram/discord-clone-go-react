package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

func MessageHandler(w http.ResponseWriter, r *http.Request) {

	// upgrade the connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR: Upgrade failed: ", err)
		http.Error(w, "Couln't initiate websocket connection", http.StatusInternalServerError)
		return
	}

	// set the connection in the list of active clients 
	clients[conn] = true

	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("ERROR: Couldn't read websocket message: ", err)
			http.Error(w, "error reading message", http.StatusInternalServerError)
			break
		}
		log.Println("Received message from : ", r.RemoteAddr)
		log.Println("Message type: ", messageType)

		input := string(message)
		log.Println("Message: ", input)



		for client := range clients {
			err = client.WriteMessage(messageType, []byte(input))
			if err != nil {
				log.Println("ERROR writing mesage :", err)
			}

		}

	}

}
