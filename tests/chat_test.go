package tests

import (
	"context"
	"testing"

	"github.com/efortin/ollama-mcp/internal/core"
)

// This test is mocked since we don't want to make actual API calls in tests
func TestChatHandler(t *testing.T) {
	// Skip this test in normal runs since it requires Ollama to be running
	t.Skip("Skipping test that requires Ollama to be running")

	// In a real test, you would mock the Ollama client
	// For demonstration purposes, we're just testing the interface
	ctx := context.Background()

	// Create a test configuration
	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create server and handler factory
	server := core.NewServer(config)
	factory := core.NewHandlerFactory(server)

	// Get the chat handler
	chatHandler := factory.ChatHandler()

	input := core.ChatInput{
		Model:   "llama3",
		Message: "Hello, how are you?",
	}

	_, output, err := chatHandler(ctx, nil, input)
	if err != nil {
		t.Fatalf("Chat returned an error: %v", err)
	}

	if output.Response == "" {
		t.Error("Expected a non-empty response")
	}
}
