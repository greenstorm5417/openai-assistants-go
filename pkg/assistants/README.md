# OpenAI Assistants API

This package provides a Go client for the OpenAI Assistants API. It allows you to create, manage, and interact with assistants that can use various tools and models.

## Installation

```bash
go get github.com/greenstorm5417/openai-assistants-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/greenstorm5417/openai-assistants-go/client"
    "github.com/greenstorm5417/openai-assistants-go/pkg/assistants"
)

func main() {
    // Create a new client
    c := client.NewClient("your-api-key")
    
    // Initialize the assistants service
    service := assistants.New(c)

    // Create a new assistant
    name := "Math Tutor"
    instructions := "You are a personal math tutor. Write and run Python code to solve math problems."
    
    assistant, err := service.Create(&assistants.CreateAssistantRequest{
        Model:        "gpt-4",
        Name:         &name,
        Instructions: &instructions,
        Tools: []assistants.Tool{
            {Type: "code_interpreter"},
        },
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created assistant: %s\n", assistant.ID)
}
```

## Features

- Create, retrieve, modify, and delete assistants
- List all assistants with pagination
- Support for all assistant tools (code_interpreter, file_search, function)
- Full type safety and error handling

## Usage Examples

### Creating an Assistant

```go
// Create a simple assistant
name := "Math Tutor"
assistant, err := service.Create(&assistants.CreateAssistantRequest{
    Model: "gpt-4",
    Name:  &name,
})

// Create an assistant with tools
instructions := "You are a coding tutor."
assistant, err := service.Create(&assistants.CreateAssistantRequest{
    Model:        "gpt-4",
    Name:         &name,
    Instructions: &instructions,
    Tools: []assistants.Tool{
        {Type: "code_interpreter"},
    },
})

// Create an assistant with custom function
functionTool := assistants.FunctionTool{
    Name:        "calculate_sum",
    Description: "Calculate the sum of two numbers",
    Parameters: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "a": map[string]interface{}{
                "type": "number",
                "description": "First number",
            },
            "b": map[string]interface{}{
                "type": "number",
                "description": "Second number",
            },
        },
        "required": []string{"a", "b"},
    },
}

assistant, err := service.Create(&assistants.CreateAssistantRequest{
    Model: "gpt-4",
    Tools: []assistants.Tool{
        {
            Type:     "function",
            Function: &functionTool,
        },
    },
})
```

### Listing Assistants

```go
// List all assistants
assistants, err := service.List(nil)

// List with pagination
limit := 10
order := "desc"
assistants, err := service.List(&assistants.ListAssistantsParams{
    Limit: &limit,
    Order: &order,
})

// Get next page using After
after := assistants.LastID
nextPage, err := service.List(&assistants.ListAssistantsParams{
    Limit: &limit,
    After: &after,
})
```

### Retrieving an Assistant

```go
assistant, err := service.Get("asst_abc123")
if err != nil {
    // Handle error
}
```

### Modifying an Assistant

```go
name := "Updated Name"
instructions := "Updated instructions"
assistant, err := service.Modify("asst_abc123", &assistants.CreateAssistantRequest{
    Name:         &name,
    Instructions: &instructions,
})
```

### Deleting an Assistant

```go
response, err := service.Delete("asst_abc123")
if err != nil {
    // Handle error
}
if response.Deleted {
    fmt.Println("Assistant deleted successfully")
}
```

### Error Handling

```go
assistant, err := service.Get("nonexistent_id")
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

## Advanced Configuration

### Response Format

```go
// Create an assistant with JSON response format
jsonFormat := &assistants.ResponseFormat{
    Type: "json_object",
}

assistant, err := service.Create(&assistants.CreateAssistantRequest{
    Model:          "gpt-4",
    ResponseFormat: jsonFormat,
})

// Create an assistant with structured output
schema := map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
        "name": map[string]interface{}{
            "type": "string",
            "description": "The name of the person",
        },
        "age": map[string]interface{}{
            "type": "number",
            "description": "The age of the person",
        },
    },
    "required": []string{"name", "age"},
}

strict := true
jsonFormat := &assistants.ResponseFormat{
    Type: "json_schema",
    JSONSchema: &assistants.JSONSchema{
        Name:        "PersonInfo",
        Description: "Information about a person",
        Schema:      schema,
        Strict:      &strict,
    },
}

assistant, err := service.Create(&assistants.CreateAssistantRequest{
    Model:          "gpt-4",
    ResponseFormat: jsonFormat,
})
```

### Tool Resources

```go
// Create an assistant with file search tool and vector store resources
assistant, err := service.Create(&assistants.CreateAssistantRequest{
    Model: "gpt-4",
    Tools: []assistants.Tool{
        {Type: "file_search"},
    },
    ToolResources: &assistants.ToolResources{
        FileSearch: &assistants.FileSearchResources{
            VectorStoreIDs: []string{"vs_abc123"},
        },
    },
})
```

## Best Practices

1. Always check for errors after API calls
2. Use pointers for optional fields in requests
3. Set appropriate temperature and top_p values for your use case
4. Include clear instructions for your assistant
5. Use appropriate tools based on your assistant's purpose

## Rate Limiting and Retries

The client does not automatically handle rate limiting or retries. You should implement appropriate retry logic in your application:

```go
import "time"

func createAssistantWithRetry(service *assistants.Service, req *assistants.CreateAssistantRequest) (*assistants.Assistant, error) {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        assistant, err := service.Create(req)
        if err == nil {
            return assistant, nil
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