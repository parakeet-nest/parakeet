package main

/*
Example from https://docs.anthropic.com/claude/page/code-clarifier
*/
import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	//ollamaUrl := "http://localhost:11434"
	// if working from a container
	ollamaUrl := "http://host.docker.internal:11434"
	//model := "deepseek-coder:instruct"
	model := "gemma:2b-instruct"

	systemContent := `Your task is to take the code snippet provided and explain it in simple, easy-to-understand language. 
    Break down the code's functionality, purpose, and key components. 
    Use analogies, examples, and plain terms to make the explanation accessible to someone with minimal coding knowledge. 
    Avoid using technical jargon unless absolutely necessary, and provide clear explanations for any jargon used. 
    The goal is to help the reader understand what the code does and how it works at a high level.`

	contextContent := `
	<code snippet>
    import random

    def bubble_sort(arr):
        n = len(arr)
        for i in range(n-1):
            for j in range(n-i-1):
                if arr[j] > arr[j+1]:
                    arr[j], arr[j+1] = arr[j+1], arr[j]
        return arr

    numbers = [random.randint(1, 100) for _ in range(10)]
    print("Unsorted array:", numbers)
    sorted_numbers = bubble_sort(numbers)
    print("Sorted array:", sorted_numbers)
    </code snippet>
	`

	userContent := `Explain the provided code snippet.`

	options := llm.Options{
		RepeatLastN:   2,   // default (64) the default value will "freeze" deepseek-coder
		RepeatPenalty: 2.0, // default (1.1)
	}

	//fmt.Println(options)

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: contextContent},
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
	fmt.Println()
}
