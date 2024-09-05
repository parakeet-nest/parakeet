package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

/*
https://github.com/01-ai/Yi-Coder/blob/main/cookbook/System_prompt/System_prompt.ipynb
*/

// This scenario demonstrates how Yi-Coder can identify errors and insert the correct code to fix them.

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "yi-coder:1.5b"

	systemContent := `You are Yi-Coder, you are exceptionally skilled in programming, coding, and any computer-related issues.`

	userContent := `
	'''python
	def quick_sort(arr):
		if len(arr) <= 1:
			return arr
		else:
			pivot = arr[len(arr) // 2]
			left = [x for x in arr if x < pivot]

			right = [x for x in arr if x > pivot]
			return quick_sort(left) + middle + quick_sort(right)

	print(quick_sort([3,6,8,10,1,2,1]))
	# Prints "[1, 1, 2, 3, 6, 8, 10]"
	'''
	Is there a problem with this code?
	`

	options := llm.Options{
		Temperature: 0.0,
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
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
