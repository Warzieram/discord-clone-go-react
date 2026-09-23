package messagerepository

import (
	message "back/internal/models/messages"
	"context"
)

// MessageRepository describes persistence operations for messages. Handlers and
// services depend on this interface rather than a concrete implementation,
// which keeps the data access swappable and unit-testable (e.g. with a fake).
type MessageRepository interface {
	// Save persists a new message and returns its generated id.
	Save(ctx context.Context, m *message.Message) (int, error)
	// GetByID returns a single message by its id.
	GetByID(ctx context.Context, id int) (*message.Message, error)
	// GetLast returns up to `limit` non-deleted messages, newest first,
	// skipping `offset` rows.
	GetLast(ctx context.Context, limit int, offset int) ([]message.Message, error)
	// MarkAsDeleted soft-deletes a message.
	MarkAsDeleted(ctx context.Context, id int) error
}
