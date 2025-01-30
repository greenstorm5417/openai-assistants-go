package messages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/types"
)

// Message represents a message within a thread
type Message struct {
	ID                string             `json:"id"`
	Object            string             `json:"object"`
	CreatedAt         int64              `json:"created_at"`
	ThreadID          string             `json:"thread_id"`
	Status            string             `json:"status,omitempty"`
	IncompleteDetails *IncompleteDetails `json:"incomplete_details,omitempty"`
	CompletedAt       *int64             `json:"completed_at,omitempty"`
	IncompleteAt      *int64             `json:"incomplete_at,omitempty"`
	Role              string             `json:"role"`
	Content           []Content          `json:"content"`
	AssistantID       *string            `json:"assistant_id,omitempty"`
	RunID             *string            `json:"run_id,omitempty"`
	Attachments       []Attachment       `json:"attachments,omitempty"`
	Metadata          types.Metadata     `json:"metadata,omitempty"`
}

// IncompleteDetails contains details about why a message is incomplete
type IncompleteDetails struct {
	Reason string `json:"reason"`
}

// Content represents a content part in a message
type Content struct {
	Type      string     `json:"type"`
	Text      *Text      `json:"text,omitempty"`
	ImageURL  *ImageURL  `json:"image_url,omitempty"`
	ImageFile *ImageFile `json:"image_file,omitempty"`
}

// Text represents text content
type Text struct {
	Value       string       `json:"value"`
	Annotations []Annotation `json:"annotations"`
}

// ImageURL represents an image URL in content
type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// ImageFile represents an image file in content
type ImageFile struct {
	FileID string `json:"file_id"`
	Detail string `json:"detail,omitempty"`
}

// Annotation represents an annotation in text content
type Annotation struct {
	Type string `json:"type"`
}

// Attachment represents a file attached to a message
type Attachment struct {
	FileID string `json:"file_id"`
	Tools  []Tool `json:"tools,omitempty"`
}

// Tool represents a tool that can use an attachment
type Tool struct {
	Type string `json:"type"`
}

// CreateMessageRequest represents the request to create a new message
type CreateMessageRequest struct {
	Role        string         `json:"role"`
	Content     interface{}    `json:"content"`
	Attachments []Attachment   `json:"attachments,omitempty"`
	Metadata    types.Metadata `json:"metadata,omitempty"`
}

// ListMessagesResponse represents the response when listing messages
type ListMessagesResponse struct {
	Object  string    `json:"object"`
	Data    []Message `json:"data"`
	FirstID string    `json:"first_id"`
	LastID  string    `json:"last_id"`
	HasMore bool      `json:"has_more"`
}

// ListMessagesParams represents the parameters for listing messages
type ListMessagesParams struct {
	Limit  *int    `json:"limit,omitempty"`
	Order  *string `json:"order,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
	RunID  *string `json:"run_id,omitempty"`
}

// DeleteMessageResponse represents the response when deleting a message
type DeleteMessageResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// Service handles communication with the messages related methods of the OpenAI API
type Service struct {
	client *client.Client
}

// New creates a new messages service using the provided client
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// Create creates a new message in a thread
func (s *Service) Create(threadID string, req *CreateMessageRequest) (*Message, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/threads/%s/messages", s.client.BaseURL, threadID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("OpenAI-Beta", "assistants=v2")

	var message Message
	if err := s.client.SendRequest(httpReq, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// List returns a list of messages for a thread
func (s *Service) List(threadID string, params *ListMessagesParams) (*ListMessagesResponse, error) {
	url := fmt.Sprintf("%s/threads/%s/messages", s.client.BaseURL, threadID)
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
		if params.RunID != nil {
			query["run_id"] = *params.RunID
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

	var response ListMessagesResponse
	if err := s.client.SendRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Get retrieves a specific message
func (s *Service) Get(threadID, messageID string) (*Message, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/threads/%s/messages/%s", s.client.BaseURL, threadID, messageID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var message Message
	if err := s.client.SendRequest(req, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// Modify modifies a message's metadata
func (s *Service) Modify(threadID, messageID string, metadata types.Metadata) (*Message, error) {
	body, err := json.Marshal(map[string]interface{}{
		"metadata": metadata,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/threads/%s/messages/%s", s.client.BaseURL, threadID, messageID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var message Message
	if err := s.client.SendRequest(req, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// Delete deletes a message
func (s *Service) Delete(threadID, messageID string) (*DeleteMessageResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/threads/%s/messages/%s", s.client.BaseURL, threadID, messageID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var response DeleteMessageResponse
	if err := s.client.SendRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
