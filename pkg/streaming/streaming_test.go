package streaming

import (
	"testing"
)

func TestEventTypes(t *testing.T) {
	if EventTypeMessage != "message" {
		t.Errorf("Expected EventTypeMessage to be 'message', got %s", EventTypeMessage)
	}
	if EventTypeError != "error" {
		t.Errorf("Expected EventTypeError to be 'error', got %s", EventTypeError)
	}
	if EventTypeDone != "done" {
		t.Errorf("Expected EventTypeDone to be 'done', got %s", EventTypeDone)
	}
}