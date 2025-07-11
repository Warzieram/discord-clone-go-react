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
		Content:  m.Content,
		CreateAt: m.CreatedAt,
		Sender:   sender.Email,
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
