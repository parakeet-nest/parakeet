package openai

// OpenAI API support
import (
	"encoding/json"
	"github.com/parakeet-nest/parakeet/llm"
)

// https://platform.openai.com/docs/api-reference/chat/create

type Query struct {
	Model    string        `json:"model"`
	Messages []llm.Message `json:"messages"`

	//Options  OpenAIOptions `json:"options"`
	//--------------------------------------------
	Stop        []string `json:"stop,omitempty"`
	Seed        int      `json:"seed,omitempty"`
	Temperature float64  `json:"temperature,omitempty"`
	TopP        float64  `json:"top_p,omitempty"`

	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	LogitBias       map[string]interface{} `json:"logit_bias,omitempty"`      // OpenAI specific
	Logprobs        bool                   `json:"logprobs,omitempty"`        // OpenAI specific
	TopLogprobs     int                    `json:"top_logprobs,omitempty"`    // OpenAI specific
	MaxTokens       int                    `json:"max_tokens,omitempty"`      // OpenAI specific
	N               int                    `json:"n,omitempty"`               // OpenAI specific
	Response_format map[string]interface{} `json:"response_format,omitempty"` // OpenAI specific
	ServiceTier     string                 `json:"service_tier,omitempty"`    // OpenAI specific

	StreamOptions map[string]interface{} `json:"stream_options,omitempty"` // OpenAI specific
	//--------------------------------------------

	Stream bool `json:"stream"`

	Tools      []llm.Tool `json:"tools,omitempty"`
	ToolChoice string `json:"tool_choice,omitempty"`

	ParallelToolCalls bool   `json:"parallel_tool_calls,omitempty"` // not used right now
	User              string `json:"user,omitempty"`                // not used right now

	//Format    string `json:"format,omitempty"` // https://github.com/ollama/ollama/blob/main/docs/api.md#request-json-mode
	//KeepAlive bool   `json:"keep_alive,omitempty"`
	//Raw       bool   `json:"raw,omitempty"`
	//System    string `json:"system,omitempty"`
	//Template  string `json:"template,omitempty"`

	//TokenHeaderName  string
	//TokenHeaderValue string
	Verbose      bool   `json:"-"`
	OpenAIAPIKey string `json:"-"`
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

type OpenAIMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
	//ToolCalls []interface{} `json:"tool_calls,omitempty"`
	ToolCalls []map[string]interface{} `json:"tool_calls,omitempty"`
}

/*
type OpenAIToolCall struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"` // "function" or "tool"
	Function FunctionTool `json:"function,omitempty"`
}
*/
//type OpenAIToolCalls []OpenAIToolCall

//! ????

/*
   "tool_calls": [
     {
       "id": "call_hkPjBb3TnBg532I7LStxDuqr",
       "type": "function",
       "function": {
         "name": "hello",
         "arguments": "{\"name\":\"Bob\"}"
       }
     }
   ],
*/

type Delta struct {
	Content string `json:"content,omitempty"`
}

type Choice struct {
	Index        int           `json:"index,omitempty"`
	Message      OpenAIMessage `json:"message,omitempty"`
	Logprobs     *string       `json:"logprobs,omitempty"` // Assuming logprobs can be null
	FinishReason string        `json:"finish_reason,omitempty"`
	Delta        Delta         `json:"delta,omitempty"`
}

// "choices":[{"index":0,"delta":{"content":" redemption"}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

/*
{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1694268190,"model":"gpt-4o-mini", "system_fingerprint": "fp_44709d6fcb", "choices":[{"index":0,"delta":{"role":"assistant","content":""},"logprobs":null,"finish_reason":null}]}

{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1694268190,"model":"gpt-4o-mini", "system_fingerprint": "fp_44709d6fcb", "choices":[{"index":0,"delta":{"content":"Hello"},"logprobs":null,"finish_reason":null}]}

....

{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1694268190,"model":"gpt-4o-mini", "system_fingerprint": "fp_44709d6fcb", "choices":[{"index":0,"delta":{},"logprobs":null,"finish_reason":"stop"}]}
*/

type Answer struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

func (answer *Answer) ToJsonString() string {
	// for the verbose mode
	// Marshal the answer into JSON
	jsonBytes, err := json.MarshalIndent(answer, "", "  ")

	if err != nil {
		return `{"error":"` + err.Error() + `"}`
	}
	// Convert JSON bytes to string
	jsonString := string(jsonBytes)
	return jsonString
}
