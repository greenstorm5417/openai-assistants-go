package messages

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/types"
)

func TestCreateMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("OpenAI-Beta") != "assistants=v2" {
			t.Errorf("Expected OpenAI-Beta header to be assistants=v2")
		}

		message := Message{
			ID:        "msg_123",
			Object:    "thread.message",
			CreatedAt: 1699000000,
			ThreadID:  "thread_123",
			Role:      "user",
			Content: []Content{
				{
					Type: "text",
					Text: &Text{
						Value:       "Hello, what is AI?",
						Annotations: []Annotation{},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(message)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	req := &CreateMessageRequest{
		Role:    "user",
		Content: "Hello, what is AI?",
	}

	message, err := service.Create("thread_123", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if message.ID != "msg_123" {
		t.Errorf("Expected ID msg_123, got %s", message.ID)
	}
}

func TestListMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		response := ListMessagesResponse{
			Object: "list",
			Data: []Message{
				{
					ID:        "msg_123",
					Object:    "thread.message",
					CreatedAt: 1699000000,
					ThreadID:  "thread_123",
					Role:      "user",
					Content: []Content{
						{
							Type: "text",
							Text: &Text{
								Value:       "Hello, what is AI?",
								Annotations: []Annotation{},
							},
						},
					},
				},
			},
			FirstID: "msg_123",
			LastID:  "msg_123",
			HasMore: false,
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

	limit := 20
	order := "desc"
	params := &ListMessagesParams{
		Limit: &limit,
		Order: &order,
	}

	response, err := service.List("thread_123", params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 message, got %d", len(response.Data))
	}
}

func TestGetMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		message := Message{
			ID:        "msg_123",
			Object:    "thread.message",
			CreatedAt: 1699000000,
			ThreadID:  "thread_123",
			Role:      "user",
			Content: []Content{
				{
					Type: "text",
					Text: &Text{
						Value:       "Hello, what is AI?",
						Annotations: []Annotation{},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(message)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	message, err := service.Get("thread_123", "msg_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if message.ID != "msg_123" {
		t.Errorf("Expected ID msg_123, got %s", message.ID)
	}
}

func TestModifyMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		message := Message{
			ID:        "msg_123",
			Object:    "thread.message",
			CreatedAt: 1699000000,
			ThreadID:  "thread_123",
			Role:      "user",
			Content: []Content{
				{
					Type: "text",
					Text: &Text{
						Value:       "Hello, what is AI?",
						Annotations: []Annotation{},
					},
				},
			},
			Metadata: types.Metadata{
				"modified": "true",
				"user":     "abc123",
			},
		}
		json.NewEncoder(w).Encode(message)
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

	message, err := service.Modify("thread_123", "msg_123", metadata)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if message.ID != "msg_123" {
		t.Errorf("Expected ID msg_123, got %s", message.ID)
	}

	if message.Metadata["modified"] != "true" {
		t.Errorf("Expected modified metadata to be true")
	}
}

func TestDeleteMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		response := DeleteMessageResponse{
			ID:      "msg_123",
			Object:  "thread.message.deleted",
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

	response, err := service.Delete("thread_123", "msg_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.Deleted {
		t.Error("Expected deleted to be true")
	}
}

func TestMessageContent(t *testing.T) {
	tests := []struct {
		name     string
		content  Content
		expected string
	}{
		{
			name: "text content",
			content: Content{
				Type: "text",
				Text: &Text{
					Value:       "Hello, world!",
					Annotations: []Annotation{},
				},
			},
			expected: "text",
		},
		{
			name: "image url content",
			content: Content{
				Type: "image_url",
				ImageURL: &ImageURL{
					URL:    "https://example.com/image.jpg",
					Detail: "high",
				},
			},
			expected: "image_url",
		},
		{
			name: "image file content",
			content: Content{
				Type: "image_file",
				ImageFile: &ImageFile{
					FileID: "file_123",
					Detail: "low",
				},
			},
			expected: "image_file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.content)
			if err != nil {
				t.Fatalf("Failed to marshal content: %v", err)
			}

			var decoded Content
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("Failed to unmarshal content: %v", err)
			}

			if decoded.Type != tt.expected {
				t.Errorf("Expected content type %s, got %s", tt.expected, decoded.Type)
			}
		})
	}
}
