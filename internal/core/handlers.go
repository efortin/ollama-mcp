package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

// HandlerFactory creates tool handlers with injected dependencies
type HandlerFactory struct {
	server *Server
}

// NewHandlerFactory creates a new handler factory
func NewHandlerFactory(server *Server) *HandlerFactory {
	return &HandlerFactory{
		server: server,
	}
}

// ChatHandler returns a handler function for the chat tool
func (h *HandlerFactory) ChatHandler() func(context.Context, *mcp.CallToolRequest, ChatInput) (*mcp.CallToolResult, ChatOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ChatInput) (*mcp.CallToolResult, ChatOutput, error) {
		// Add timeout to context
		timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// Validate input
		if err := h.validateChatInput(input); err != nil {
			return nil, ChatOutput{}, fmt.Errorf("invalid input: %w", err)
		}

		// Get configuration
		config := h.server.GetConfig()
		if config == nil {
			return nil, ChatOutput{}, fmt.Errorf("server configuration not found")
		}

		// Determine which tool is being used
		toolName := "chat" // Default
		if input.ToolName != "" {
			toolName = strings.ToLower(input.ToolName)
			if toolName != "code" && toolName != "chat" {
				toolName = "chat"
			}
		}

		// Use default model if not specified
		modelToUse := input.Model
		if modelToUse == "" {
			modelToUse = config.GetModel(toolName)
		}

		// Validate model name
		if err := h.validateModelName(modelToUse); err != nil {
			return nil, ChatOutput{}, err
		}

		// Build the chat request
		chatRequest := &api.ChatRequest{
			Model: modelToUse,
			Messages: []api.Message{
				{
					Role:    "user",
					Content: input.Message,
				},
			},
			Stream:  new(bool), // Set to false (non-streaming)
			Options: make(map[string]interface{}),
		}

		// Add system prompt if provided
		if input.SystemPrompt != "" {
			systemMsg := api.Message{
				Role:    "system",
				Content: input.SystemPrompt,
			}
			chatRequest.Messages = append([]api.Message{systemMsg}, chatRequest.Messages...)
		}

		// Set context size
		if input.ContextSize != nil {
			chatRequest.Options["num_ctx"] = *input.ContextSize
		} else {
			chatRequest.Options["num_ctx"] = config.ContextSize
		}

		// Set temperature
		if input.Temperature != nil {
			chatRequest.Options["temperature"] = *input.Temperature
		}

		// Set top_p
		if input.TopP != nil {
			chatRequest.Options["top_p"] = *input.TopP
		}

		// Set top_k
		if input.TopK != nil {
			chatRequest.Options["top_k"] = *input.TopK
		}

		// Add any additional options
		if input.Options != nil {
			for k, v := range input.Options {
				chatRequest.Options[k] = v
			}
		}

		// Set keep_alive
		if input.KeepAlive != nil {
			chatRequest.Options["keep_alive"] = *input.KeepAlive
		} else {
			chatRequest.Options["keep_alive"] = config.KeepAlive
		}

		// Variable to store the final response
		var finalResponse string

		// Use the official client's Chat method with timeout context
		err := config.Client.Chat(timeoutCtx, chatRequest, func(response api.ChatResponse) error {
			if response.Message.Content != "" {
				finalResponse = response.Message.Content
			}
			return nil
		})

		if err != nil {
			return nil, ChatOutput{}, fmt.Errorf("failed to chat with Ollama: %w", err)
		}

		return nil, ChatOutput{Response: finalResponse}, nil
	}
}

// CodeHandler returns a handler function for the code tool
func (h *HandlerFactory) CodeHandler() func(context.Context, *mcp.CallToolRequest, CodeInput) (*mcp.CallToolResult, CodeOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input CodeInput) (*mcp.CallToolResult, CodeOutput, error) {
		// Convert CodeInput to ChatInput and use the chat handler
		chatInput := ChatInput{
			Model:        "", // Always use the environment variable for code tool
			Message:      input.Message,
			ContextSize:  input.ContextSize,
			SystemPrompt: input.SystemPrompt,
			Options:      input.Options,
			KeepAlive:    input.KeepAlive,
			ToolName:     "code", // Specify that this is the code tool
		}

		// If no system prompt is provided, use a default one for code
		if chatInput.SystemPrompt == "" {
			chatInput.SystemPrompt = "You are a helpful coding assistant. Provide clear, concise, and well-commented code solutions."
		}

		// Call the chat handler
		result, chatOutput, err := h.ChatHandler()(ctx, req, chatInput)
		if err != nil {
			return nil, CodeOutput{}, err
		}

		// Convert ChatOutput to CodeOutput
		return result, CodeOutput(chatOutput), nil
	}
}

// ListModelsHandler returns a handler function for the list-models tool
func (h *HandlerFactory) ListModelsHandler() func(context.Context, *mcp.CallToolRequest, ListModelsInput) (*mcp.CallToolResult, ListModelsOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ListModelsInput) (*mcp.CallToolResult, ListModelsOutput, error) {
		// Add timeout to context
		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// Get the Ollama client from server
		client := h.server.GetClient()
		if client == nil {
			return nil, ListModelsOutput{}, fmt.Errorf("ollama client not initialized")
		}

		// Get the list of models
		response, err := client.List(timeoutCtx)
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
}

// ModelInfoHandler returns a handler function for the model-info tool
func (h *HandlerFactory) ModelInfoHandler() func(context.Context, *mcp.CallToolRequest, ModelInfoInput) (*mcp.CallToolResult, ModelInfoOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ModelInfoInput) (*mcp.CallToolResult, ModelInfoOutput, error) {
		// Add timeout to context
		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// Validate input
		if err := h.validateModelName(input.Name); err != nil {
			return nil, ModelInfoOutput{}, err
		}

		// Get the Ollama client from server
		client := h.server.GetClient()
		if client == nil {
			return nil, ModelInfoOutput{}, fmt.Errorf("ollama client not initialized")
		}

		// Get the model information
		response, err := client.Show(timeoutCtx, &api.ShowRequest{Name: input.Name})
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
}

// PullModelHandler returns a handler function for the pull-model tool
func (h *HandlerFactory) PullModelHandler() func(context.Context, *mcp.CallToolRequest, PullModelInput) (*mcp.CallToolResult, PullModelOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input PullModelInput) (*mcp.CallToolResult, PullModelOutput, error) {
		// Input validation
		if err := h.validateModelName(input.Name); err != nil {
			return nil, PullModelOutput{}, err
		}

		// Get the Ollama client from server
		client := h.server.GetClient()
		if client == nil {
			return nil, PullModelOutput{}, fmt.Errorf("ollama client not initialized")
		}

		// Set up the pull request
		pullRequest := &api.PullRequest{
			Model:    input.Name,
			Insecure: input.Insecure,
		}

		// Pull the model - let the library handle everything
		// Note: Pull operations can take a long time, so we don't add a timeout here
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
}

// Validation helper methods

func (h *HandlerFactory) validateChatInput(input ChatInput) error {
	if input.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	// Validate context size if provided
	if input.ContextSize != nil && *input.ContextSize <= 0 {
		return fmt.Errorf("context size must be positive")
	}

	// Validate temperature if provided
	if input.Temperature != nil && (*input.Temperature < 0 || *input.Temperature > 2.0) {
		return fmt.Errorf("temperature must be between 0 and 2.0")
	}

	// Validate top_p if provided
	if input.TopP != nil && (*input.TopP < 0 || *input.TopP > 1.0) {
		return fmt.Errorf("top_p must be between 0 and 1.0")
	}

	// Validate top_k if provided
	if input.TopK != nil && *input.TopK < 0 {
		return fmt.Errorf("top_k must be non-negative")
	}

	return nil
}

func (h *HandlerFactory) validateModelName(model string) error {
	if model == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	// Prevent path traversal attacks
	if strings.Contains(model, "..") || strings.Contains(model, "/") || strings.Contains(model, "\\") {
		return fmt.Errorf("invalid model name: %s", model)
	}

	// Check for suspicious characters
	if strings.ContainsAny(model, "<>|&;`$") {
		return fmt.Errorf("model name contains invalid characters: %s", model)
	}

	return nil
}
