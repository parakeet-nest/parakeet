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

	/*
	{
		"input": "Give me a list of all containers, indicating their status as well.",
		"instruction": "translate this sentence in docker command",
		"output": "docker ps -a"
	},
	{
		"input": "List all containers with Ubuntu as their ancestor.",
		"instruction": "translate this sentence in docker command",
		"output": "docker ps --filter 'ancestor=ubuntu'"
	},
	{
		"input": "Give me a list of all the local Docker images.",
		"instruction": "translate this sentence in docker command",
		"output": "docker images"
	},
	*/

	systemContent := `instruction: 
	translate the user question in docker command using the given context.
	Stay brief.`

	contextContent := `<context>
		<doc>
		input: Give me a list of all containers, indicating their status as well.
		output: docker ps -a
		</doc>
		<doc>
		input: List all containers with Ubuntu as their ancestor.
		output: docker ps --filter 'ancestor=ubuntu'
		</doc>
		<doc>
		input: Give me a list of all the local Docker images.
		output: docker images
		</doc>
	</context>
	`

	for {
		question := input(smallChatModel)
		if question == "bye" {
			break
		}

		// Prepare the query
		query := llm.Query{
			Model: smallChatModel,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: contextContent},
				{Role: "user", Content: question},
			},
			Options: llm.Options{
				Temperature:   0.0,
				RepeatLastN:   2,
				RepeatPenalty: 3.0,
				TopK:          10,
				TopP:          0.5,
			},
		}

		// Answer the question
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
