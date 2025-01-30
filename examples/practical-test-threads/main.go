package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/threads"
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
	service := threads.New(c)

	// Example 1: Create a simple thread
	fmt.Println("\n=== Creating Simple Thread ===")
	simpleThread, err := createSimpleThread(service)
	if err != nil {
		log.Fatalf("Failed to create simple thread: %v", err)
	}
	printThread(simpleThread)

	// Example 2: Create a thread with initial messages
	fmt.Println("\n=== Creating Thread with Messages ===")
	threadWithMessages, err := createThreadWithMessages(service)
	if err != nil {
		log.Fatalf("Failed to create thread with messages: %v", err)
	}
	printThread(threadWithMessages)

	// Example 3: Create a thread with tool resources
	fmt.Println("\n=== Creating Thread with Tool Resources ===")
	threadWithTools, err := createThreadWithTools(service)
	if err != nil {
		log.Fatalf("Failed to create thread with tools: %v", err)
	}
	printThread(threadWithTools)

	// Example 4: Modify a thread
	fmt.Println("\n=== Modifying Thread ===")
	modifiedThread, err := modifyThread(service, simpleThread.ID)
	if err != nil {
		log.Fatalf("Failed to modify thread: %v", err)
	}
	printThread(modifiedThread)

	// Example 5: Retrieve a thread
	fmt.Println("\n=== Retrieving Thread ===")
	retrievedThread, err := service.Get(simpleThread.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve thread: %v", err)
	}
	printThread(retrievedThread)

	// Example 6: Delete threads
	fmt.Println("\n=== Deleting Threads ===")
	threadsToDelete := []string{
		simpleThread.ID,
		threadWithMessages.ID,
		threadWithTools.ID,
	}

	for _, id := range threadsToDelete {
		if err := deleteThread(service, id); err != nil {
			log.Printf("Failed to delete thread %s: %v", id, err)
		}
	}
}

func createSimpleThread(service *threads.Service) (*threads.Thread, error) {
	return service.Create(&threads.CreateThreadRequest{
		Metadata: types.Metadata{
			"purpose": "general_discussion",
			"topic":   "AI and Machine Learning",
		},
	})
}

func createThreadWithMessages(service *threads.Service) (*threads.Thread, error) {
	messages := []threads.Message{
		{
			Role:    "user",
			Content: "Hello! I'd like to learn about artificial intelligence.",
		},
		{
			Role:    "user",
			Content: "Can you explain how neural networks work?",
		},
	}

	return service.Create(&threads.CreateThreadRequest{
		Messages: messages,
		Metadata: types.Metadata{
			"purpose": "educational",
			"topic":   "neural_networks",
		},
	})
}

func createThreadWithTools(service *threads.Service) (*threads.Thread, error) {
	return service.Create(&threads.CreateThreadRequest{
		ToolResources: &threads.ToolResources{
			CodeInterpreter: &threads.CodeInterpreterResources{
				FileIDs: []string{},
			},
			// FileSearch resources can be added when you have a vector store ID
		},
		Metadata: types.Metadata{
			"purpose": "code_analysis",
			"project": "data_science",
		},
	})
}

func modifyThread(service *threads.Service, threadID string) (*threads.Thread, error) {
	toolResources := &threads.ToolResources{
		CodeInterpreter: &threads.CodeInterpreterResources{
			FileIDs: []string{},
		},
	}

	metadata := types.Metadata{
		"purpose":  "advanced_analysis",
		"priority": "high",
		"status":   "in_progress",
	}

	return service.Modify(threadID, toolResources, metadata)
}

func deleteThread(service *threads.Service, threadID string) error {
	response, err := service.Delete(threadID)
	if err != nil {
		return fmt.Errorf("failed to delete thread: %w", err)
	}

	if response.Deleted {
		fmt.Printf("Successfully deleted thread %s\n", threadID)
	} else {
		fmt.Printf("Failed to delete thread %s\n", threadID)
	}

	return nil
}

func printThread(t *threads.Thread) {
	// Convert thread to pretty JSON
	jsonData, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal thread: %v", err)
		return
	}

	fmt.Printf("Thread Details:\n%s\n", string(jsonData))
}
