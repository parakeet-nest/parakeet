<!-- TOPIC: Parakeet - A Go Library for Creating GenAI Apps SUMMARY: Parakeet is a simple Go library used to create text-based GenAI apps, allowing users to generate new content based on training data. KEYWORDS: Parakeet, GenAI, Go, Library, Text Generation, AI -->

# 🦜🪺 Parakeet

Parakeet is the simplest Go library to create **GenAI apps** with **[Ollama](https://ollama.com/)**.

> A GenAI app is an application that uses generative AI technology. Generative AI can create new text, images, or other content based on what it's been trained on. So a GenAI app could help you write a poem, design a logo, or even compose a song! These are still under development, but they have the potential to be creative tools for many purposes. - [Gemini](https://gemini.google.com)

> ✋ Parakeet is only for creating GenAI apps generating **text** (not image, music,...).

## Install

```bash
go get github.com/parakeet-nest/parakeet
```
<!-- split -->

<!-- TOPIC: Simple Completion in Golang using Parakeet and LLaMA SUMMARY: This code snippet demonstrates the use of simple completion in Golang to generate a response for a given prompt with a provided model, specifically using Parakeet and LLaMA. KEYWORDS: Golang, Parakeet, LLaMA, Simple Completion, AI-powered Text Generation -->
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

	question := llm.GenQuery{
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
<!-- split -->


<!-- TOPIC: Golang programming and Stream completion with LLaMA model SUMMARY: This code snippet demonstrates the use of LLaMA model for generating a response to a given question using Go language. The code sets up an LLaMA connection, defines a query with a prompt and options, and then generates a stream of answers. KEYWORDS: Golang, LLaMA, Stream completion, Natural Language Processing -->
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

	question := llm.GenQuery{
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
<!-- split -->

<!-- TOPIC: Contextual completion with Ollama SUMMARY: The code demonstrates the use of Ollama's API to generate completions in a conversational context. KEYWORDS: Ollama, contextual completion, conversation, API, tinydolphin, James T Kirk, best friend -->
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

	firstQuestion := llm.GenQuery{
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

	secondQuestion := llm.GenQuery{
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
<!-- split -->

<!-- TOPIC: Chat Completion SUMMARY: The chat completion feature is used to generate a conversational response for a given set of messages with a provided model. KEYWORDS: Go, Golang, Parakeet, Conversational AI, Chat Completion, BBolt -->
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

<!-- split -->

<!-- TOPIC: Chat Completion with Stream using Golang and LLaMA API SUMMARY: This Go program uses the LLaMA API to create a chat completion stream, generating responses based on user input and system content. It provides a basic "hello world" example in Golang. KEYWORDS: Golang, LLaMA API, Chat Completion, Stream, Programming -->
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
<!-- split -->

## Protected endpoint
If your Ollama endpoint is protected with a header token, you can specify the token like this:

```golang
query := llm.Query{
    Model: model,
    Messages: []llm.Message{
        {Role: "system", Content: systemContent},
        {Role: "user", Content: userContent},
    },
    Options: options,
    TokenHeaderName: "X-TOKEN",
    TokenHeaderValue: "john_doe",
}
```


## Verbose mode

You can activate the "verbose mode" with all kinds of completions.


```golang
options := llm.Options{
    Temperature: 0.5,
    RepeatLastN: 2,
    RepeatPenalty: 2.0,
    Verbose: true,
}
```
You will get an output like this (with the query and the completion):

```json
[llm/query] {
  "model": "deepseek-coder",
  "messages": [
    {
      "role": "system",
      "content": "You are an expert in computer programming.\n\tPlease make friendly answer for the noobs.\n\tAdd source code examples if you can."
    },
    {
      "role": "user",
      "content": "I need a clear explanation regarding the following question:\n\tCan you create a \"hello world\" program in Golang?\n\tAnd, please, be structured with bullet points"
    }
  ],
  "options": {
    "repeat_last_n": 2,
    "temperature": 0.5,
    "repeat_penalty": 2,
    "Verbose": true
  },
  "stream": false,
  "prompt": "",
  "context": null,
  "tools": null,
  "TokenHeaderName": "",
  "TokenHeaderValue": ""
}

[llm/completion] {
  "model": "deepseek-coder",
  "message": {
    "role": "assistant",
    "content": "Sure, here's a simple \"Hello, World!\" program in Golang.\n\t1. First, you need to have Golang installed on your machine.\n\t2. Open your text editor, and write the following code:\n\t```go\n\tpackage main\n\timport \"fmt\"\n\tfunc main() {\n\t    fmt.Println(\"Hello, World!\")\n\t} \n\t```\n\t3. Save the file with a `.go` extension (like `hello.go`).\n\t4. In your terminal, navigate to the directory containing the `.go` file.\n\t5. Run the program with the command:\n\t```\n\tgo run hello.go\n\t```\n\t6. If everything goes well, you should see \"Hello, World!\" printed in your terminal.\n\t7. If there's an error, you will see the error message.\n\t8. If everything is correct, you'll see \"Hello, World!\" printed in your terminal.\n"
  },
  "done": true,
  "response": "",
  "context": null,
  "created_at": "2024-08-19T05:57:23.979361Z",
  "total_duration": 3361191958,
  "load_duration": 2044932125,
  "prompt_eval_count": 79,
  "prompt_eval_duration": 95034000,
  "eval_count": 222,
  "eval_duration": 1216689000
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

> 👀 you will find a complete example in `examples/11-chat-conversational-bbolt`
> - `examples/11-chat-conversational-bbolt/begin`: start a conversation and save the history
> - `examples/11-chat-conversational-bbolt/resume`: load the messages from the history bucket and resue the conversation
<!-- split -->

<!-- TOPIC: Embeddings and Vector Stores SUMMARY: This document discusses creating embeddings, storing them in an efficient way using a vector store, and searching for similar embeddings. KEYWORDS: Embeddings, Vector Store, MemoryVectorStore, BBolt -->
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
<!-- split -->

### Redis vector store

**Create a store, and open an existing store**:
```golang
redisStore := embeddings.RedisVectorStore{}
err := redisStore.Initialize("localhost:6379", "", "chronicles-bucket")

if err != nil {
	log.Fatalln("😡:", err)
}
```

> 👀 you will find a complete example in `examples/32-rag-with-redis`
> - `examples/32-rag-with-redis/create-embeddings`: create and populate the vector store
> - `examples/32-rag-with-redis/use-embeddings`: search similarities in the vector store

### Elasticsearch vector store

**Create a store, and open an existing store**:
```golang
cert, _ := os.ReadFile(os.Getenv("ELASTIC_CERT_PATH"))

elasticStore := embeddings.ElasticSearchStore{}
err := elasticStore.Initialize(
	[]string{
		os.Getenv("ELASTIC_ADDRESS"),
	},
	os.Getenv("ELASTIC_USERNAME"),
	os.Getenv("ELASTIC_PASSWORD"),
	cert,
	"chronicles-index",
)
```

> 👀 you will find a complete example in `examples/33-rag-with-elastic`
> - `examples/33-rag-with-elastic/create-embeddings`: create and populate the vector store
> - `examples/33-rag-with-elastic/use-embeddings`: search similarities in the vector store

### Additional data

you can add additional data to a vector record (embedding):

```golang
embedding.Text()
embedding.Reference()
embedding.MetaData()
```


<!-- TOPIC: Natural Language Processing, Computer Programming, Golang SUMMARY: This document discusses the process of creating embeddings from text files and performing similarity searches using a Bolt Vector Store. The document also demonstrates how to use these embeddings to generate context for a chat system. KEYWORDS: Embeddings, Similarity Search, Bolt Vector Store, Natural Language Processing, Computer Programming, Golang -->
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

### Other similarity search methods

`SearchMaxSimilarity` searches for the vector record in the `BboltVectorStore` that has the maximum **cosine distance similarity** to the given `embeddingFromQuestion`:
```golang
similarity, _ := store.SearchMaxSimilarity(embeddingFromQuestion)
```

`SearchTopNSimilarities` searches for vector records in the `BboltVectorStore` that have a **cosine distance similarity** greater than or equal to the given `limit` and returns the top `n` records:
```golang
similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, limit, n)
```

<!-- split -->
## Chunkers and Splitters

There are three methods in the `content` package to help you chunk and split text:

- `ChunkText` takes a text string and divides it into chunks of a specified size with a given overlap. It returns a slice of strings, where each string represents a chunk of the original text.

```golang
chunks := content.ChunkText(documentContent, 900, 400)
```

- `SplitTextWithDelimiter` splits the given text using the specified delimiter and returns a slice of strings.

```golang
chunks := content.SplitTextWithDelimiter(documentContent, "<!-- SPLIT -->")
```

- `SplitTextWithRegex` splits the given text using the provided regular expression delimiter. It returns a slice of strings containing the split parts of the text.

```golang
chunks := content.SplitTextWithRegex(documentContent, `## *`)
```

- `SplitMarkdownBySections` splits the given markdown text using the title sections (`#, ##, etc.`) and returns a slice of strings.

```golang
chunks := content.SplitMarkdownBySections(documentContent)
```

- `SplitAsciiDocBySections` splits the given asciidoc text using the title sections (`=, ==, etc.`) and returns a slice of strings.

```golang
chunks := content.SplitAsciiDocBySections(documentContent)
```

- `SplitHTMLBySections` splits the given html text using the title sections (`h1, h2, h3, h4, h5, h6`) and returns a slice of strings.

```golang
chunks := content.SplitHTMLBySections(documentContent)
```

<!-- split -->


<!-- TOPIC: Function Calling SUMMARY: A feature in LLMs that allows them to provide a specific output with the same format (predictable output format). KEYWORDS: function calling, predictable output format, LLMs. -->
## Function Calling (before tool support)

What is **"Function Calling"**? First, it's not a feature where a LLM can call and execute a function. "Function Calling" is the ability for certain LLMs to provide a specific output with the same format (we could say: "a predictable output format").

So, the principle is simple:

- You (or your GenAI application) will create a prompt with a delimited list of tools (the functions) composed by name, descriptions, and parameters: `SayHello`, `AddNumbers`, etc.
- Then, you will add your question ("Hey, say 'hello' to Bob!") to the prompt and send all of this to the LLM.
- If the LLM "understand" that the `SayHello` function can be used to say "hello" to Bob, then the LLM will answer with only the name of the function with the parameter(s). For example: `{"name":"SayHello","arguments":{"name":"Bob"}}`.

Then, it will be up to you to implement the call of the function.

The [latest version (v0.3) of Mistral 7b](https://ollama.com/library/mistral:7b) supports function calling and is available for Ollama.

### Define a list of tools

First, you have to provide the LLM with a list of tools with the following format:

```golang
toolsList := []llm.Tool{
	{
		Type: "function",
		Function: llm.Function{
			Name:        "hello",
			Description: "Say hello to a given person with his name",
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
```

### Generate a prompt from the tools list and the user instructions

The `tools.GenerateContent` method generates a string with the tools in JSON format surrounded by `[AVAILABLE_TOOLS]` and `[/AVAILABLE_TOOLS]`:
```golang
toolsContent, err := tools.GenerateContent(toolsList)
if err != nil {
	log.Fatal("😡:", err)
}
```


The `tools.GenerateInstructions` method generates a string with the user instructions surrounded by `[INST]` and `[/INST]`:
```golang
userContent := tools.GenerateInstructions(`say "hello" to Bob`)
```

Then, you can add these two strings to the messages list:
```golang
messages := []llm.Message{
	{Role: "system", Content: toolsContent},
	{Role: "user", Content: userContent},
}
```

### Send the prompt (messages) to the LLM

It's important to set the `Temperature` to `0.0`:
```golang
options := llm.Options{
	Temperature:   0.0,
	RepeatLastN:   2,
	RepeatPenalty: 2.0,
}

You must set the `Format` to `json` and `Raw` to `true`:
query := llm.Query{
	Model: model,
	Messages: messages,
	Options: options,
	Format:  "json",
	Raw:     true,
}
```
> When building the payload to be sent to Ollama, we need to set the `Raw` field to true, thanks to that, no formatting will be applied to the prompt (we override the prompt template of Mistral), and we need to set the `Format` field to `"json"`.

No you can call the `Chat` method. The answer of the LLM will be in JSON format:
```golang
answer, err := completion.Chat(ollamaUrl, query)
if err != nil {
	log.Fatal("😡:", err)
}
// PrettyString is a helper that prettyfies the JSON string
result, err := gear.PrettyString(answer.Message.Content)
if err != nil {
	log.Fatal("😡:", err)
}
fmt.Println(result)
```

You should get this answer:
```json
{
  "name": "hello",
  "arguments": {
    "name": "Bob"
  }
}
```

You can try with the other tool (or function):
```golang
userContent := tools.GenerateInstructions(`add 2 and 40`)
```

You should get this answer:
```json
{
  "name": "addNumbers",
  "arguments": {
    "a": 2,
    "b": 40
  }
}
```

> **Remark**: always test the format of the output, even if Mistral is trained for "function calling", the result are not entirely predictable.

Look at this sample for a complete sample: [examples/15-mistral-function-calling](examples/15-mistral-function-calling)
<!-- split -->

<!-- TOPIC: Function Calling with LLMs that do not implement Function Calling SUMMARY: A technique to reproduce function calling feature in LLMs without native support by adding specific messages at the beginning and end of the conversation. KEYWORDS: phi3, mini, golang, JSON, tool calling, argument passing -->
## Function Calling with LLMs that do not implement Function Calling

It is possible to reproduce this feature with some LLMs that do not implement the "Function Calling" feature natively, but we need to supervise them and explain precisely what we need. The result (the output) will be less predictable, so you will need to add some tests before using the output, but with some "clever" LLMs, you will obtain correct results. I did my experiments with **[phi3:mini](https://ollama.com/library/phi3:mini)**.

The trick is simple:

Add this message at the begining of the list of messages:
```golang
systemContentIntroduction := `You have access to the following tools:`
```

Add this message at the end of the list of messages, just before the user message:
```golang
systemContentInstructions := `If the question of the user matched the description of a tool, the tool will be called.
To call a tool, respond with a JSON object with the following structure: 
{
	"name": <name of the called tool>,
	"arguments": {
	<name of the argument>: <value of the argument>
	}
}

search the name of the tool in the list of tools with the Name field
`
```

At the end, you will have this:
```golang
messages := []llm.Message{
	{Role: "system", Content: systemContentIntroduction},
	{Role: "system", Content: toolsContent},
	{Role: "system", Content: systemContentInstructions},
	{Role: "user", Content: `say "hello" to Bob`},
}
```

Look at this sample for a complete sample: [examples/17-fake-function-calling](examples/17-fake-function-calling)
<!-- split -->


## Function Calling with tool support
> Ollama API: chat request with tools https://github.com/ollama/ollama/blob/main/docs/api.md#chat-request-with-tools

Since Ollama `0.3.0`, Ollama supports **tools calling**, blog post: https://ollama.com/blog/tool-support.
A list of supported models can be found under the Tools category on the models page: https://ollama.com/search?c=tools


### Define a list of tools

> use a supported model
```golang
model := "mistral:7b"

toolsList := []llm.Tool{
	{
		Type: "function",
		Function: llm.Function{
			Name:        "hello",
			Description: "Say hello to a given person with his name",
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
```

### Set the Tools property of the query

> - set the `Temperature` to `0.0`
> - you don't need to set the row mode to true
> - set `query.Tools` with `toolsList`
```golang
messages := []llm.Message{
	{Role: "user", Content: `say "hello" to Bob`},
}

options := llm.Options{
	Temperature:   0.0,
	RepeatLastN:   2,
	RepeatPenalty: 2.0,
}

query := llm.Query{
	Model:    model,
	Messages: messages,
	Tools:    toolsList,
	Options:  options,
	Format:   "json",
}
```

### Run the completion

```go
answer, err := completion.Chat(ollamaUrl, query)
if err != nil {
	log.Fatal("😡:", err)
}

// It's a []map[string]interface{}
toolCalls := answer.Message.ToolCalls

// Convert toolCalls into a JSON string
jsonBytes, err := json.Marshal(toolCalls)
if err != nil {
	log.Fatal("😡:", err)
}
// Convert JSON bytes to string
result := string(jsonBytes)

fmt.Println(result)
```

The result will look like this:
```json
[{"function":{"arguments":{"name":"Bob"},"name":"hello"}}]
```

#### Or you can use the `ToolCallsToJSONString` helper

```golang
answer, err = completion.Chat(ollamaUrl, query)
if err != nil {
	log.Fatal("😡:", err)
}

result, err = answer.Message.ToolCallsToJSONString()
if err != nil {
	log.Fatal("😡:", err)
}
fmt.Println(result)
```
The result will look like this:
```json
[{"function":{"arguments":{"name":"Bob"},"name":"hello"}}]
```

Look here for a complete sample: [examples/19-mistral-function-calling-tool-support](examples/19-mistral-function-calling-tool-support)

#### Or (better) you can use the `ToolCallsToJSONString` helper

```golang
answer, err := completion.Chat(ollamaUrl, query)
if err != nil {
	log.Fatal("😡:", err)
}

result, err := answer.Message.ToolCalls[0].Function.ToJSONString()
if err != nil {
	log.Fatal("😡:", err)
}
fmt.Println(result)
```

The result will look like this:
```json
{"name":"hello","arguments":{"name":"Bob"}}
```

Look at these samples:
- [examples/43-function-calling/01-xlam](examples/43-function-calling/01-xlam)
- [examples/43-function-calling/02-qwen2tools](examples/43-function-calling/02-qwen2tools)

<!-- split -->

<!-- TOPIC: WebAssembly plugins for Parakeet SUMMARY: The release of Parakeet's version 0.0.6 brings support for WebAssembly, allowing users to write their own wasm plugins in various languages (Rust, Go, C, etc.) and use them with the "Function Calling" feature. KEYWORDS: Parakeet, WebAssembly, Wasm plugins, Extism project, TinyGo, Function Calling -->
## Wasm plugins

The release `0.0.6` of Parakeet brings the support of **WebAssembly** thanks to the **[Extism project](https://extism.org/)**. That means you can write your own wasm plugins for Parakeet to add new features (for example, a chunking helper for doing RAG) with various languages (Rust, Go, C, ...).

Or you can use the Wasm plugins with the "Function Calling" feature, which is implemented in Parakeet.

You can find an example of "Wasm Function Calling" in [examples/18-call-functions-for-real](examples/18-call-functions-for-real) - the wasm plugin is located in the `wasm` folder and it is built with **[TinyGo](https://tinygo.org/)**.

🚧 more samples to come.
<!-- split -->

<!-- TOPIC: Other Parakeet methods for interacting with models SUMMARY: This document describes two methods used to interact with models in Parakeet, including retrieving information about a model and pulling a model. KEYWORDS: Parakeet, models, API, ShowModelInformation, PullModel -->
## Other Parakeet methods

### Get Information about a model

```golang
llm.ShowModelInformation(url, model string) (llm.ModelInformation, int, error)
```

`ShowModelInformation` retrieves information about a model from the specified URL.

**Parameters**:
  - url: the base URL of the API.
  - model: the name of the model to retrieve information for.

**Returns**:
  - ModelInformation: the information about the model.
  - int: the HTTP status code of the response.
  - error: an error if the request fails.

**✋ Remark**: if the model does not exist, it will return an error with a status code of 404.

### Pull a model

```golang
llm.PullModel(url, model string) (llm.PullResult, int, error)
```

`PullModel` sends a POST request to the specified URL to pull a model with the given name.

**Parameters**:
  - url: The URL to send the request to.
  - model: The name of the model to pull.

**Returns**:
  - PullResult: The result of the pull operation.
  - int: The HTTP status code of the response.
  - error: An error if the request fails.
<!-- split -->

<!-- TOPIC: Prompt helpers and meta prompts SUMMARY: A collection of special instructions, known as meta-prompts, to guide language models in generating specific kinds of responses. KEYWORDS: Meta prompts, prompt helpers, AI, LLM, natural language processing, NLP -->
### Prompt helpers

#### Meta prompts
> package: `prompt`

Meta-prompts are special instructions embedded within a prompt to guide a language model in generating a specific kind of response.

|  Meta-Prompt   |  Purpose  |
| :------------  | :-------- |
|[Brief] What is AI? | For a concise answer
|[In Layman’s Terms] Explain LLM | For a simplified explanation
|[As a Story] Describe the evolution of cars | To get the information in story form
|[Pros and Cons] Is AI useful? | For a balanced view with advantages and disadvantages
|[Step-by-Step] How to do a smart prompt? | For a detailed, step-by-step guide
|[Factual] What is the best pizza of the world? | For a straightforward, factual answer
|[Opinion] What is the best pizza of the world? | To get an opinion-based answer
|[Comparison] Compare pineapple pizza to pepperoni pizza | For a comparative analysis
|[Timeline] What are the key milestones to develop a WebApp? | For a chronological account of key events
|[As a Poem] How to cook a cake? | For a poetic description
|[For Kids] How to cook a cake? | For a child-friendly explanation
|[Advantages Only] What are the benefits of AI? | To get a list of only the advantages
|[As a Recipe] How to cook a cake? | To receive the information in the form of a recipe

##### Meta prompts methods

- `prompt.Brief(s string) string`
- `prompt.InLaymansTerms(s string) string`
- `prompt.AsAStory(s string) string`
- `prompt.ProsAndCons(s string) string`
- `prompt.StepByStep(s string) string`
- `prompt.Factual(s string) string`
- `prompt.Opinion(s string) string`
- `prompt.Comparison(s string) string`
- `prompt.Timeline(s string) string`
- `prompt.AsAPoem(s string) string`
- `prompt.ForKids(s string) string`
- `prompt.AdvantagesOnly(s string) string`
- `prompt.AsARecipe(s string) string`
<!-- split -->

<!-- TOPIC: Parakeet Demos and Blog Posts SUMMARY: A collection of Parakeet demos and blog posts showcasing the ease of creating GenAI applications with Ollama, Golang, and other tools. KEYWORDS: Parakeet, GenAI, Ollama, Golang, function calling, RAG, Mistral 7B, Bash, Jq, LLMs -->
## Parakeet Demos

- https://github.com/parakeet-nest/parakeet-demo
- https://github.com/parakeet-nest/tiny-genai-stack

## Blog Posts

- [Parakeet, an easy way to create GenAI applications with Ollama and Golang](https://k33g.hashnode.dev/parakeet-an-easy-way-to-create-genai-applications-with-ollama-and-golang)
- [Understand RAG with Parakeet](https://k33g.hashnode.dev/understand-rag-with-parakeet)
-[Function Calling with Ollama, Mistral 7B, Bash and Jq](https://k33g.hashnode.dev/function-calling-with-ollama-mistral-7b-bash-and-jq)
- [Function Calling with Ollama and LLMs that do not support function calling](https://k33g.hashnode.dev/function-calling-with-ollama-and-llms-that-do-not-support-function-calling)

