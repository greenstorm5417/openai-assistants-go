# OpenAI Threads API

This package provides a Go client for the OpenAI Threads API. Threads are conversations that can contain messages and be processed by assistants.

## Installation

```bash
go get github.com/greenstorm5417/openai-assistants-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/greenstorm5417/openai-assistants-go/internal/client"
    "github.com/greenstorm5417/openai-assistants-go/pkg/threads"
)

func main() {
    // Create a new client
    c := client.NewClient("your-api-key")
    service := threads.New(c)

    // Create a new thread
    thread, err := service.Create(nil)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created thread: %s\n", thread.ID)
}
```

## Features

- Create, retrieve, modify, and delete threads
- Support for tool resources (code interpreter, file search)
- Full metadata support
- Comprehensive error handling

## Usage Examples

### Creating Threads

```go
// Create an empty thread
thread, err := service.Create(nil)

// Create a thread with initial messages
messages := []Message{
    {
        Role:    "user",
        Content: "Hello! I'd like to learn about AI.",
    },
    {
        Role:    "user",
        Content: "Can you explain machine learning?",
    },
}

thread, err := service.Create(&threads.CreateThreadRequest{
    Messages: messages,
    Metadata: types.Metadata{
        "purpose": "educational",
        "topic":   "ai_basics",
    },
})
```

### Creating a Thread with Tool Resources

```go
// Create a thread with code interpreter
thread, err := service.Create(&threads.CreateThreadRequest{
    ToolResources: &threads.ToolResources{
        CodeInterpreter: &threads.CodeInterpreterResources{
            FileIDs: []string{"file-123", "file-456"},
        },
    },
})

// Create a thread with file search
thread, err := service.Create(&threads.CreateThreadRequest{
    ToolResources: &threads.ToolResources{
        FileSearch: &threads.FileSearchResources{
            VectorStoreIDs: []string{"vs-789"},
        },
    },
})
```

### Modifying Threads

```go
// Update thread metadata
metadata := types.Metadata{
    "status":   "in_progress",
    "priority": "high",
}

thread, err := service.Modify(threadID, &threads.ToolResources{
    CodeInterpreter: &threads.CodeInterpreterResources{
        FileIDs: []string{"file-updated"},
    },
}, metadata)
```

### Retrieving Threads

```go
// Get a specific thread
thread, err := service.Get("thread_abc123")
if err != nil {
    // Handle error
}

// Access thread properties
fmt.Printf("Thread ID: %s\n", thread.ID)
fmt.Printf("Created At: %d\n", thread.CreatedAt)
fmt.Printf("Metadata: %v\n", thread.Metadata)
```

### Deleting Threads

```go
response, err := service.Delete("thread_abc123")
if err != nil {
    // Handle error
}

if response.Deleted {
    fmt.Println("Thread deleted successfully")
}
```

## Error Handling

```go
thread, err := service.Get("nonexistent_id")
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
2. Clean up threads when they're no longer needed
3. Use meaningful metadata to organize threads
4. Handle pagination when listing threads
5. Set appropriate tool resources based on thread purpose

## Rate Limiting and Retries

The client does not automatically handle rate limiting or retries. You should implement appropriate retry logic in your application:

```go
func createThreadWithRetry(service *threads.Service) (*threads.Thread, error) {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        thread, err := service.Create(nil)
        if err == nil {
            return thread, nil
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

## Thread Lifecycle Management

```go
// Create a thread with cleanup
thread, err := service.Create(&threads.CreateThreadRequest{
    Metadata: types.Metadata{
        "auto_delete": "true",
        "ttl":        "24h",
    },
})
if err != nil {
    panic(err)
}
defer func() {
    if _, err := service.Delete(thread.ID); err != nil {
        log.Printf("Failed to delete thread %s: %v", thread.ID, err)
    }
}()

// Use the thread...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.