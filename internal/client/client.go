package client

import (
        "encoding/json"
        "fmt"
        "io"
        "net/http"
)

const (
        defaultBaseURL = "https://api.openai.com/v1"
)

// Client is the OpenAI API client
type Client struct {
        BaseURL    string
        APIKey     string
        HTTPClient *http.Client
}

// APIError represents an error response from the OpenAI API
type APIError struct {
        ErrorInfo struct {
                Message string `json:"message"`
                Type    string `json:"type"`
                Param   string `json:"param"`
                Code    string `json:"code"`
        } `json:"error"`
}

func (e *APIError) Error() string {
        return fmt.Sprintf("OpenAI API error: %s (type: %s, code: %s)", e.ErrorInfo.Message, e.ErrorInfo.Type, e.ErrorInfo.Code)
}

// NewClient creates a new OpenAI API client
func NewClient(apiKey string) *Client {
        return &Client{
                BaseURL:    defaultBaseURL,
                APIKey:     apiKey,
                HTTPClient: &http.Client{},
        }
}

// SendRequest sends an HTTP request and decodes the response into v
func (c *Client) SendRequest(req *http.Request, v interface{}) error {
        // Set common headers
        req.Header.Set("Authorization", "Bearer "+c.APIKey)
        req.Header.Set("Content-Type", "application/json")

        // Send request
        resp, err := c.HTTPClient.Do(req)
        if err != nil {
                return fmt.Errorf("failed to send request: %w", err)
        }
        defer resp.Body.Close()

        // Read response body
        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return fmt.Errorf("failed to read response body: %w", err)
        }

        // Check for error response
        if resp.StatusCode != http.StatusOK {
                var apiErr APIError
                if err := json.Unmarshal(body, &apiErr); err != nil {
                        return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
                }
                return &apiErr
        }

        // Decode response
        if err := json.Unmarshal(body, v); err != nil {
                return fmt.Errorf("failed to decode response: %w", err)
        }

        return nil
}