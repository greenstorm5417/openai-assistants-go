package runs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/greenstorm5417/openai-assistants-go/client"
)

func TestCreateRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("OpenAI-Beta") != "assistants=v2" {
			t.Errorf("Expected OpenAI-Beta header to be assistants=v2")
		}

		run := Run{
			ID:          "run_123",
			Object:      "thread.run",
			CreatedAt:   1699000000,
			ThreadID:    "thread_123",
			AssistantID: "asst_123",
			Status:      "queued",
			Model:       "gpt-4",
		}
		json.NewEncoder(w).Encode(run)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	req := &CreateRunRequest{
		AssistantID:  "asst_123",
		Model:        stringPtr("gpt-4"),
		Instructions: stringPtr("You are a helpful assistant."),
	}

	run, err := service.Create("thread_123", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if run.ID != "run_123" {
		t.Errorf("Expected ID run_123, got %s", run.ID)
	}
}

func TestCreateAndStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Error("Expected ResponseWriter to be a Flusher")
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		events := []string{
			`event: thread.run.created
data: {"id":"run_123","object":"thread.run","created_at":1699000000,"status":"queued"}`,
			`event: thread.run.queued
data: {"id":"run_123","object":"thread.run","created_at":1699000000,"status":"queued"}`,
			`event: thread.run.in_progress
data: {"id":"run_123","object":"thread.run","created_at":1699000000,"status":"in_progress"}`,
			`event: thread.run.completed
data: {"id":"run_123","object":"thread.run","created_at":1699000000,"status":"completed"}`,
			"data: [DONE]",
		}

		for _, event := range events {
			_, err := w.Write([]byte(event + "\n\n"))
			if err != nil {
				t.Errorf("Error writing event: %v", err)
				return
			}
			flusher.Flush()
		}
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	req := &CreateRunRequest{
		AssistantID: "asst_123",
		Stream:      true,
	}

	events, err := service.CreateAndStream("thread_123", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedEvents := []string{
		"thread.run.created",
		"thread.run.queued",
		"thread.run.in_progress",
		"thread.run.completed",
		"done",
	}

	i := 0
	for event := range events {
		if i >= len(expectedEvents) {
			t.Errorf("Received more events than expected")
			break
		}

		if event.Event != expectedEvents[i] {
			t.Errorf("Expected event %s, got %s", expectedEvents[i], event.Event)
		}
		i++
	}

	if i != len(expectedEvents) {
		t.Errorf("Expected %d events, got %d", len(expectedEvents), i)
	}
}

func TestSubmitToolOutputs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		run := Run{
			ID:          "run_123",
			Object:      "thread.run",
			CreatedAt:   1699000000,
			ThreadID:    "thread_123",
			AssistantID: "asst_123",
			Status:      "completed",
			Model:       "gpt-4",
		}
		json.NewEncoder(w).Encode(run)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	req := &SubmitToolOutputsRequest{
		ToolOutputs: []ToolOutput{
			{
				ToolCallID: "call_123",
				Output:     "The weather is sunny and 72°F",
			},
		},
	}

	run, err := service.SubmitToolOutputs("thread_123", "run_123", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if run.ID != "run_123" {
		t.Errorf("Expected ID run_123, got %s", run.ID)
	}

	if run.Status != "completed" {
		t.Errorf("Expected status completed, got %s", run.Status)
	}
}

func TestSubmitToolOutputsStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Error("Expected ResponseWriter to be a Flusher")
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		events := []string{
			`event: thread.run.queued
data: {"id":"run_123","object":"thread.run","created_at":1699000000,"status":"queued"}`,
			`event: thread.run.in_progress
data: {"id":"run_123","object":"thread.run","created_at":1699000000,"status":"in_progress"}`,
			`event: thread.message.created
data: {"id":"msg_123","object":"thread.message","role":"assistant","content":"The weather in San Francisco is sunny and 72°F."}`,
			`event: thread.run.completed
data: {"id":"run_123","object":"thread.run","created_at":1699000000,"status":"completed"}`,
			"data: [DONE]",
		}

		for _, event := range events {
			_, err := w.Write([]byte(event + "\n\n"))
			if err != nil {
				t.Errorf("Error writing event: %v", err)
				return
			}
			flusher.Flush()
		}
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	req := &SubmitToolOutputsRequest{
		ToolOutputs: []ToolOutput{
			{
				ToolCallID: "call_123",
				Output:     "The weather is sunny and 72°F",
			},
		},
		Stream: true,
	}

	events, err := service.SubmitToolOutputsStream("thread_123", "run_123", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedEvents := []string{
		"thread.run.queued",
		"thread.run.in_progress",
		"thread.message.created",
		"thread.run.completed",
		"done",
	}

	i := 0
	for event := range events {
		if i >= len(expectedEvents) {
			t.Errorf("Received more events than expected")
			break
		}

		if event.Event != expectedEvents[i] {
			t.Errorf("Expected event %s, got %s", expectedEvents[i], event.Event)
		}
		i++
	}

	if i != len(expectedEvents) {
		t.Errorf("Expected %d events, got %d", len(expectedEvents), i)
	}
}

func TestCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		run := Run{
			ID:          "run_123",
			Object:      "thread.run",
			CreatedAt:   1699000000,
			ThreadID:    "thread_123",
			AssistantID: "asst_123",
			Status:      "cancelled",
			Model:       "gpt-4",
		}
		json.NewEncoder(w).Encode(run)
	}))
	defer server.Close()

	c := &client.Client{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		HTTPClient: server.Client(),
	}

	service := New(c)

	run, err := service.Cancel("thread_123", "run_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if run.ID != "run_123" {
		t.Errorf("Expected ID run_123, got %s", run.ID)
	}

	if run.Status != "cancelled" {
		t.Errorf("Expected status cancelled, got %s", run.Status)
	}
}

func stringPtr(s string) *string {
	return &s
}
