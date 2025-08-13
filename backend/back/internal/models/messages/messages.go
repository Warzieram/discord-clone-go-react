package message

import (
	"back/internal/database"
	"back/internal/models/user"
	"errors"
	"log"
	"time"
)

type Message struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	SenderID  int       `json:"sender_id"`
	RoomID    int       `json:"room_id"`
}

type MessageResponse struct {
	Id       int       `json:"id"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"created_at"`
	Sender   string    `json:"sender"`
}

func CreateMessage(content string, senderId int, roomId int) (*Message, error) {
	message := &Message{Content: content, SenderID: senderId, RoomID: roomId}

	if content == "" {
		return nil, errors.New("message content can't be null")
	}

	return message, nil
}

func (m *Message) ToSendFormat() (*MessageResponse, error) {
	sender, err := user.GetUserById(m.SenderID)
	if err != nil {
		log.Println("ERROR getting sender: ", err)
		return nil, err
	}

	response := &MessageResponse{
		Id:       m.Id,
		Content:  m.Content,
		CreateAt: m.CreatedAt,
		Sender:   sender.Username,
	}

	return response, nil

}

func (m Message) Save() (int, error) {

	log.Printf("Saving message with content %s ans senderID %v", m.Content, m.SenderID)

	id := 0
	query := `INSERT INTO messages (content, sender_id, room_id) VALUES ($1, $2, $3) RETURNING id`
	err := database.DbInstance.DB.QueryRow(query, m.Content, m.SenderID, m.RoomID).Scan(&id)
	if err != nil {
		log.Println("ERROR while saving message: ", err)
		return 0, err
	}

	return id, nil
}

func GetLastMessages(room_id int, limit int, offset int) ([]MessageResponse, error) {
	query := `SELECT id, content, created_at, sender_id, room_id FROM messages WHERE deleted = false AND room_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := database.DbInstance.DB.Query(query, room_id, limit, offset)
	if err != nil {
		log.Printf("[ERROR] Couldn't get last %v messages: %s", limit, err)
		return nil, err
	}

	defer rows.Close()

	var messages []MessageResponse

	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.Id, &message.Content, &message.CreatedAt, &message.SenderID, &message.RoomID); err != nil {
			log.Println("[ERROR] Couldn't map messages into memory: ", err)
			return messages, err
		}
		sendFormat, err := message.ToSendFormat()
		if err != nil {
			log.Println("[ERROR] Couldn't convert message to send format: ", err)
			return nil, err
		}
		messages = append(messages, *sendFormat)
	}
	log.Println(messages)

	if err = rows.Err(); err != nil {
		return messages, err
	}

	return messages, nil

}

func MarkAsDeleted(id int) error {
	query := `UPDATE messages SET deleted=true WHERE id=$1 `

	_, err := database.DbInstance.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func GetMessageById(id int) (*Message, error) {
	query := `SELECT id, content, created_at, sender_id, room_id
	FROM messages WHERE id = $1`

	message := &Message{}
	err := database.DbInstance.DB.QueryRow(query, id).Scan(
		&message.Id, &message.Content, &message.CreatedAt, &message.SenderID, &message.RoomID,
	)
	if err != nil {
		log.Println("Error retrieving message in database: ", err)
		return nil, err
	}

	return message, nil
}
