// C:\Users\sduss\Desktop\Projects\ai-club-v3-go-lang\pkg\runsteps\runsteps.go
package runsteps

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/types"
)

// RunStep represents a step in the execution of a run.
type RunStep struct {
	ID          string         `json:"id"`
	Object      string         `json:"object"`
	CreatedAt   int64          `json:"created_at"`
	AssistantID string         `json:"assistant_id"`
	ThreadID    string         `json:"thread_id"`
	RunID       string         `json:"run_id"`
	Type        string         `json:"type"`
	Status      string         `json:"status"`
	StepDetails StepDetails    `json:"step_details"`
	LastError   *ErrorObject   `json:"last_error,omitempty"`
	ExpiresAt   *int64         `json:"expires_at,omitempty"`
	CancelledAt *int64         `json:"cancelled_at,omitempty"`
	FailedAt    *int64         `json:"failed_at,omitempty"`
	CompletedAt *int64         `json:"completed_at,omitempty"`
	Metadata    types.Metadata `json:"metadata"`
	Usage       *Usage         `json:"usage,omitempty"`
}

// StepDetails contains specific details about the run step based on its type.
type StepDetails struct {
	Type            string           `json:"type"`
	MessageCreation *MessageCreation `json:"message_creation,omitempty"`
	ToolCalls       []ToolCall       `json:"tool_calls,omitempty"`
}

// MessageCreation holds details for message creation steps.
type MessageCreation struct {
	MessageID string `json:"message_id"`
}

// ToolCall represents a call to a tool within a run step.
type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents the function call details within a tool call.
type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
	Output    string `json:"output"`
}

// ErrorObject represents an error that occurred during the run step.
type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Usage holds the token usage statistics for the run step.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ListRunStepsResponse represents the response from listing run steps.
type ListRunStepsResponse struct {
	Object  string    `json:"object"`
	Data    []RunStep `json:"data"`
	FirstID string    `json:"first_id"`
	LastID  string    `json:"last_id"`
	HasMore bool      `json:"has_more"`
}

// ListRunStepsParams represents the query parameters for listing run steps.
type ListRunStepsParams struct {
	Limit   *int     `json:"limit,omitempty"`
	Order   *string  `json:"order,omitempty"` // "asc" or "desc"
	After   *string  `json:"after,omitempty"`
	Before  *string  `json:"before,omitempty"`
	Include []string `json:"include,omitempty"` // e.g., "step_details.tool_calls[*].file_search.results[*].content"
}

// GetRunStepParams represents the query parameters for retrieving a run step.
type GetRunStepParams struct {
	Include []string `json:"include,omitempty"` // e.g., "step_details.tool_calls[*].file_search.results[*].content"
}

// Service handles communication with the run steps related methods of the OpenAI API.
type Service struct {
	client *client.Client
}

// New creates a new runsteps service using the provided client.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// List retrieves a list of run steps belonging to a specific run.
func (s *Service) List(threadID, runID string, params *ListRunStepsParams) (*ListRunStepsResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/runs/%s/steps", s.client.BaseURL, threadID, runID)
	if params != nil {
		query := make([]string, 0)
		if params.Limit != nil {
			query = append(query, fmt.Sprintf("limit=%d", *params.Limit))
		}
		if params.Order != nil {
			query = append(query, fmt.Sprintf("order=%s", *params.Order))
		}
		if params.After != nil {
			query = append(query, fmt.Sprintf("after=%s", *params.After))
		}
		if params.Before != nil {
			query = append(query, fmt.Sprintf("before=%s", *params.Before))
		}
		for _, include := range params.Include {
			query = append(query, fmt.Sprintf("include[]=%s", include))
		}
		if len(query) > 0 {
			url += "?" + strings.Join(query, "&")
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var response ListRunStepsResponse
	if err := s.client.SendRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Get retrieves a specific run step by its ID.
func (s *Service) Get(threadID, runID, stepID string, params *GetRunStepParams) (*RunStep, error) {
	url := fmt.Sprintf("%s/threads/%s/runs/%s/steps/%s", s.client.BaseURL, threadID, runID, stepID)
	if params != nil && len(params.Include) > 0 {
		query := make([]string, 0)
		for _, include := range params.Include {
			query = append(query, fmt.Sprintf("include[]=%s", include))
		}
		if len(query) > 0 {
			url += "?" + strings.Join(query, "&")
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var runStep RunStep
	if err := s.client.SendRequest(req, &runStep); err != nil {
		return nil, err
	}

	return &runStep, nil
}
