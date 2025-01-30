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
	"github.com/greenstorm5417/openai-assistants-go/pkg/threads"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create a new client
	c := client.NewClient(apiKey)
	assistantService := assistants.New(c)
	threadService := threads.New(c)
	messageService := messages.New(c)
	runService := runs.New(c)

	// First, create an assistant
	fmt.Println("\n=== Creating Assistant ===")
	assistant, err := createAssistant(assistantService)
	if err != nil {
		log.Fatalf("Failed to create assistant: %v", err)
	}
	fmt.Printf("Created assistant: %s\n", assistant.ID)

	// Create a thread
	fmt.Println("\n=== Creating Thread ===")
	thread, err := threadService.Create(nil)
	if err != nil {
		log.Fatalf("Failed to create thread: %v", err)
	}
	fmt.Printf("Created thread: %s\n", thread.ID)

	// Add a message to the thread
	fmt.Println("\n=== Adding Message ===")
	message, err := messageService.Create(thread.ID, &messages.CreateMessageRequest{
		Role:    "user",
		Content: "What is the weather like in San Francisco?",
	})
	if err != nil {
		log.Fatalf("Failed to create message: %v", err)
	}
	fmt.Printf("Added message: %s\n", message.ID)

	// Example 1: Create and run normally
	fmt.Println("\n=== Creating Normal Run ===")
	run, err := createNormalRun(runService, thread.ID, assistant.ID)
	if err != nil {
		log.Fatalf("Failed to create normal run: %v", err)
	}
	printRun(run)

	// Wait for run to complete or require action
	fmt.Println("\n=== Waiting for Run ===")
	run, err = waitForRun(runService, thread.ID, run.ID)
	if err != nil {
		log.Fatalf("Failed waiting for run: %v", err)
	}
	printRun(run)

	// Handle tool outputs if required
	if run.RequiredAction != nil && run.RequiredAction.Type == "submit_tool_outputs" {
		fmt.Println("\n=== Submitting Tool Outputs ===")
		if err := submitToolOutputs(runService, thread.ID, run.ID, run.RequiredAction.SubmitToolOutputs.ToolCalls); err != nil {
			log.Fatalf("Failed to submit tool outputs: %v", err)
		}

		// Wait for run to complete after submitting tool outputs
		fmt.Println("\n=== Waiting for Run After Tool Outputs ===")
		run, err = waitForRun(runService, thread.ID, run.ID)
		if err != nil {
			log.Fatalf("Failed waiting for run: %v", err)
		}
		printRun(run)
	}

	// Example 2: Create and stream run
	fmt.Println("\n=== Creating Streaming Run ===")
	if err := createStreamingRun(runService, thread.ID, assistant.ID); err != nil {
		log.Fatalf("Failed to create streaming run: %v", err)
	}

	// Example 3: Create run with function calling
	fmt.Println("\n=== Creating Function Calling Run ===")
	run, err = createFunctionRun(runService, thread.ID, assistant.ID)
	if err != nil {
		log.Fatalf("Failed to create function run: %v", err)
	}
	printRun(run)

	// Wait for run to require action
	fmt.Println("\n=== Waiting for Function Run ===")
	run, err = waitForRun(runService, thread.ID, run.ID)
	if err != nil {
		log.Fatalf("Failed waiting for run: %v", err)
	}
	printRun(run)

	// Handle tool outputs if required
	if run.RequiredAction != nil && run.RequiredAction.Type == "submit_tool_outputs" {
		fmt.Println("\n=== Submitting Tool Outputs ===")
		if err := submitToolOutputs(runService, thread.ID, run.ID, run.RequiredAction.SubmitToolOutputs.ToolCalls); err != nil {
			log.Fatalf("Failed to submit tool outputs: %v", err)
		}

		// Wait for run to complete after submitting tool outputs
		fmt.Println("\n=== Waiting for Run After Tool Outputs ===")
		run, err = waitForRun(runService, thread.ID, run.ID)
		if err != nil {
			log.Fatalf("Failed waiting for run: %v", err)
		}
		printRun(run)
	}

	// Example 5: List runs
	fmt.Println("\n=== Listing Runs ===")
	if err := listRuns(runService, thread.ID); err != nil {
		log.Fatalf("Failed to list runs: %v", err)
	}

	// Example 6: Cancel a run
	fmt.Println("\n=== Creating Run to Cancel ===")
	run, err = runService.Create(thread.ID, &runs.CreateRunRequest{
		AssistantID: assistant.ID,
	})
	if err != nil {
		log.Fatalf("Failed to create run: %v", err)
	}

	fmt.Println("\n=== Cancelling Run ===")
	run, err = runService.Cancel(thread.ID, run.ID)
	if err != nil {
		log.Printf("Failed to cancel run %s: %v", run.ID, err)
	} else {
		fmt.Printf("Successfully cancelled run %s\n", run.ID)
	}

	// Clean up
	fmt.Println("\n=== Cleaning Up ===")
	if _, err := assistantService.Delete(assistant.ID); err != nil {
		log.Printf("Failed to delete assistant %s: %v", assistant.ID, err)
	} else {
		fmt.Printf("Successfully deleted assistant %s\n", assistant.ID)
	}

	if _, err := threadService.Delete(thread.ID); err != nil {
		log.Printf("Failed to delete thread %s: %v", thread.ID, err)
	} else {
		fmt.Printf("Successfully deleted thread %s\n", thread.ID)
	}
}

func createAssistant(service *assistants.Service) (*assistants.Assistant, error) {
	name := "Weather Assistant"
	instructions := "You are a helpful assistant that provides weather information."
	model := "gpt-4"
	temp := 0.7

	return service.Create(&assistants.CreateAssistantRequest{
		Model:        model,
		Name:         &name,
		Instructions: &instructions,
		Tools: []assistants.Tool{
			{
				Type: "function",
				Function: &assistants.FunctionTool{
					Name:        "get_current_weather",
					Description: "Get the current weather in a given location",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]interface{}{
								"type":        "string",
								"description": "The city and state, e.g. San Francisco, CA",
							},
							"unit": map[string]interface{}{
								"type": "string",
								"enum": []string{"celsius", "fahrenheit"},
							},
						},
						"required": []string{"location"},
					},
				},
			},
		},
		Temperature: &temp,
	})
}

func createNormalRun(service *runs.Service, threadID, assistantID string) (*runs.Run, error) {
	return service.Create(threadID, &runs.CreateRunRequest{
		AssistantID: assistantID,
	})
}

func createStreamingRun(service *runs.Service, threadID, assistantID string) error {
	// Create and stream the run
	events, err := service.CreateAndStream(threadID, &runs.CreateRunRequest{
		AssistantID: assistantID,
		Stream:      true,
	})
	if err != nil {
		return err
	}

	var currentRun *runs.Run
	for event := range events {
		fmt.Printf("Event: %s\n", event.Event)
		if len(event.Data) > 0 {
			fmt.Printf("Data: %s\n", string(event.Data))

			// Parse run data to check for required actions
			if event.Event == "thread.run.requires_action" {
				var run runs.Run
				if err := json.Unmarshal(event.Data, &run); err != nil {
					return fmt.Errorf("failed to parse run data: %w", err)
				}
				currentRun = &run
			}
		}
	}

	// If the run requires tool outputs, submit them
	if currentRun != nil && currentRun.RequiredAction != nil && currentRun.RequiredAction.Type == "submit_tool_outputs" {
		fmt.Println("\n=== Submitting Tool Outputs (Streaming) ===")
		events, err := service.SubmitToolOutputsStream(threadID, currentRun.ID, &runs.SubmitToolOutputsRequest{
			ToolOutputs: []runs.ToolOutput{
				{
					ToolCallID: currentRun.RequiredAction.SubmitToolOutputs.ToolCalls[0].ID,
					Output:     "The weather in San Francisco is currently 72°F and sunny.",
				},
			},
			Stream: true,
		})
		if err != nil {
			return fmt.Errorf("failed to submit tool outputs: %w", err)
		}

		for event := range events {
			fmt.Printf("Event: %s\n", event.Event)
			if len(event.Data) > 0 {
				fmt.Printf("Data: %s\n", string(event.Data))
			}
		}
	}

	return nil
}

func createFunctionRun(service *runs.Service, threadID, assistantID string) (*runs.Run, error) {
	return service.Create(threadID, &runs.CreateRunRequest{
		AssistantID: assistantID,
		ToolChoice: map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name": "get_current_weather",
			},
		},
	})
}

func submitToolOutputs(service *runs.Service, threadID, runID string, toolCalls []runs.ToolCall) error {
	var toolOutputs []runs.ToolOutput
	for _, call := range toolCalls {
		if call.Function != nil {
			toolOutputs = append(toolOutputs, runs.ToolOutput{
				ToolCallID: call.ID,
				Output:     "The weather in San Francisco is currently 72°F and sunny.",
			})
		}
	}

	events, err := service.SubmitToolOutputsStream(threadID, runID, &runs.SubmitToolOutputsRequest{
		ToolOutputs: toolOutputs,
		Stream:      true,
	})
	if err != nil {
		return err
	}

	for event := range events {
		fmt.Printf("Event: %s\n", event.Event)
		if len(event.Data) > 0 {
			fmt.Printf("Data: %s\n", string(event.Data))
		}
	}

	return nil
}

func listRuns(service *runs.Service, threadID string) error {
	limit := 10
	order := "desc"
	response, err := service.List(threadID, &runs.ListRunsParams{
		Limit: &limit,
		Order: &order,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Found %d runs:\n", len(response.Data))
	for _, run := range response.Data {
		fmt.Printf("- %s (Status: %s)\n", run.ID, run.Status)
		if len(run.Tools) > 0 {
			fmt.Printf("  Tools: %d\n", len(run.Tools))
			for _, tool := range run.Tools {
				fmt.Printf("  - %s\n", tool.Type)
			}
		}
	}

	return nil
}

func waitForRun(service *runs.Service, threadID, runID string) (*runs.Run, error) {
	maxAttempts := 60
	for i := 0; i < maxAttempts; i++ {
		run, err := service.Get(threadID, runID)
		if err != nil {
			return nil, err
		}

		switch run.Status {
		case "completed", "failed", "cancelled":
			return run, nil
		case "requires_action":
			return run, nil
		case "queued", "in_progress":
			time.Sleep(time.Second)
			continue
		default:
			return nil, fmt.Errorf("unexpected run status: %s", run.Status)
		}
	}

	return nil, fmt.Errorf("timeout waiting for run to complete")
}

func printRun(r *runs.Run) {
	fmt.Printf("Run ID: %s\n", r.ID)
	fmt.Printf("Status: %s\n", r.Status)
	fmt.Printf("Created At: %d\n", r.CreatedAt)
	if r.RequiredAction != nil {
		fmt.Printf("Required Action: %s\n", r.RequiredAction.Type)
		if r.RequiredAction.SubmitToolOutputs != nil {
			fmt.Printf("Tool Calls: %d\n", len(r.RequiredAction.SubmitToolOutputs.ToolCalls))
			for _, call := range r.RequiredAction.SubmitToolOutputs.ToolCalls {
				fmt.Printf("- Tool Call ID: %s\n", call.ID)
				fmt.Printf("  Type: %s\n", call.Type)
				if call.Function != nil {
					fmt.Printf("  Function: %s\n", call.Function.Name)
					fmt.Printf("  Arguments: %s\n", call.Function.Arguments)
				}
			}
		}
	}
}
