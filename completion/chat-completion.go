package completion

import (
	"encoding/json"
	"fmt"

	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
)

func getProvider(options ...string) string {
	return gear.GetOptionString(0, "", options...)
}
func getOpenAIKey(options ...string) string {
	return gear.GetOptionString(1, "", options...)
}

func Chat(url string, query llm.Query, options ...string) (llm.Answer, error) {
	selectedProvider := getProvider(options...)
	// ? should I test error instead of ""
	switch selectedProvider {
	case provider.Ollama:
		return ollamaChat(url, query)
	case provider.DockerModelRunner:
		return modelRunnerChat(url, query)
	case provider.OpenAI:
		openAIKey := getOpenAIKey(options...)
		return openAIChat(url, query, openAIKey)

	default: // if no provider is specified or empty, use the default one
		return ollamaChat(url, query)
	}

}

func ChatStream(url string, query llm.Query, onChunk func(llm.Answer) error, options ...string) (llm.Answer, error) {
	selectedProvider := getProvider(options...)
	switch selectedProvider {
	case provider.Ollama:
		return ollamaChatStream(url, query, onChunk)
	case provider.DockerModelRunner:
		return modelRunnerChatStream(url, query, onChunk)
	case provider.OpenAI:
		openAIKey := getOpenAIKey(options...)
		return openAIChatStream(url, query, onChunk, openAIKey)
	default:
		return ollamaChatStream(url, query, onChunk)
	}
}

func convertOpenAIAnswerToAnswer(openAIAnswer llm.OpenAIAnswer) (llm.Answer, error) {

	if openAIAnswer.Choices[0].Message.ToolCalls != nil {
		// Conversion of openAI ToolCalls  to toolcalls for llm.Answer
		var toolCalls = llm.ToolCalls{}
		for _, toolCall := range openAIAnswer.Choices[0].Message.ToolCalls {
			functionMame := toolCall["function"].(map[string]interface{})["name"].(string)
			arguments := toolCall["function"].(map[string]interface{})["arguments"].(string)
			// Create a map to hold the JSON data
			var result map[string]interface{}
			// Unmarshal the JSON string into the map
			err := json.Unmarshal([]byte(arguments), &result)
			if err != nil {
				return llm.Answer{}, fmt.Errorf("Error unmarshaling JSON arguments tool calls: %w", err)
			}

			tc := llm.ToolCall{
				Function: llm.FunctionTool{
					Name:      functionMame,
					Arguments: result,
				},
			}
			toolCalls = append(toolCalls, tc)
		}

		// Conversion of openAIAnswer  to llm.Answer
		answer := llm.Answer{
			Model: openAIAnswer.Model,
			Message: llm.Message{
				Role:      openAIAnswer.Choices[0].Message.Role,
				Content:   openAIAnswer.Choices[0].Message.Content,
				ToolCalls: toolCalls,
			},
		}
		return answer, nil
	} else {
		// Conversion of openAIAnswer  to llm.Answer
		answer := llm.Answer{
			Model: openAIAnswer.Model,
			Message: llm.Message{
				Role:    openAIAnswer.Choices[0].Message.Role,
				Content: openAIAnswer.Choices[0].Message.Content,
			},
		}
		return answer, nil
	}

}
