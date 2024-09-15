# Function Calling with tool support
> Ollama API: chat request with tools https://github.com/ollama/ollama/blob/main/docs/api.md#chat-request-with-tools

Since Ollama `0.3.0`, Ollama supports **tools calling**, blog post: https://ollama.com/blog/tool-support.
A list of supported models can be found under the Tools category on the models page: https://ollama.com/search?c=tools

## Define a list of tools

> use a supported model

```golang
model := "mistral:7b"

toolsList := []llm.Tool{
    {
        Type: "function",
        Function: llm.Function{
            Name:        "hello",
            Description: "Say hello to a given person with his name",
            Parameters: llm.Parameters{
                Type: "object",
                Properties: map[string]llm.Property{
                    "name": {
                        Type:        "string",
                        Description: "The name of the person",
                    },
                },
                Required: []string{"name"},
            },
        },
    },
    {
        Type: "function",
        Function: llm.Function{
            Name:        "addNumbers",
            Description: "Make an addition of the two given numbers",
            Parameters: llm.Parameters{
                Type: "object",
                Properties: map[string]llm.Property{
                    "a": {
                        Type:        "number",
                        Description: "first operand",
                    },
                    "b": {
                        Type:        "number",
                        Description: "second operand",
                    },
                },
                Required: []string{"a", "b"},
            },
        },
    },
}
```

## Set the Tools property of the query

> - set the `Temperature` to `0.0`
> - you don't need to set the row mode to true
> - set `query.Tools` with `toolsList`

```golang
messages := []llm.Message{
    {Role: "user", Content: `say "hello" to Bob`},
}

options := llm.SetOptions(map[string]interface{}{
    option.Temperature: 0.0,
    option.RepeatLastN: 2,
    option.RepeatPenalty: 2.0,
})

query := llm.Query{
    Model:    model,
    Messages: messages,
    Tools:    toolsList,
    Options:  options,
    Format:   "json",
}
```

## Run the completion

```go
answer, err := completion.Chat(ollamaUrl, query)
if err != nil {
    log.Fatal("😡:", err)
}

// It's a []map[string]interface{}
toolCalls := answer.Message.ToolCalls

// Convert toolCalls into a JSON string
jsonBytes, err := json.Marshal(toolCalls)
if err != nil {
    log.Fatal("😡:", err)
}
// Convert JSON bytes to string
result := string(jsonBytes)

fmt.Println(result)
```

The result will look like this:
```json
[{"function":{"arguments":{"name":"Bob"},"name":"hello"}}]
```

### Or you can use the `ToolCallsToJSONString` helper

```golang
answer, err = completion.Chat(ollamaUrl, query)
if err != nil {
    log.Fatal("😡:", err)
}

result, err = answer.Message.ToolCallsToJSONString()
if err != nil {
    log.Fatal("😡:", err)
}
fmt.Println(result)
```
The result will look like this:
```json
[{"function":{"arguments":{"name":"Bob"},"name":"hello"}}]
```

!!! note
	Look here for a complete sample: [examples/19-mistral-function-calling-tool-support](https://github.com/parakeet-nest/parakeet/tree/main/examples/19-mistral-function-calling-tool-support)

### Or (better) you can use the `ToolCallsToJSONString` helper

```golang
answer, err := completion.Chat(ollamaUrl, query)
if err != nil {
    log.Fatal("😡:", err)
}

result, err := answer.Message.ToolCalls[0].Function.ToJSONString()
if err != nil {
    log.Fatal("😡:", err)
}
fmt.Println(result)
```

The result will look like this:
```json
{"name":"hello","arguments":{"name":"Bob"}}
```

!!! note
	Look at these samples:

    - [examples/43-function-calling/01-xlam](https://github.com/parakeet-nest/parakeet/tree/main/examples/43-function-calling/01-xlam)
    - [examples/43-function-calling/02-qwen2tools](https://github.com/parakeet-nest/parakeet/tree/main/examples/43-function-calling/02-qwen2tools)
