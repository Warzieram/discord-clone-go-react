package handlers
/*
import (
	"back/internal/models/user"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByVerificationToken(token string) (*user.User, error) {
	args := m.Called(token)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) VerifyEmail() error {
	args := m.Called()
	return args.Error(0)
}

type User struct {
	ID                    int            `json:"id"`
	Email                 string         `json:"email"`
	Username							string				 `json:"username"`
	PasswordHash          string         `json:"-"`
	CreatedAt             time.Time      `json:"created_at"`
	EmailVerified         bool           `json:"email_verified"`
	VerificationToken     sql.NullString `json:"verification_token"`
	VerificationExpiresAt sql.NullTime   `json:"verification_expires_at"`
}

func TestVerifyEmail(t *testing.T) {
	mockUserService := new(MockUserService)

	// case with no token
	t.Run("Missing token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/verify-email", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			VerifyEmail(w, r)
		})

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Missing token param\n", rec.Body.String())
	})

	// case with valid token
	t.Run("Valid token", func(t *testing.T) {
		mockUser := &User{Email: "test@example.com"}
		mockUserService.On("GetUserByVerificationToken", "valid-token").Return(mockUser, nil)
		mockUserService.On("VerifyEmail").Return(nil)
		
		req, err := http.NewRequest("GET", "/verify-email?token=valid-token", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			VerifyEmail(w, r)
		})

		handler.ServeHTTP(rec, req)


		assert.Equal(t, http.StatusMovedPermanently, rec.Code)
		assert.Equal(t, "http://192.168.1.151:5173/login", rec.Header().Get("Location"))
	})
}
*/
