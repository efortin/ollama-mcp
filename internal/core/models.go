package core

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

// Model represents an Ollama model
type Model struct {
	Name        string `json:"name" jsonschema:"name of the model"`
	Size        int64  `json:"size" jsonschema:"size of the model in bytes"`
	ModifiedAt  string `json:"modified_at" jsonschema:"timestamp when the model was last modified"`
	Digest      string `json:"digest" jsonschema:"digest of the model"`
	Description string `json:"description" jsonschema:"description of the model"`
}

// ListModelsInput represents the input for the ListModels function
type ListModelsInput struct {
	// No input parameters needed for listing models
}

// ListModelsOutput represents the output from the ListModels function
type ListModelsOutput struct {
	Models []Model `json:"models" jsonschema:"list of available models"`
}

// ListModels lists all available models from Ollama
func ListModels(ctx context.Context, req *mcp.CallToolRequest, input ListModelsInput) (
	*mcp.CallToolResult,
	ListModelsOutput,
	error,
) {
	// Get the Ollama client from configuration
	config := GetConfig()
	if config == nil {
		return nil, ListModelsOutput{}, fmt.Errorf("failed to get configuration")
	}

	client := config.Client

	// Get the list of models
	response, err := client.List(ctx)
	if err != nil {
		return nil, ListModelsOutput{}, fmt.Errorf("failed to list models: %w", err)
	}

	// Convert the response to our output format
	models := make([]Model, len(response.Models))
	for i, model := range response.Models {
		models[i] = Model{
			Name:        model.Name,
			Size:        model.Size,
			ModifiedAt:  model.ModifiedAt.Format(time.RFC3339),
			Digest:      model.Digest,
			Description: "", // API doesn't provide description in List response
		}
	}

	return nil, ListModelsOutput{Models: models}, nil
}

// ModelInfoInput represents the input for the ModelInfo function
type ModelInfoInput struct {
	Name string `json:"name" jsonschema:"name of the model to get information about"`
}

// ModelInfoOutput represents the output from the ModelInfo function
type ModelInfoOutput struct {
	Name       string `json:"name" jsonschema:"name of the model"`
	License    string `json:"license" jsonschema:"license of the model"`
	Modelfile  string `json:"modelfile" jsonschema:"modelfile content"`
	Parameters string `json:"parameters" jsonschema:"model parameters"`
	Template   string `json:"template" jsonschema:"model template"`
	System     string `json:"system" jsonschema:"system prompt"`
	ModifiedAt string `json:"modified_at" jsonschema:"timestamp when the model was last modified"`
}

// ModelInfo gets detailed information about a specific model
func ModelInfo(ctx context.Context, req *mcp.CallToolRequest, input ModelInfoInput) (
	*mcp.CallToolResult,
	ModelInfoOutput,
	error,
) {
	// Get the Ollama client from configuration
	config := GetConfig()
	if config == nil {
		return nil, ModelInfoOutput{}, fmt.Errorf("failed to get configuration")
	}

	client := config.Client

	// Get the model information
	response, err := client.Show(ctx, &api.ShowRequest{Name: input.Name})
	if err != nil {
		return nil, ModelInfoOutput{}, fmt.Errorf("failed to get model info: %w", err)
	}

	// Convert the response to our output format
	output := ModelInfoOutput{
		Name:       input.Name,
		License:    response.License,
		Modelfile:  response.Modelfile,
		Parameters: response.Parameters,
		Template:   response.Template,
		System:     response.System,
		ModifiedAt: response.ModifiedAt.Format(time.RFC3339),
	}

	return nil, output, nil
}
