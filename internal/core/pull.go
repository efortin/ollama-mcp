package core

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

// PullModelInput represents the input for the PullModel function
type PullModelInput struct {
	Name       string `json:"name" jsonschema:"name of the model to pull"`
	Insecure   bool   `json:"insecure,omitempty" jsonschema:"allow insecure connections to the Ollama library"`
	NoProgress bool   `json:"no_progress,omitempty" jsonschema:"do not show progress"`
}

// PullModelOutput represents the output from the PullModel function
type PullModelOutput struct {
	Status  string `json:"status" jsonschema:"status of the pull operation"`
	Message string `json:"message" jsonschema:"message from the pull operation"`
}

// PullModel pulls a model from the Ollama library
func PullModel(ctx context.Context, req *mcp.CallToolRequest, input PullModelInput) (
	*mcp.CallToolResult,
	PullModelOutput,
	error,
) {
	// Input validation
	if input.Name == "" {
		return nil, PullModelOutput{}, fmt.Errorf("model name is required")
	}

	// Get the Ollama client from configuration
	config := GetConfig()
	if config == nil {
		return nil, PullModelOutput{}, fmt.Errorf("configuration not found")
	}
	if config.Client == nil {
		return nil, PullModelOutput{}, fmt.Errorf("ollama client not initialized")
	}

	client := config.Client

	// Set up the pull request
	pullRequest := &api.PullRequest{
		Model:    input.Name,
		Insecure: input.Insecure,
	}

	// Pull the model - let the library handle everything
	err := client.Pull(ctx, pullRequest, func(progress api.ProgressResponse) error {
		// Simply pass through - the library handles progress internally
		return nil
	})

	if err != nil {
		return nil, PullModelOutput{}, fmt.Errorf("failed to pull model %s: %w", input.Name, err)
	}

	return nil, PullModelOutput{
		Status:  "success",
		Message: fmt.Sprintf("Successfully pulled model %s", input.Name),
	}, nil
}
