# Squawk: Getting Started

!!! info
    **Squawk is like the "jQuery" of Generative AI**

    Squawk is a DSL for Parakeet designed to simplify interactions with language models, making generative AI more accessible and easier to work with - similar to how jQuery simplified JavaScript development. 

## Small chat with Squawk

```golang
ollamaBaseUrl := "http://localhost:11434"
model := "qwen2.5:1.5b"

options := llm.SetOptions(map[string]interface{}{
    option.Temperature:   0.5,
    option.RepeatLastN:   2,
    option.RepeatPenalty: 2.2,
})

ollamaParrot := squawk.New().
    Model(model).
    BaseURL(ollamaBaseUrl).
    Provider(provider.Ollama)

ollamaParrot.
    Options(options).
    System("You are a useful AI agent, you are a Star Trek expert.").
    User("Who is James T Kirk?").
    Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
        fmt.Println(answer.Message.Content)
    }).
    SaveAssistantAnswer()

ollamaParrot.
    User("Who is his best friend?").
    Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
        fmt.Println(answer.Message.Content)
    })
```

## You can chain everything

```golang
squawk.New().
    Model(model).
    BaseURL(ollamaBaseUrl).
    Provider(provider.Ollama).
    Options(options).
    System("You are a useful AI agent, you are a Star Trek expert.").
    User("Who is James T Kirk?").
    Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
        fmt.Println(answer.Message.Content)
    }).
    SaveAssistantAnswer().
    User("Who is his best friend?").
    Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
        fmt.Println(answer.Message.Content)
    })
```
