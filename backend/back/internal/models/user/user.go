package user

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"back/internal/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                    int            `json:"id"`
	Email                 string         `json:"email"`
	Username              string         `json:"username"`
	PasswordHash          string         `json:"-"`
	CreatedAt             time.Time      `json:"created_at"`
	EmailVerified         bool           `json:"email_verified"`
	VerificationToken     sql.NullString `json:"verification_token"`
	VerificationExpiresAt sql.NullTime   `json:"verification_expires_at"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (u *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	return nil

}

func (u *User) CheckPassword(password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) CreateAccountVerificationToken() error {
	token := uuid.New()
	expirationTime := time.Now().Add(24 * time.Hour)
	u.VerificationToken = sql.NullString{String: token.String(), Valid: true}
	u.VerificationExpiresAt = sql.NullTime{Time: expirationTime, Valid: true}
	query := `
	UPDATE users 
	SET verification_token = $1, verification_expires_at = $2 
	WHERE email = $3 `
	_, err := database.DbInstance.DB.Exec(query, token, expirationTime, u.Email)
	if err != nil {
		log.Println("ERROR: could'create verification token in db:", err)
		return err
	}
	return nil
}

func (u *User) Save() error {
	query := `INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at`
	_, err := database.DbInstance.DB.Exec(query, u.Email, u.Username, u.PasswordHash)
	if err != nil {
		return err
	}
	err = u.CreateAccountVerificationToken()
	if err != nil {
		log.Println("Couldn't create verification token", err)
		return err
	}
	SendCreationEmail(u)

	return nil
}

func (u *User) ResendVerification() error {
	err := u.CreateAccountVerificationToken()
	if err != nil {
		return err
	}

	return nil
}

func (u *User) VerifyEmail() error {
	if u.EmailVerified {
		return errors.New("email already verified")
	}

	// if expiration date is already passed
	if u.VerificationExpiresAt.Time.Before(time.Now()) {
		return errors.New("verification token expired")
	}
	query := `UPDATE users SET email_verified = true WHERE id = $1`

	_, err := database.DbInstance.DB.Exec(query, u.ID)
	if err != nil {
		return err
	}
	u.EmailVerified = true
	return nil
}

func CreateUser(email string, username string, password string) (*User, error) {
	user := &User{Email: email, Username: username}

	if err := user.HashPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserById(id int) (*User, error) {
	user := &User{}
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

func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, email, username, password_hash, created_at, email_verified FROM users WHERE email=$1`

	err := database.DbInstance.DB.QueryRow(query, email).Scan(
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

func GetUserByVerificationToken(token string) (*User, error) {
	user := &User{}
	query := `SELECT id, email, username, created_at, verification_expires_at FROM users WHERE verification_token=$1`

	err := database.DbInstance.DB.QueryRow(query, token).Scan(
		&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.VerificationExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
