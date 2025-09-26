package core

// CodeInput represents the input for code generation operations
type CodeInput struct {
	Message      string         `json:"message" jsonschema:"the message to send to the model"`
	ContextSize  *int           `json:"context_size,omitempty" jsonschema:"maximum context size in tokens (optional)"`
	SystemPrompt string         `json:"system_prompt,omitempty" jsonschema:"system prompt to use (optional)"`
	Options      map[string]any `json:"options,omitempty" jsonschema:"additional model options (optional)"`
	KeepAlive    *string        `json:"keep_alive,omitempty" jsonschema:"duration to keep the model loaded in memory (optional)"`
}

type CodeOutput struct {
	Response string `json:"response" jsonschema:"the response from the model"`
}

// Note: Code is deprecated. Use HandlerFactory.CodeHandler() instead.
// This function is kept for backward compatibility but should not be used directly.
