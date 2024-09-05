package completion

import "github.com/parakeet-nest/parakeet/llm"

func ChatWithOpenAI(url string, query llm.OpenAIQuery) (llm.OpenAIAnswer, error) {
	return llm.OpenAIAnswer{}, nil
}


func ChatWithOpenAIStream(url string, query llm.OpenAIQuery, onChunk func(llm.OpenAIAnswer) error) (llm.OpenAIAnswer, error) {
	return llm.OpenAIAnswer{}, nil
}
