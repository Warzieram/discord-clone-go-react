package room

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// MAX_NAME_LENGTH matches the rooms.name VARCHAR(20) column
const MAX_NAME_LENGTH = 20

type Room struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	CreatorID int `json:"creator_id"`
}

func CreateRoom(name string, creator_id int) (*Room, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("room name can't be empty")
	}
	if utf8.RuneCountInString(name) > MAX_NAME_LENGTH {
		return nil, errors.New("room name is too long")
	}

	if creator_id < 1{
		return nil, errors.New("Invalid user Id for creator")
	}

	return &Room{Name: name, CreatorID: creator_id}, nil
}
