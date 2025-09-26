package tests

import (
	"context"
	"testing"

	"mcp-hello/internal/core"
)

// This test is mocked since we don't want to make actual API calls in tests
func TestChatWithOllama(t *testing.T) {
	// Skip this test in normal runs since it requires Ollama to be running
	t.Skip("Skipping test that requires Ollama to be running")

	// In a real test, you would mock the Ollama client
	// For demonstration purposes, we're just testing the interface
	ctx := context.Background()
	input := core.ChatInput{
		Model:   "llama3",
		Message: "Hello, how are you?",
	}

	_, output, err := core.ChatWithOllama(ctx, nil, input)
	if err != nil {
		t.Fatalf("Code returned an error: %v", err)
	}

	if output.Response == "" {
		t.Error("Expected a non-empty response")
	}
}
