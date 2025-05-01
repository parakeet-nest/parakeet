# Squawk Tutorial: The "jQuery" of Generative AI

!!! note
    This tutorial covers the core functionality of Squawk. For more advanced use cases and detailed API documentation, refer to the official Squawk documentation and source code.


Squawk is a DSL for Parakeet designed to simplify interactions with language models, making generative AI more accessible and easier to work with - similar to how jQuery simplified JavaScript development. This tutorial covers how to use Squawk for chat completions, retrieval-augmented generation (RAG), structured outputs, and function calling.

## Table of Contents

- [Installation](#installation)
- [Basic Chat Completions](#basic-chat-completions)
- [Retrieval-Augmented Generation (RAG)](#retrieval-augmented-generation-rag)
- [Structured Output](#structured-output)
- [Function Calling with Tools](#function-calling-with-tools)
- [Additional Features](#additional-features)

## Basic Chat Completions

### Simple Chat Completion

```go
squawk.New().
  Model("qwen2.5:3b").                    // Set the model
  Provider(provider.Ollama).              // Set the provider
  BaseURL("http://localhost:11434").      // Set the API URL
  System("You are a Go expert").          // Add a system message
  User("Explain concurrency in Go").      // Add a user message
  Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
    if err != nil {
      fmt.Println("Error:", err)
      return
    }
    fmt.Println(answer.Message.Content)   // Print response
  })
```

### Streaming Chat Completion

For real-time responses as they are generated:

```go
squawk.New().
  Model("qwen2.5:3b").
  Provider(provider.Ollama).
  System("You are a Go expert").
  User("Explain channels in Go").
  ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
    fmt.Print(answer.Message.Content)     // Print each chunk
    return nil
  })
```

### Conversation Management

You can maintain a conversation by saving the assistant's responses:

```go
sq := squawk.New().
  Model("mistral:latest").
  Provider(provider.Ollama).
  System("You are a Go expert").
  User("What is a goroutine?").
  Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
    if err != nil {
      return
    }
    fmt.Println(answer.Message.Content)
  }).
  SaveAssistantAnswer()                   // Save response to conversation history

// Continue the conversation
sq.User("How do they differ from threads?").
  Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
    if err != nil {
      return
    }
    fmt.Println(answer.Message.Content)
  })
```

### Message Labeling and Management

Squawk allows labeling messages and selectively removing them:

```go
squawk.New().
  Model("qwen2.5:3b").
  Provider(provider.Ollama).
  System("You are a Go expert", "role").
  User("What is a channel?", "question-1").
  Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
    // Handle response
  }).
  SaveAssistantAnswer("answer-1").
  RemoveMessageByLabel("question-1").     // Remove specific message
  User("How do goroutines work?", "question-2")
```

## Retrieval-Augmented Generation (RAG)

RAG enhances responses by retrieving relevant information from a knowledge base.

### Setting Up a Vector Store

```go
// Initialize a memory vector store
store := embeddings.MemoryVectorStore{
    Records: make(map[string]llm.VectorRecord),
}
```

### Generating Embeddings

```go
// Sample documents
docs := []string{
  "Go is a statically typed language",
  "Go supports concurrent programming",
}

sq := squawk.New().
  EmbeddingsModel("mxbai-embed-large:latest").     // Set embeddings model
  Model("qwen2.5:3b").
  Provider(provider.Ollama).
  BaseURL("http://localhost:11434").
  Store(&store).                            // Set the vector store
  GenerateEmbeddings(docs, true)           // Generate embeddings with logging
```

### Similarity Search

```go
// Search by content
sq.SimilaritySearch(
    "concurrent programming", // Query text
    0.7,                      // Similarity threshold (0-1)
    2,                        // Max results
)

// Get similarities
similarities := sq.Similarities()
for _, sim := range similarities {
    fmt.Println("Similar Content: ",sim.Prompt)
}
```

### Search from User Message

```go
sq.User("Tell me about Go concurrency", "question-1").
  SimilaritySearchFromUserMessage("question-1", 0.8, 3)
```

### Adding Context to Conversation

```go
sq.SimilaritySearch("concurrent Go features", 0.8, 3).
  AddSimilaritiesToMessages().
  User("Explain these concepts").
  Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
    fmt.Println(answer.Message.Content)
  })

// With custom prefix
sq.SimilaritySearch("concurrent Go features", 0.8, 3).
AddSimilaritiesToMessagesWithPrefix("Use this documentation as context: \n\n").    
User("Explain these concepts").
Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
    fmt.Println(answer.Message.Content)
})
```

### Complete RAG Example

```go
// Initialize a memory vector store
store := embeddings.MemoryVectorStore{
    Records: make(map[string]llm.VectorRecord),
}

// Star Wars characters example
starWarsChars := []string{
  `Luke Skywalker is the main protagonist of the original Star Wars trilogy...`,
  `Princess Leia Organa is a leader of the Rebel Alliance...`,
}

squawk.New().
    EmbeddingsModel("mxbai-embed-large:latest").
    Model("qwen2.5:3b").
    Provider(provider.Ollama).
    Store(&store).
    GenerateEmbeddings(starWarsChars, true).
    System("You are a Star Wars expert").
    User("Who is Luke Skywalker?", "q1").
    SimilaritySearchFromUserMessage("q1", 0.6, 1).
    AddSimilaritiesToMessages("sim1").
    ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
        fmt.Print(answer.Message.Content)
        return nil
    }).
    SaveAssistantAnswer("a1").
    RemoveMessageByLabel("q1").
    RemoveMessageByLabel("sim1").
    User("Who is Leia?", "q2").
    SimilaritySearchFromUserMessage("q2", 0.6, 1).
    AddSimilaritiesToMessages().
    ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
        fmt.Print(answer.Message.Content)
        return nil
    })
```

## Structured Output

Structured output allows you to get responses in a specific JSON format.

### Basic Structured Output

```go
squawk.New().
    Model("qwen2.5:3b").
    Provider(provider.Ollama).
    Schema(map[string]any{
        "type": "object",
        "properties": map[string]any{
            "name": map[string]any{
                "type":        "string",
                "description": "The name of the person",
            },
            "age": map[string]any{
                "type":        "integer",
                "description": "The age of the person",
            },
        },
        "required": []string{"name", "age"},
    }).
    User("Extract name and age from: John is 25 years old").
    StructuredOutput(func(answer llm.Answer, self *squawk.Squawk, err error) {
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        fmt.Printf("%+v\n", answer.Message.Content)
    })
```

### Complex Schema Example

```go
squawk.New().
    Model("qwen2.5:3b").
    Provider(provider.Ollama).
    Schema(map[string]any{
        "type": "object",
        "properties": map[string]any{
            "name": map[string]any{
                "type": "string",
            },
            "capital": map[string]any{
                "type": "string",
            },
            "languages": map[string]any{
                "type": "array",
                "items": map[string]any{
                    "type": "string",
                },
            },
        },
        "required": []string{"name", "capital", "languages"},
    }).
    User("Tell me about Canada").
    StructuredOutput(func(answer llm.Answer, self *squawk.Squawk, err error) {
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        fmt.Println("Structured Output:", answer.Message.Content)
    })
```

## Function Calling with Tools

Function calling allows the model to call predefined functions or tools.

### Defining Tools

```go
toolsList := []llm.Tool{
  {
    Type: "function",
    Function: llm.Function{
      Name:        "get_weather",
      Description: "Get current weather for a location",
      Parameters: llm.Parameters{
        Type: "object",
        Properties: map[string]llm.Property{
          "location": {
            Type:        "string",
            Description: "City name",
          },
        },
        Required: []string{"location"},
      },
    },
  },
  {
    Type: "function",
    Function: llm.Function{
      Name:        "calculate",
      Description: "Perform a calculation",
      Parameters: llm.Parameters{
        Type: "object",
        Properties: map[string]llm.Property{
          "operation": {
            Type:        "string",
            Description: "Math operation",
          },
          "a": {
            Type:        "number",
            Description: "First number",
          },
          "b": {
            Type:        "number",
            Description: "Second number",
          },
        },
        Required: []string{"operation", "a", "b"},
      },
    },
  },
}
```

### Using Function Calling

```go
squawk.New().
    Model("qwen2.5:3b").
    Provider(provider.Ollama).
    Tools(toolsList).
    User("What's the weather in Paris and calculate 2 + 2").
    FunctionCalling(func(answer llm.Answer, self *squawk.Squawk, err error) {
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        for _, toolCall := range self.ToolCalls() {
            fmt.Printf("Tool: %s\n", toolCall.Function.Name)
            fmt.Printf("Arguments: %v\n", toolCall.Function.Arguments)

            // Implement function logic
            switch toolCall.Function.Name {
            case "get_weather":
                location := toolCall.Function.Arguments["location"]
                // Call weather API
                fmt.Printf("Getting weather for %s\n", location)
            case "calculate":
                a := toolCall.Function.Arguments["a"].(float64)
                b := toolCall.Function.Arguments["b"].(float64)
                op := toolCall.Function.Arguments["operation"].(string)
                // Perform calculation
                fmt.Printf("Calculating %v %s %v\n", a, op, b)
            }
        }
    })
```

### Complete Function Calling Example

```go
toolsList := []llm.Tool{
    {
        Type: "function",
        Function: llm.Function{
            Name:        "hello",
            Description: "Say hello to a given person with their name",
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

squawk.New().
    Model("qwen2.5:3b").
    Provider(provider.Ollama).
    Tools(toolsList).
    User(`say "hello" to Bob, say "hello" to Sam`).
    User(`add 2 and 40`).
    FunctionCalling(func(answer llm.Answer, self *squawk.Squawk, err error) {
        var results string
        for _, toolCall := range self.ToolCalls() {
            switch toolCall.Function.Name {
            case "hello":
                results += fmt.Sprintf("Hello %s\n", toolCall.Function.Arguments["name"])
            case "addNumbers":
                a := toolCall.Function.Arguments["a"]
                b := toolCall.Function.Arguments["b"]
                results += fmt.Sprintf("Addition of %v and %v is %v\n",
                    a, b, a.(float64)+b.(float64))
            }
        }
        self.System("RESULTS:\n" + results) // Add results to conversation
    }).
    User("Use the results and format the output with fancy emojis").
    ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
        fmt.Print(answer.Message.Content)
        return nil
    })
```

## Additional Features

### Provider Configuration

Squawk supports multiple providers:

```go
// Ollama
squawk.New().
  Model("mistral:latest").
  BaseURL("http://localhost:11434").
  Provider(provider.Ollama)

// OpenAI
squawk.New().
  Model("gpt-4").
  Provider(provider.OpenAI, "your-api-key")

// Docker Model Runner
squawk.New().
  Model("llama:13b").
  BaseURL("http://localhost:12434").
  Provider(provider.DockerModelRunner)
```

### Model Options

Configure model behavior:

```go
squawk.New().
  Model("mistral:latest").
  Provider(provider.Ollama).
  Options(llm.SetOptions(map[string]interface{}{
    option.Temperature:   0.7,
    option.TopP:          0.9,
    option.RepeatLastN:   64,
    option.RepeatPenalty: 1.1,
    option.MaxTokens:     2000,
  }))
```

### Prompting Templates

Squawk includes built-in prompting templates:

```go
// Explain for kids
squawk.New().
    Model("qwen2.5:3b").
    Provider(provider.Ollama).
    ForKids("Explain Docker").
    Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
        fmt.Println(answer.Message.Content)
    })

// Brief explanation
squawk.New().
    Model("qwen2.5:3b").
    Provider(provider.Ollama).
    Brief("Explain Docker").
    Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
        fmt.Println(answer.Message.Content)
    })


// Other templates include:
// - AsAPoem()
// - AsARecipe()
// - Timeline()
// - Comparison()
// - StepByStep()
// - ProsAndCons()
// - InLaymansTerms()
// - Summarize()
```

### Error Handling

```go
sq := squawk.New().
  Model("mistral:latest").
  Provider(provider.Ollama).
  User("Generate some code").
  Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
    if err != nil {
      self.LastError(err)  // Set error
      return
    }
  })

// Check for errors
if err := sq.LastError(); err != nil {
  fmt.Printf("Operation failed: %v\n", err)
}
```

### Command Execution

Execute custom commands on a Squawk instance:

```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    System("You are a Go expert").
    Cmd(func(self *squawk.Squawk) {
        fmt.Println(self.Messages())
    })
```
