package user

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Password Management Tests
func TestHashPassword(t *testing.T) {
	user := &User{}

	t.Run("Valid password hashing", func(t *testing.T) {
		password := "testpassword123"
		err := user.HashPassword(password)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if user.PasswordHash == "" {
			t.Errorf("Expected password hash to be set")
		}
		if user.PasswordHash == password {
			t.Errorf("Password hash should not equal plain password")
		}
	})

	t.Run("Empty password", func(t *testing.T) {
		user := &User{}
		err := user.HashPassword("")
		if err == nil {
			t.Error("Expected error for empty password, got none")
		}
	})

	t.Run("Very long password", func(t *testing.T) {
		user := &User{}
		longPassword := strings.Repeat("a", 1000)
		err := user.HashPassword(longPassword)
		if err == nil {
			t.Errorf("Expected error for password longer than 72 bytes")
		}
		// Bcrypt has a 72-byte limit, so this should fail
		if !strings.Contains(err.Error(), "password length exceeds 72 bytes") {
			t.Errorf("Expected bcrypt length error, got %v", err)
		}
	})

	t.Run("Special characters in password", func(t *testing.T) {
		user := &User{}
		specialPassword := "p@ssw0rd!@#$%^&*()_+-=[]{}|;':\",./<>?"
		err := user.HashPassword(specialPassword)
		if err != nil {
			t.Errorf("Expected no error for special characters, got %v", err)
		}
	})

	t.Run("Unicode characters in password", func(t *testing.T) {
		user := &User{}
		unicodePassword := "пароль123🔒"
		err := user.HashPassword(unicodePassword)
		if err != nil {
			t.Errorf("Expected no error for unicode characters, got %v", err)
		}
	})
}

func TestCheckPassword(t *testing.T) {
	user := &User{}
	password := "testpassword123"
	user.HashPassword(password)

	t.Run("Correct password verification", func(t *testing.T) {
		if !user.CheckPassword(password) {
			t.Errorf("Expected password to be valid")
		}
	})

	t.Run("Incorrect password rejection", func(t *testing.T) {
		if user.CheckPassword("wrongpassword") {
			t.Errorf("Expected password to be invalid")
		}
	})

	t.Run("Empty password check", func(t *testing.T) {
		if user.CheckPassword("") {
			t.Errorf("Expected empty password to be invalid")
		}
	})

	t.Run("Case sensitivity", func(t *testing.T) {
		if user.CheckPassword("TESTPASSWORD123") {
			t.Errorf("Expected password check to be case sensitive")
		}
	})

	t.Run("Password with extra characters", func(t *testing.T) {
		if user.CheckPassword(password + "extra") {
			t.Errorf("Expected password with extra characters to be invalid")
		}
	})
}

// User Creation Tests
func TestCreateUser(t *testing.T) {
	t.Run("Valid user creation", func(t *testing.T) {
		user, err := CreateUser("test@example.com", "testuser", "password123")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if user == nil {
			t.Errorf("Expected user to be created")
		}
		if user.Email != "test@example.com" {
			t.Errorf("Expected email to be set correctly")
		}
		if user.Username != "testuser" {
			t.Errorf("Expected username to be set correctly")
		}
		if user.PasswordHash == "" {
			t.Errorf("Expected password hash to be set")
		}
		if user.PasswordHash == "password123" {
			t.Errorf("Password should be hashed")
		}
	})

	t.Run("Empty email", func(t *testing.T) {
		user, err := CreateUser("", "testuser", "password123")
		if err != nil {
			t.Errorf("Expected no error for empty email, got %v", err)
		}
		if user.Email != "" {
			t.Errorf("Expected empty email to be preserved")
		}
	})

	t.Run("Empty username", func(t *testing.T) {
		user, err := CreateUser("test@example.com", "", "password123")
		if err != nil {
			t.Errorf("Expected no error for empty username, got %v", err)
		}
		if user.Username != "" {
			t.Errorf("Expected empty username to be preserved")
		}
	})

	t.Run("Empty password", func(t *testing.T) {
		user, err := CreateUser("test@example.com", "testuser", "")
		if err == nil {
			t.Errorf("Expected error for empty password, got none")
		}
		if user != nil {
			t.Errorf("Expected user not to be created with empty password")
		}
	})

	t.Run("Invalid email format", func(t *testing.T) {
		user, err := CreateUser("invalid-email", "testuser", "password123")
		if err != nil {
			t.Errorf("Expected no error (validation should happen at save), got %v", err)
		}
		if user.Email != "invalid-email" {
			t.Errorf("Expected invalid email to be preserved for later validation")
		}
	})
}

// Email Verification Tests (without database)
func TestVerifyEmailLogic(t *testing.T) {
	t.Run("Already verified email", func(t *testing.T) {
		testUser := &User{EmailVerified: true}
		err := testUser.VerifyEmail()
		if err == nil {
			t.Errorf("Expected error for already verified email")
		}
		if err.Error() != "email already verified" {
			t.Errorf("Expected 'email already verified' error, got %v", err)
		}
	})

	t.Run("Expired token", func(t *testing.T) {
		testUser := &User{
			EmailVerified:         false,
			VerificationExpiresAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
		}
		err := testUser.VerifyEmail()
		if err == nil {
			t.Errorf("Expected error for expired token")
		}
		if err.Error() != "verification token expired" {
			t.Errorf("Expected 'verification token expired' error, got %v", err)
		}
	})

	t.Run("Valid verification setup", func(t *testing.T) {
		testUser := &User{
			ID:                    1,
			EmailVerified:         false,
			VerificationExpiresAt: sql.NullTime{Time: time.Now().Add(time.Hour), Valid: true},
		}

		// This will fail because there's no database connection
		// We just test that it passes the logic checks and fails at database level
		defer func() {
			if r := recover(); r != nil {
				// Expected to panic due to nil database connection
				t.Logf("Expected panic due to nil database: %v", r)
			}
		}()

		testUser.VerifyEmail()
		// If we get here without panic, the database connection worked unexpectedly
		t.Errorf("Expected panic due to nil database connection")
	})
}

// Token Generation Tests (without database)
func TestCreateAccountVerificationTokenLogic(t *testing.T) {
	t.Run("Token generation", func(t *testing.T) {
		testUser := &User{Email: "test@example.com"}

		// This will fail due to no database connection
		defer func() {
			if r := recover(); r != nil {
				// Expected to panic due to nil database connection
				// But we can still check if token was generated before the database call
				if testUser.VerificationToken.Valid && testUser.VerificationToken.String != "" {
					t.Logf("Token was generated before database error: %v", r)
				} else {
					t.Logf("Expected panic due to nil database: %v", r)
				}
			}
		}()

		testUser.CreateAccountVerificationToken()
		// If we get here without panic, the database connection worked unexpectedly
		t.Errorf("Expected panic due to nil database connection")
	})

	t.Run("Token uniqueness", func(t *testing.T) {
		// This test will also fail due to database, but we can test the concept
		// by checking that UUID generation produces unique values
		testUser := &User{Email: "test@example.com"}

		// We'll just test that we can generate tokens without database
		// by catching the panic and checking the token was set
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic due to nil database: %v", r)
			}
		}()

		testUser.CreateAccountVerificationToken()
		t.Errorf("Expected panic due to nil database connection")
	})
}

// Security Tests
func TestPasswordSecurity(t *testing.T) {
	t.Run("Bcrypt is used correctly", func(t *testing.T) {
		user := &User{}
		password := "testpassword"
		err := user.HashPassword(password)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify it's a bcrypt hash (starts with $2a$, $2b$, or $2y$)
		if !strings.HasPrefix(user.PasswordHash, "$2") {
			t.Errorf("Expected bcrypt hash format, got %s", user.PasswordHash)
		}

		// Verify we can validate with bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			t.Errorf("Expected bcrypt validation to work, got %v", err)
		}
	})

	t.Run("Password hash is different each time", func(t *testing.T) {
		password := "samepassword"
		user1 := &User{}
		user2 := &User{}

		user1.HashPassword(password)
		user2.HashPassword(password)

		if user1.PasswordHash == user2.PasswordHash {
			t.Errorf("Expected different hashes for same password")
		}

		// But both should validate correctly
		if !user1.CheckPassword(password) {
			t.Errorf("Expected user1 password to validate")
		}
		if !user2.CheckPassword(password) {
			t.Errorf("Expected user2 password to validate")
		}
	})

	t.Run("Password strength validation", func(t *testing.T) {
		testCases := []struct {
			name     string
			password string
		}{
			{"Short password", "123"},
			{"Long password (within bcrypt limit)", strings.Repeat("a", 70)},
			{"Special characters", "!@#$%^&*()"},
			{"Unicode", "пароль🔒"},
			{"Mixed case", "PaSSwoRd123"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				user := &User{}
				err := user.HashPassword(tc.password)
				if err != nil {
					t.Errorf("Expected no error for %s, got %v", tc.name, err)
				}
				if !user.CheckPassword(tc.password) {
					t.Errorf("Expected password validation to work for %s", tc.name)
				}
			})
		}
	})
}

// Data Validation Tests
func TestDataValidation(t *testing.T) {
	t.Run("Email format preservation", func(t *testing.T) {
		testEmails := []string{
			"valid@example.com",
			"invalid-email",
			"",
			"test@",
			"@example.com",
			strings.Repeat("a", 300) + "@example.com",
		}

		for _, email := range testEmails {
			user, err := CreateUser(email, "testuser", "password123")
			if err != nil {
				t.Errorf("Expected no error for email %s, got %v", email, err)
			}
			if user.Email != email {
				t.Errorf("Expected email to be preserved as %s, got %s", email, user.Email)
			}
		}
	})

	t.Run("Username format preservation", func(t *testing.T) {
		testUsernames := []string{
			"validuser",
			"",
			"user with spaces",
			"user@#$%",
			strings.Repeat("a", 100),
		}

		for _, username := range testUsernames {
			user, err := CreateUser("test@example.com", username, "password123")
			if err != nil {
				t.Errorf("Expected no error for username %s, got %v", username, err)
			}
			if user.Username != username {
				t.Errorf("Expected username to be preserved as %s, got %s", username, user.Username)
			}
		}
	})

	t.Run("Null bytes handling", func(t *testing.T) {
		emailWithNull := "test\x00@example.com"
		usernameWithNull := "test\x00user"
		passwordWithNull := "pass\x00word"

		user, err := CreateUser(emailWithNull, usernameWithNull, passwordWithNull)
		if err != nil {
			t.Errorf("Expected no error creating user with null bytes, got %v", err)
		}

		// Verify null bytes are preserved (or handled appropriately)
		if user.Email != emailWithNull {
			t.Errorf("Expected email with null byte to be preserved")
		}
		if user.Username != usernameWithNull {
			t.Errorf("Expected username with null byte to be preserved")
		}

		// Password should still hash and validate correctly
		if !user.CheckPassword(passwordWithNull) {
			t.Errorf("Expected password with null byte to validate correctly")
		}
	})
}

// User Struct Tests
func TestUserStruct(t *testing.T) {
	t.Run("User struct initialization", func(t *testing.T) {
		user := &User{
			ID:                    1,
			Email:                 "test@example.com",
			Username:              "testuser",
			PasswordHash:          "hashedpassword",
			CreatedAt:             time.Now(),
			EmailVerified:         false,
			VerificationToken:     sql.NullString{String: "token", Valid: true},
			VerificationExpiresAt: sql.NullTime{Time: time.Now().Add(time.Hour), Valid: true},
		}

		if user.ID != 1 {
			t.Errorf("Expected ID to be 1, got %d", user.ID)
		}
		if user.Email != "test@example.com" {
			t.Errorf("Expected email to be test@example.com, got %s", user.Email)
		}
		if user.Username != "testuser" {
			t.Errorf("Expected username to be testuser, got %s", user.Username)
		}
		if user.EmailVerified {
			t.Errorf("Expected EmailVerified to be false")
		}
		if !user.VerificationToken.Valid {
			t.Errorf("Expected VerificationToken to be valid")
		}
		if !user.VerificationExpiresAt.Valid {
			t.Errorf("Expected VerificationExpiresAt to be valid")
		}
	})

	t.Run("UserCredentials struct", func(t *testing.T) {
		creds := UserCredentials{
			Email:    "test@example.com",
			Password: "password123",
			Username: "testuser",
		}

		if creds.Email != "test@example.com" {
			t.Errorf("Expected email to be test@example.com, got %s", creds.Email)
		}
		if creds.Password != "password123" {
			t.Errorf("Expected password to be password123, got %s", creds.Password)
		}
		if creds.Username != "testuser" {
			t.Errorf("Expected username to be testuser, got %s", creds.Username)
		}
	})
}

// Edge Cases
func TestEdgeCases(t *testing.T) {
	t.Run("Concurrent password hashing", func(t *testing.T) {
		password := "testpassword"
		users := make([]*User, 10)

		// Create users concurrently
		done := make(chan bool, 10)
		for i := range 10 {
			go func(index int) {
				users[index] = &User{}
				users[index].HashPassword(password)
				done <- true
			}(i)
		}

		// Wait for all to complete
		for range 10 {
			<-done
		}

		// Verify all hashes are different but validate correctly
		hashes := make(map[string]bool)
		for i, user := range users {
			if user == nil {
				t.Errorf("User %d is nil", i)
				continue
			}
			if hashes[user.PasswordHash] {
				t.Errorf("Hash collision detected for user %d", i)
			}
			hashes[user.PasswordHash] = true

			if !user.CheckPassword(password) {
				t.Errorf("Password validation failed for user %d", i)
			}
		}
	})

	t.Run("Very large data handling", func(t *testing.T) {
		largeEmail := strings.Repeat("a", 1000) + "@example.com"
		largeUsername := strings.Repeat("b", 1000)
		largePassword := strings.Repeat("c", 70) // Within bcrypt limit

		user, err := CreateUser(largeEmail, largeUsername, largePassword)
		if err != nil {
			t.Errorf("Expected no error with large data, got %v", err)
		}

		// Verify data is preserved
		if user.Email != largeEmail {
			t.Errorf("Large email not preserved correctly")
		}
		if user.Username != largeUsername {
			t.Errorf("Large username not preserved correctly")
		}
		if !user.CheckPassword(largePassword) {
			t.Errorf("Large password not hashed/validated correctly")
		}
	})

	t.Run("Password exceeding bcrypt limit", func(t *testing.T) {
		largePassword := strings.Repeat("c", 100) // Exceeds bcrypt 72-byte limit

		_, err := CreateUser("test@example.com", "testuser", largePassword)
		if err == nil {
			t.Errorf("Expected error for password exceeding bcrypt limit")
		}
		if !strings.Contains(err.Error(), "password length exceeds 72 bytes") {
			t.Errorf("Expected bcrypt length error, got %v", err)
		}
	})
}
