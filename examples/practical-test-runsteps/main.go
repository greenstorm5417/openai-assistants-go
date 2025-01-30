package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/greenstorm5417/openai-assistants-go/client"
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

// Helper functions to create pointers for test parameters
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
