# OpenAI Assistants Go Client (Unofficial)

[![Go Reference](https://pkg.go.dev/badge/github.com/greenstorm5417/openai-assistants-go.svg)](https://pkg.go.dev/github.com/greenstorm5417/openai-assistants-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/greenstorm5417/openai-assistants-go)](https://goreportcard.com/report/github.com/greenstorm5417/openai-assistants-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive Go client for the OpenAI Assistants API (Beta). This library provides full support for the OpenAI Assistants API, including assistants, threads, messages, runs, and run steps.

## Features

- 🤖 Complete support for all OpenAI Assistants API endpoints
- 🧵 Full thread management capabilities
- 📝 Message handling with support for text and files
- ⚡ Real-time streaming support
- 🔧 Tool integration (Code Interpreter, File Search, Function Calling)
- 🔄 Comprehensive run and run steps management
- ✅ 100% test coverage
- 📚 Detailed documentation for each package

## Installation

```bash
go get github.com/greenstorm5417/openai-assistants-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/greenstorm5417/openai-assistants-go/internal/client"
    "github.com/greenstorm5417/openai-assistants-go/pkg/assistants"
    "github.com/greenstorm5417/openai-assistants-go/pkg/messages"
    "github.com/greenstorm5417/openai-assistants-go/pkg/threads"
)

func main() {
    // Initialize client with your API key
    c := client.NewClient(os.Getenv("OPENAI_API_KEY"))

    // Create services
    assistantService := assistants.New(c)
    threadService := threads.New(c)
    messageService := messages.New(c)

    // Create an assistant
    name := "Math Tutor"
    assistant, err := assistantService.Create(&assistants.CreateAssistantRequest{
        Model: "gpt-4",
        Name:  &name,
        Tools: []assistants.Tool{
            {Type: "code_interpreter"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a thread
    thread, err := threadService.Create(nil)
    if err != nil {
        log.Fatal(err)
    }

    // Add a message to the thread
    message, err := messageService.Create(thread.ID, &messages.CreateMessageRequest{
        Role:    "user",
        Content: "Can you help me understand how to calculate the derivative of x²?",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created assistant: %s\n", assistant.ID)
    fmt.Printf("Created thread: %s\n", thread.ID)
    fmt.Printf("Added message: %s\n", message.ID)
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
└── examples/           # Example implementations
```

## Documentation

Each package contains its own detailed documentation in its respective directory:

- [Assistants](pkg/assistants/README.md)
- [Messages](pkg/messages/README.md)
- [Runs](pkg/runs/README.md)
- [Run Steps](pkg/runsteps/README.md)
- [Threads](pkg/threads/README.md)

## Examples

Complete working examples can be found in the [examples](examples) directory:

- Basic Assistant Creation and Management
- Thread and Message Handling
- Run Management with Tool Outputs
- Streaming Responses
- Function Calling
- File Handling

## Error Handling

The client provides structured error handling for API errors:

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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This is an unofficial library and is not affiliated with OpenAI. The OpenAI Assistants API is currently in Beta, and this library may need updates as the API evolves.