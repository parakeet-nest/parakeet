# Squawk Documentation

!!! note "🚧 work in progress"

> "Squawk is the jQuery of generative AI"

Squawk simplifies common tasks with generative AI, making the technology more accessible and easier to work with.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Core Concepts](#core-concepts)
- [Basic Methods](#basic-methods)
- [Message Management](#message-management)
- [Chat Completions](#chat-completions)
- [Streaming Responses](#streaming-responses)
- [Structured Output](#structured-output)
- [Embeddings and RAG](#embeddings-and-rag)
- [Function Calling with Tools](#function-calling-with-tools)
- [Meta Prompts](#meta-prompts)
- [Complete Examples](#complete-examples)

## Overview

Squawk provides a simplified interface for interacting with generative AI models and embeddings. It encapsulates various configurations, states, and tools required for processing and generating responses from language models. It supports functionalities such as setting models, managing conversation contexts, handling embeddings, and executing structured or streaming outputs.

Key features:
- Simplified interaction with language models and embeddings
- Support for multiple providers (OpenAI, Ollama, Docker Model Runner)
- Flexible configuration options for chat and structured outputs
- Tools for managing conversation contexts and embeddings
- Callback-based execution for chat and streaming responses

## Installation

To use Squawk in your Go project:

```go
go get github.com/parakeet-nest/parakeet/squawk
```

## Core Concepts

### Squawk Structure

The Squawk struct contains all the configurations, states, and tools required for interacting with language models:

```go
type Squawk struct {
    setOfMessages   []llm.Message
    baseUrl         string
    apiUrl          string
    provider        string
    chatModel       string
    embeddingsModel string
    options         llm.Options
    openAPIKey      string
    lastAnswer      llm.Answer
    lastError       error
    schema          map[string]any

    // embeddings
    vectorStore  embeddings.VectorStore
    similarities []llm.VectorRecord

    // tools
    tools     []llm.Tool
    toolCalls []llm.ToolCall
}
```

### Creating a New Squawk Instance

Start by creating a new Squawk instance:

```go
sq := squawk.New()
```

### Method Chaining

Squawk follows the method chaining pattern, allowing for a fluent interface:

```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    System("You are a Go expert").
    User("Explain channels in Go").
    Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
        fmt.Println(answer.Message.Content)
    })
```

## Basic Methods

### Model

Sets the chat model identifier for the Squawk instance.

```go
func (s *Squawk) Model(model string) *Squawk
```

Example:
```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama)
```

### EmbeddingsModel

Sets the embeddings model identifier for the Squawk instance.

```go
func (s *Squawk) EmbeddingsModel(model string) *Squawk
```

Example:
```go
squawk.New().
    EmbeddingsModel("nomic-embed-text").
    Provider(provider.Ollama)
```

### BaseURL

Sets the base URL for the API endpoint of the chosen provider.

```go
func (s *Squawk) BaseURL(url string) *Squawk
```

Example:
```go
squawk.New().
    Model("mistral:latest").
    BaseURL("http://localhost:11434").
    Provider(provider.Ollama)
```

### Provider

Sets the LLM provider for the Squawk instance and configures the API URL and other parameters.

```go
func (s *Squawk) Provider(llmProvider string, parameters ...string) *Squawk
```

Examples:

```go
// Ollama
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama)

// OpenAI
squawk.New().
    Model("gpt-4").
    Provider(provider.OpenAI, "your-api-key")

// Docker Model Runner
squawk.New().
    Model("llama:13b").
    Provider(provider.DockerModelRunner)
```

### Options

Sets the configuration options for the language model interactions.

```go
func (s *Squawk) Options(options llm.Options) *Squawk
```

Example:
```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    Options(llm.SetOptions(map[string]interface{}{
        option.Temperature:   0.7,
        option.TopP:          0.9,
        option.RepeatLastN:   64,
        option.RepeatPenalty: 1.1,
    }))
```

Common options:
- Temperature: Controls randomness (0.0 to 1.0)
- TopP: Nucleus sampling parameter
- RepeatLastN: Number of tokens to look back for repetitions
- RepeatPenalty: Penalty for repeated tokens
- TopK: Limits vocabulary to top K tokens
- MaxTokens: Maximum tokens to generate

## Message Management

### System

Adds a system message to the conversation context.

```go
func (s *Squawk) System(message string, optionalParameters ...string) *Squawk
```

Examples:
```go
// Basic system message
squawk.New().
    System("You are a helpful AI assistant specializing in Go programming")

// Labeled system message
squawk.New().
    System("You are a code review expert", "code-reviewer")
```

### User

Adds a user message to the conversation context.

```go
func (s *Squawk) User(message string, optionalParameters ...string) *Squawk
```

Examples:
```go
// Basic user message
squawk.New().
    User("What is the capital of France?")

// Labeled user message
squawk.New().
    User("Review this Go code...", "code-review-1")
```

### Assistant

Adds an assistant message to the conversation context.

```go
func (s *Squawk) Assistant(message string, optionalParameters ...string) *Squawk
```

Examples:
```go
// Basic assistant message
squawk.New().
    Assistant("A goroutine is a lightweight thread of execution in Go")

// Labeled assistant message
squawk.New().
    Assistant("The function looks good but needs error handling", "review-1")
```

### AddSetOfMessages

Creates or extends a conversation with provided messages.

```go
func (s *Squawk) AddSetOfMessages(messages ...interface{}) *Squawk
```

Examples:
```go
// Adding single messages
squawk.New().
    AddSetOfMessages(
        llm.Message{Role: "system", Content: "You are a Go expert"},
        llm.Message{Role: "user", Content: "Explain channels"},
    )

// Adding a slice of messages
previousMessages := []llm.Message{
    {Role: "system", Content: "You are a Go expert"},
    {Role: "user", Content: "What is concurrency?"},
    {Role: "assistant", Content: "Concurrency is..."},
}

squawk.New().
    AddSetOfMessages(previousMessages).
    User("How does this relate to goroutines?")
```

### Messages

Returns the current conversation context as a slice of messages.

```go
func (s *Squawk) Messages() []llm.Message
```

Example:
```go
messages := squawk.Messages()
for _, msg := range messages {
    fmt.Printf("Role: %s, Content: %s\n", msg.Role, msg.Content)
}
```

### SaveAssistantAnswer

Adds the last model response to the conversation history as an assistant message.

```go
func (s *Squawk) SaveAssistantAnswer(optionalParameters ...string) *Squawk
```

Examples:
```go
// Save without label
squawk.New().
    Chat(/* ... */).
    SaveAssistantAnswer()

// Save with label
squawk.New().
    Chat(/* ... */).
    SaveAssistantAnswer("answer-1")
```

### LastAnswer

Gets or sets the most recent answer from the language model.

```go
func (s *Squawk) LastAnswer(optionalAnswer ...llm.Answer) llm.Answer
```

Example:
```go
// Get last answer
lastAnswer := squawk.LastAnswer()
fmt.Println("Last response:", lastAnswer.Message.Content)

// Set custom answer
customAnswer := llm.Answer{
    Message: llm.Message{
        Role: "assistant",
        Content: "Custom response",
    },
}
squawk.LastAnswer(customAnswer)
```

### RemoveMessageByLabel

Removes messages with the specified label from the conversation context.

```go
func (s *Squawk) RemoveMessageByLabel(label string) *Squawk
```

Example:
```go
squawk.New().
    System("You are a Go expert", "role").
    User("What is a channel?", "question-1").
    Chat(/* ... */).
    SaveAssistantAnswer("answer-1").
    RemoveMessageByLabel("question-1")
```

## Chat Completions

### Chat

Executes a non-streaming chat completion request and handles the response through a callback function.

```go
func (s *Squawk) Chat(callBack func(answer llm.Answer, self *Squawk, err error)) *Squawk
```

Example:
```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    System("You are a Go expert").
    User("What is a channel?").
    Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }
        fmt.Println(answer.Message.Content)
    })
```

## Streaming Responses

### ChatStream

Executes a streaming chat completion request and handles the response through a callback function.

```go
func (s *Squawk) ChatStream(callBack func(answer llm.Answer, self *Squawk) error) *Squawk
```

Example:
```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    System("You are a Go expert").
    User("Explain channels in detail").
    ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
        if answer.Error != nil {
            fmt.Printf("Stream error: %v\n", answer.Error)
            return answer.Error
        }

        fmt.Print(answer.Message.Content)
        return nil
    })
```

## Structured Output

### Schema

Sets a schema for structured output from the language model.

```go
func (s *Squawk) Schema(schema map[string]any) *Squawk
```

Example:
```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    Schema(map[string]any{
        "type": "object",
        "properties": map[string]any{
            "name": map[string]any{
                "type": "string",
                "description": "The name of the person",
            },
            "age": map[string]any{
                "type": "integer",
                "description": "The age of the person",
            },
        },
    })
```

### StructuredOutput

Processes a structured chat completion request and handles the response through a callback function.

```go
func (s *Squawk) StructuredOutput(callBack func(answer llm.Answer, self *Squawk, err error)) *Squawk
```

Example:
```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    Schema(map[string]any{
        "type": "object",
        "properties": map[string]any{
            "functionName": map[string]any{
                "type": "string",
                "description": "Name of the function",
            },
            "complexity": map[string]any{
                "type": "string",
                "enum": ["O(1)", "O(n)", "O(n^2)"],
                "description": "Time complexity of the function",
            },
        },
    }).
    User("Analyze this function: func BubbleSort(arr []int) []int { ... }").
    StructuredOutput(func(answer llm.Answer, self *squawk.Squawk, err error) {
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        result := answer.StructuredOutput
        fmt.Printf("Function: %s\nComplexity: %s\n",
            result["functionName"],
            result["complexity"])
    })
```

## Embeddings and RAG

### Store

Sets the vector store for embeddings in the Squawk instance.

```go
func (s *Squawk) Store(store embeddings.VectorStore, optionalParameters ...string) *Squawk
```

Example:
```go
memStore := embeddings.NewMemoryVectorStore()
squawk.New().
    EmbeddingsModel("nomic-embed-text").
    Provider(provider.Ollama).
    Store(memStore)
```

### GenerateEmbeddings

Creates vector embeddings for a list of documents and stores them in the configured vector store.

```go
func (s *Squawk) GenerateEmbeddings(docs []string, optionalParameters ...any) *Squawk
```

Example:
```go
squawk.New().
    EmbeddingsModel("nomic-embed-text").
    Provider(provider.Ollama).
    Store(embeddings.NewMemoryVectorStore()).
    GenerateEmbeddings(
        []string{
            "Go is a statically typed language",
            "Go has built-in concurrency support",
        },
        true  // Enable logging
    )
```

### SimilaritySearch

Performs a semantic search against stored embeddings to find similar documents.

```go
func (s *Squawk) SimilaritySearch(content string, limit float64, max int, optionalParameters ...any) *Squawk
```

Example:
```go
squawk.New().
    EmbeddingsModel("nomic-embed-text").
    Provider(provider.Ollama).
    Store(embeddings.NewMemoryVectorStore()).
    GenerateEmbeddings([]string{
        "Go is a statically typed language",
        "Go supports concurrent programming",
    }).
    SimilaritySearch("concurrent Go features", 0.7, 3, true)
```

### SimilaritySearchFromUserMessage

Performs a semantic search using the content of a labeled user message as the search query.

```go
func (s *Squawk) SimilaritySearchFromUserMessage(userMessageLabel string, limit float64, max int, optionalParameters ...any) *Squawk
```

Example:
```go
squawk.New().
    EmbeddingsModel("nomic-embed-text").
    Provider(provider.Ollama).
    Store(embeddings.NewMemoryVectorStore()).
    GenerateEmbeddings([]string{
        "Go is a statically typed language",
        "Go supports concurrent programming",
    }).
    User("Tell me about Go concurrency", "question-1").
    SimilaritySearchFromUserMessage("question-1", 0.7, 3, true)
```

### AddSimilaritiesToMessages

Adds the context generated from similarity search results to the conversation as a system message.

```go
func (s *Squawk) AddSimilaritiesToMessages(optionalParameters ...string) *Squawk
```

Example:
```go
squawk.New().
    SimilaritySearch("concurrent Go features", 0.7, 3).
    AddSimilaritiesToMessages("context-1").
    User("Explain these concepts")
```

### AddSimilaritiesToMessagesWithPrefix

Adds the context generated from similarity search results to the conversation as a system message, with a custom prefix.

```go
func (s *Squawk) AddSimilaritiesToMessagesWithPrefix(prefix string, optionalParameters ...string) *Squawk
```

Example:
```go
squawk.New().
    SimilaritySearch("concurrent Go features", 0.7, 3).
    AddSimilaritiesToMessagesWithPrefix(
        "Use this documentation as context: \n\n",
        "context-1",
    ).
    User("Explain these concepts")
```

### Similarities

Returns the vector records from the most recent similarity search.

```go
func (s *Squawk) Similarities() []llm.VectorRecord
```

Example:
```go
similarities := squawk.Similarities()
for _, record := range similarities {
    fmt.Printf("Score: %.2f, Content: %s\n", 
        record.Score, 
        record.Content,
    )
}
```

### ContextFromSimilarities

Generates a formatted string containing the content from similarity search results.

```go
func (s *Squawk) ContextFromSimilarities() string
```

Example:
```go
context := squawk.ContextFromSimilarities()
fmt.Println("Retrieved context:", context)
```

## Function Calling with Tools

### Tools

Sets the list of available tools for function calling capabilities.

```go
func (s *Squawk) Tools(toolsList []llm.Tool) *Squawk
```

Example:
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
}

squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    Tools(toolsList)
```

### FunctionCalling

Executes a function calling request and processes the results through a callback function.

```go
func (s *Squawk) FunctionCalling(callBack func(answer llm.Answer, self *Squawk, err error)) *Squawk
```

Example:
```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    Tools(toolsList).
    User("Find all records about Go programming").
    FunctionCalling(func(answer llm.Answer, s *squawk.Squawk, err error) {
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        
        for _, call := range s.ToolCalls() {
            fmt.Printf("Function: %s\nArguments: %v\n", 
                call.Name, 
                call.Arguments,
            )
        }
    })
```

### ToolCalls

Returns the tool calls from the most recent function calling interaction.

```go
func (s *Squawk) ToolCalls() []llm.ToolCall
```

Example:
```go
calls := s.ToolCalls()
for _, call := range calls {
    fmt.Printf("Called: %s with %v\n", call.Name, call.Arguments)
}
```

## Meta Prompts

Squawk provides several methods that modify the input message to the LLM using predefined prompting templates.

### ForKids

Formats a message to be kid-friendly.

```go
func (s *Squawk) ForKids(message string, optionalParameters ...string) *Squawk
```

### Brief

Formats a message to get a brief response.

```go
func (s *Squawk) Brief(message string, optionalParameters ...string) *Squawk
```

### AsAPoem

Formats a message to get a response in the form of a poem.

```go
func (s *Squawk) AsAPoem(message string, optionalParameters ...string) *Squawk
```

### AdvantagesOnly

Formats a message to get only advantages in the response.

```go
func (s *Squawk) AdvantagesOnly(message string, optionalParameters ...string) *Squawk
```

### AsARecipe

Formats a message to get a response in the form of a recipe.

```go
func (s *Squawk) AsARecipe(message string, optionalParameters ...string) *Squawk
```

### Timeline

Formats a message to get a timeline in the response.

```go
func (s *Squawk) Timeline(message string, optionalParameters ...string) *Squawk
```

### Comparison

Formats a message to get a comparison in the response.

```go
func (s *Squawk) Comparison(message string, optionalParameters ...string) *Squawk
```

### Opinion

Formats a message to get an opinion in the response.

```go
func (s *Squawk) Opinion(message string, optionalParameters ...string) *Squawk
```

### Factual

Formats a message to get factual information in the response.

```go
func (s *Squawk) Factual(message string, optionalParameters ...string) *Squawk
```

### StepByStep

Formats a message to get a step-by-step explanation in the response.

```go
func (s *Squawk) StepByStep(message string, optionalParameters ...string) *Squawk
```

### ProsAndCons

Formats a message to get pros and cons in the response.

```go
func (s *Squawk) ProsAndCons(message string, optionalParameters ...string) *Squawk
```

### AsAStory

Formats a message to get a response in the form of a story.

```go
func (s *Squawk) AsAStory(message string, optionalParameters ...string) *Squawk
```

### InLaymansTerms

Formats a message to get a response in simple terms.

```go
func (s *Squawk) InLaymansTerms(message string, optionalParameters ...string) *Squawk
```

### Summarize

Formats a message to get a summary in the response.

```go
func (s *Squawk) Summarize(message string, optionalParameters ...string) *Squawk
```

### SummarizeLastAnswer

Formats a message to get a summary of the last answer in the response.

```go
func (s *Squawk) SummarizeLastAnswer(optionalParameters ...string) *Squawk
```

## Error Handling

### LastError

Gets or sets the most recent error encountered during processing.

```go
func (s *Squawk) LastError(optionalError ...error) error
```

Example:
```go
// Check for errors
if err := squawk.LastError(); err != nil {
    fmt.Printf("Last operation failed: %v\n", err)
}

// Set an error
squawk.LastError(errors.New("custom error"))
```

## Miscellaneous

### Cmd

Executes a custom command function on the Squawk instance.

```go
func (s *Squawk) Cmd(callBack func(self *Squawk)) *Squawk
```

Example:
```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    System("You are a Go expert").
    Cmd(func(self *squawk.Squawk) {
        fmt.Printf("Current model: %s\n", self.chatModel)
        fmt.Printf("Messages count: %d\n", len(self.Messages()))
    }).
    User("What is a channel?")
```

## Complete Examples

### Basic Chat Example

```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    BaseURL("http://localhost:11434").
    System("You are a Go expert").
    User("Explain concurrency in Go").
    Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        fmt.Println(answer.Message.Content)
    })
```

### Streaming Chat Example

```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    System("You are a Go expert").
    User("Explain channels in Go").
    ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
        if answer.Error != nil {
            fmt.Println("Error:", answer.Error)
            return answer.Error
        }
        fmt.Print(answer.Message.Content)
        return nil
    })
```

### RAG Example

```go
// Sample documents
docs := []string{
    "Go is a statically typed language",
    "Go supports concurrent programming",
}

// Initialize vector store
store := embeddings.NewMemoryVectorStore()

squawk.New().
    EmbeddingsModel("nomic-embed-text").
    Model("mistral:latest").
    Provider(provider.Ollama).
    Store(store).
    GenerateEmbeddings(docs, true).
    System("You are a Go expert").
    User("Tell me about Go concurrency", "question-1").
    SimilaritySearchFromUserMessage("question-1", 0.7, 3).
    AddSimilaritiesToMessages("context").
    ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
        fmt.Print(answer.Message.Content)
        return nil
    })
```

### Structured Output Example

```go
squawk.New().
    Model("mistral:latest").
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
        fmt.Println("Structured Output:", answer.StructuredOutput)
    })
```

### Function Calling Example

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
    Model("mistral:latest").
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
        self.System("RESULTS:\n"+results)
    }).
    User("Use the results and format the output with fancy emojis").
    ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
        fmt.Print(answer.Message.Content)
        return nil
    })
```

### Meta Prompts Example

```go
squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    ForKids("Explain Docker").
    ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
        fmt.Print(answer.Message.Content)
        return nil
    })

squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    Brief("Explain Docker").
    ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
        fmt.Print(answer.Message.Content)
        return nil
    })

squawk.New().
    Model("mistral:latest").
    Provider(provider.Ollama).
    InLaymansTerms("Explain Docker").
    Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
        fmt.Print(answer.Message.Content)
    }).
    SummarizeLastAnswer().
    ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
        fmt.Print(answer.Message.Content)
        return nil
    })
```