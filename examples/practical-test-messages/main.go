package main

import (
	"fmt"
	"log"
	"os"

	"github.com/greenstorm5417/openai-assistants-go/internal/client"
	"github.com/greenstorm5417/openai-assistants-go/pkg/messages"
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
	messageService := messages.New(c)
	threadService := threads.New(c)

	// First, create a thread to work with
	fmt.Println("\n=== Creating Thread ===")
	thread, err := threadService.Create(nil)
	if err != nil {
		log.Fatalf("Failed to create thread: %v", err)
	}
	fmt.Printf("Created thread: %s\n", thread.ID)

	// Example 1: Create a simple text message
	fmt.Println("\n=== Creating Text Message ===")
	textMessage, err := createTextMessage(messageService, thread.ID)
	if err != nil {
		log.Fatalf("Failed to create text message: %v", err)
	}
	printMessage(textMessage)

	// Example 2: Create a message with metadata
	fmt.Println("\n=== Creating Message with Metadata ===")
	metadataMessage, err := createMessageWithMetadata(messageService, thread.ID)
	if err != nil {
		log.Fatalf("Failed to create message with metadata: %v", err)
	}
	printMessage(metadataMessage)

	// Example 3: List messages in the thread
	fmt.Println("\n=== Listing Messages ===")
	if err := listMessages(messageService, thread.ID); err != nil {
		log.Fatalf("Failed to list messages: %v", err)
	}

	// Example 4: Modify a message
	fmt.Println("\n=== Modifying Message ===")
	modifiedMessage, err := modifyMessage(messageService, thread.ID, textMessage.ID)
	if err != nil {
		log.Fatalf("Failed to modify message: %v", err)
	}
	printMessage(modifiedMessage)

	// Example 5: Retrieve a specific message
	fmt.Println("\n=== Retrieving Message ===")
	retrievedMessage, err := messageService.Get(thread.ID, textMessage.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve message: %v", err)
	}
	printMessage(retrievedMessage)

	// Example 6: Delete messages
	fmt.Println("\n=== Deleting Messages ===")
	messagesToDelete := []string{
		textMessage.ID,
		metadataMessage.ID,
	}

	for _, id := range messagesToDelete {
		if err := deleteMessage(messageService, thread.ID, id); err != nil {
			log.Printf("Failed to delete message %s: %v", id, err)
		}
	}

	// Clean up: Delete the thread
	fmt.Println("\n=== Cleaning Up ===")
	if _, err := threadService.Delete(thread.ID); err != nil {
		log.Printf("Failed to delete thread %s: %v", thread.ID, err)
	} else {
		fmt.Printf("Successfully deleted thread %s\n", thread.ID)
	}
}

func createTextMessage(service *messages.Service, threadID string) (*messages.Message, error) {
	return service.Create(threadID, &messages.CreateMessageRequest{
		Role:    "user",
		Content: "Hello! I'd like to learn about artificial intelligence.",
	})
}

func createMessageWithMetadata(service *messages.Service, threadID string) (*messages.Message, error) {
	return service.Create(threadID, &messages.CreateMessageRequest{
		Role:    "user",
		Content: "What are the main branches of AI?",
		Metadata: types.Metadata{
			"importance": "high",
			"category":   "ai_fundamentals",
			"topic":      "overview",
		},
	})
}

func listMessages(service *messages.Service, threadID string) error {
	// List with pagination
	limit := 2
	order := "desc"
	firstPage, err := service.List(threadID, &messages.ListMessagesParams{
		Limit: &limit,
		Order: &order,
	})
	if err != nil {
		return fmt.Errorf("failed to get first page: %w", err)
	}

	fmt.Printf("First page - %d messages:\n", len(firstPage.Data))
	for _, msg := range firstPage.Data {
		fmt.Printf("- %s (Role: %s)\n", msg.ID, msg.Role)
		if len(msg.Content) > 0 {
			if msg.Content[0].Text != nil {
				fmt.Printf("  Content: %s\n", msg.Content[0].Text.Value)
			}
		}
	}

	// If there are more pages, get the next page
	if firstPage.HasMore {
		after := firstPage.LastID
		secondPage, err := service.List(threadID, &messages.ListMessagesParams{
			Limit: &limit,
			Order: &order,
			After: &after,
		})
		if err != nil {
			fmt.Printf("\nFailed to get second page: %v\n", err)
			return nil
		}

		fmt.Printf("\nSecond page - %d messages:\n", len(secondPage.Data))
		for _, msg := range secondPage.Data {
			fmt.Printf("- %s (Role: %s)\n", msg.ID, msg.Role)
			if len(msg.Content) > 0 {
				if msg.Content[0].Text != nil {
					fmt.Printf("  Content: %s\n", msg.Content[0].Text.Value)
				}
			}
		}
	}

	return nil
}

func modifyMessage(service *messages.Service, threadID, messageID string) (*messages.Message, error) {
	metadata := types.Metadata{
		"modified":   "true",
		"importance": "medium",
		"reviewed":   "true",
	}

	return service.Modify(threadID, messageID, metadata)
}

func deleteMessage(service *messages.Service, threadID, messageID string) error {
	response, err := service.Delete(threadID, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if response.Deleted {
		fmt.Printf("Successfully deleted message %s\n", messageID)
	} else {
		fmt.Printf("Failed to delete message %s\n", messageID)
	}

	return nil
}

func printMessage(m *messages.Message) {
	fmt.Printf("Message ID: %s\n", m.ID)
	fmt.Printf("Role: %s\n", m.Role)
	if len(m.Content) > 0 {
		if m.Content[0].Text != nil {
			fmt.Printf("Content: %s\n", m.Content[0].Text.Value)
		}
	}
	if len(m.Metadata) > 0 {
		fmt.Printf("Metadata: %+v\n", m.Metadata)
	}
	fmt.Printf("Created At: %d\n", m.CreatedAt)
}
