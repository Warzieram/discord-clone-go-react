package messagerepository

import (
	message "back/internal/models/messages"
	"context"
	"database/sql"
)

// PostgresMessageRepository is the Postgres-backed implementation of MessageRepository.
type PostgresMessageRepository struct {
	db *sql.DB
}

// NewPostgresMessageRepository builds a repository around an injected DB handle.
func NewPostgresMessageRepository(db *sql.DB) *PostgresMessageRepository {
	return &PostgresMessageRepository{db: db}
}

func (r *PostgresMessageRepository) Save(ctx context.Context, m *message.Message) (int, error) {
	const query = `INSERT INTO messages (content, sender_id, room_id) VALUES ($1, $2, $3) RETURNING id`

	id := 0
	if err := r.db.QueryRowContext(ctx, query, m.Content, m.SenderID, m.RoomID).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresMessageRepository) GetByID(ctx context.Context, id int) (*message.Message, error) {
	const query = `SELECT id, content, created_at, sender_id, COALESCE(room_id, 0) FROM messages WHERE id = $1`

	m := &message.Message{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Content, &m.CreatedAt, &m.SenderID, &m.RoomID,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *PostgresMessageRepository) GetLast(ctx context.Context, roomID int, limit int, offset int) ([]message.Message, error) {
	const query = `SELECT id, content, created_at, sender_id, room_id
	FROM messages
	WHERE deleted = false AND room_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []message.Message
	for rows.Next() {
		var m message.Message
		if err := rows.Scan(&m.Id, &m.Content, &m.CreatedAt, &m.SenderID, &m.RoomID); err != nil {
			return messages, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return messages, err
	}
	return messages, nil
}

func (r *PostgresMessageRepository) MarkAsDeleted(ctx context.Context, id int) error {
	const query = `UPDATE messages SET deleted = true WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
