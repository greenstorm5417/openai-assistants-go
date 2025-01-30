package assistants

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
)

func TestCreateAssistant(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check Beta header
		if r.Header.Get("OpenAI-Beta") != "assistants=v2" {
			t.Errorf("Expected OpenAI-Beta header to be assistants=v2")
		}

		// Return mock response
		assistant := Assistant{
			ID:        "asst_123",
			Object:    "assistant",
			CreatedAt: 1699000000,
			Model:     "gpt-4",
			Tools:     []Tool{{Type: "code_interpreter"}},
		}
		json.NewEncoder(w).Encode(assistant)
	}))
	defer server.Close()

	// Create client with mock server URL
	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	// Create test request
	name := "Test Assistant"
	req := &CreateAssistantRequest{
		Model: "gpt-4",
		Name:  &name,
		Tools: []Tool{{Type: "code_interpreter"}},
	}

	// Make request
	assistant, err := service.Create(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if assistant.ID != "asst_123" {
		t.Errorf("Expected ID asst_123, got %s", assistant.ID)
	}
}

func TestListAssistants(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		response := ListAssistantsResponse{
			Object:  "list",
			Data:    []Assistant{{ID: "asst_123", Object: "assistant"}},
			FirstID: "asst_123",
			LastID:  "asst_123",
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
	params := &ListAssistantsParams{
		Limit: &limit,
		Order: &order,
	}

	response, err := service.List(params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 assistant, got %d", len(response.Data))
	}
}

func TestGetAssistant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		assistant := Assistant{
			ID:     "asst_123",
			Object: "assistant",
			Model:  "gpt-4",
		}
		json.NewEncoder(w).Encode(assistant)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	assistant, err := service.Get("asst_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if assistant.ID != "asst_123" {
		t.Errorf("Expected ID asst_123, got %s", assistant.ID)
	}
}

func TestModifyAssistant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		assistant := Assistant{
			ID:     "asst_123",
			Object: "assistant",
			Model:  "gpt-4",
		}
		json.NewEncoder(w).Encode(assistant)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	name := "Updated Assistant"
	req := &CreateAssistantRequest{
		Model: "gpt-4",
		Name:  &name,
	}

	assistant, err := service.Modify("asst_123", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if assistant.ID != "asst_123" {
		t.Errorf("Expected ID asst_123, got %s", assistant.ID)
	}
}

func TestDeleteAssistant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		response := DeleteAssistantResponse{
			ID:      "asst_123",
			Object:  "assistant.deleted",
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

	response, err := service.Delete("asst_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.Deleted {
		t.Error("Expected deleted to be true")
	}
}
