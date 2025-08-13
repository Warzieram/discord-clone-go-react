package rooms

import (
	"back/internal/database"
	"back/internal/models/user"
	"errors"
	"fmt"
	"log"
)

type Room struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatorID int    `json:"creator_id"`
	Deleted   bool   `json:"deleted"`
}

func GetRooms() ([]Room, error) {
	query := `SELECT id, name, creator_id FROM chatroom.rooms WHERE deleted=false`
	rows, err := database.DbInstance.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var retrievedRooms []Room

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name, &room.CreatorID); err != nil {
			log.Println("[ERROR] Couldn't map rooms into memory: ", err)
			return nil, err
		}
		retrievedRooms = append(retrievedRooms, room)
	}
	log.Println(retrievedRooms)

	return retrievedRooms, nil
}

func GetRoomByID(ID int) (*Room, error) {
	query := `SELECT id, name, creator_id
	FROM rooms WHERE id = $1`

	room := &Room{}
	err := database.DbInstance.DB.QueryRow(query, ID).Scan(
		&room.ID, &room.Name, &room.CreatorID,
	)
	if err != nil {
		log.Println("Error retrieving room in database: ", err)
		return nil, err
	}

	return room, nil

}

func CreateRoom(name string, creator_id int) (*Room, error) {
	room := &Room{Name: name, CreatorID: creator_id}

	if room.Name == "" {
		return nil, errors.New("empty room name")
	}
	_, err := user.GetUserById(creator_id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", err)
	}

	return room, nil
}

func (r *Room) Save() error {
	query := `INSERT into rooms (name, creator_id) VALUES ($1, $2) RETURNING id`
	err := database.DbInstance.DB.QueryRow(query, r.Name, r.CreatorID).Scan(&r.ID)
	if err != nil {
		return err
	}
	return nil
}

func MarkAsDeleted(id int) error {
	query := `UPDATE rooms SET deleted=true WHERE id=$1 `
	_, err := database.DbInstance.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func GetCreatorID(id int) (int, error) {
	var creatorID int
	query := `SELECT creator_id FROM rooms WHERE id=$1`
	err := database.DbInstance.DB.QueryRow(query, id).Scan(&creatorID)
	if err != nil {
		return 0, err
	}
	return creatorID, nil
}
