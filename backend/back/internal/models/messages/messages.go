package message

import (
	"back/internal/models/user"
	"errors"
	"log"
	"time"
	"unicode/utf8"
)

// MAX_CONTENT_LENGTH matches the messages.content VARCHAR(500) column
const MAX_CONTENT_LENGTH = 500

type Message struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	SenderID  int       `json:"sender_id"`
	RoomID    int       `json:"room_id"`
	Edited    bool      `json:"edited"`
	Deleted bool `json:"-"`
}

type MessageResponse struct {
	Id       int       `json:"id"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"created_at"`
	Sender   string    `json:"sender"`
	RoomID   int       `json:"room_id"`
	Edited   bool      `json:"edited"`
}

func CreateMessage(content string, senderId int, roomId int) (*Message, error) {
	message := &Message{Content: content, SenderID: senderId, RoomID: roomId}

	if content == "" {
		return nil, errors.New("message content can't be null")
	}

	if utf8.RuneCountInString(content) > MAX_CONTENT_LENGTH {
		return nil, errors.New("message content is too long")
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
		RoomID:   m.RoomID,
		Edited:   m.Edited,
	}

	return response, nil

}
