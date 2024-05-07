# 🦜🪺 Parakeet

Parakeet is the simplest Go library to create **GenAI apps** with **[Ollma](https://ollama.com/)**.

> A GenAI app is an application that uses generative AI technology. Generative AI can create new text, images, or other content based on what it's been trained on. So a GenAI app could help you write a poem, design a logo, or even compose a song! These are still under development, but they have the potential to be creative tools for many purposes. - [Gemini](https://gemini.google.com)

> ✋ Parakeet is only for creating GenAI apps generating **text** (not image, music,...).

## Install

```bash
go get github.com/parakeet-nest/parakeet
```

## Simple completion

The simple completion can be used to generate a response for a given prompt with a provided model.

```golang
package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "tinydolphin"

	options := llm.Options{
		Temperature: 0.5,  // default (0.8)
	}

	question := llm.Query{
		Model: model,
		Prompt: "Who is James T Kirk?",
		Options: options,
	}

	answer, err := completion.Generate(ollamaUrl, question)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)
}
```

### Simple completion with stream

```golang
package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "tinydolphin"

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
	}

	question := llm.Query{
		Model: model,
		Prompt: "Who is James T Kirk?",
		Options: options,
	}
	
	answer, err := completion.GenerateStream(ollamaUrl, question,
		func(answer llm.Answer) error {
			fmt.Print(answer.Response)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}
}

```

### Completion with context
> see: https://github.com/ollama/ollama/blob/main/docs/api.md#generate-a-completion

> The context can be used to keep a short conversational memory for the next completion.

```golang
package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "tinydolphin"

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
	}

	firstQuestion := llm.Query{
		Model: model,
		Prompt: "Who is James T Kirk?",
		Options: options,
	}

	answer, err := completion.Generate(ollamaUrl, firstQuestion)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)

	fmt.Println()

	secondQuestion := llm.Query{
		Model: model,
		Prompt: "Who is his best friend?",
		Context: answer.Context,
		Options: options,
	}

	answer, err = completion.Generate(ollamaUrl, secondQuestion)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)
}
```

## Chat completion

The chat completion can be used to generate a conversational response for a given set of messages with a provided model.

```golang
package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "deepseek-coder"

	systemContent := `You are an expert in computer programming.
	Please make friendly answer for the noobs.
	Add source code examples if you can.`

	userContent := `I need a clear explanation regarding the following question:
	Can you create a "hello world" program in Golang?
	And, please, be structured with bullet points`

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
		RepeatLastN: 2, // default (64)
		RepeatPenalty: 2.0, // default (1.1)
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
		Stream: false,
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.Content)
}
```

✋ **To keep a conversational memory** for the next chat completion, update the list of messages with the previous question and answer.
> I plan to add the support of [bbolt](https://github.com/etcd-io/bbolt) in the incoming v0.0.1 of Parakeet to store the conversational memory.

### Chat completion with stream

```golang
package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "deepseek-coder"

	systemContent := `You are an expert in computer programming.
	Please make friendly answer for the noobs.
	Add source code examples if you can.`

	userContent := `I need a clear explanation regarding the following question:
	Can you create a "hello world" program in Golang?
	And, please, be structured with bullet points`

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
		RepeatLastN: 2, // default (64) 
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
		Stream:  false,
	}

	_, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}
}
```

## Embeddings

### Create embeddings

```golang
embedding, err := embeddings.CreateEmbedding(
	ollamaUrl,
	llm.Query4Embedding{
		Model:  "all-minilm",
		Prompt: "Jean-Luc Picard is a fictional character in the Star Trek franchise.",
	},
	"Picard", // identifier
)
```

## Vector stores

A vector store allows to store and search for embeddings in an efficient way.

### In memory vector store

**Create a store**:
```golang
store := embeddings.MemoryVectorStore{
	Records: make(map[string]llm.VectorRecord),
}
```

**Save embeddings**:
```golang
store.Save(embedding)
```

**Search embeddings**:
```golang
embeddingFromQuestion, err := embeddings.CreateEmbedding(
	ollamaUrl,
	llm.Query4Embedding{
		Model:  "all-minilm",
		Prompt: "Who is Jean-Luc Picard?",
	},
	"question",
)
similarity, _ := store.SearchMaxSimilarity(embeddingFromQuestion)

documentsContent := `<context><doc>` + similarity.Prompt + `</doc></context>`
```

> 👀 you will find a complete example in `examples/08-embeddings`

### Bbolt vector store

**[Bbolt](https://github.com/etcd-io/bbolt)** is an embedded key/value database for Go.

**Create a store, and open an existing store**:
```golang
store := embeddings.BboltVectorStore{}
store.Initialize("../embeddings.db")
```

> 👀 you will find a complete example in `examples/09-embeddings`
> - `examples/09-embeddings/create-embeddings`: create and populate the vector store
> - `examples/09-embeddings/use-embeddings`: search similarities in the vector store

## Demos

- https://github.com/parakeet-nest/parakeet-demo
- https://github.com/parakeet-nest/tiny-genai-stack
