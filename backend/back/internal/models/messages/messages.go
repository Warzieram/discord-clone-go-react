package message

import (
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
