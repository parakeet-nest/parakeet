package tools

import (
	"encoding/json"

	"github.com/parakeet-nest/parakeet/llm"
)

// GenerateContent generates a string containing the JSON representation of the provided tools array,
// wrapped in [AVAILABLE_TOOLS] and [/AVAILABLE_TOOLS] tags.
//
// Parameters:
// - tools: an array of llm.Tool objects representing the tools to be converted to JSON.
//
// Returns:
// - string: the JSON representation of the tools array, wrapped in [AVAILABLE_TOOLS] and [/AVAILABLE_TOOLS] tags.
// - error: an error if the JSON marshaling fails.
func GenerateContent(tools []llm.Tool) (string, error) {
	toolsJSON, err := json.Marshal(&tools)
	if err != nil {
		return "", err
	}
	return "[AVAILABLE_TOOLS] " + string(toolsJSON) + " [/AVAILABLE_TOOLS]", nil
}

// GenerateInstructions generates a string containing the user message wrapped in [INST] and [/INST] tags.
//
// Parameters:
// - userMessage: a string representing the user message to be wrapped.
//
// Returns:
// - string: the user message wrapped in [INST] and [/INST] tags.
func GenerateInstructions(userMessage string) string {
	return "[INST] " + userMessage + " [/INST]"
}
