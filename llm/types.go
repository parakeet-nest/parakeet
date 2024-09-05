package llm

import (
	"encoding/json"
)

type LLM struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// not used
/*
type _Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	ToolCalls []struct {
		Function struct {
			Name      string `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		} `json:"function"`
	} `json:"tool_calls"`
}
*/

type FunctionTool struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"` // used for the ToolCalls list
}

func (ft *FunctionTool) ToJSONString() (string, error) {
	// Marshal the data into JSON
	jsonBytes, err := json.Marshal(ft)
	if err != nil {
		return "", err
	}
	// Convert JSON bytes to string
	jsonString := string(jsonBytes)
	return jsonString, nil
}

type Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	ToolCalls []struct {
		Function FunctionTool `json:"function"`
	} `json:"tool_calls"`
}

func (m *Message) ToolCallsToJSONString() (string, error) {
	// Marshal the data into JSON
	jsonBytes, err := json.Marshal(m.ToolCalls)
	if err != nil {
		return "", err
	}
	// Convert JSON bytes to string
	jsonString := string(jsonBytes)
	return jsonString, nil
}

func (m *Message) FirstToolCallToJSONString() (string, error) {
	// Marshal the data into JSON
	jsonBytes, err := json.Marshal(m.ToolCalls[0])
	if err != nil {
		return "", err
	}
	// Convert JSON bytes to string
	jsonString := string(jsonBytes)
	return jsonString, nil
}

type MessageRecord struct {
	Id      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

/*
type AnswerGenerate struct {
	Model    string  `json:"model"`
	Done     bool    `json:"done"`
	Response string  `json:"response"` // For "Simple" Completion
	Context  []int   `json:"context"`  // For "Simple" Completion
}

type AnswerChat struct {
	Model    string  `json:"model"`
	Message  Message `json:"message"` // For Chat Completion
	Done     bool    `json:"done"`
}
*/

/*
- https://github.com/ollama/ollama/blob/main/docs/api.md#generate-a-completion
- https://github.com/ollama/ollama/blob/main/docs/api.md#generate-a-chat-completion
- https://github.com/ollama/ollama/blob/main/api/types.go
- https://github.com/ollama/ollama/blob/main/docs/modelfile.md

- https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values
*/

type Options struct {
	RepeatLastN   int      `json:"repeat_last_n,omitempty"`
	Temperature   float64  `json:"temperature,omitempty"`
	Seed          int      `json:"seed,omitempty"`
	RepeatPenalty float64  `json:"repeat_penalty,omitempty"`
	Stop          []string `json:"stop,omitempty"`

	NumKeep          int     `json:"num_keep,omitempty"`
	NumPredict       int     `json:"num_predict,omitempty"`
	TopK             int     `json:"top_k,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	TFSZ             float64 `json:"tfs_z,omitempty"`
	TypicalP         float64 `json:"typical_p,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	Mirostat         int     `json:"mirostat,omitempty"`
	MirostatTau      float64 `json:"mirostat_tau,omitempty"`
	MirostatEta      float64 `json:"mirostat_eta,omitempty"`
	PenalizeNewline  bool    `json:"penalize_newline,omitempty"`

	Verbose bool
}

/* Embeddings */

type VectorRecord struct {
	Id             string    `json:"id"`
	Prompt         string    `json:"prompt"`
	Embedding      []float64 `json:"embedding"`
	CosineDistance float64
	Score          float64 // ElasticSearch

	Reference string `json:"reference"`
	MetaData  string `json:"metaData"`
	Text      string `json:"text"`
}

type Query4Embedding struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`

	TokenHeaderName  string
	TokenHeaderValue string
}
