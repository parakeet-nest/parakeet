package tools

import (
	"encoding/json"

	"github.com/parakeet-nest/parakeet/llm"
)

// GenerateAvailableToolsContent generates a string containing the JSON representation of the provided tools array,
// wrapped in [AVAILABLE_TOOLS] and [/AVAILABLE_TOOLS] tags.
//
// ✋ it works for mistral:7b
//
// Parameters:
// - tools: an array of llm.Tool objects representing the tools to be converted to JSON.
//
// Returns:
// - string: the JSON representation of the tools array, wrapped in [AVAILABLE_TOOLS] and [/AVAILABLE_TOOLS] tags.
// - error: an error if the JSON marshaling fails.
func GenerateAvailableToolsContent(tools []llm.Tool) (string, error) {
	toolsJSON, err := json.Marshal(&tools)
	if err != nil {
		return "", err
	}
	return "[AVAILABLE_TOOLS] " + string(toolsJSON) + " [/AVAILABLE_TOOLS]", nil
}

// GenerateToolsInstructions generates a string containing the user message wrapped in [INST] and [/INST] tags.
//
// ✋ it works for mistral:7b
//
// Parameters:
// - userMessage: a string representing the user message to be wrapped.
//
// Returns:
// - string: the user message wrapped in [INST] and [/INST] tags.
func GenerateToolsInstructions(userMessage string) string {
	return "[INST] " + userMessage + " [/INST]"
}

// GenerateSystemInstructions generates a string containing the system content instructions for using "function calling".
// ✋ Use it only if the LLM does not implement function calling.
func GenerateSystemInstructions() string {
	systemContentInstructions := `If the question of the user matched the description of a tool, the tool will be called.
	To call a tool, respond with a JSON object with the following structure: 
	{
	  "name": <name of the called tool>,
	  "arguments": {
	    <name of the argument>: <value of the argument>
	  }
	}
	
	search the name of the tool in the list of tools with the Name field
	`
	return systemContentInstructions
}
