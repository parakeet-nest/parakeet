package squawk

// Model sets the chat model identifier for the Squawk instance.
//
// The method allows specifying which language model should be used for chat
// interactions. It follows the method chaining pattern, allowing for fluent
// interface usage.
//
// Parameters:
//   - model: A string representing the model identifier (e.g., "gpt-4", "mistral:latest")
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example:
//
//	squawk := New().
//	  Model("gpt-4").
//	  Provider(provider.OpenAI)
func (s *Squawk) Model(model string) *Squawk {
	s.chatModel = model
	return s
}

// EmbeddingsModel sets the embeddings model identifier for the Squawk instance.
//
// The method specifies which model should be used for generating embeddings.
// It follows the method chaining pattern, allowing for fluent interface usage.
//
// Parameters:
//   - model: A string representing the embeddings model identifier for Ollama
//     (e.g., "nomic-embed-text", "all-minilm")
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example using Ollama provider:
//
//	squawk := New().
//	  EmbeddingsModel("nomic-embed-text").
//	  BaseURL("http://localhost:11434").
//	  Provider(provider.Ollama).
//	  Store(embeddings.NewMemoryVectorStore())
func (s *Squawk) EmbeddingsModel(model string) *Squawk {
	s.embeddingsModel = model
	return s
}
