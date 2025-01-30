// C:\Users\sduss\Desktop\Projects\ai-club-v3-go-lang\pkg\runsteps\runsteps_test.go
package runsteps

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/greenstorm5417/openai-assistants-go/client"

	"testing"
)

func TestListRunSteps(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		responseBody   interface{}
		params         *ListRunStepsParams
	}{
		{
			name:           "successful list run steps",
			method:         "GET",
			expectedStatus: http.StatusOK,
			responseBody: ListRunStepsResponse{
				Object: "list",
				Data: []RunStep{
					{
						ID:          "step_abc123",
						Object:      "thread.run.step",
						CreatedAt:   1699063291,
						RunID:       "run_abc123",
						AssistantID: "asst_abc123",
						ThreadID:    "thread_abc123",
						Type:        "message_creation",
						Status:      "completed",
						StepDetails: StepDetails{
							Type: "message_creation",
							MessageCreation: &MessageCreation{
								MessageID: "msg_abc123",
							},
						},
						Usage: &Usage{
							PromptTokens:     123,
							CompletionTokens: 456,
							TotalTokens:      579,
						},
					},
				},
				FirstID: "step_abc123",
				LastID:  "step_abc456",
				HasMore: false,
			},
			params: &ListRunStepsParams{
				Limit:   intPtr(20),
				Order:   stringPtr("desc"),
				After:   stringPtr("step_abc123"),
				Include: []string{"step_details.tool_calls[*].file_search.results[*].content"},
			},
		},
		{
			name:           "api error",
			method:         "GET",
			expectedStatus: http.StatusBadRequest,
			responseBody: client.APIError{
				ErrorInfo: struct {
					Message string `json:"message"`
					Type    string `json:"type"`
					Param   string `json:"param"`
					Code    string `json:"code"`
				}{
					Message: "Invalid run ID",
					Type:    "invalid_request_error",
					Code:    "invalid_run_id",
				},
			},
			params: &ListRunStepsParams{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify HTTP method
				if r.Method != tt.method {
					t.Errorf("Expected method %s, got %s", tt.method, r.Method)
				}

				// Verify headers
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Error("Authorization header not set correctly")
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("Content-Type header not set correctly")
				}
				if r.Header.Get("OpenAI-Beta") != "assistants=v2" {
					t.Error("OpenAI-Beta header not set correctly")
				}

				// Respond with the desired status code and body
				w.WriteHeader(tt.expectedStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Initialize the client with the mock server URL
			c := &client.Client{
				BaseURL:    server.URL,
				APIKey:     "test-key",
				HTTPClient: server.Client(),
			}

			service := New(c)

			// Call the List method
			response, err := service.List("thread_abc123", "run_abc123", tt.params)

			if tt.expectedStatus == http.StatusOK {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				if response.Object != "list" {
					t.Errorf("Expected object 'list', got '%s'", response.Object)
				}
				if len(response.Data) != 1 {
					t.Errorf("Expected 1 run step, got %d", len(response.Data))
				}
			} else {
				if err == nil {
					t.Fatal("Expected error, got none")
				}
				// Optionally, check the error message
				apiErr, ok := err.(*client.APIError)
				if !ok {
					t.Fatalf("Expected APIError, got %T", err)
				}
				if apiErr.ErrorInfo.Code != "invalid_run_id" {
					t.Errorf("Expected error code 'invalid_run_id', got '%s'", apiErr.ErrorInfo.Code)
				}
			}
		})
	}
}

func TestGetRunStep(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		responseBody   interface{}
		params         *GetRunStepParams
	}{
		{
			name:           "successful get run step",
			method:         "GET",
			expectedStatus: http.StatusOK,
			responseBody: RunStep{
				ID:          "step_abc123",
				Object:      "thread.run.step",
				CreatedAt:   1699063291,
				RunID:       "run_abc123",
				AssistantID: "asst_abc123",
				ThreadID:    "thread_abc123",
				Type:        "message_creation",
				Status:      "completed",
				StepDetails: StepDetails{
					Type: "message_creation",
					MessageCreation: &MessageCreation{
						MessageID: "msg_abc123",
					},
				},
				Usage: &Usage{
					PromptTokens:     123,
					CompletionTokens: 456,
					TotalTokens:      579,
				},
			},
			params: &GetRunStepParams{
				Include: []string{"step_details.tool_calls[*].file_search.results[*].content"},
			},
		},
		{
			name:           "run step not found",
			method:         "GET",
			expectedStatus: http.StatusNotFound,
			responseBody: client.APIError{
				ErrorInfo: struct {
					Message string `json:"message"`
					Type    string `json:"type"`
					Param   string `json:"param"`
					Code    string `json:"code"`
				}{
					Message: "Run step not found",
					Type:    "not_found_error",
					Code:    "run_step_not_found",
				},
			},
			params: &GetRunStepParams{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify HTTP method
				if r.Method != tt.method {
					t.Errorf("Expected method %s, got %s", tt.method, r.Method)
				}

				// Verify headers
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Error("Authorization header not set correctly")
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("Content-Type header not set correctly")
				}
				if r.Header.Get("OpenAI-Beta") != "assistants=v2" {
					t.Error("OpenAI-Beta header not set correctly")
				}

				// Respond with the desired status code and body
				w.WriteHeader(tt.expectedStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Initialize the client with the mock server URL
			c := &client.Client{
				BaseURL:    server.URL,
				APIKey:     "test-key",
				HTTPClient: server.Client(),
			}

			service := New(c)

			// Call the Get method
			runStep, err := service.Get("thread_abc123", "run_abc123", "step_abc123", tt.params)

			if tt.expectedStatus == http.StatusOK {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				if runStep.ID != "step_abc123" {
					t.Errorf("Expected run step ID 'step_abc123', got '%s'", runStep.ID)
				}
				if runStep.Status != "completed" {
					t.Errorf("Expected status 'completed', got '%s'", runStep.Status)
				}
			} else {
				if err == nil {
					t.Fatal("Expected error, got none")
				}
				// Optionally, check the error message
				apiErr, ok := err.(*client.APIError)
				if !ok {
					t.Fatalf("Expected APIError, got %T", err)
				}
				if apiErr.ErrorInfo.Code != "run_step_not_found" {
					t.Errorf("Expected error code 'run_step_not_found', got '%s'", apiErr.ErrorInfo.Code)
				}
			}
		})
	}
}

// Helper functions to create pointers for test parameters
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
