package rooms

import (
	"back/internal/database"
	"log"
)

type Room struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

func GetRooms() ([]Room, error) {
	query := `SELECT id, name FROM chatroom.rooms`
	rows, err := database.DbInstance.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var retrievedRooms []Room

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name); err != nil {
			log.Println("[ERROR] Couldn't map rooms into memory: ", err)
			return nil, err
		}
		retrievedRooms = append(retrievedRooms, room)
	}
	log.Println(retrievedRooms)

	return retrievedRooms, nil
}
