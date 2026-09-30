package message

import (
	"strings"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	t.Run("Valid message creation", func(t *testing.T) {
		message, err := CreateMessage("Hello Worlders ", 1, 1)
		if err != nil {
			t.Errorf("Expected no error, got : %v", err)
		}
		if message == nil {
			t.Error("Expected message to be created")
		}
	})

	t.Run("Empty content", func(t *testing.T) {
		_, err := CreateMessage("", 1, 1)
		if err == nil {
			t.Error("Expected an error for empty content")
		}
	})

	t.Run("Content too long", func(t *testing.T) {
		_, err := CreateMessage(strings.Repeat("a", MAX_CONTENT_LENGTH+1), 1, 1)
		if err == nil {
			t.Error("Expected an error for content over MAX_CONTENT_LENGTH")
		}
	})

	t.Run("Content at max length", func(t *testing.T) {
		_, err := CreateMessage(strings.Repeat("a", MAX_CONTENT_LENGTH), 1, 1)
		if err != nil {
			t.Errorf("Expected no error at exactly MAX_CONTENT_LENGTH, got: %v", err)
		}
	})

}
