<!-- TOPIC:
Experimental OpenAI API support
-->
# OpenAI API support

!!! info "✋ Only tested with the `gpt-4o-mini` model"


!!! note "Ollama provides **experimental** [compatibility](https://github.com/ollama/ollama/blob/main/docs/openai.md ) with parts of the [OpenAI API](https://platform.openai.com/docs/api-reference). As it's experimental, I prefer to keep the completion methods of Ollama and OpenAI "separated"."


## Chat completion

```golang
openAIUrl := "https://api.openai.com/v1"
model := "gpt-4o-mini"

systemContent := `You are an expert in Star Trek.`
userContent := `Who is Jean-Luc Picard?`

query := llm.OpenAIQuery{
	Model: model,
	Messages: []llm.Message{
		{Role: "system", Content: systemContent},
		{Role: "user", Content: userContent},
	},
	//Verbose: true,
	OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
}

answer, err := completion.ChatWithOpenAI(openAIUrl, query)
if err != nil {
	log.Fatal("😡:", err)
}
fmt.Println(answer.Choices[0].Message.Content)
```

## Chat completion with stream

```golang
openAIUrl := "https://api.openai.com/v1"
model := "gpt-4o-mini"

systemContent := `You are an expert in Star Trek.`
userContent := `Who is Jean-Luc Picard?`

query := llm.OpenAIQuery{
	Model: model,
	Messages: []llm.Message{
		{Role: "system", Content: systemContent},
		{Role: "user", Content: userContent},
	},
	//Verbose: true,
	OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
}

err = completion.ChatWithOpenAIStream(openAIUrl, query,
	func(answer llm.OpenAIAnswer) error {
		fmt.Print(answer.Choices[0].Delta.Content)
		return nil
	})

if err != nil {
	log.Fatal("😡:", err)
}
```

## Chat completion with tools
> 🚧 in progress