package llm

import "encoding/json"


// https://github.com/ollama/ollama/blob/main/docs/api.md#parameters
type GenQuery struct {
	Model    string    `json:"model"`
	Options  Options   `json:"options"`
	Stream   bool      `json:"stream"`
	Prompt   string    `json:"prompt"`  // For "Simple" Completion
	Context  []int     `json:"context"` // For "Simple" Completion
	//Tools    []Tool    `json:"tools"`
	Suffix string `json:"suffix,omitempty"`

	Format    string `json:"format,omitempty"` // https://github.com/ollama/ollama/blob/main/docs/api.md#request-json-mode
	KeepAlive bool   `json:"keep_alive,omitempty"`
	Raw       bool   `json:"raw,omitempty"`
	System    string `json:"system,omitempty"`
	Template  string `json:"template,omitempty"`

	TokenHeaderName  string
	TokenHeaderValue string
}


func (query *GenQuery) ToJsonString() string {
	// for the verbose mode
	// Marshal the answer into JSON
	jsonBytes, err := json.MarshalIndent(query, "", "  ")

	if err != nil {
		return `{"error":"` + err.Error() + `"}`
	}
	// Convert JSON bytes to string
	jsonString := string(jsonBytes)
	return jsonString
}

// --------

type Query struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"` // For Chat Completion
	Options  Options   `json:"options"`
	Stream   bool      `json:"stream"`
	Tools    []Tool    `json:"tools"`

	//Format    string `json:"format"` // https://github.com/ollama/ollama/blob/main/docs/api.md#request-json-mode
	// change the type string to any to accept any type of format in case of structured outputs
	// https://ollama.com/blog/structured-outputs
	Format    any `json:"format,omitempty"` 
	KeepAlive bool   `json:"keep_alive,omitempty"`
	Raw       bool   `json:"raw,omitempty"`
	System    string `json:"system,omitempty"`
	Template  string `json:"template,omitempty"`

	TokenHeaderName  string
	TokenHeaderValue string

	// *** OpenAI specific fields ***
	// I will movve them to the options
	ResponseFormat map[string]interface{} `json:"response_format,omitempty"` // "text" or "json"

	// *** End of OpenAI specific fields ***

}

func (query *Query) ToJsonString() string {
	// for the verbose mode
	// Marshal the answer into JSON
	jsonBytes, err := json.MarshalIndent(query, "", "  ")

	if err != nil {
		return `{"error":"` + err.Error() + `"}`
	}
	// Convert JSON bytes to string
	jsonString := string(jsonBytes)
	return jsonString
}