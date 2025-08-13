package message

import "testing"

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

}
