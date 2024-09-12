<!-- TOPIC: Chat Completion SUMMARY: The chat completion feature is used to generate a conversational response for a given set of messages with a provided model. KEYWORDS: Go, Golang, Parakeet, Conversational AI, Chat Completion, BBolt -->
# Chat completion

## Completion

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
		Temperature: 0.5, 
		RepeatLastN: 2, 
		RepeatPenalty: 2.0,
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

<!-- split -->

<!-- TOPIC: Chat Completion with Stream using Golang and LLaMA API SUMMARY: This Go program uses the LLaMA API to create a chat completion stream, generating responses based on user input and system content. It provides a basic "hello world" example in Golang. KEYWORDS: Golang, LLaMA API, Chat Completion, Stream, Programming -->
## Completion with stream

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

<!-- TOPIC: Chat completion with conversational memory SUMMARY: A Go program that uses the Parakeet library to store messages in memory and complete conversations using a conversational memory history. KEYWORDS: Parakeet, conversational memory, chat completion, Go programming language -->
## Chat completion with conversational memory

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
<!-- split -->

<!-- TOPIC: Bbolt history and usage in Go SUMMARY: A brief introduction to using bbolt, an embedded key-value database for Go, for storing message histories. KEYWORDS: bbolt, Go, key-value database, message history, Golang -->
### Bbolt history

**[Bbolt](https://github.com/etcd-io/bbolt)** is an embedded key/value database for Go.

To store the messages in a bbolt bucket, use `history.BboltMessages`

```golang
conversation := history.BboltMessages{}
conversation.Initialize("../conversation.db")
```

!!! note
	👀 you will find a complete example in:

    - [examples/11-chat-conversational-bbolt](https://github.com/parakeet-nest/parakeet/tree/main/examples/11-chat-conversational-bbolt)
    - [examples/11-chat-conversational-bbolt/begin](https://github.com/parakeet-nest/parakeet/tree/main/examples/11-chat-conversational-bbolt/begin): start a conversation and save the history
    - [examples/11-chat-conversational-bbolt/resume](https://github.com/parakeet-nest/parakeet/tree/main/examples/11-chat-conversational-bbolt/resume): load the messages from the history bucket and resue the conversation

