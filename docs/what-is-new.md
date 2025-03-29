# What's new with Parakeet 

## 🦜 Parakeet `v0.2.7` 🐳 [spouting whale]

Addition of **Docker Model Runner** support (and OpenAI at the same time) allowing easy development of generative AI applications in Docker containers.

```golang
modelRunnerURL := "http://model-runner.docker.internal/engines/llama.cpp/v1/"
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

answer, err := completion.Chat(
    modelRunnerURL, 
    query, 
    provider.DockerModelRunner,
)
if err != nil {
    log.Fatal("🫢 Oops!", err)
}
fmt.Println(answer.Message.Content)
```

!!! note
	👋 Look at this **Docker Compose** sample [examples/82-web-chat-bot-model-runner](https://github.com/parakeet-nest/parakeet/tree/main/examples/82-web-chat-bot-model-runner)


