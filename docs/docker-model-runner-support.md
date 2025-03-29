<!-- TOPIC:
Experimental Docker Model Runner support
-->
# Docker Model Runner support

!!! note "Since the release `0.2.7` I unified the completion methods"

When you use **Docker Model Runner**, if
- the application runs outside a container, use the following url: `http://localhost:12434/engines/llama.cpp/v1`
- the application runs into a container, use the following url: `http://model-runner.docker.internal/engines/llama.cpp/v1/`

## Chat completion

Use: `completion.Chat(modelRunnerURL, query, provider.DockerModelRunner)`

```golang
modelRunnerURL := "http://localhost:12434/engines/llama.cpp/v1"
model := "ai/qwen2.5:latest" 

systemContent := `You are an expert in Star Trek.`
userContent := `Who is Jean-Luc Picard?`

options := llm.SetOptions(map[string]interface{}{
    option.Temperature: 0.5,
    option.RepeatPenalty: 2.0,
})

query := llm.Query{
    Model: model,
    Messages: []llm.Message{
        {Role: "system", Content: systemContent},
        {Role: "user", Content: userContent},
    },
    Options: options,
}

answer, err := completion.Chat(modelRunnerURL, query, provider.DockerModelRunner)
if err != nil {
    log.Fatal("😡:", err)
}
fmt.Println(answer.Message.Content)
```

## Chat completion with stream

Use: `completion.ChatStream(openAIURL, query, function, provider.DockerModelRunner)`

```golang
options := llm.SetOptions(map[string]interface{}{
    option.Temperature:   0.5,
    option.RepeatPenalty: 3.0,
})

query := llm.Query{
    Model: "ai/mistral:latest",
    Messages: []llm.Message{
        {Role: "system", Content: `You are a Borg in Star Trek. Speak like a Borg`},
        {Role: "user", Content: `Who is Jean-Luc Picard?`},
    },
    Options: options,
}

_, err := completion.ChatStream("http://localhost:12434/engines/llama.cpp/v1", query,
    func(answer llm.Answer) error {
        fmt.Print(answer.Message.Content)
        return nil
    }, provider.DockerModelRunner)

if err != nil {
    log.Fatal("😡:", err)
}
```

!!! note
	You can find examples in 
    
	  - [examples/77-chat-model-runner](https://github.com/parakeet-nest/parakeet/tree/main/77-chat-model-runner)
	  - [examples/78-chat-stream-model-runner](https://github.com/parakeet-nest/parakeet/tree/main/examples/78-chat-stream-model-runner)

## Chat completion with tools

Use: `completion.Chat(openAIURL, query, provider.DockerModelRunner) and set `query.Tools`.

```golang
modelRunnerURL := "http://localhost:12434/engines/llama.cpp/v1"
model := "ai/smollm2"

toolsList := []llm.Tool{
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

messages := []llm.Message{
    {Role: "user", Content: `add 2 and 40`},
}

query := llm.Query{
    Model:    model,
    Messages: messages,
    Tools:    toolsList,
    Options:  options,
}

answer, err := completion.Chat(modelRunnerURL, query, provider.DockerModelRunner)
if err != nil {
    log.Fatal("😡 completion bis:", err)
}

// Search tool to call in the answer
tool, err := answer.Message.ToolCalls.Find("addNumbers")
if err != nil {
    log.Fatal("😡 ToolCalls.Find bis:", err)
}
result, _ := tool.Function.ToJSONString()
fmt.Println(result)
```

!!! note
	You can find an example in [examples/80-tools-model-runner](https://github.com/parakeet-nest/parakeet/tree/main/examples/80-tools-model-runner)


## Create embeddings

```golang
// Create an embedding from a content
embedding, err := embeddings.CreateEmbedding(
	openAIURL,
	llm.Query4Embedding{
		Model:  "ai/mxbai-embed-large",
		Prompt: "thi is the content of the document",				
	},
	"unique identifier",
	provider.DockerModelRunner,
)
```

!!! note
	You can find an example in [examples/79-embeddings-memory-model-runner](https://github.com/parakeet-nest/parakeet/tree/main/examples/79-embeddings-memory-model-runner)
