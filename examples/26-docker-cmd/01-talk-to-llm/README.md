# Make a Small Language Model Smarter: Teach Docker Commands

## Introduction

The topic of this blog post is to make a SLM smarter (SLM for Small Language model). 

Why would you want to use a SLM?
- You can run it on a Raspberry Pi (so, you don't need a GPU)
- The hosting costs are lower (no need for a GPU, size of the model is smaller)

I would say that running a SLM is more ecological and democratic.

So, today, I will try to teach some Docker commands to a SLM. My tools are:

- **Ollama** and its API
- For chatting: the [qwen2:0.5b LLM](https://ollama.com/library/qwen2:0.5b), a small language model of 352MB
- For the embedding: the [All-minilm:33m LLM](https://ollama.com/library/all-minilm:33m), a tiny language model of 67MB
- **Parakeet**, a simple and easy to use Golang wrapper around the Ollama API. I developed it to make it easier to use the Ollama API in my Golang projects instead of using Langchain for Go wich is more complex to use.

**As Parakeet is only a Golang wrapper around the Ollama API, you can reproduce this experiment in any language that you want and with other frameworks like LangChain.**

## Let's check the Docker knowledge of the qwen2:0.5b LLM

First, you will create a new Golang project and add the Parakeet package to your project. This project is a simple command-line application that interacts with a language model to translate user input into Docker commands.


```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	smallChatModel := "qwen2:0.5b"

    // 1️⃣ Prepare the system content
	systemContent := `instruction: translate this sentence in docker command - stay brief`

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 3.0,
		option.TopK: 10,
		option.TopP: 0.5,
	})

    // 2️⃣ Start the conversation
	for {
		question := input(smallChatModel)
		if question == "bye" {
			break
		}

		// 3️⃣ Prepare the query
		query := llm.Query{
			Model: smallChatModel,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "user", Content: question},
			},
			Options: options,
		}

		// 4️⃣ Answer the question (stream mode)
		_, err := completion.ChatStream(ollamaUrl, query,
			func(answer llm.Answer) error {
				fmt.Print(answer.Message.Content)
				return nil
			})

		if err != nil {
			log.Fatal("😡:", err)
		}

		fmt.Println()
	}
}

func input(smallChatModel string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("🐳 [%s] ask me something> ", smallChatModel)
	question, _ := reader.ReadString('\n')
	return strings.TrimSpace(question)
}
```


1. `systemContent` is a string that sets the instruction for the language model.
2. The program enters an infinite loop where it repeatedly prompts the user for input. `question := input()` calls the `input` function to read user input from the command line. If the user types `bye`, the loop breaks and the program exits.
3. Query Preparation: a query of type `llm.Query` is prepared with the following fields:
  - `Model`: The model identifier (`smallChatModel`).
  - `Messages`: A slice of `llm.Message` containing two messages:
    - A `system` message with the content `systemContent`.
    - A `user` message with the content of the user's `question`.
  - `Options`: Various options for the language model, such as `Temperature`, `RepeatLastN`, `RepeatPenalty`, `TopK`, and `TopP`. (more info about these options can be found in the [Ollama documentation](https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values))
4. Query Execution:
  - The `completion.ChatStream` function is called with `ollamaUrl` and the `query`.
  - A callback function is provided to handle the response (`answer` `llm.Answer`), which prints the content of the answer message to the console.


Ok, now, let's run the program and see how well the Qwen2:0.5b LLM can translate user input into Docker commands.

```bash
go run main.go
```

Then, ask the following questions:

- "Give me a list of all the local Docker images."
- "Give me a list of all containers, indicating their status as well."
- "List all containers with Ubuntu as their ancestor."

Well ..., the LLM is completely off the mark and is talking nonsense. We can say that the LLM is terrible at Docker.

The right answers to the questions are:

- "Give me a list of all the local Docker images."
    - `docker images`
- "Give me a list of all containers, indicating their status as well."
    - `docker ps -a`
- "List all containers with Ubuntu as their ancestor."
    - `docker ps --filter 'ancestor=ubuntu'`
