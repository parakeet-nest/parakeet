package tools

import (
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
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

// GenerateUserToolsInstructions generates a string containing the user message wrapped in [INST] and [/INST] tags.
//
// ✋ it works for mistral:7b
//
// Parameters:
// - userMessage: a string representing the user message to be wrapped.
//
// Returns:
// - string: the user message wrapped in [INST] and [/INST] tags.
func GenerateUserToolsInstructions(userMessage string) string {
	return "[INST] " + userMessage + " [/INST]"
}

// GenerateSystemToolsInstructions generates a string containing the system content instructions for using "function calling".
//
// ✋ Use it only if the LLM does not implement function calling.
func GenerateSystemToolsInstructions() string {
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

// ConvertToLLMTools converts mcp.Tool to llm.Tool format
func ConvertMCPTools(tools []mcp.Tool) []llm.Tool {
	llmTools := make([]llm.Tool, len(tools))
	for i, tool := range tools {
		llmTools[i] = llm.Tool{
			Type: "function",
			Function: llm.Function{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters: llm.Parameters{
					Type:       tool.InputSchema.Type,
					Required:   tool.InputSchema.Required,
					Properties: convertToLLMProperties(tool.InputSchema.Properties),
				},
			},
		}
	}
	return llmTools
}

// Helper function to convert properties to llm.Property format
func convertToLLMProperties(props map[string]interface{}) map[string]llm.Property {
	result := make(map[string]llm.Property)

	for name, prop := range props {
		if propMap, ok := prop.(map[string]interface{}); ok {
			property := llm.Property{
				Type:        getString(propMap, "type"),
				Description: getString(propMap, "description"),
			}
			result[name] = property
		}
	}

	return result
}

// Helper function to safely get string values from map (unchanged)
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
