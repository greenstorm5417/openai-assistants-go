package client

import (
        "encoding/json"
        "net/http"
        "net/http/httptest"
        "testing"
)

func TestNewClient(t *testing.T) {
        apiKey := "test-key"
        client := NewClient(apiKey)

        if client.APIKey != apiKey {
                t.Errorf("Expected API key %s, got %s", apiKey, client.APIKey)
        }

        if client.BaseURL != defaultBaseURL {
                t.Errorf("Expected base URL %s, got %s", defaultBaseURL, client.BaseURL)
        }

        if client.HTTPClient == nil {
                t.Error("HTTP client should not be nil")
        }
}

func TestSendRequest(t *testing.T) {
        type TestResponse struct {
                Message string `json:"message"`
        }

        tests := []struct {
                name           string
                statusCode     int
                responseBody   interface{}
                expectedError  bool
                expectedResult *TestResponse
        }{
                {
                        name:       "successful request",
                        statusCode: http.StatusOK,
                        responseBody: TestResponse{
                                Message: "success",
                        },
                        expectedError: false,
                        expectedResult: &TestResponse{
                                Message: "success",
                        },
                },
                {
                        name:       "api error",
                        statusCode: http.StatusBadRequest,
                        responseBody: APIError{
                                ErrorInfo: struct {
                                        Message string `json:"message"`
                                        Type    string `json:"type"`
                                        Param   string `json:"param"`
                                        Code    string `json:"code"`
                                }{
                                        Message: "error message",
                                        Type:    "invalid_request_error",
                                        Code:    "invalid_api_key",
                                },
                        },
                        expectedError:  true,
                        expectedResult: nil,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                                // Check headers
                                if r.Header.Get("Authorization") != "Bearer test-key" {
                                        t.Error("Authorization header not set correctly")
                                }
                                if r.Header.Get("Content-Type") != "application/json" {
                                        t.Error("Content-Type header not set correctly")
                                }

                                // Set response
                                w.WriteHeader(tt.statusCode)
                                json.NewEncoder(w).Encode(tt.responseBody)
                        }))
                        defer server.Close()

                        client := &Client{
                                BaseURL:    server.URL,
                                APIKey:     "test-key",
                                HTTPClient: server.Client(),
                        }

                        req, err := http.NewRequest("GET", server.URL+"/test", nil)
                        if err != nil {
                                t.Fatalf("Failed to create request: %v", err)
                        }

                        var result TestResponse
                        err = client.SendRequest(req, &result)

                        if tt.expectedError && err == nil {
                                t.Error("Expected error but got none")
                        }
                        if !tt.expectedError && err != nil {
                                t.Errorf("Expected no error but got: %v", err)
                        }
                        if tt.expectedResult != nil && result.Message != tt.expectedResult.Message {
                                t.Errorf("Expected message %s, got %s", tt.expectedResult.Message, result.Message)
                        }
                })
        }
}