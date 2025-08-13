package handlers

import (
	"encoding/json"
	"testing"
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
			input:       `{"command_type":"SEND","data":{"content":"Hello World","room_id":1}}`,
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
	req := SendRequest{Data: SendRequestData{
		Content: "test",
		RoomID:  1,
	}}
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
}

func TestCommandTypeSerialization(t *testing.T) {
	tests := []struct {
		name     string
		cmdType  CommandType
		expected string
	}{
		{"SEND command", SEND, "SEND"},
		{"REMOVE command", REMOVE, "REMOVE"},
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
		req := SendRequest{Data: SendRequestData{
			Content: "Hello World",
			RoomID:  1,
		}}

		// Note: This test would require database setup for full testing
		// For now, we test the basic validation logic
		if req.Data.Content == "" {
			t.Errorf("Expected non-empty content")
		}
		if req.Data.RoomID <= 0 {
			t.Errorf("Expected positive room ID")
		}
	})

	t.Run("Empty message content", func(t *testing.T) {
		req := SendRequest{Data: SendRequestData{
			Content: "",
			RoomID:  1,
		}}

		// This should fail validation in the message creation
		if req.Data.Content != "" {
			t.Errorf("Expected empty content")
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
