package squawk

import (
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
)

// BaseURL sets the base URL for the API endpoint of the chosen provider.
//
// Parameters:
//   - url: A string representing the base URL for the API endpoint
//     (e.g., "http://localhost:11434" for Ollama)
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  BaseURL("http://localhost:11434").
//	  Provider(provider.Ollama)
//
// Common use cases:
//   - Configuring custom API endpoints
//   - Setting up local development environments
//   - Connecting to remote model servers
func (s *Squawk) BaseURL(url string) *Squawk {
	s.baseUrl = url
	return s
}

// Provider sets the LLM (Large Language Model) provider for the Squawk instance
// and configures the API URL and other parameters based on the selected provider.
//
// Parameters:
//   - llmProvider: A string representing the LLM provider. Supported values include:
//   - provider.Ollama: Configures the API URL for the Ollama provider.
//   - provider.DockerModelRunner: Configures the API URL for the Docker Model Runner provider.
//   - provider.OpenAI: Configures the API URL for OpenAI and sets the OpenAI API key.
//   - parameters: Optional additional parameters. For provider.OpenAI, the first parameter
//     should be the OpenAI API key.
//
// Behavior:
//   - If the baseUrl is not set, a default API URL is assigned based on the provider.
//   - For provider.OpenAI, the first parameter in the `parameters` slice is used as the API key.
//   - If an unsupported provider is specified, the baseUrl is used as the API URL.
//
// Returns:
//   - A pointer to the updated Squawk instance.
func (s *Squawk) Provider(llmProvider string, parameters ...string) *Squawk {
	s.provider = llmProvider
	switch llmProvider {
	case provider.Ollama:
		//log.Println("🦙", llmProvider)
		if s.baseUrl == "" {
			s.apiUrl = "http://localhost:11434"
		} else {
			s.apiUrl = s.baseUrl
		}
	case provider.DockerModelRunner:
		//log.Println("🐳", llmProvider)
		if s.baseUrl == "" {
			s.apiUrl = "http://localhost:12434/engines/llama.cpp/v1"
		} else {
			s.apiUrl = s.baseUrl + "/engines/llama.cpp/v1"
		}

	case provider.OpenAI:
		//log.Println("🔵", llmProvider, parameters[0])
		if s.baseUrl == "" {
			s.apiUrl = "https://api.openai.com/v1"
		} else {
			s.apiUrl = s.baseUrl + "/v1"
		}
		s.openAPIKey = parameters[0]
	default: // Ollama
		s.apiUrl = s.baseUrl
	}
	return s
}

// Options sets the configuration options for the language model interactions.
//
// Parameters:
//   - options: An llm.Options struct containing configuration parameters such as
//     temperature, top_p, repeat_penalty, etc.
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  Options(llm.SetOptions(map[string]interface{}{
//	      option.Temperature:   0.7,
//	      option.TopP:         0.9,
//	      option.RepeatLastN:  64,
//	      option.RepeatPenalty: 1.1,
//	  }))
//
// Common options for Ollama:
//   - Temperature: Controls randomness (0.0 to 1.0)
//   - TopP: Nucleus sampling parameter
//   - RepeatLastN: Number of tokens to look back for repetitions
//   - RepeatPenalty: Penalty for repeated tokens
//   - TopK: Limits vocabulary to top K tokens
//   - MaxTokens: Maximum tokens to generate
func (s *Squawk) Options(options llm.Options) *Squawk {
	s.options = options
	return s
}
