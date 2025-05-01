package embeddings

import (
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
)

type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

func getProvider(options ...string) string {
	return gear.GetOptionString(0, "", options...)
}
func getOpenAIKey(options ...string) string {
	return gear.GetOptionString(1, "", options...)
}

func CreateEmbedding(engineURL string, query llm.Query4Embedding, id string, options ...string) (llm.VectorRecord, error) {
	selectedProvider := getProvider(options...)
	// ? should I test error instead of ""

	switch selectedProvider {
	case provider.Ollama:
		return ollamaCreateEmbedding(engineURL, query, id)
	case provider.DockerModelRunner:
		return modelRunnerCreateEmbedding(engineURL, query, id)
	case provider.OpenAI:
		openAIKey := getOpenAIKey(options...)
		return openAICreateEmbedding(engineURL, query, id, openAIKey)

	default: // if no provider is specified or empty, use the default one
		return ollamaCreateEmbedding(engineURL, query, id)
	}

}
