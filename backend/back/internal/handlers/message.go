package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"back/internal/models/messages"
	messagerepository "back/internal/repositories/message_repository"

	"github.com/gorilla/websocket"
)

type MessageHandlers struct {
	repo messagerepository.MessageRepository
}

func NewMessageHandlers(repo messagerepository.MessageRepository) *MessageHandlers {
	return &MessageHandlers{repo: repo}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Println("WebSocket origin: ", origin)
		// return strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1")
		return true
	},
}

// clients maps each active connection to the room it is connected to
var clients = make(map[*websocket.Conn]int)
var broadcast = make(chan roomBroadcast)
var mutex = &sync.Mutex{}

const IDLE_TIMEOUT = 1800

type CommandType string

const (
	SEND   CommandType = "SEND"
	REMOVE CommandType = "REMOVE"
)

type Command struct {
	Type CommandType `json:"command_type"`
}

type Request interface {
	GetType() CommandType
	Execute(ctx context.Context, repo messagerepository.MessageRepository, userID int, roomID int) error
}

type SendRequest struct {
	Data string `json:"data"`
}

func (s SendRequest) GetType() CommandType {
	return SEND
}

type RemoveRequest struct {
	Data int `json:"data"`
}

func (r RemoveRequest) GetType() CommandType {
	return REMOVE
}

type BroadcastData interface {
	int | message.MessageResponse
}

type Broadcast[T BroadcastData] struct {
	Type CommandType `json:"command_type"`
	Data T           `json:"data"`
}

// roomBroadcast is a payload to send to every client connected to RoomID
type roomBroadcast struct {
	RoomID  int
	Payload []byte
}

func (s SendRequest) Execute(ctx context.Context, repo messagerepository.MessageRepository, userID int, roomID int) error {
	log.Println("Executing request: ", s)
	m, err := message.CreateMessage(s.Data, userID, roomID)
	if err != nil {
		return err
	}

	log.Println("Created message: ", m)

	id, errSave := repo.Save(ctx, m)
	if errSave != nil {
		return errSave
	}
	log.Println("Saved message ID: ", id)

	retrievedMessage, errRetrieved := repo.GetByID(ctx, id)
	if errRetrieved != nil {
		return errRetrieved
	}
	log.Println("Retrieved message: ", retrievedMessage)

	response, err := retrievedMessage.ToSendFormat()
	if err != nil {
		log.Println("ERROR converting to send format: ", err)
		return err
	}

	b := Broadcast[message.MessageResponse]{
		"SEND",
		*response,
	}

	output, jsonErr := json.Marshal(b)
	if jsonErr != nil {
		log.Println("ERROR converting message to json: ", err)
		return jsonErr
	}

	broadcast <- roomBroadcast{roomID, output}

	return nil
}

func (r RemoveRequest) Execute(ctx context.Context, repo messagerepository.MessageRepository, userID int, roomID int) error {
	log.Println("Executing request: ", r)

	m, err := repo.GetByID(ctx, r.Data)
	if err != nil {
		return err
	}

	//Proprietary check
	if m.SenderID == userID && m.RoomID == roomID {
		err = repo.MarkAsDeleted(ctx, r.Data)
		if err != nil {
			return err
		}
		b := Broadcast[int]{
			"REMOVE",
			r.Data,
		}
		output, err := json.Marshal(b)
		if err != nil {
			return err
		}

		broadcast <- roomBroadcast{roomID, output}
	}

	return nil
}

func parseReq(s string) (Request, error) {

	log.Println("PARSING request: ", s)
	c := &Command{}

	err := json.Unmarshal([]byte(s), c)
	if err != nil {
		return nil, err
	}
	log.Println(string(c.Type))

	switch string(c.Type) {
	case string(SEND):
		req := &SendRequest{}
		err := json.Unmarshal([]byte(s), req)
		if err != nil {
			return nil, err
		}
		return req, nil
	case string(REMOVE):
		req := &RemoveRequest{}
		err := json.Unmarshal([]byte(s), req)
		if err != nil {
			return nil, err
		}
		return req, nil
	default:
		return nil, errors.New("unknown request type")
	}
}

func (h *MessageHandlers) MessageHandler(w http.ResponseWriter, r *http.Request) {

	// Echo back the client's requested subprotocol (e.g. "auth.<token>").
	// Browsers fail the handshake if a subprotocol was offered but the server
	// doesn't select one in the Sec-WebSocket-Protocol response header.
	// We create a separate ws connection for each room
	// Route: /message?roomID={roomID}

	roomID, err := strconv.Atoi(r.URL.Query().Get("roomID"))
	if err != nil {
		http.Error(w, "missing or invalid room ID", http.StatusBadRequest)
		return
	}

	var respHeader http.Header
	if proto := r.Header.Get("Sec-WebSocket-Protocol"); proto != "" {
		respHeader = http.Header{"Sec-WebSocket-Protocol": {proto}}
	}

	// upgrade the connection to websocket
	conn, err := upgrader.Upgrade(w, r, respHeader)
	if err != nil {
		log.Println("ERROR: Upgrade failed: ", err)
		http.Error(w, "Couln't initiate websocket connection", http.StatusInternalServerError)
		return
	}

	// set the connection in the list of active clients
	mutex.Lock()
	clients[conn] = roomID
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
			continue
		}

		userID := r.Context().Value("user_id").(int)
		log.Println("User ID: ", userID)

		err = request.Execute(r.Context(), h.repo, userID, roomID)
		if err != nil {
			log.Println("ERROR while executing message request: ", err)
			continue
		}

	}

}

func SendMessage() {
	for {
		b := <-broadcast

		log.Println("sending message to room ", b.RoomID)

		mutex.Lock()
		for client, roomID := range clients {
			if roomID != b.RoomID {
				continue
			}
			err := client.WriteMessage(1, b.Payload)
			if err != nil {
				log.Println("ERROR writing mesage :", err)
				delete(clients, client)
			}

		}
		mutex.Unlock()

	}
}
