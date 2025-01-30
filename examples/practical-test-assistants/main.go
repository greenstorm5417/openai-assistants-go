package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/greenstorm5417/openai-assistants-go/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/assistants"
	"github.com/greenstorm5417/openai-assistants-go/pkg/types"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create a new client
	c := client.NewClient(apiKey)
	service := assistants.New(c)

	// Example 1: Create a simple code interpreter assistant
	fmt.Println("\n=== Creating Code Interpreter Assistant ===")
	codeAssistant, err := createCodeInterpreterAssistant(service)
	if err != nil {
		log.Fatalf("Failed to create code interpreter assistant: %v", err)
	}
	printAssistant(codeAssistant)

	// Example 2: Create an assistant with function calling
	fmt.Println("\n=== Creating Function Calling Assistant ===")
	functionAssistant, err := createFunctionAssistant(service)
	if err != nil {
		log.Fatalf("Failed to create function assistant: %v", err)
	}
	printAssistant(functionAssistant)

	// Example 3: Create an assistant with file search and JSON response format
	fmt.Println("\n=== Creating File Search Assistant with JSON Response ===")
	fileSearchAssistant, err := createFileSearchAssistant(service)
	if err != nil {
		log.Fatalf("Failed to create file search assistant: %v", err)
	}
	printAssistant(fileSearchAssistant)

	// Example 4: List all assistants with pagination
	fmt.Println("\n=== Listing Assistants ===")
	if err := listAssistantsWithPagination(service); err != nil {
		log.Fatalf("Failed to list assistants: %v", err)
	}

	// Example 5: Modify an assistant
	fmt.Println("\n=== Modifying Assistant ===")
	modifiedAssistant, err := modifyAssistant(service, codeAssistant.ID)
	if err != nil {
		log.Fatalf("Failed to modify assistant: %v", err)
	}
	printAssistant(modifiedAssistant)

	// Example 6: Retrieve an assistant
	fmt.Println("\n=== Retrieving Assistant ===")
	retrievedAssistant, err := service.Get(codeAssistant.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve assistant: %v", err)
	}
	printAssistant(retrievedAssistant)

	// Example 7: Delete assistants
	fmt.Println("\n=== Deleting Assistants ===")
	assistantsToDelete := []string{
		codeAssistant.ID,
		functionAssistant.ID,
		fileSearchAssistant.ID,
	}

	for _, id := range assistantsToDelete {
		if err := deleteAssistant(service, id); err != nil {
			log.Printf("Failed to delete assistant %s: %v", id, err)
		}
	}
}

func createCodeInterpreterAssistant(service *assistants.Service) (*assistants.Assistant, error) {
	name := "Python Math Tutor"
	instructions := "You are a math tutor who helps students understand mathematical concepts by writing and running Python code to demonstrate solutions."
	temp := 0.7
	topP := 0.9

	return service.Create(&assistants.CreateAssistantRequest{
		Model:        "gpt-4-1106-preview",
		Name:         &name,
		Instructions: &instructions,
		Tools: []assistants.Tool{
			{Type: "code_interpreter"},
		},
		Temperature: &temp,
		TopP:        &topP,
		Metadata: types.Metadata{
			"expertise": "mathematics",
			"level":     "advanced",
		},
	})
}

func createFunctionAssistant(service *assistants.Service) (*assistants.Assistant, error) {
	name := "Calendar Assistant"
	instructions := "You are a calendar management assistant that helps users schedule and manage their appointments."

	// Define a custom function for scheduling appointments
	scheduleFn := assistants.FunctionTool{
		Name:        "schedule_appointment",
		Description: "Schedule a new appointment in the calendar",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Title of the appointment",
				},
				"start_time": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "Start time of the appointment in ISO 8601 format",
				},
				"duration_minutes": map[string]interface{}{
					"type":        "integer",
					"description": "Duration of the appointment in minutes",
				},
				"attendees": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "string"},
					"description": "List of attendee email addresses",
				},
			},
			"required": []string{"title", "start_time", "duration_minutes"},
		},
	}

	return service.Create(&assistants.CreateAssistantRequest{
		Model:        "gpt-4-1106-preview",
		Name:         &name,
		Instructions: &instructions,
		Tools: []assistants.Tool{
			{
				Type:     "function",
				Function: &scheduleFn,
			},
		},
	})
}

func createFileSearchAssistant(service *assistants.Service) (*assistants.Assistant, error) {
	name := "Document Assistant"
	instructions := "You are a document assistant that helps users find and analyze information in their files."

	jsonFormat := assistants.ResponseFormat("auto")

	return service.Create(&assistants.CreateAssistantRequest{
		Model:        "gpt-4-1106-preview",
		Name:         &name,
		Instructions: &instructions,
		Tools: []assistants.Tool{
			{Type: "file_search"},
		},
		// ToolResources can be added when you have a vector store ID
		ResponseFormat: jsonFormat,
	})
}

func listAssistantsWithPagination(service *assistants.Service) error {
	// First page
	limit := 2
	order := "desc"
	firstPage, err := service.List(&assistants.ListAssistantsParams{
		Limit: &limit,
		Order: &order,
	})
	if err != nil {
		return fmt.Errorf("failed to get first page: %w", err)
	}

	fmt.Printf("First page - %d assistants:\n", len(firstPage.Data))
	for _, asst := range firstPage.Data {
		fmt.Printf("- %s (%s)\n", asst.ID, *asst.Name)
	}

	// If there are more pages, get the next page
	if firstPage.HasMore {
		after := firstPage.LastID
		secondPage, err := service.List(&assistants.ListAssistantsParams{
			Limit: &limit,
			Order: &order,
			After: &after,
		})
		if err != nil {
			fmt.Printf("\nFailed to get second page: %v\n", err)
			return nil
		}

		fmt.Printf("\nSecond page - %d assistants:\n", len(secondPage.Data))
		for _, asst := range secondPage.Data {
			fmt.Printf("- %s (%s)\n", asst.ID, *asst.Name)
		}
	}

	return nil
}

func modifyAssistant(service *assistants.Service, assistantID string) (*assistants.Assistant, error) {
	name := "Enhanced Math Tutor"
	instructions := "You are an advanced math tutor who helps students understand complex mathematical concepts through interactive Python code examples and visualizations."
	temp := 0.5 // More focused responses

	return service.Modify(assistantID, &assistants.CreateAssistantRequest{
		Name:         &name,
		Instructions: &instructions,
		Temperature:  &temp,
		Metadata: types.Metadata{
			"expertise": "mathematics",
			"level":     "expert",
			"updated":   time.Now().Format(time.RFC3339),
		},
	})
}

func deleteAssistant(service *assistants.Service, assistantID string) error {
	response, err := service.Delete(assistantID)
	if err != nil {
		return fmt.Errorf("failed to delete assistant: %w", err)
	}

	if response.Deleted {
		fmt.Printf("Successfully deleted assistant %s\n", assistantID)
	} else {
		fmt.Printf("Failed to delete assistant %s\n", assistantID)
	}

	return nil
}

func printAssistant(a *assistants.Assistant) {
	// Convert assistant to pretty JSON
	jsonData, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal assistant: %v", err)
		return
	}

	fmt.Printf("Assistant Details:\n%s\n", string(jsonData))
}
