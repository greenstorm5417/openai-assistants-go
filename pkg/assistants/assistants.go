package assistants

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/types"
)

type Assistant struct {
	ID             string         `json:"id"`
	Object         string         `json:"object"`
	CreatedAt      int64          `json:"created_at"`
	Name           *string        `json:"name,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Model          string         `json:"model"`
	Instructions   *string        `json:"instructions,omitempty"`
	Tools          []Tool         `json:"tools"`
	ToolResources  *ToolResources `json:"tool_resources,omitempty"`
	Metadata       types.Metadata `json:"metadata,omitempty"`
	Temperature    *float64       `json:"temperature,omitempty"`
	TopP           *float64       `json:"top_p,omitempty"`
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`
}

type Tool struct {
	Type     string        `json:"type"`
	Function *FunctionTool `json:"function,omitempty"`
}

type FunctionTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

type ToolResources struct {
	FileSearch *FileSearchResources `json:"file_search,omitempty"`
}

type FileSearchResources struct {
	VectorStoreIDs []string `json:"vector_store_ids"`
}

type ResponseFormat string

type JSONSchema struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Schema      any    `json:"schema"`
	Strict      *bool  `json:"strict,omitempty"`
}

type CreateAssistantRequest struct {
	Model          string         `json:"model"`
	Name           *string        `json:"name,omitempty"`
	Description    *string        `json:"description,omitempty"`
	Instructions   *string        `json:"instructions,omitempty"`
	Tools          []Tool         `json:"tools,omitempty"`
	ToolResources  *ToolResources `json:"tool_resources,omitempty"`
	Metadata       types.Metadata `json:"metadata,omitempty"`
	Temperature    *float64       `json:"temperature,omitempty"`
	TopP           *float64       `json:"top_p,omitempty"`
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`
}

type ListAssistantsResponse struct {
	Object  string      `json:"object"`
	Data    []Assistant `json:"data"`
	FirstID string      `json:"first_id"`
	LastID  string      `json:"last_id"`
	HasMore bool        `json:"has_more"`
}

type ListAssistantsParams struct {
	Limit  *int    `json:"limit,omitempty"`
	Order  *string `json:"order,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
}

type DeleteAssistantResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// Service handles communication with the assistants related methods of the OpenAI API.
type Service struct {
	client *client.Client
}

// New creates a new assistants service using the provided client.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// Create creates a new assistant.
func (s *Service) Create(req *CreateAssistantRequest) (*Assistant, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", s.client.BaseURL+"/assistants", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("OpenAI-Beta", "assistants=v2")

	var assistant Assistant
	if err := s.client.SendRequest(httpReq, &assistant); err != nil {
		return nil, err
	}

	return &assistant, nil
}

// List returns a list of assistants.
func (s *Service) List(params *ListAssistantsParams) (*ListAssistantsResponse, error) {
	url := s.client.BaseURL + "/assistants"
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
		// Add query parameters to URL
		if len(query) > 0 {
			url += "?"
			for k, v := range query {
				url += fmt.Sprintf("%s=%s&", k, v)
			}
			url = url[:len(url)-1] // Remove trailing &
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var response ListAssistantsResponse
	if err := s.client.SendRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Get retrieves an assistant.
func (s *Service) Get(assistantID string) (*Assistant, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/assistants/%s", s.client.BaseURL, assistantID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var assistant Assistant
	if err := s.client.SendRequest(req, &assistant); err != nil {
		return nil, err
	}

	return &assistant, nil
}

// Modify modifies an existing assistant.
func (s *Service) Modify(assistantID string, req *CreateAssistantRequest) (*Assistant, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/assistants/%s", s.client.BaseURL, assistantID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("OpenAI-Beta", "assistants=v2")

	var assistant Assistant
	if err := s.client.SendRequest(httpReq, &assistant); err != nil {
		return nil, err
	}

	return &assistant, nil
}

// Delete deletes an assistant.
func (s *Service) Delete(assistantID string) (*DeleteAssistantResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/assistants/%s", s.client.BaseURL, assistantID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var response DeleteAssistantResponse
	if err := s.client.SendRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
