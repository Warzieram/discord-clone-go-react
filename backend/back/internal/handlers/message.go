package handlers

// TODO: Order this hot mess

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"back/internal/models/messages"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Println("WebSocket origin: ", origin)
		// return strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1")
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var mutex = &sync.Mutex{}

const IDLE_TIMEOUT = 1800

type CommandType string

const (
	SEND   CommandType = "SEND"
	REMOVE CommandType = "REMOVE"
)

type Request interface {
	GetType() CommandType
	Execute(userID int) error
}

type SendRequest struct {
	Data string
}

func (s SendRequest) GetType() CommandType {
	return SEND
}

func (s SendRequest) Execute(userID int) error {
	log.Println("Executing request: ", s)
	m, err := message.CreateMessage(s.Data, userID)
	if err != nil {
		return err
	}

	log.Println("Created message: ", m)

	id, errSave := m.Save()
	if errSave != nil {
		return errSave
	}
	log.Println("Saved message ID: ", id)

	retrievedMessage, errRetrieved := message.GetMessageById(id)
	if errRetrieved != nil {
		return errRetrieved
	}
	log.Println("Retrieved message: ", retrievedMessage)

	response, err := retrievedMessage.ToSendFormat()
	if err != nil {
		log.Println("ERROR converting to send format: ", err)
		return err
	}

	output, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		log.Println("ERROR converting message to json: ", err)
		return jsonErr
	}

	broadcast <- output

	return nil
}

type RemoveRequest struct {
	Data int
}

func (r RemoveRequest) GetType() CommandType {
	return REMOVE
}

func (r RemoveRequest) Execute(userID int) error {
	log.Println("Executing request: ", r)

	err := message.MarkAsDeleted(r.Data)
	return err
}

func parseReq(s string) (Request, error) {

	log.Println("PARSING request: ", s)

	trimedString := strings.Trim(s, " ")
	splitedString := strings.Split(trimedString, " ")

	cmd := splitedString[0]
	data := strings.Join(splitedString[1:], " ")

	log.Println("COMMAND: ", cmd)
	log.Println("DATA: ", data)

	if cmd == string(SEND) {
		return SendRequest{Data: data}, nil
	} else if cmd == string(REMOVE) {
		message_id, err := strconv.Atoi(data)
		if err != nil {
			return nil, err
		}
		return RemoveRequest{Data: message_id}, nil
	}
	err := errors.New("invalid command")
	return nil, err
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {

	// upgrade the connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR: Upgrade failed: ", err)
		http.Error(w, "Couln't initiate websocket connection", http.StatusInternalServerError)
		return
	}

	// set the connection in the list of active clients
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()
	}()

	log.Println("Initiated websocket connection with: ", r.RemoteAddr)

	for {
		messageType, message, err := conn.ReadMessage()
		conn.SetReadDeadline(time.Now().Add(IDLE_TIMEOUT * time.Second))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(IDLE_TIMEOUT * time.Second))
			return nil
		})
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Println("WebSocket error: ", err)
			} else {
				log.Println("WebSocket closed normally", err)
			}
			break
		}
		log.Println("Received message from : ", r.RemoteAddr)
		log.Println("Message type: ", messageType)

		input := string(message)
		log.Println("Message: ", input)

		request, err := parseReq(input)
		if err != nil {
			log.Println("Error while parsing the request: ", err)
		}

		userID := r.Context().Value("user_id").(int)
		log.Println("User ID: ", userID)

		err = request.Execute(userID)
		if err != nil {
			log.Println("ERROR while executing message request: ", err)
			continue
		}

	}

}

func SendMessage() {
	for {
		message := <-broadcast

		log.Println("sending message")

		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(1, message)
			if err != nil {
				log.Println("ERROR writing mesage :", err)
				delete(clients, client)
			}

		}
		mutex.Unlock()

	}
}
