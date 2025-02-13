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
		Function FunctionTool //`json:"function"`
		Result   interface{}
		Error    error
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
	SessionId string `json:"sessionId"`
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
	NumCtx int `json:"num_ctx,omitempty"`

	RepeatLastN   int      `json:"repeat_last_n"`
	Temperature   float64  `json:"temperature"`
	Seed          int      `json:"seed"`
	RepeatPenalty float64  `json:"repeat_penalty"`
	Stop          []string `json:"stop,omitempty"`

	NumKeep          int     `json:"num_keep"`
	NumPredict       int     `json:"num_predict"`
	TopK             int     `json:"top_k"`
	TopP             float64 `json:"top_p"`
	TFSZ             float64 `json:"tfs_z"`
	TypicalP         float64 `json:"typical_p"`
	PresencePenalty  float64 `json:"presence_penalty"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	Mirostat         int     `json:"mirostat"`
	MirostatTau      float64 `json:"mirostat_tau"`
	MirostatEta      float64 `json:"mirostat_eta"`
	PenalizeNewline  bool    `json:"penalize_newline"`

	Verbose bool
}

/* Default Ollama Options
https://github.com/ollama/ollama/blob/main/api/types.go
*/

func DefaultOptions() Options {
	return Options{
		NumPredict: -1,

		NumKeep:          4,
		Temperature:      0.8,
		TopK:             40,
		TopP:             0.9,
		TFSZ:             1.0,
		TypicalP:         1.0,
		RepeatLastN:      64,
		RepeatPenalty:    1.1,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
		Mirostat:         0,
		MirostatTau:      5.0,
		MirostatEta:      0.1,
		PenalizeNewline:  true,
		Seed:             -1,
	}
}

func SetOptions(options map[string]interface{}) Options {
	defaultOptions := DefaultOptions()
	for key, value := range options {
		switch key {
		case "NumCtx":
			defaultOptions.NumCtx = value.(int)
		case "NumPredict":
			defaultOptions.NumPredict = value.(int)
		case "NumKeep":
			defaultOptions.NumKeep = value.(int)
		case "Temperature":
			defaultOptions.Temperature = value.(float64)
		case "TopK":
			defaultOptions.TopK = value.(int)
		case "TopP":
			defaultOptions.TopP = value.(float64)
		case "TFSZ":
			defaultOptions.TFSZ = value.(float64)
		case "TypicalP":
			defaultOptions.TypicalP = value.(float64)
		case "RepeatLastN":
			defaultOptions.RepeatLastN = value.(int)
		case "RepeatPenalty":
			defaultOptions.RepeatPenalty = value.(float64)
		case "PresencePenalty":
			defaultOptions.PresencePenalty = value.(float64)
		case "FrequencyPenalty":
			defaultOptions.FrequencyPenalty = value.(float64)
		case "Mirostat":
			defaultOptions.Mirostat = value.(int)
		case "MirostatTau":
			defaultOptions.MirostatTau = value.(float64)
		case "MirostatEta":
			defaultOptions.MirostatEta = value.(float64)
		case "PenalizeNewline":
			defaultOptions.PenalizeNewline = value.(bool)
		case "Seed":
			defaultOptions.Seed = value.(int)
		case "Verbose":
			defaultOptions.Verbose = value.(bool)
		}
	}
	return defaultOptions
}

/* Embeddings */

type VectorRecord struct {
	Id               string    `json:"id"`
	Prompt           string    `json:"prompt"`
	Embedding        []float64 `json:"embedding"`
	CosineSimilarity float64
	Score            float64 // ElasticSearch

	Reference      string                 `json:"reference"`
	SimpleMetaData string                 `json:"metaData"`
	Metadata       map[string]interface{} `json:"metadata"` // additional metadata
	Text           string                 `json:"text"`
}

type Query4Embedding struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`

	TokenHeaderName  string
	TokenHeaderValue string
}
