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

## Chat completion with conversational memeory

### In memory history

To store the messages in memory, use `history.MemoryMessages`

```golang
package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/history"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "tinydolphin" // fast, and perfect answer (short, brief)

	conversation := history.MemoryMessages{
		Messages: make(map[string]llm.MessageRecord),
	}

	systemContent := `You are an expert with the Star Trek series. use the history of the conversation to answer the question`

	userContent := `Who is James T Kirk?`

	options := llm.Options{
		Temperature: 0.5,
		RepeatLastN: 2,  
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	// Ask the question
	answer, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		},
	)
	if err != nil {
		log.Fatal("😡:", err)
	}

	// Save the conversation
	_, err = conversation.SaveMessage("1", llm.Message{
		Role:    "user",
		Content: userContent,
	})
	if err != nil {
		log.Fatal("😡:", err)
	}

	_, err = conversation.SaveMessage("2", llm.Message{
		Role:    "system",
		Content: answer.Message.Content,
	})

	if err != nil {
		log.Fatal("😡:", err)
	}

	// New question
	userContent = `Who is his best friend ?`

	previousMessages, _ := conversation.GetAllMessages()

	// (Re)Create the conversation
	conversationMessages := []llm.Message{}
	// instruction
	conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: systemContent})
	// history
	conversationMessages = append(conversationMessages, previousMessages...)
	// last question
	conversationMessages = append(conversationMessages, llm.Message{Role: "user", Content: userContent})

	query = llm.Query{
		Model:    model,
		Messages: conversationMessages,
		Options:  options,
	}

	answer, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		},
	)
	fmt.Println()
	if err != nil {
		log.Fatal("😡:", err)
	}

}
```

### Bbolt history

**[Bbolt](https://github.com/etcd-io/bbolt)** is an embedded key/value database for Go.

To store the messages in a bbolt bucket, use `history.BboltMessages`

```golang
conversation := history.BboltMessages{}
conversation.Initialize("../conversation.db")
```

> 👀 you will find a complete example in `examples/11-chat-conversational-bbolt`
> - `examples/11-chat-conversational-bbolt/begin`: start a conversation and save the history
> - `examples/11-chat-conversational-bbolt/resume`: load the messages from the history bucket and resue the conversation


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
// find the nearest vector
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

> 👀 you will find a complete example in `examples/09-embeddings-bbolt`
> - `examples/09-embeddings-bbolt/create-embeddings`: create and populate the vector store
> - `examples/09-embeddings-bbolt/use-embeddings`: search similarities in the vector store


## Create embeddings from text files and Similarity search

### Create embeddings
```golang
ollamaUrl := "http://localhost:11434"
embeddingsModel := "all-minilm"

store := embeddings.BboltVectorStore{}
store.Initialize("../embeddings.db")

// Parse all golang source code of the examples
// Create embeddings from documents and save them in the store
counter := 0
_, err := content.ForEachFile("../../examples", ".go", func(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fmt.Println("📝 Creating embedding from:", path)
	counter++
	embedding, err := embeddings.CreateEmbedding(
		ollamaUrl,
		llm.Query4Embedding{
			Model:  embeddingsModel,
			Prompt: string(data),
		},
		strconv.Itoa(counter), // don't forget the id (unique identifier)
	)
	fmt.Println("📦 Created: ", len(embedding.Embedding))

	if err != nil {
		fmt.Println("😡:", err)
	} else {
		_, err := store.Save(embedding)
		if err != nil {
			fmt.Println("😡:", err)
		}
	}
	return nil
})
if err != nil {
	log.Fatalln("😡:", err)
}
```

### Similarity search

```golang
ollamaUrl := "http://localhost:11434"
embeddingsModel := "all-minilm"
chatModel := "magicoder:latest"

store := embeddings.BboltVectorStore{}
store.Initialize("../embeddings.db")

systemContent := `You are a Golang developer and an expert in computer programming.
Please make friendly answer for the noobs. Use the provided context and doc to answer.
Add source code examples if you can.`

// Question for the Chat system
userContent := `How to create a stream chat completion with Parakeet?`

// Create an embedding from the user question
embeddingFromQuestion, err := embeddings.CreateEmbedding(
	ollamaUrl,
	llm.Query4Embedding{
		Model:  embeddingsModel,
		Prompt: userContent,
	},
	"question",
)
if err != nil {
	log.Fatalln("😡:", err)
}
fmt.Println("🔎 searching for similarity...")

similarities, _ := store.SearchSimilarities(embeddingFromQuestion, 0.3)

// Generate the context from the similarities
// This will generate a string with a content like this one:
// `<context><doc>...<doc><doc>...<doc></context>`
documentsContent := embeddings.GenerateContextFromSimilarities(similarities)

fmt.Println("🎉 similarities", len(similarities))

query := llm.Query{
	Model: chatModel,
	Messages: []llm.Message{
		{Role: "system", Content: systemContent},
		{Role: "system", Content: documentsContent},
		{Role: "user", Content: userContent},
	},
	Options: llm.Options{
		Temperature: 0.4,
		RepeatLastN: 2,
	},
	Stream: false,
}

fmt.Println("")
fmt.Println("🤖 answer:")

// Answer the question
_, err = completion.ChatStream(ollamaUrl, query,
	func(answer llm.Answer) error {
		fmt.Print(answer.Message.Content)
		return nil
	})

if err != nil {
	log.Fatal("😡:", err)
}
```

## Demos

- https://github.com/parakeet-nest/parakeet-demo
- https://github.com/parakeet-nest/tiny-genai-stack
