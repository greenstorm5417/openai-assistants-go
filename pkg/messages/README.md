# OpenAI Messages API

This package provides a Go client for the OpenAI Messages API. Messages are the building blocks of conversations within threads.

## Installation

```bash
go get github.com/workspace
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/greenstorm5417/openai-assistants-go/internal/client"
    "github.com/greenstorm5417/openai-assistants-go/pkg/messages"
    "github.com/greenstorm5417/openai-assistants-go/pkg/threads"
)

func main() {
    // Create a new client
    c := client.NewClient("your-api-key")
    messageService := messages.New(c)
    threadService := threads.New(c)

    // First, create a thread
    thread, err := threadService.Create(nil)
    if err != nil {
        panic(err)
    }

    // Create a message in the thread
    message, err := messageService.Create(thread.ID, &messages.CreateMessageRequest{
        Role:    "user",
        Content: "Hello! I'd like to learn about AI.",
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created message: %s\n", message.ID)
}
```

## Features

- Create, retrieve, modify, and delete messages
- Support for different content types (text, images)
- File attachments
- Metadata management
- Comprehensive error handling

## Usage Examples

### Creating Messages

```go
// Create a simple text message
message, err := service.Create(threadID, &messages.CreateMessageRequest{
    Role:    "user",
    Content: "Hello! I'd like to learn about AI.",
})

// Create a message with metadata
message, err := service.Create(threadID, &messages.CreateMessageRequest{
    Role:    "user",
    Content: "What are the main branches of AI?",
    Metadata: types.Metadata{
        "importance": "high",
        "category":   "ai_fundamentals",
        "topic":      "overview",
    },
})

// Create a message with file attachments
message, err := service.Create(threadID, &messages.CreateMessageRequest{
    Role:    "user",
    Content: "Here's my dataset for analysis.",
    Attachments: []messages.Attachment{
        {
            FileID: "file-123",
            Tools: []messages.Tool{
                {Type: "code_interpreter"},
            },
        },
    },
})
```

### Working with Message Content

```go
// Text content with annotations
content := []messages.Content{
    {
        Type: "text",
        Text: &messages.Text{
            Value: "Check out this example code.",
            Annotations: []messages.Annotation{
                {Type: "code"},
            },
        },
    },
}

// Image content via URL
content := []messages.Content{
    {
        Type: "image_url",
        ImageURL: &messages.ImageURL{
            URL:    "https://example.com/image.jpg",
            Detail: "high",
        },
    },
}

// Image content via file
content := []messages.Content{
    {
        Type: "image_file",
        ImageFile: &messages.ImageFile{
            FileID: "file-123",
            Detail: "low",
        },
    },
}
```

### Listing Messages

```go
// List all messages in a thread
messages, err := service.List(threadID, nil)

// List with pagination
limit := 20
order := "desc"
messages, err := service.List(threadID, &messages.ListMessagesParams{
    Limit: &limit,
    Order: &order,
})

// Get next page
after := messages.LastID
nextPage, err := service.List(threadID, &messages.ListMessagesParams{
    Limit: &limit,
    After: &after,
})

// Filter by run ID
runID := "run_123"
messages, err := service.List(threadID, &messages.ListMessagesParams{
    RunID: &runID,
})
```

### Modifying Messages

```go
// Update message metadata
metadata := types.Metadata{
    "reviewed": "true",
    "status":   "approved",
    "tags":     "important,follow-up",
}

message, err := service.Modify(threadID, messageID, metadata)
```

### Retrieving Messages

```go
// Get a specific message
message, err := service.Get(threadID, messageID)
if err != nil {
    // Handle error
}

// Access message properties
fmt.Printf("Message ID: %s\n", message.ID)
fmt.Printf("Role: %s\n", message.Role)
if len(message.Content) > 0 {
    if message.Content[0].Text != nil {
        fmt.Printf("Content: %s\n", message.Content[0].Text.Value)
    }
}
```

### Deleting Messages

```go
response, err := service.Delete(threadID, messageID)
if err != nil {
    // Handle error
}

if response.Deleted {
    fmt.Println("Message deleted successfully")
}
```

## Error Handling

```go
message, err := service.Get(threadID, "nonexistent_id")
if err != nil {
    switch e := err.(type) {
    case *client.APIError:
        fmt.Printf("API error: %s (type: %s, code: %s)\n", 
            e.ErrorInfo.Message, 
            e.ErrorInfo.Type, 
            e.ErrorInfo.Code)
    default:
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## Best Practices

1. Always check for errors after API calls
2. Use appropriate content types for different kinds of data
3. Set meaningful metadata for message organization
4. Handle pagination when listing messages
5. Clean up messages when they're no longer needed

## Message Content Guidelines

1. Text Content:
   - Keep messages concise and clear
   - Use annotations when appropriate
   - Consider using markdown for formatting

2. Image Content:
   - Use appropriate detail level (high/low) based on needs
   - Consider bandwidth and token usage
   - Prefer file references over URLs for persistence

3. File Attachments:
   - Specify appropriate tools for file processing
   - Clean up files when no longer needed
   - Consider file size limits

## Rate Limiting and Retries

The client does not automatically handle rate limiting or retries. You should implement appropriate retry logic in your application:

```go
func createMessageWithRetry(service *messages.Service, threadID string, req *messages.CreateMessageRequest) (*messages.Message, error) {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        message, err := service.Create(threadID, req)
        if err == nil {
            return message, nil
        }
        
        if apiErr, ok := err.(*client.APIError); ok && apiErr.ErrorInfo.Type == "rate_limit_exceeded" {
            time.Sleep(time.Second * time.Duration(i+1))
            continue
        }
        
        return nil, err
    }
    return nil, fmt.Errorf("max retries exceeded")
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.