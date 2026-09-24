package room

import (
	"strings"
	"testing"
)

func TestCreateRoom(t *testing.T) {
	t.Run("Valid room creation", func(t *testing.T) {
		room, err := CreateRoom("  general ")
		if err != nil {
			t.Errorf("Expected no error, got : %v", err)
		}
		if room == nil || room.Name != "general" {
			t.Errorf("Expected room named \"general\", got : %v", room)
		}
	})

	t.Run("Empty name", func(t *testing.T) {
		_, err := CreateRoom("   ")
		if err == nil {
			t.Error("Expected an error for an empty name")
		}
	})

	t.Run("Name too long", func(t *testing.T) {
		_, err := CreateRoom(strings.Repeat("a", MAX_NAME_LENGTH+1))
		if err == nil {
			t.Error("Expected an error for a name that is too long")
		}
	})
}
