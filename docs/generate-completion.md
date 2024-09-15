<!-- TOPIC: Simple Completion in Golang using Parakeet and LLaMA SUMMARY: This code snippet demonstrates the use of simple completion in Golang to generate a response for a given prompt with a provided model, specifically using Parakeet and LLaMA. KEYWORDS: Golang, Parakeet, LLaMA, Simple Completion, AI-powered Text Generation -->
# Generate completion

## Completion

The simple completion can be used to generate a response for a given prompt with a provided model.

```golang
package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"

	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "tinydolphin"

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.5,
	})

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
## Completion with stream

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
		Temperature: 0.5,
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
## Completion with context
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
		Temperature: 0.5,
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