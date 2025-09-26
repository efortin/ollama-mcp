package tests

import (
	"context"
	"os"
	"testing"

	"github.com/efortin/ollama-mcp/internal/core"
)

func TestListModels(t *testing.T) {
	// Skip if OLLAMA_TEST is not set to true
	if os.Getenv("OLLAMA_TEST") != "true" {
		t.Skip("Skipping test that requires Ollama to be running. Set OLLAMA_TEST=true to run this test.")
	}

	// Create a context
	ctx := context.Background()

	// Load config and create handler
	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	server := core.NewServer(config)
	factory := core.NewHandlerFactory(server)
	listModelsHandler := factory.ListModelsHandler()

	// Call the ListModels function
	_, output, err := listModelsHandler(ctx, nil, core.ListModelsInput{})
	if err != nil {
		t.Fatalf("ListModels returned an error: %v", err)
	}

	// Check that we got some models
	t.Logf("Found %d models", len(output.Models))
	for i, model := range output.Models {
		t.Logf("Model %d: %s", i+1, model.Name)
	}
}

func TestModelInfo(t *testing.T) {
	// Skip if OLLAMA_TEST is not set to true
	if os.Getenv("OLLAMA_TEST") != "true" {
		t.Skip("Skipping test that requires Ollama to be running. Set OLLAMA_TEST=true to run this test.")
	}

	// Skip if OLLAMA_MODEL is not set
	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		t.Skip("Skipping test that requires a model name. Set OLLAMA_MODEL to run this test.")
	}

	// Create a context
	ctx := context.Background()

	// Load config and create handler
	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	server := core.NewServer(config)
	factory := core.NewHandlerFactory(server)
	modelInfoHandler := factory.ModelInfoHandler()

	// Call the ModelInfo function
	_, output, err := modelInfoHandler(ctx, nil, core.ModelInfoInput{Name: modelName})
	if err != nil {
		t.Fatalf("ModelInfo returned an error: %v", err)
	}

	// Check that we got the model info
	t.Logf("Model: %s", output.Name)
	t.Logf("License: %s", output.License)
	t.Logf("ModifiedAt: %s", output.ModifiedAt)
	t.Logf("Template: %s", output.Template)
}

func TestPullModel(t *testing.T) {
	// Skip if OLLAMA_TEST is not set to true
	if os.Getenv("OLLAMA_TEST") != "true" {
		t.Skip("Skipping test that requires Ollama to be running. Set OLLAMA_TEST=true to run this test.")
	}

	// Skip if OLLAMA_PULL_MODEL is not set
	modelName := os.Getenv("OLLAMA_PULL_MODEL")
	if modelName == "" {
		t.Skip("Skipping test that requires a model name to pull. Set OLLAMA_PULL_MODEL to run this test.")
	}

	// Create a context
	ctx := context.Background()

	// Load config and create handler
	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	server := core.NewServer(config)
	factory := core.NewHandlerFactory(server)
	pullModelHandler := factory.PullModelHandler()

	// Call the PullModel function with no progress to make the test faster
	_, output, err := pullModelHandler(ctx, nil, core.PullModelInput{
		Name:       modelName,
		NoProgress: true,
	})
	if err != nil {
		t.Fatalf("PullModel returned an error: %v", err)
	}

	// Check that we got a status
	t.Logf("Status: %s", output.Status)
	t.Logf("Message: %s", output.Message)
}

func TestOllamaChatCalculation(t *testing.T) {
	// Skip if OLLAMA_TEST is not set to true
	if os.Getenv("OLLAMA_TEST") != "true" {
		t.Skip("Skipping test that requires Ollama to be running. Set OLLAMA_TEST=true to run this test.")
	}

	// Skip if OLLAMA_MODEL is not set
	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		t.Skip("Skipping test that requires a model name. Set OLLAMA_MODEL to run this test.")
	}

	// Create a context
	ctx := context.Background()

	// Load config and create handler
	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	server := core.NewServer(config)
	factory := core.NewHandlerFactory(server)
	chatHandler := factory.ChatHandler()

	// Call the ChatWithOllama function
	_, output, err := chatHandler(ctx, nil, core.ChatInput{
		Model:   modelName,
		Message: "Calculate 1+1 and give me the result.",
	})
	if err != nil {
		t.Fatalf("ChatWithOllama returned an error: %v", err)
	}

	// Check that we got a response
	t.Logf("Response: %s", output.Response)
}
