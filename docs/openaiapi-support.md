<!-- TOPIC:
Experimental OpenAI API support
-->
# OpenAI API support

!!! info "✋ Only tested with the `gpt-4o-mini` model, and `text-embedding-3-large` for the embeddings model."

!!! note "Since the release `0.2.7` I unified the completion methods"

## Chat completion

Use: `completion.Chat(openAIURL, query, provider.OpenAI, <OPENAI-API-KEY>)`

```golang
openAIURL := "https://api.openai.com/v1"
model := "gpt-4o-mini"

systemContent := `You are an expert in Star Trek.`
userContent := `Who is Jean-Luc Picard?`

query := llm.Query{
	Model: model,
	Messages: []llm.Message{
		{Role: "system", Content: systemContent},
		{Role: "user", Content: userContent},
	},
}

answer, err := completion.Chat(openAIURL, query, provider.OpenAI, os.Getenv("OPENAI_API_KEY"))
if err != nil {
	log.Fatal("😡:", err)
}
fmt.Println(answer.Message.Content)
```

## Chat completion with stream

Use: `completion.ChatStream(openAIURL, query, function, provider.OpenAI, <OPENAI-API-KEY>)`


```golang
openAIURL := "https://api.openai.com/v1"
model := "gpt-4o-mini"

systemContent := `You are an expert in Star Trek.`
userContent := `Who is Jean-Luc Picard?`

query := llm.Query{
	Model: model,
	Messages: []llm.Message{
		{Role: "system", Content: systemContent},
		{Role: "user", Content: userContent},
	},
}

_, err = completion.ChatStream(openAIURL, query,
	func(answer llm.Answer) error {
		fmt.Print(answer.Message.Content)
		return nil
	}, provider.OpenAI, os.Getenv("OPENAI_API_KEY"))

if err != nil {
	log.Fatal("😡:", err)
}
```

!!! note
	You can find examples in 
	
	  - [examples/44-chat-openai](https://github.com/parakeet-nest/parakeet/tree/main/examples/44-chat-openai)
	  - [examples/45-chat-stream-openai](https://github.com/parakeet-nest/parakeet/tree/main/examples/45-chat-stream-openai)

## Chat completion with tools

Use: `completion.Chat(openAIURL, query, provider.OpenAI, <OPENAI-API-KEY>)` and set `query.Tools`.

```golang
openAIURL := "https://api.openai.com/v1"
model := "gpt-4o-mini"
openAIKey := os.Getenv("OPENAI_API_KEY")

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

options := llm.SetOptions(map[string]interface{}{
	option.Temperature: 0.0,
})

messages := []llm.Message{
	{Role: "user", Content: `add 2 and 40`},
}

query := llm.Query{
	Model:    model,
	Messages: messages,
	Tools:    toolsList,
	Options:  options,
}

answer, err := completion.Chat(openAIURL, query, provider.OpenAI, openAIKey)
if err != nil {
	log.Fatal("😡:", err)
}

// Search tool to call in the answer
tool, err := answer.Message.ToolCalls.Find("addNumbers")
if err != nil {
	log.Fatal("😡:", err)
}
result, _ := tool.Function.ToJSONString()
fmt.Println(result)
```

!!! note
	You can find an example in [examples/81-tools-openai](https://github.com/parakeet-nest/parakeet/tree/main/examples/81-tools-openai)


## Create embeddings

```golang
// Create an embedding from a content
embedding, err := embeddings.CreateEmbedding(
	openAIURL,
	llm.Query4Embedding{
		Model:  "text-embedding-3-large",
		Prompt: "thi is the content of the document",				
	},
	"unique identifier",
	provider.OpenAI,
	os.Getenv("OPENAI_API_KEY"),
)
```

!!! note
	You can find an example in [examples/49-embeddings-memory-openai](https://github.com/parakeet-nest/parakeet/tree/main/examples/49-embeddings-memory-openai)
