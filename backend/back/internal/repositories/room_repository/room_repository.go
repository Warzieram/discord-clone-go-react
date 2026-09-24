package roomrepository

import (
	"back/internal/models/room"
	"context"
)

// RoomRepository describes persistence operations for rooms.
type RoomRepository interface {
	// Save persists a new room and returns its generated id.
	Save(ctx context.Context, r *room.Room) (int, error)
	// GetAll returns every room, oldest first.
	GetAll(ctx context.Context) ([]room.Room, error)
}
