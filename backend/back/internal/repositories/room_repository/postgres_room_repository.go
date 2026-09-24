package roomrepository

import (
	"back/internal/models/room"
	"context"
	"database/sql"
)

// PostgresRoomRepository is the Postgres-backed implementation of RoomRepository.
type PostgresRoomRepository struct {
	db *sql.DB
}

// NewPostgresRoomRepository builds a repository around an injected DB handle.
func NewPostgresRoomRepository(db *sql.DB) *PostgresRoomRepository {
	return &PostgresRoomRepository{db: db}
}

func (r *PostgresRoomRepository) Save(ctx context.Context, rm *room.Room) (int, error) {
	const query = `INSERT INTO rooms (name) VALUES ($1) RETURNING id`

	id := 0
	if err := r.db.QueryRowContext(ctx, query, rm.Name).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresRoomRepository) GetAll(ctx context.Context) ([]room.Room, error) {
	const query = `SELECT id, name FROM rooms ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := []room.Room{}
	for rows.Next() {
		var rm room.Room
		if err := rows.Scan(&rm.Id, &rm.Name); err != nil {
			return rooms, err
		}
		rooms = append(rooms, rm)
	}
	if err := rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}
