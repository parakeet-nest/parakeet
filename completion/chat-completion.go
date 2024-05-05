package completion

import (
	"github.com/parakeet-nest/parakeet/llm"
)

func Chat(url string, query llm.Query) (llm.Answer, error) {
	return completion(url, "chat", query)
}

func ChatStream(url string, query llm.Query, onChunk func(llm.Answer) error) (llm.Answer, error) {
	return completionStream(url, "chat", query, onChunk)
}
