package completion

import (
	"github.com/parakeet-nest/parakeet/llm"
)

func Generate(url string, query llm.Query) (llm.Answer, error) {
	return completion(url, "generate", query)
}

func GenerateStream(url string, query llm.Query, onChunk func(llm.Answer) error) (llm.Answer, error) {
	return completionStream(url, "generate", query, onChunk)
}