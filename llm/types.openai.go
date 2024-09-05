package llm

import "encoding/json"

// https://platform.openai.com/docs/api-reference/chat/create

// Default chat query for gpt-4o-mini
/*
curl https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [
      {
        "role": "system",
        "content": "You are a helpful assistant."
      },
      {
        "role": "user",
        "content": "Hello!"
      }
    ]
  }'

### Response

{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "gpt-4o-mini",
  "system_fingerprint": "fp_44709d6fcb",
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "\n\nHello there, how may I assist you today?",
    },
    "logprobs": null,
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 9,
    "completion_tokens": 12,
    "total_tokens": 21
  }
}

*/


// Tool chat query for gpt-4o-mini
/*
curl https://api.openai.com/v1/chat/completions \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $OPENAI_API_KEY" \
-d '{
  "model": "gpt-4o",
  "messages": [
    {
      "role": "user",
      "content": "What'\''s the weather like in Boston today?"
    }
  ],
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_current_weather",
        "description": "Get the current weather in a given location",
        "parameters": {
          "type": "object",
          "properties": {
            "location": {
              "type": "string",
              "description": "The city and state, e.g. San Francisco, CA"
            },
            "unit": {
              "type": "string",
              "enum": ["celsius", "fahrenheit"]
            }
          },
          "required": ["location"]
        }
      }
    }
  ],
  "tool_choice": "auto"
}'

## Response
{
  "id": "chatcmpl-abc123",
  "object": "chat.completion",
  "created": 1699896916,
  "model": "gpt-4o-mini",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": null,
        "tool_calls": [
          {
            "id": "call_abc123",
            "type": "function",
            "function": {
              "name": "get_current_weather",
              "arguments": "{\n\"location\": \"Boston, MA\"\n}"
            }
          }
        ]
      },
      "logprobs": null,
      "finish_reason": "tool_calls"
    }
  ],
  "usage": {
    "prompt_tokens": 82,
    "completion_tokens": 17,
    "total_tokens": 99
  }
}

*/


/*
```bash
curl "https://api.openai.com/v1/chat/completions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $OPENAI_API_KEY" \
    -d '{
        "model": "gpt-4o-mini",
        "messages": [
            {
                "role": "system",
                "content": "You are a helpful assistant."
            },
            {
                "role": "user",
                "content": "who is Jean-Luc Picard."
            }
        ]
    }'
```

*/

type OpenAIQuery struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`

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

	Stream            bool   `json:"stream"`
	Tools             []Tool `json:"tools,omitempty"`                         // not used right now
	ToolChoices       string `json:"tool_choices,omitempty"`        // not used right now
	ParallelToolCalls bool   `json:"parallel_tool_calls,omitempty"` // not used right now
	User              string `json:"user,omitempty"`                // not used right now

	//Format    string `json:"format,omitempty"` // https://github.com/ollama/ollama/blob/main/docs/api.md#request-json-mode
	//KeepAlive bool   `json:"keep_alive,omitempty"`
	//Raw       bool   `json:"raw,omitempty"`
	//System    string `json:"system,omitempty"`
	//Template  string `json:"template,omitempty"`

	//TokenHeaderName  string
	//TokenHeaderValue string
	Verbose bool `json:"-"`  
	OpenAIAPIKey string `json:"-"`  
}

func (query *OpenAIQuery) ToJsonString() string {
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
}

type Choice struct {
    Index        int     `json:"index,omitempty"`
    Message      OpenAIMessage `json:"message,omitempty"`
    Logprobs     *string `json:"logprobs,omitempty"` // Assuming logprobs can be null
    FinishReason string  `json:"finish_reason,omitempty"`
}

type Usage struct {
    PromptTokens     int `json:"prompt_tokens,omitempty"`
    CompletionTokens int `json:"completion_tokens,omitempty"`
    TotalTokens      int `json:"total_tokens,omitempty"`
}


type OpenAIAnswer struct {
    ID               string   `json:"id"`
    Object           string   `json:"object"`
    Created          int64    `json:"created"`
    Model            string   `json:"model"`
    SystemFingerprint string  `json:"system_fingerprint"`
    Choices          []Choice `json:"choices"`
    Usage            Usage    `json:"usage"`
}

func (answer *OpenAIAnswer) ToJsonString() string {
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
