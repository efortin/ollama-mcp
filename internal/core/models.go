package core

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

// Note: ListModels is deprecated. Use HandlerFactory.ListModelsHandler() instead.
// This function is kept for backward compatibility but should not be used directly.

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

// Note: ModelInfo is deprecated. Use HandlerFactory.ModelInfoHandler() instead.
// This function is kept for backward compatibility but should not be used directly.
