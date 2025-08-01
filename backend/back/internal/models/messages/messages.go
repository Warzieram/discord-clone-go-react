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
}

type MessageResponse struct {
	Id       int       `json:"id"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"created_at"`
	Sender   string    `json:"sender"`
}

func CreateMessage(content string, senderId int) (*Message, error) {
	message := &Message{Content: content, SenderID: senderId}

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
	query := `INSERT INTO messages (content, sender_id) VALUES ($1, $2) RETURNING id`
	err := database.DbInstance.DB.QueryRow(query, m.Content, m.SenderID).Scan(&id)
	if err != nil {
		log.Println("ERROR while saving message: ", err)
		return 0, err
	}

	log.Println("ID: ", id)

	return id, nil
}

func GetLastMessages(limit int, offset int) ([]MessageResponse, error) {
	query := `SELECT id, content, created_at, sender_id FROM messages WHERE deleted = false ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := database.DbInstance.DB.Query(query, limit, offset)
	if err != nil {
		log.Printf("[ERROR] Couldn't get last %v messages: %s", limit, err)
		return nil, err
	}

	defer rows.Close()

	var messages []MessageResponse

	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.Id, &message.Content, &message.CreatedAt, &message.SenderID); err != nil {
			log.Println("[ERROR] Couldn't map messages into memory: ", err)
			return messages, err
		}
		sendFormat, err := message.ToSendFormat()
		if err != nil {
			log.Println("[ERROR] Couldn't convert message to send format: ", err)
			return messages, err
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
	query := `SELECT id, content, created_at, sender_id 
	FROM messages WHERE id = $1`

	message := &Message{}
	err := database.DbInstance.DB.QueryRow(query, id).Scan(
		&message.Id, &message.Content, &message.CreatedAt, &message.SenderID,
	)
	if err != nil {
		log.Println("Error creating message in database: ", err)
		return nil, err
	}

	return message, nil
}
