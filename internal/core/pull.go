package core

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

// Note: PullModel is deprecated. Use HandlerFactory.PullModelHandler() instead.
// This function is kept for backward compatibility but should not be used directly.
