package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"back/internal/models/room"
)

type fakeRoomRepository struct {
	rooms []room.Room
	err   error
}

func (f *fakeRoomRepository) Save(ctx context.Context, r *room.Room) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	id := len(f.rooms) + 1
	f.rooms = append(f.rooms, room.Room{Id: id, Name: r.Name})
	return id, nil
}

func (f *fakeRoomRepository) GetAll(ctx context.Context) ([]room.Room, error) {
	return f.rooms, f.err
}

func TestCreateRoom(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		repoErr        error
		expectedStatus int
	}{
		{"Valid room", `{"name": "general"}`, nil, http.StatusCreated},
		{"Invalid JSON", `{"name":`, nil, http.StatusBadRequest},
		{"Empty name", `{"name": ""}`, nil, http.StatusBadRequest},
		{"Name too long", `{"name": "` + strings.Repeat("a", room.MAX_NAME_LENGTH+1) + `"}`, nil, http.StatusBadRequest},
		{"Repository error", `{"name": "general"}`, errors.New("db down"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewRoomHandlers(&fakeRoomRepository{err: tt.repoErr})
			req := httptest.NewRequest(http.MethodPost, "/api/rooms", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			h.CreateRoom(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if tt.expectedStatus == http.StatusCreated {
				var created room.Room
				if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
					t.Fatalf("Couldn't decode response: %v", err)
				}
				if created.Id != 1 || created.Name != "general" {
					t.Errorf("Expected room {1 general}, got %v", created)
				}
			}
		})
	}
}

func TestListRooms(t *testing.T) {
	t.Run("Returns every room", func(t *testing.T) {
		repo := &fakeRoomRepository{rooms: []room.Room{{Id: 1, Name: "general"}, {Id: 2, Name: "random"}}}
		h := NewRoomHandlers(repo)
		w := httptest.NewRecorder()

		h.ListRooms(w, httptest.NewRequest(http.MethodGet, "/api/rooms", nil))

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}
		var rooms []room.Room
		if err := json.NewDecoder(w.Body).Decode(&rooms); err != nil {
			t.Fatalf("Couldn't decode response: %v", err)
		}
		if len(rooms) != 2 {
			t.Errorf("Expected 2 rooms, got %d", len(rooms))
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		h := NewRoomHandlers(&fakeRoomRepository{err: errors.New("db down")})
		w := httptest.NewRecorder()

		h.ListRooms(w, httptest.NewRequest(http.MethodGet, "/api/rooms", nil))

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})
}
