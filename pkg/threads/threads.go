package threads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/types"
)

// Thread represents a thread that contains messages
type Thread struct {
	ID            string         `json:"id"`
	Object        string         `json:"object"`
	CreatedAt     int64          `json:"created_at"`
	ToolResources *ToolResources `json:"tool_resources,omitempty"`
	Metadata      types.Metadata `json:"metadata,omitempty"`
}

// ToolResources represents resources available to assistant tools in a thread
type ToolResources struct {
	CodeInterpreter *CodeInterpreterResources `json:"code_interpreter,omitempty"`
	FileSearch      *FileSearchResources      `json:"file_search,omitempty"`
}

// CodeInterpreterResources represents resources for the code interpreter tool
type CodeInterpreterResources struct {
	FileIDs []string `json:"file_ids"`
}

// FileSearchResources represents resources for the file search tool
type FileSearchResources struct {
	VectorStoreIDs []string `json:"vector_store_ids"`
}

// Message represents a message in a thread
type Message struct {
	Role        string         `json:"role"`
	Content     string         `json:"content"`
	Attachments []string       `json:"attachments,omitempty"`
	Metadata    types.Metadata `json:"metadata,omitempty"`
}

// TextContent represents text content in a message
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ImageFileContent represents an image file in a message
type ImageFileContent struct {
	Type      string `json:"type"`
	ImageFile struct {
		FileID string `json:"file_id"`
		Detail string `json:"detail,omitempty"`
	} `json:"image_file"`
}

// ImageURLContent represents an image URL in a message
type ImageURLContent struct {
	Type     string `json:"type"`
	ImageURL struct {
		URL    string `json:"url"`
		Detail string `json:"detail,omitempty"`
	} `json:"image_url"`
}

// VectorStore represents a vector store configuration
type VectorStore struct {
	FileIDs          []string          `json:"file_ids"`
	ChunkingStrategy *ChunkingStrategy `json:"chunking_strategy,omitempty"`
	Metadata         types.Metadata    `json:"metadata,omitempty"`
}

// ChunkingStrategy represents the strategy for chunking files
type ChunkingStrategy struct {
	Type   string        `json:"type,omitempty"`
	Static *StaticConfig `json:"static,omitempty"`
}

// StaticConfig represents static chunking configuration
type StaticConfig struct {
	MaxChunkSizeTokens int `json:"max_chunk_size_tokens"`
	ChunkOverlapTokens int `json:"chunk_overlap_tokens"`
}

// CreateThreadRequest represents the request to create a new thread
type CreateThreadRequest struct {
	Messages      []Message      `json:"messages,omitempty"`
	ToolResources *ToolResources `json:"tool_resources,omitempty"`
	Metadata      types.Metadata `json:"metadata,omitempty"`
}

// DeleteThreadResponse represents the response when deleting a thread
type DeleteThreadResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// Service handles communication with the threads related methods of the OpenAI API
type Service struct {
	client *client.Client
}

// New creates a new threads service using the provided client
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// Create creates a new thread
func (s *Service) Create(req *CreateThreadRequest) (*Thread, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", s.client.BaseURL+"/threads", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("OpenAI-Beta", "assistants=v2")

	var thread Thread
	if err := s.client.SendRequest(httpReq, &thread); err != nil {
		return nil, err
	}

	return &thread, nil
}

// Get retrieves a thread
func (s *Service) Get(threadID string) (*Thread, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/threads/%s", s.client.BaseURL, threadID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var thread Thread
	if err := s.client.SendRequest(req, &thread); err != nil {
		return nil, err
	}

	return &thread, nil
}

// Modify modifies a thread
func (s *Service) Modify(threadID string, toolResources *ToolResources, metadata types.Metadata) (*Thread, error) {
	body, err := json.Marshal(map[string]interface{}{
		"tool_resources": toolResources,
		"metadata":       metadata,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/threads/%s", s.client.BaseURL, threadID), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var thread Thread
	if err := s.client.SendRequest(req, &thread); err != nil {
		return nil, err
	}

	return &thread, nil
}

// Delete deletes a thread
func (s *Service) Delete(threadID string) (*DeleteThreadResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/threads/%s", s.client.BaseURL, threadID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OpenAI-Beta", "assistants=v2")

	var response DeleteThreadResponse
	if err := s.client.SendRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
