package threads

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/types"
)

func TestCreateThread(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("OpenAI-Beta") != "assistants=v2" {
			t.Errorf("Expected OpenAI-Beta header to be assistants=v2")
		}

		thread := Thread{
			ID:        "thread_123",
			Object:    "thread",
			CreatedAt: 1699000000,
			Metadata:  types.Metadata{"key": "value"},
		}
		json.NewEncoder(w).Encode(thread)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	// Test creating a thread with messages
	messages := []Message{
		{
			Role:    "user",
			Content: "Hello, what is AI?",
		},
		{
			Role:    "user",
			Content: "How does AI work? Explain it in simple terms.",
		},
	}

	thread, err := service.Create(&CreateThreadRequest{
		Messages: messages,
		Metadata: types.Metadata{"purpose": "test"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if thread.ID != "thread_123" {
		t.Errorf("Expected ID thread_123, got %s", thread.ID)
	}
}

func TestGetThread(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		thread := Thread{
			ID:        "thread_123",
			Object:    "thread",
			CreatedAt: 1699000000,
			ToolResources: &ToolResources{
				CodeInterpreter: &CodeInterpreterResources{
					FileIDs: []string{},
				},
			},
		}
		json.NewEncoder(w).Encode(thread)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	thread, err := service.Get("thread_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if thread.ID != "thread_123" {
		t.Errorf("Expected ID thread_123, got %s", thread.ID)
	}
}

func TestModifyThread(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		thread := Thread{
			ID:        "thread_123",
			Object:    "thread",
			CreatedAt: 1699000000,
			Metadata: types.Metadata{
				"modified": "true",
				"user":     "abc123",
			},
		}
		json.NewEncoder(w).Encode(thread)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	metadata := types.Metadata{
		"modified": "true",
		"user":     "abc123",
	}

	thread, err := service.Modify("thread_123", nil, metadata)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if thread.ID != "thread_123" {
		t.Errorf("Expected ID thread_123, got %s", thread.ID)
	}

	if thread.Metadata["modified"] != "true" {
		t.Errorf("Expected modified metadata to be true")
	}
}

func TestDeleteThread(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		response := DeleteThreadResponse{
			ID:      "thread_123",
			Object:  "thread.deleted",
			Deleted: true,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	response, err := service.Delete("thread_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.Deleted {
		t.Error("Expected deleted to be true")
	}
}

func TestMessage(t *testing.T) {
	tests := []struct {
		name    string
		message Message
	}{
		{
			name: "simple text message",
			message: Message{
				Role:    "user",
				Content: "Hello, world!",
			},
		},
		{
			name: "message with metadata",
			message: Message{
				Role:    "user",
				Content: "Message with metadata",
				Metadata: types.Metadata{
					"importance": "high",
					"category":   "test",
				},
			},
		},
		{
			name: "message with attachments",
			message: Message{
				Role:        "user",
				Content:     "Message with attachments",
				Attachments: []string{"file1", "file2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.message)
			if err != nil {
				t.Fatalf("Failed to marshal message: %v", err)
			}

			var decoded Message
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("Failed to unmarshal message: %v", err)
			}

			if decoded.Role != tt.message.Role {
				t.Errorf("Expected role %s, got %s", tt.message.Role, decoded.Role)
			}

			if decoded.Content != tt.message.Content {
				t.Errorf("Expected content %s, got %s", tt.message.Content, decoded.Content)
			}

			if len(tt.message.Attachments) > 0 {
				if len(decoded.Attachments) != len(tt.message.Attachments) {
					t.Errorf("Expected %d attachments, got %d", len(tt.message.Attachments), len(decoded.Attachments))
				}
			}

			if len(tt.message.Metadata) > 0 {
				if len(decoded.Metadata) != len(tt.message.Metadata) {
					t.Errorf("Expected %d metadata entries, got %d", len(tt.message.Metadata), len(decoded.Metadata))
				}
			}
		})
	}
}
