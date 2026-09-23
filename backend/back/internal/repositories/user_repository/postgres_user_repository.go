package userrepository

import (
	"back/internal/database"
	"back/internal/models/user"
	"context"
	"database/sql"
	"errors"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db}
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*user.User, error) {

	user := &user.User{}
	query := `SELECT id, email, username, password_hash, created_at, email_verified FROM users WHERE id=$1`

	err := database.DbInstance.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.EmailVerified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil

}
