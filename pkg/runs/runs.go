package runs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/greenstorm5417/openai-assistants-go/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/types"
)

// Run represents an execution run on a thread
type Run struct {
	ID                  string              `json:"id"`
	Object              string              `json:"object"`
	CreatedAt           int64               `json:"created_at"`
	ThreadID            string              `json:"thread_id"`
	AssistantID         string              `json:"assistant_id"`
	Status              string              `json:"status"`
	RequiredAction      *RequiredAction     `json:"required_action,omitempty"`
	LastError           *ErrorObject        `json:"last_error,omitempty"`
	ExpiresAt           *int64              `json:"expires_at,omitempty"`
	StartedAt           *int64              `json:"started_at,omitempty"`
	CancelledAt         *int64              `json:"cancelled_at,omitempty"`
	FailedAt            *int64              `json:"failed_at,omitempty"`
	CompletedAt         *int64              `json:"completed_at,omitempty"`
	Model               string              `json:"model"`
	Instructions        *string             `json:"instructions,omitempty"`
	Tools               []Tool              `json:"tools"`
	ToolResources       *ToolResources      `json:"tool_resources,omitempty"`
	Metadata            types.Metadata      `json:"metadata,omitempty"`
	Usage               *Usage              `json:"usage,omitempty"`
	Temperature         *float64            `json:"temperature,omitempty"`
	TopP                *float64            `json:"top_p,omitempty"`
	MaxPromptTokens     *int                `json:"max_prompt_tokens,omitempty"`
	MaxCompletionTokens *int                `json:"max_completion_tokens,omitempty"`
	TruncationStrategy  *TruncationStrategy `json:"truncation_strategy,omitempty"`
	ResponseFormat      interface{}         `json:"response_format,omitempty"`
	ToolChoice          interface{}         `json:"tool_choice,omitempty"`
	ParallelToolCalls   bool                `json:"parallel_tool_calls"`
}

// RequiredAction represents an action required to continue the run
type RequiredAction struct {
	Type              string             `json:"type"`
	SubmitToolOutputs *SubmitToolOutputs `json:"submit_tool_outputs,omitempty"`
}

// SubmitToolOutputs represents the tool outputs that need to be submitted
type SubmitToolOutputs struct {
	ToolCalls []ToolCall `json:"tool_calls"`
}

// ToolCall represents a tool call that needs output
type ToolCall struct {
	ID       string        `json:"id"`
	Type     string        `json:"type"`
	Function *FunctionCall `json:"function,omitempty"`
}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
	Output    string `json:"output,omitempty"`
}

// ErrorObject represents an error that occurred during the run
type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Usage represents the token usage for the run
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Tool represents a tool that can be used by the assistant
type Tool struct {
	Type     string        `json:"type"`
	Function *FunctionTool `json:"function,omitempty"`
}

// FunctionTool represents a function tool
type FunctionTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

// ToolResources represents resources available to tools
type ToolResources struct {
	CodeInterpreter *CodeInterpreterResources `json:"code_interpreter,omitempty"`
	FileSearch      *FileSearchResources      `json:"file_search,omitempty"`
}

// CodeInterpreterResources represents resources for code interpreter
type CodeInterpreterResources struct {
	FileIDs []string `json:"file_ids"`
}

// FileSearchResources represents resources for file search
type FileSearchResources struct {
	VectorStoreIDs []string `json:"vector_store_ids"`
}

// TruncationStrategy represents how a thread will be truncated
type TruncationStrategy struct {
	Type         string `json:"type"`
	LastMessages *int   `json:"last_messages,omitempty"`
}

// CreateRunRequest represents the request to create a run
type CreateRunRequest struct {
	AssistantID            string              `json:"assistant_id"`
	Model                  *string             `json:"model,omitempty"`
	Instructions           *string             `json:"instructions,omitempty"`
	AdditionalInstructions *string             `json:"additional_instructions,omitempty"`
	Tools                  []Tool              `json:"tools,omitempty"`
	ToolResources          *ToolResources      `json:"tool_resources,omitempty"`
	Metadata               types.Metadata      `json:"metadata,omitempty"`
	Temperature            *float64            `json:"temperature,omitempty"`
	TopP                   *float64            `json:"top_p,omitempty"`
	Stream                 bool                `json:"stream,omitempty"`
	MaxPromptTokens        *int                `json:"max_prompt_tokens,omitempty"`
	MaxCompletionTokens    *int                `json:"max_completion_tokens,omitempty"`
	TruncationStrategy     *TruncationStrategy `json:"truncation_strategy,omitempty"`
	ResponseFormat         interface{}         `json:"response_format,omitempty"`
	ToolChoice             interface{}         `json:"tool_choice,omitempty"`
	ParallelToolCalls      *bool               `json:"parallel_tool_calls,omitempty"`
}

// CreateThreadAndRunRequest represents the request to create a thread and run
type CreateThreadAndRunRequest struct {
	AssistantID string         `json:"assistant_id"`
	Thread      *ThreadRequest `json:"thread,omitempty"`
	CreateRunRequest
}

// ThreadRequest represents the thread creation part of CreateThreadAndRunRequest
type ThreadRequest struct {
	Messages []Message      `json:"messages,omitempty"`
	Metadata types.Metadata `json:"metadata,omitempty"`
}

// Message represents a message in a thread
type Message struct {
	Role     string         `json:"role"`
	Content  string         `json:"content"`
	Metadata types.Metadata `json:"metadata,omitempty"`
}

// ListRunsResponse represents the response when listing runs
type ListRunsResponse struct {
	Object  string `json:"object"`
	Data    []Run  `json:"data"`
	FirstID string `json:"first_id"`
	LastID  string `json:"last_id"`
	HasMore bool   `json:"has_more"`
}

// ListRunsParams represents the parameters for listing runs
type ListRunsParams struct {
	Limit  *int    `json:"limit,omitempty"`
	Order  *string `json:"order,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
}

// SubmitToolOutputsRequest represents the request to submit tool outputs
type SubmitToolOutputsRequest struct {
	ToolOutputs []ToolOutput `json:"tool_outputs"`
	Stream      bool         `json:"stream,omitempty"`
}

// ToolOutput represents a tool output
type ToolOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}

// RunEvent represents an event in a streaming response
type RunEvent struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// Service handles communication with the runs related methods of the OpenAI API
type Service struct {
	client *client.Client
}

// New creates a new runs service using the provided client
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// Create creates a new run
func (s *Service) Create(threadID string, req *CreateRunRequest) (*Run, error) {
	return s.createRun(fmt.Sprintf("%s/threads/%s/runs", s.client.BaseURL, threadID), req)
}

// CreateAndStream creates a new run and returns a channel of events
func (s *Service) CreateAndStream(threadID string, req *CreateRunRequest) (<-chan RunEvent, error) {
	req.Stream = true
	return s.createRunStream(fmt.Sprintf("%s/threads/%s/runs", s.client.BaseURL, threadID), req)
}

// CreateThreadAndRun creates a thread and run in one request
func (s *Service) CreateThreadAndRun(req *CreateThreadAndRunRequest) (*Run, error) {
	return s.createRun(fmt.Sprintf("%s/threads/runs", s.client.BaseURL), req)
}

// prepareRequest sets the necessary headers for a request
func (s *Service) prepareRequest(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.client.APIKey))
	req.Header.Set("OpenAI-Beta", "assistants=v2")
	req.Header.Set("Content-Type", "application/json")
}

// CreateThreadAndRunStream creates a thread and run in one request and returns a channel of events
func (s *Service) CreateThreadAndRunStream(req *CreateThreadAndRunRequest) (<-chan RunEvent, error) {
	req.Stream = true
	return s.createRunStream(fmt.Sprintf("%s/threads/runs", s.client.BaseURL), req)
}

func (s *Service) createRun(url string, req interface{}) (*Run, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("OpenAI-Beta", "assistants=v2")

	var run Run
	if err := s.client.SendRequest(httpReq, &run); err != nil {
		return nil, err
	}

	return &run, nil
}

func (s *Service) createRunStream(url string, req interface{}) (<-chan RunEvent, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Set necessary headers
	s.prepareRequest(httpReq)

	resp, err := s.client.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	events := make(chan RunEvent)
	go func() {
		defer resp.Body.Close()
		defer close(events)

		reader := bufio.NewReader(resp.Body)
		var currentEvent string

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					events <- RunEvent{Event: "error", Data: json.RawMessage(fmt.Sprintf(`{"error":"%s"}`, err.Error()))}
				}
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Parse event type
			if strings.HasPrefix(line, "event: ") {
				currentEvent = strings.TrimPrefix(line, "event: ")
				continue
			}

			// Parse data
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					events <- RunEvent{Event: "done"}
					return
				}

				events <- RunEvent{
					Event: currentEvent,
					Data:  json.RawMessage(data),
				}
			}
		}
	}()

	return events, nil
}

// List returns a list of runs for a thread
func (s *Service) List(threadID string, params *ListRunsParams) (*ListRunsResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/runs", s.client.BaseURL, threadID)
	if params != nil {
		query := make(map[string]string)
		if params.Limit != nil {
			query["limit"] = fmt.Sprintf("%d", *params.Limit)
		}
		if params.Order != nil {
			query["order"] = *params.Order
		}
		if params.After != nil {
			query["after"] = *params.After
		}
		if params.Before != nil {
			query["before"] = *params.Before
		}
		if len(query) > 0 {
			url += "?"
			for k, v := range query {
				url += fmt.Sprintf("%s=%s&", k, v)
			}
			url = url[:len(url)-1]
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var response ListRunsResponse
	if err := s.client.SendRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Get retrieves a specific run
func (s *Service) Get(threadID, runID string) (*Run, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/threads/%s/runs/%s", s.client.BaseURL, threadID, runID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var run Run
	if err := s.client.SendRequest(req, &run); err != nil {
		return nil, err
	}

	return &run, nil
}

// Modify modifies a run
func (s *Service) Modify(threadID, runID string, metadata types.Metadata) (*Run, error) {
	body, err := json.Marshal(map[string]interface{}{
		"metadata": metadata,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/threads/%s/runs/%s", s.client.BaseURL, threadID, runID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var run Run
	if err := s.client.SendRequest(req, &run); err != nil {
		return nil, err
	}

	return &run, nil
}

// SubmitToolOutputs submits outputs for tool calls
func (s *Service) SubmitToolOutputs(threadID, runID string, req *SubmitToolOutputsRequest) (*Run, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/threads/%s/runs/%s/submit_tool_outputs", s.client.BaseURL, threadID, runID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("OpenAI-Beta", "assistants=v2")

	var run Run
	if err := s.client.SendRequest(httpReq, &run); err != nil {
		return nil, err
	}

	return &run, nil
}

// SubmitToolOutputsStream submits outputs for tool calls and returns a channel of events
func (s *Service) SubmitToolOutputsStream(threadID, runID string, req *SubmitToolOutputsRequest) (<-chan RunEvent, error) {
	req.Stream = true
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/threads/%s/runs/%s/submit_tool_outputs", s.client.BaseURL, threadID, runID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Set necessary headers
	s.prepareRequest(httpReq)

	resp, err := s.client.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	events := make(chan RunEvent)
	go func() {
		defer resp.Body.Close()
		defer close(events)

		reader := bufio.NewReader(resp.Body)
		var currentEvent string

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					events <- RunEvent{Event: "error", Data: json.RawMessage(fmt.Sprintf(`{"error":"%s"}`, err.Error()))}
				}
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Parse event type
			if strings.HasPrefix(line, "event: ") {
				currentEvent = strings.TrimPrefix(line, "event: ")
				continue
			}

			// Parse data
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					events <- RunEvent{Event: "done"}
					return
				}

				events <- RunEvent{
					Event: currentEvent,
					Data:  json.RawMessage(data),
				}
			}
		}
	}()

	return events, nil
}

// Cancel cancels a run
func (s *Service) Cancel(threadID, runID string) (*Run, error) {
	fmt.Printf("Canceling run: threadID=%s, runID=%s\n", threadID, runID)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/threads/%s/runs/%s/cancel", s.client.BaseURL, threadID, runID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var run Run
	if err := s.client.SendRequest(req, &run); err != nil {
		return nil, fmt.Errorf("SendRequest failed: %w", err)
	}

	fmt.Printf("Cancel response run: %+v\n", run)

	return &run, nil
}
