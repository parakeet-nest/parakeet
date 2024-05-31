package llm

import (
	"bytes"
	"encoding/json"
)

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

func GenerateToolsContent(tools []Tool) (string, error) {
	toolsJSON, err := json.Marshal(&tools)
	if err != nil {
		return "", err
	}
	return "[AVAILABLE_TOOLS] " + string(toolsJSON) + " [/AVAILABLE_TOOLS]", nil
}

func GenerateToolsInstruction(userMessage string) string {
	return "[INST] " + userMessage + " [/INST]"
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
