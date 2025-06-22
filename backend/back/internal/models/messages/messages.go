package message

import (
	"back/internal/database"
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

func CreateMessage(content string, senderId int) (*Message, error) {
	message := &Message{Content: content, SenderID: senderId}
	
	if content == "" {
		return nil, errors.New("message content can't be null")
	}

	return message, nil
}

func (m Message) Save() error {
	query := `INSERT INTO messages (content, sender_id) VALUES ($1, $2)`
	_, err := database.DbInstance.DB.Exec(query, m.Content, m.SenderID)
	if err != nil {
		log.Println("ERROR while saving message: ", err)
	}
	return nil
}
