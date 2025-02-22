# Error handling

!!! info "🚧 work in progress"

## ModelNotFoundError

```golang
// package completion
type ModelNotFoundError struct {
  Code    int
  Message string
  Model   string
}
```

**Usage**:
```golang
answer, err := completion.Chat(ollamaUrl, query)
if err != nil {
  // test if the model is not found
  if modelErr, ok := err.(*completion.ModelNotFoundError); ok {
    fmt.Printf("💥 Got Model Not Found error: %s\n", modelErr.Message)
    fmt.Printf("😡 Error code: %d\n", modelErr.Code)
    fmt.Printf("🧠 Expected Model: %s\n", modelErr.Model)
  } else {
    log.Fatal("😡:", err)
  }
}
```

## NoSuchOllamaHostError

```golang
// package completion
type NoSuchOllamaHostError struct {
	Host string
	Message string
}
```

**Usage**:
```golang
if noHostErr, ok := err.(*completion.NoSuchOllamaHostError); ok {
  fmt.Printf("🦙 Got No Such Ollama Host error: %s\n", noHostErr.Message)
  fmt.Printf("🌍 Expected Host: %s\n", noHostErr.Host)
}
```


!!! note
	👀 you will find a complete example in:

    - [examples/04-chat-stream](https://github.com/parakeet-nest/parakeet/tree/main/examples/04-chat-stream)
    - [examples/66-structured-outputs](https://github.com/parakeet-nest/parakeet/tree/main/examples/66-structured-outputs)

