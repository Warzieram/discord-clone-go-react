package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type JWTInput struct {
	ID int
	Email string
	CreatedAt time.Time
}

func TestGenerateAndValidateJWT(t *testing.T) {
	testUser := &JWTInput{
		ID: 1,
		Email: "test@example.com",
		CreatedAt: time.Now(),
	}

	jwt, err :=  GenerateJWT(testUser.ID, testUser.Email, testUser.CreatedAt)
	if err != nil {
		t.Errorf("Couldn't generate jwt: %v", err)
		return
	}

	createdClaim, validateErr := ValidateJWT(jwt)
	if validateErr != nil {
		t.Errorf("Couldn't validate jwt: %v", validateErr)
		return
	}


	assert.Equal(t, createdClaim.UserID, testUser.ID)
	assert.Equal(t, createdClaim.Email, testUser.Email)
	assert.WithinDuration(t, createdClaim.CreatedAt, testUser.CreatedAt, 0)
}
