# OpenAI Assistants API Go Client

A Go client library for the OpenAI Assistants API (Beta). This package provides a comprehensive implementation for interacting with OpenAI's Assistants API, including support for assistants, threads, messages, runs, and run steps.

## Features

- Complete implementation of OpenAI Assistants API endpoints
- Support for all major features:
  - Assistants management (CRUD operations)
  - Threads handling
  - Messages handling
  - Runs management
  - Run Steps tracking
- Streaming support for real-time updates
- Comprehensive error handling
- Full test coverage
- Type-safe implementation

## Installation

```bash
go get github.com/yourusername/your-repo-name
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/yourusername/your-repo-name/internal/client"
    "github.com/yourusername/your-repo-name/pkg/assistants"
)

func main() {
    // Initialize client
    c := client.NewClient("your-api-key")
    
    // Create assistants service
    service := assistants.New(c)
    
    // Create an assistant
    name := "Math Tutor"
    assistant, err := service.Create(&assistants.CreateAssistantRequest{
        Model:        "gpt-4",
        Name:         &name,
        Instructions: stringPtr("You are a helpful math tutor."),
        Tools: []assistants.Tool{
            {Type: "code_interpreter"},
        },
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Created assistant: %s\n", assistant.ID)
}
```

## Package Structure

```
.
├── internal/
│   └── client/         # Base HTTP client implementation
├── pkg/
│   ├── assistants/     # Assistants API implementation
│   ├── messages/       # Messages API implementation
│   ├── runs/           # Runs API implementation
│   ├── runsteps/       # Run Steps API implementation
│   ├── threads/        # Threads API implementation
│   ├── streaming/      # Streaming support
│   ├── types/          # Shared types
│   └── vectorstores/   # Vector stores implementation
└── examples/           # Example implementations
```

## Testing

To run tests:

```bash
go test ./...
```

## Examples

Check the `examples/` directory for complete examples of:
- Creating and managing assistants
- Working with threads and messages
- Handling runs and run steps
- Using tool functions
- Streaming responses

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Disclaimer

This is an unofficial library and is not affiliated with OpenAI. The OpenAI Assistants API is currently in Beta, and this library may need updates as the API evolves.