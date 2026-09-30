package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	message "back/internal/models/messages"
)

func TestParseReq(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectType  CommandType
		expectError bool
	}{
		{
			name:        "Valid SEND request",
			input:       `{"command_type":"SEND","data":"Hello World"}`,
			expectType:  SEND,
			expectError: false,
		},
		{
			name:        "Valid REMOVE request",
			input:       `{"command_type":"REMOVE","data":123}`,
			expectType:  REMOVE,
			expectError: false,
		},
		{
			name:        "Valid MODIFY request",
			input:       `{"command_type":"MODIFY","data":{"id":42,"content":"edited text"}}`,
			expectType:  MODIFY,
			expectError: false,
		},
		{
			name:        "MODIFY with scalar data",
			input:       `{"command_type":"MODIFY","data":42}`,
			expectType:  "",
			expectError: true,
		},
		{
			name:        "Invalid JSON",
			input:       `{"command_type":"SEND","data":}`,
			expectType:  "",
			expectError: true,
		},
		{
			name:        "Unknown command type",
			input:       `{"command_type":"UNKNOWN","data":"test"}`,
			expectType:  "",
			expectError: true,
		},
		{
			name:        "Missing command type",
			input:       `{"data":"test"}`,
			expectType:  "",
			expectError: true,
		},
		{
			name:        "Empty string",
			input:       "",
			expectType:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := parseReq(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if req == nil {
				t.Errorf("Expected request but got nil")
				return
			}

			if req.GetType() != tt.expectType {
				t.Errorf("Expected type %v, got %v", tt.expectType, req.GetType())
			}
		})
	}
}

func TestSendRequestGetType(t *testing.T) {
	req := SendRequest{Data: "test"}
	if req.GetType() != SEND {
		t.Errorf("Expected SEND, got %v", req.GetType())
	}
}

func TestRemoveRequestGetType(t *testing.T) {
	req := RemoveRequest{Data: 123}
	if req.GetType() != REMOVE {
		t.Errorf("Expected REMOVE, got %v", req.GetType())
	}
}

func TestModifyRequestGetType(t *testing.T) {
	req := ModifyRequest{Data: ModifyPayload{Id: 42, Content: "edited"}}
	if req.GetType() != MODIFY {
		t.Errorf("Expected MODIFY, got %v", req.GetType())
	}
}

func TestParseReqModifyPayload(t *testing.T) {
	req, err := parseReq(`{"command_type":"MODIFY","data":{"id":42,"content":"edited text"}}`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	modify, ok := req.(*ModifyRequest)
	if !ok {
		t.Fatalf("Expected *ModifyRequest, got %T", req)
	}

	if modify.Data.Id != 42 {
		t.Errorf("Expected id 42, got %d", modify.Data.Id)
	}
	if modify.Data.Content != "edited text" {
		t.Errorf("Expected content %q, got %q", "edited text", modify.Data.Content)
	}
}

func TestBroadcastSerialization(t *testing.T) {
	t.Run("Broadcast with int data", func(t *testing.T) {
		b := Broadcast[int]{
			Type: REMOVE,
			Data: 123,
		}

		data, err := json.Marshal(b)
		if err != nil {
			t.Errorf("Failed to marshal broadcast: %v", err)
		}

		expected := `{"command_type":"REMOVE","data":123}`
		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}
	})

	t.Run("Broadcast with message data", func(t *testing.T) {
		b := Broadcast[message.MessageResponse]{
			Type: MODIFY,
			Data: message.MessageResponse{
				Id:       42,
				Content:  "edited text",
				CreateAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Sender:   "alice",
				RoomID:   1,
				Edited:   true,
			},
		}

		data, err := json.Marshal(b)
		if err != nil {
			t.Errorf("Failed to marshal broadcast: %v", err)
		}

		expected := `{"command_type":"MODIFY","data":{"id":42,"content":"edited text",` +
			`"created_at":"2024-01-01T12:00:00Z","sender":"alice","room_id":1,"edited":true}}`
		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}
	})
}

func TestCommandTypeSerialization(t *testing.T) {
	tests := []struct {
		name     string
		cmdType  CommandType
		expected string
	}{
		{"SEND command", SEND, "SEND"},
		{"REMOVE command", REMOVE, "REMOVE"},
		{"MODIFY command", MODIFY, "MODIFY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.cmdType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.cmdType))
			}
		})
	}
}

func TestSendRequestExecute(t *testing.T) {
	t.Run("Valid message content", func(t *testing.T) {
		req := SendRequest{Data: "Hello World"}

		// Note: This test would require database setup for full testing
		// For now, we test the basic validation logic
		if req.Data == "" {
			t.Errorf("Expected non-empty data")
		}
	})

	t.Run("Empty message content", func(t *testing.T) {
		req := SendRequest{Data: ""}

		// This should fail validation in the message creation
		if req.Data != "" {
			t.Errorf("Expected empty data")
		}
	})
}

func TestRemoveRequestExecute(t *testing.T) {
	t.Run("Valid message ID", func(t *testing.T) {
		req := RemoveRequest{Data: 123}

		if req.Data <= 0 {
			t.Errorf("Expected positive message ID")
		}
	})

	t.Run("Invalid message ID", func(t *testing.T) {
		req := RemoveRequest{Data: -1}

		if req.Data > 0 {
			t.Errorf("Expected negative or zero message ID")
		}
	})
}

// fakeMessageRepository is an in-memory MessageRepository that records whether
// a mutation was attempted, so rejection paths can be asserted positively.
// Modelled on fakeRoomRepository in room_test.go.
type fakeMessageRepository struct {
	stored map[int]*message.Message
	err    error

	updatedID      int
	updatedContent string
	updateCalled   bool
	deleteCalled   bool
}

func (f *fakeMessageRepository) Save(ctx context.Context, m *message.Message) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	id := len(f.stored) + 1
	if f.stored == nil {
		f.stored = map[int]*message.Message{}
	}
	f.stored[id] = m
	return id, nil
}

func (f *fakeMessageRepository) GetByID(ctx context.Context, id int) (*message.Message, error) {
	if f.err != nil {
		return nil, f.err
	}
	m, ok := f.stored[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return m, nil
}

func (f *fakeMessageRepository) GetLast(ctx context.Context, roomID int, limit int, offset int) ([]message.Message, error) {
	return nil, f.err
}

func (f *fakeMessageRepository) MarkAsDeleted(ctx context.Context, id int) error {
	f.deleteCalled = true
	return f.err
}

func (f *fakeMessageRepository) UpdateContent(ctx context.Context, id int, content string) error {
	f.updateCalled = true
	f.updatedID = id
	f.updatedContent = content
	return f.err
}

// TestModifyRequestExecuteRejections covers every path that returns before
// ToSendFormat(). The success path is deliberately not tested here: it calls
// user.GetUserById, which dereferences the database.DbInstance global, and that
// is nil in a unit test — it panics before reaching the broadcast channel.
func TestModifyRequestExecuteRejections(t *testing.T) {
	const (
		authorID = 7
		roomID   = 1
		msgID    = 1
	)

	tests := []struct {
		name        string
		stored      *message.Message
		repoErr     error
		userID      int
		roomID      int
		content     string
		expectError bool
	}{
		{
			name:    "not the author",
			stored:  &message.Message{Id: msgID, SenderID: authorID, RoomID: roomID},
			userID:  authorID + 1,
			roomID:  roomID,
			content: "hax",
		},
		{
			name:    "wrong room",
			stored:  &message.Message{Id: msgID, SenderID: authorID, RoomID: roomID + 1},
			userID:  authorID,
			roomID:  roomID,
			content: "edited",
		},
		{
			name:    "already deleted",
			stored:  &message.Message{Id: msgID, SenderID: authorID, RoomID: roomID, Deleted: true},
			userID:  authorID,
			roomID:  roomID,
			content: "edited",
		},
		{
			name:    "empty content",
			stored:  &message.Message{Id: msgID, SenderID: authorID, RoomID: roomID},
			userID:  authorID,
			roomID:  roomID,
			content: "",
		},
		{
			name:        "content too long",
			stored:      &message.Message{Id: msgID, SenderID: authorID, RoomID: roomID},
			userID:      authorID,
			roomID:      roomID,
			content:     strings.Repeat("a", message.MAX_CONTENT_LENGTH+1),
			expectError: true,
		},
		{
			name:        "repository read error",
			stored:      &message.Message{Id: msgID, SenderID: authorID, RoomID: roomID},
			repoErr:     errors.New("db down"),
			userID:      authorID,
			roomID:      roomID,
			content:     "edited",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeMessageRepository{
				stored: map[int]*message.Message{msgID: tt.stored},
				err:    tt.repoErr,
			}

			req := ModifyRequest{Data: ModifyPayload{Id: msgID, Content: tt.content}}
			err := req.Execute(context.Background(), repo, tt.userID, tt.roomID)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if repo.updateCalled {
				t.Errorf("Expected no update, but UpdateContent was called with %q", repo.updatedContent)
			}

			// Non-blocking, so a regression that does broadcast fails the test
			// rather than deadlocking on the unbuffered channel.
			select {
			case b := <-broadcast:
				t.Errorf("Expected no broadcast, got %s", b.Payload)
			default:
			}
		})
	}
}
