# Function Calling SLM Examples

> Ollama API: https://github.com/ollama/ollama/blob/main/docs/api.md#chat-request-with-tools

## xlam

https://ollama.com/allenporter/xlam

```bash
curl http://localhost:11434/api/chat -d '{
  "model": "allenporter/xlam:1b",
  "messages": [
    {
      "role": "user",
      "content": "What is the weather today in Paris?"
    }
  ],
  "stream": false,
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_current_weather",
        "description": "Get the current weather for a location",
        "parameters": {
          "type": "object",
          "properties": {
            "location": {
              "type": "string",
              "description": "The location to get the weather for, e.g. San Francisco, CA"
            },
            "format": {
              "type": "string",
              "description": "The format to return the weather in, e.g. 'celsius' or 'fahrenheit'",
              "enum": ["celsius", "fahrenheit"]
            }
          },
          "required": ["location", "format"]
        }
      }
    }
  ]
}' | jq '.message.tool_calls'
```

## qwen2tools

https://ollama.com/sam4096/qwen2tools

```bash
curl http://localhost:11434/api/chat -d '{
  "model": "sam4096/qwen2tools:0.5b",
  "messages": [
    {
      "role": "user",
      "content": "What is the weather today in Paris?"
    }
  ],
  "stream": false,
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_current_weather",
        "description": "Get the current weather for a location",
        "parameters": {
          "type": "object",
          "properties": {
            "location": {
              "type": "string",
              "description": "The location to get the weather for, e.g. San Francisco, CA"
            },
            "format": {
              "type": "string",
              "description": "The format to return the weather in, e.g. 'celsius' or 'fahrenheit'",
              "enum": ["celsius", "fahrenheit"]
            }
          },
          "required": ["location", "format"]
        }
      }
    }
  ]
}' | jq '.message.tool_calls'
```

```bash
curl http://localhost:11434/api/chat -d '{
  "model": "qwen2.5:latest",
  "messages": [
    {
      "role": "user",
      "content": "What is the weather today in Paris?"
    }
  ],
  "stream": false,
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_current_weather",
        "description": "Get the current weather for a location",
        "parameters": {
          "type": "object",
          "properties": {
            "location": {
              "type": "string",
              "description": "The location to get the weather for, e.g. San Francisco, CA"
            },
            "format": {
              "type": "string",
              "description": "The format to return the weather in, e.g. 'celsius' or 'fahrenheit'",
              "enum": ["celsius", "fahrenheit"]
            }
          },
          "required": ["location", "format"]
        }
      }
    }
  ]
}'
```

🟣 JSON Query: {
  "model": "qwen2.5:latest",
  "messages": [
    {
      "role": "user",
      "content": "say \"hello\" to Bob",
      "Label": ""
    }
  ],
  "options": {
    "repeat_last_n": 2,
    "temperature": 0,
    "seed": -1,
    "repeat_penalty": 2,
    "num_keep": 4,
    "num_predict": -1,
    "top_k": 40,
    "top_p": 0.9,
    "tfs_z": 1,
    "typical_p": 1,
    "presence_penalty": 0,
    "frequency_penalty": 0,
    "mirostat": 0,
    "mirostat_tau": 5,
    "mirostat_eta": 0.1,
    "penalize_newline": true,
    "min_p": 0.05
  },
  "stream": false,
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "hello",
        "description": "Say hello to a given person with his name",
        "parameters": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string",
              "description": "The name of the person"
            }
          },
          "required": [
            "name"
          ]
        }
      }
    },
    {
      "type": "function",
      "function": {
        "name": "addNumbers",
        "description": "Make an addition of the two given numbers",
        "parameters": {
          "type": "object",
          "properties": {
            "a": {
              "type": "number",
              "description": "first operand"
            },
            "b": {
              "type": "number",
              "description": "second operand"
            }
          },
          "required": [
            "a",
            "b"
          ]
        }
      }
    }
  ],
  "format": "json",
  "TokenHeaderName": "",
  "TokenHeaderValue": ""
}