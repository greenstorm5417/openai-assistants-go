# Runs Package

The `runs` package is a comprehensive Go wrapper for managing execution runs within threads using the OpenAI API. It provides functionalities to create, manage, list, and cancel runs, as well as handle required actions such as submitting tool outputs. This package abstracts the complexities of interacting directly with the OpenAI API, offering a streamlined and developer-friendly interface.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Initialization](#initialization)
  - [Creating a Run](#creating-a-run)
  - [Listing Runs](#listing-runs)
  - [Retrieving a Specific Run](#retrieving-a-specific-run)
  - [Modifying a Run](#modifying-a-run)
  - [Submitting Tool Outputs](#submitting-tool-outputs)
  - [Cancelling a Run](#cancelling-a-run)
- [Example](#example)
- [Error Handling](#error-handling)
- [Contributing](#contributing)
- [License](#license)

---

## Installation

To integrate the `runs` package into your Go project, ensure you have Go installed and set up properly. Then, use `go get` to fetch the package:

```bash
go get github.com/greenstorm5417/openai-assistants-go/pkg/runs
```

**Note**: Replace `github.com/greenstorm5417/openai-assistants-go/pkg/runs` with the actual repository path if different.

---

## Usage

### Initialization

Before utilizing the `runs` package, initialize the API client and the runs service.

```go
package main

import (
	"log"
	"os"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/runs"
)

func main() {
	// Retrieve API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Initialize the API client
	c := client.NewClient(apiKey)

	// Initialize the runs service
	runService := runs.New(c)

	// Now, runService can be used to manage runs
}
```

### Creating a Run

To create a new run within a specific thread:

```go
import (
	"log"
	"github.com/greenstorm5417/openai-assistants-go/pkg/runs"
)

// Assuming runService is already initialized as shown above

func createRun(runService *runs.Service, threadID, assistantID string) {
	run, err := runService.Create(threadID, &runs.CreateRunRequest{
		AssistantID: assistantID,
		Model:       stringPtr("gpt-4"),
		Instructions: func() *string {
			instr := "Provide a detailed summary of the run steps involved in managing run steps."
			return &instr
		}(),
		Tools: []runs.Tool{
			{
				Type: "function",
				Function: &runs.FunctionTool{
					Name:        "summarize_steps",
					Description: "Summarize the run steps for testing purposes.",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"detail_level": map[string]interface{}{
								"type":        "string",
								"description": "Level of detail for the summary",
								"enum":        []string{"brief", "detailed"},
							},
						},
						"required": []string{"detail_level"},
					},
				},
			},
		},
		Temperature: func() *float64 {
			temp := 0.5
			return &temp
		}(),
		TopP: func() *float64 {
			tp := 0.9
			return &tp
		}(),
		Metadata: map[string]interface{}{
			"test_step": "create_run",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create run: %v", err)
	}
	log.Printf("Run Created: ID=%s, Status=%s\n", run.ID, run.Status)
}

func stringPtr(s string) *string {
	return &s
}
```

### Listing Runs

Retrieve a list of runs associated with a thread.

```go
func listRuns(runService *runs.Service, threadID string) {
	limit := 10
	order := "desc"
	response, err := runService.List(threadID, &runs.ListRunsParams{
		Limit: &limit,
		Order: &order,
	})
	if err != nil {
		log.Fatalf("Failed to list runs: %v", err)
	}

	log.Printf("Found %d runs:\n", len(response.Data))
	for _, run := range response.Data {
		log.Printf("- %s (Status: %s)\n", run.ID, run.Status)
		if len(run.Tools) > 0 {
			log.Printf("  Tools: %d\n", len(run.Tools))
			for _, tool := range run.Tools {
				log.Printf("  - %s\n", tool.Type)
			}
		}
	}
}
```

### Retrieving a Specific Run

Fetch details of a specific run using its ID.

```go
func getRun(runService *runs.Service, threadID, runID string) {
	run, err := runService.Get(threadID, runID)
	if err != nil {
		log.Fatalf("Failed to get run: %v", err)
	}
	log.Printf("Run ID: %s, Status: %s\n", run.ID, run.Status)
}
```

### Modifying a Run

Update the metadata of an existing run.

```go
func modifyRun(runService *runs.Service, threadID, runID string) {
	metadata := map[string]interface{}{
		"updated_field": "new_value",
	}
	run, err := runService.Modify(threadID, runID, metadata)
	if err != nil {
		log.Fatalf("Failed to modify run: %v", err)
	}
	log.Printf("Run Modified: ID=%s, New Metadata: %v\n", run.ID, run.Metadata)
}
```

### Submitting Tool Outputs

When a run requires additional tool outputs, submit them using the `SubmitToolOutputs` method.

```go
func submitToolOutputs(runService *runs.Service, threadID, runID string, toolCalls []runs.ToolCall) {
	var toolOutputs []runs.ToolOutput
	for _, call := range toolCalls {
		if call.Function != nil {
			toolOutputs = append(toolOutputs, runs.ToolOutput{
				ToolCallID: call.ID,
				Output:     "The weather in San Francisco is currently 72°F and sunny.",
			})
		}
	}

	updatedRun, err := runService.SubmitToolOutputs(threadID, runID, &runs.SubmitToolOutputsRequest{
		ToolOutputs: toolOutputs,
		Stream:      false, // Set to true if you want to handle streaming responses
	})
	if err != nil {
		log.Fatalf("Failed to submit tool outputs: %v", err)
	}

	log.Printf("Submitted Tool Outputs. Updated Run Status: %s\n", updatedRun.Status)
}
```

### Cancelling a Run

Cancel an ongoing run.

```go
func cancelRun(runService *runs.Service, threadID, runID string) {
	run, err := runService.Cancel(threadID, runID)
	if err != nil {
		log.Fatalf("Failed to cancel run: %v", err)
	}
	log.Printf("Run Cancelled: ID=%s, Status=%s\n", run.ID, run.Status)
}
```

---

## Example

Below is a comprehensive example demonstrating how to use the `runs` package in a practical scenario. This example covers creating an assistant, setting up a thread, adding messages, initiating a run, handling required actions, listing run steps, and performing cleanup.

### **File Path:** `practical-test-runsteps/main.go`

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"encoding/json"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/assistants"
	"github.com/greenstorm5417/openai-assistants-go/pkg/messages"
	"github.com/greenstorm5417/openai-assistants-go/pkg/runs"
	"github.com/greenstorm5417/openai-assistants-go/pkg/runsteps"
	"github.com/greenstorm5417/openai-assistants-go/pkg/threads"
	"github.com/greenstorm5417/openai-assistants-go/pkg/types"
)

func main() {
	// Retrieve API key from environment variable for security
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Initialize the API client
	c := client.NewClient(apiKey)

	// Initialize services
	assistantService := assistants.New(c)
	threadService := threads.New(c)
	messageService := messages.New(c)
	runService := runs.New(c)
	runStepService := runsteps.New(c)

	// Step 1: Create a New Assistant
	fmt.Println("\n=== Creating a New Assistant ===")
	assistant, err := createAssistant(assistantService)
	if err != nil {
		log.Fatalf("Failed to create assistant: %v", err)
	}
	fmt.Printf("Assistant Created: ID=%s, Name=%s\n", assistant.ID, *assistant.Name)

	// Step 2: Create a New Thread
	fmt.Println("\n=== Creating a New Thread ===")
	thread, err := threadService.Create(&threads.CreateThreadRequest{
		Metadata: types.Metadata{
			"purpose": "practical_test",
			"topic":   "Run Steps Testing",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create thread: %v", err)
	}
	fmt.Printf("Thread Created: ID=%s\n", thread.ID)

	// Step 3: Add Messages to the Thread
	fmt.Println("\n=== Adding Messages to the Thread ===")
	messagesToAdd := []struct {
		Role    string
		Content string
	}{
		{
			Role:    "user",
			Content: "Hello Assistant! Let's test the run steps.",
		},
		{
			Role:    "user",
			Content: "Can you summarize the steps involved in creating and managing run steps?",
		},
	}

	for _, msg := range messagesToAdd {
		createdMsg, err := messageService.Create(thread.ID, &messages.CreateMessageRequest{
			Role:    msg.Role,
			Content: msg.Content,
			Metadata: types.Metadata{
				"test_case": "run_steps_practical_test",
			},
		})
		if err != nil {
			log.Fatalf("Failed to create message: %v", err)
		}
		fmt.Printf("Message Added: ID=%s, Role=%s\n", createdMsg.ID, createdMsg.Role)
	}

	// Step 4: Create a Run
	fmt.Println("\n=== Creating a Run ===")
	run, err := runService.Create(thread.ID, &runs.CreateRunRequest{
		AssistantID: assistant.ID,
		Model:       stringPtr("gpt-4"),
		Instructions: func() *string {
			instr := "Provide a detailed summary of the run steps involved in managing run steps."
			return &instr
		}(),
		Tools: []runs.Tool{
			{
				Type: "function",
				Function: &runs.FunctionTool{
					Name:        "summarize_steps",
					Description: "Summarize the run steps for testing purposes.",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"detail_level": map[string]interface{}{
								"type":        "string",
								"description": "Level of detail for the summary",
								"enum":        []string{"brief", "detailed"},
							},
						},
						"required": []string{"detail_level"},
					},
				},
			},
		},
		Temperature: func() *float64 {
			temp := 0.5
			return &temp
		}(),
		TopP: func() *float64 {
			tp := 0.9
			return &tp
		}(),
		Metadata: types.Metadata{
			"test_step": "create_run",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create run: %v", err)
	}
	fmt.Printf("Run Created: ID=%s, Status=%s\n", run.ID, run.Status)

	// Step 5: Wait for Run Completion or Requires Action
	fmt.Println("\n=== Waiting for Run to Complete ===")
	run, err = waitForRunCompletion(runService, thread.ID, run.ID, 60*time.Second)
	if err != nil {
		log.Fatalf("Error waiting for run completion: %v", err)
	}
	fmt.Printf("Run Status after Waiting: ID=%s, Status=%s\n", run.ID, run.Status)

	// Step 6: Handle Requires Action (If Any)
	if run.Status == "requires_action" && run.RequiredAction != nil && run.RequiredAction.Type == "submit_tool_outputs" {
		fmt.Println("\n=== Handling 'requires_action': Submitting Tool Outputs ===")
		err = handleRequiresAction(runService, thread.ID, run.ID, run.RequiredAction)
		if err != nil {
			log.Fatalf("Failed to handle requires_action: %v", err)
		}

		// Wait again for the run to complete after submitting tool outputs
		fmt.Println("\n=== Waiting for Run to Complete After Submitting Tool Outputs ===")
		run, err = waitForRunCompletion(runService, thread.ID, run.ID, 60*time.Second)
		if err != nil {
			log.Fatalf("Error waiting for run completion after submitting tool outputs: %v", err)
		}
		fmt.Printf("Run Completed: ID=%s, Status=%s\n", run.ID, run.Status)
	} else {
		fmt.Println("\n=== Run Completed Without Requiring Action ===")
	}

	// Step 7: List Run Steps
	fmt.Println("\n=== Listing Run Steps ===")
	runSteps, err := runStepService.List(thread.ID, run.ID, &runsteps.ListRunStepsParams{
		Limit: intPtr(10),
		Order: stringPtr("asc"),
		// Include is omitted or set to a supported value
		Include: []string{"step_details.tool_calls[*].file_search.results[*].content"}, // Ensure this is supported
	})
	if err != nil {
		log.Fatalf("Failed to list run steps: %v", err)
	}

	fmt.Printf("Total Run Steps: %d\n", len(runSteps.Data))
	for _, step := range runSteps.Data {
		fmt.Printf("- Step ID: %s, Type: %s, Status: %s\n", step.ID, step.Type, step.Status)
		if step.StepDetails.ToolCalls != nil {
			for _, toolCall := range step.StepDetails.ToolCalls {
				fmt.Printf("  - Tool Call ID: %s, Function: %s\n", toolCall.ID, toolCall.Function.Name)
			}
		}
	}

	// Step 8: Delete the Assistant (Cleanup)
	fmt.Println("\n=== Deleting the Assistant ===")
	deleteResp, err := assistantService.Delete(assistant.ID)
	if err != nil {
		log.Fatalf("Failed to delete assistant: %v", err)
	}
	if deleteResp.Deleted {
		fmt.Printf("Assistant Deleted: ID=%s\n", deleteResp.ID)
	} else {
		fmt.Printf("Assistant Deletion Failed: ID=%s\n", deleteResp.ID)
	}

	// Step 9: Delete the Thread (Cleanup)
	fmt.Println("\n=== Deleting the Thread ===")
	_, err = threadService.Delete(thread.ID)
	if err != nil {
		log.Fatalf("Failed to delete thread: %v", err)
	}
	fmt.Printf("Thread Deleted: ID=%s\n", thread.ID)
}

// createAssistant initializes a new assistant with predefined settings.
func createAssistant(service *assistants.Service) (*assistants.Assistant, error) {
	name := "Run Steps Tester"
	description := "Assistant for testing run steps functionalities."
	instructions := "You assist in testing run steps by summarizing and managing them effectively."
	temperature := 0.7
	topP := 0.9

	req := &assistants.CreateAssistantRequest{
		Model:        "gpt-4",
		Name:         &name,
		Description:  &description,
		Instructions: &instructions,
		Tools: []assistants.Tool{
			{
				Type: "function",
				Function: &assistants.FunctionTool{
					Name:        "test_function",
					Description: "A test function for run steps.",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"param1": map[string]interface{}{
								"type":        "string",
								"description": "A test parameter.",
							},
						},
						"required": []string{"param1"},
					},
				},
			},
		},
		Temperature:    &temperature,
		TopP:           &topP,
		ResponseFormat: "auto",
		Metadata: types.Metadata{
			"test_case": "run_steps_practical_test",
		},
	}

	assistant, err := service.Create(req)
	if err != nil {
		return nil, err
	}
	return assistant, nil
}

// waitForRunCompletion polls the run status until it completes, fails, is cancelled, or times out.
func waitForRunCompletion(service *runs.Service, threadID, runID string, timeout time.Duration) (*runs.Run, error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			run, err := service.Get(threadID, runID)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Current Run Status: %s\n", run.Status)
			if run.Status == "completed" || run.Status == "failed" || run.Status == "cancelled" || run.Status == "requires_action" {
				return run, nil
			}
		case <-timeoutChan:
			return nil, fmt.Errorf("timeout waiting for run to complete")
		}
	}
}

// handleRequiresAction handles the 'requires_action' status by submitting tool outputs.
func handleRequiresAction(service *runs.Service, threadID, runID string, action *runs.RequiredAction) error {
	if action.SubmitToolOutputs == nil {
		return fmt.Errorf("no submit_tool_outputs found in required_action")
	}

	var toolOutputs []runs.ToolOutput
	for _, toolCall := range action.SubmitToolOutputs.ToolCalls {
		// Prepare the output for each tool call.
		// In a real scenario, you would generate or retrieve the appropriate output.
		// Here, we'll use a dummy output for demonstration purposes.
		output := fmt.Sprintf("Dummy output for tool call %s", toolCall.ID)
		toolOutputs = append(toolOutputs, runs.ToolOutput{
			ToolCallID: toolCall.ID,
			Output:     output,
		})
	}

	req := &runs.SubmitToolOutputsRequest{
		ToolOutputs: toolOutputs,
		Stream:      false, // Set to true if you want to handle streaming responses
	}

	updatedRun, err := service.SubmitToolOutputs(threadID, runID, req)
	if err != nil {
		return fmt.Errorf("failed to submit tool outputs: %w", err)
	}

	fmt.Printf("Submitted Tool Outputs. Updated Run Status: %s\n", updatedRun.Status)
	return nil
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
```

**Explanation:**

1. **Assistant Creation**: Initializes a new assistant with predefined settings, including tools and metadata.
2. **Thread Creation**: Sets up a new thread with specific metadata.
3. **Adding Messages**: Simulates user interactions by adding messages to the thread.
4. **Run Creation**: Initiates a run associated with the assistant and thread.
5. **Polling for Completion**: Waits for the run to complete or require action.
6. **Handling Required Actions**: If the run requires additional tool outputs, it submits them.
7. **Listing Run Steps**: Retrieves and displays the steps involved in the run.
8. **Cleanup**: Deletes the created assistant and thread to ensure no residual data remains.

---

## Error Handling

The `runs` package provides robust error handling mechanisms. Errors returned by the package can be due to:

- **API Errors**: Issues related to the OpenAI API, such as invalid parameters or authentication failures.
- **Network Errors**: Connectivity issues preventing communication with the API.
- **Parsing Errors**: Failures in decoding API responses.

**Best Practices:**

- **Check for Specific Errors**: Use type assertions to handle different error types distinctly.
  
  ```go
  if err != nil {
      if apiErr, ok := err.(*client.APIError); ok {
          log.Fatalf("API Error: %s", apiErr.Error())
      } else {
          log.Fatalf("Unexpected Error: %v", err)
      }
  }
  ```

- **Implement Retries for Transient Errors**: Incorporate retry logic with exponential backoff for handling temporary failures.

- **Validate Inputs**: Ensure that all required parameters are provided and correctly formatted before making API calls.

---

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

1. **Fork the Repository**: Click the "Fork" button on the repository's GitHub page.
2. **Clone Your Fork**:
   ```bash
   git clone https://github.com/your_username/runs.git
   ```
3. **Create a Branch**:
   ```bash
   git checkout -b feature/YourFeature
   ```
4. **Commit Your Changes**:
   ```bash
   git commit -m "Add feature XYZ"
   ```
5. **Push to Your Fork**:
   ```bash
   git push origin feature/YourFeature
   ```
6. **Open a Pull Request**: Navigate to the original repository and click "New Pull Request".

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Contact

For any questions or feedback, please reach out to [your.email@example.com](mailto:your.email@example.com).

---

# Detailed Explanation

### Introduction

The README begins by introducing the `runs` package, explaining its purpose as a Go wrapper for managing runs within threads using the OpenAI API. It highlights the functionalities it provides and the abstraction it offers over the direct API interactions.

### Table of Contents

A table of contents is included for easy navigation, listing all the major sections of the README.

### Installation

Provides instructions on how to install the `runs` package using `go get`. It also includes a note to replace the repository path if it's different.

### Usage

This section is broken down into multiple subsections, each detailing a specific functionality of the `runs` package.

- **Initialization**: Demonstrates how to initialize the API client and the runs service. It emphasizes the importance of retrieving the API key securely from environment variables.

- **Creating a Run**: Shows how to create a new run within a specific thread. It includes setting up the `CreateRunRequest` with necessary parameters like `AssistantID`, `Model`, `Instructions`, `Tools`, `Temperature`, `TopP`, and `Metadata`. Helper functions like `stringPtr` are used to create pointers for string values.

- **Listing Runs**: Explains how to retrieve a list of runs associated with a thread. It covers setting parameters like `Limit` and `Order` and iterating over the returned runs to display their details.

- **Retrieving a Specific Run**: Details how to fetch a specific run using its `runID`.

- **Modifying a Run**: Provides guidance on how to update the metadata of an existing run.

- **Submitting Tool Outputs**: Addresses scenarios where a run requires additional tool outputs. It demonstrates how to prepare and submit these outputs using the `SubmitToolOutputs` method.

- **Cancelling a Run**: Illustrates how to cancel an ongoing run.

### Example

This comprehensive example is adapted from the practical test code provided by the user. It walks through the entire lifecycle of managing a run, including:

1. **Assistant Creation**: Setting up an assistant with predefined tools and instructions.
2. **Thread Creation**: Establishing a new thread for message exchanges.
3. **Adding Messages**: Simulating user interactions by adding messages to the thread.
4. **Run Creation**: Initiating a run and setting up instructions and tools.
5. **Polling for Completion**: Waiting for the run to complete or require action.
6. **Handling Required Actions**: Submitting tool outputs if the run requires further actions.
7. **Listing Run Steps**: Retrieving and displaying the steps involved in the run.
8. **Cleanup**: Deleting the created assistant and thread to ensure no residual data remains.

The example includes code snippets and explanations for each step, providing a clear and practical guide on how to use the `runs` package effectively.

### Error Handling

Discusses the types of errors that can occur when using the `runs` package, including API errors, network issues, and parsing failures. It offers best practices for handling these errors, such as checking for specific error types, implementing retries, and validating inputs.

### Contributing

Encourages contributions to the project and outlines the steps to contribute, including forking the repository, cloning it, creating a branch, committing changes, pushing to the fork, and opening a pull request.

### License

Specifies the licensing for the project, directing users to the LICENSE file.

### Contact

Provides contact information for users to reach out with questions or feedback.

### Conclusion

The README is designed to be comprehensive yet concise, offering all the necessary information someone would need to understand, install, and use the `runs` package effectively. It incorporates example code from the practical test, ensuring that users have a real-world reference to guide their implementation.