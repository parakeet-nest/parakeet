package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
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
	//model := "yi-coder:9b"
	model := "yi-coder:1.5b"

	systemContent := `SYSTEM:
	You are Yi-Coder, you are exceptionally skilled in programming, coding, devops and any computer-related issues.
	`

	allSourceCode, err := content.GetMapOfContentFiles("./", ".yml")
	if err != nil {
		log.Fatal(err)
	}

	codebase := "CODEBASE:\n"
	for _, pipeline := range allSourceCode {
		codebase += "<>\n```yaml\n" + pipeline + "```\n<>\n"
	}


	userContent := `[Step-by-Step] Using the above codebase, explain the GitLab CI pipeline.
	Make a detqiled section per CI job.`

	options := llm.Options{
		Temperature: 0.0,
	}

	query := llm.GenQuery{
		Model: model,
		Prompt: systemContent + codebase + userContent,
		Options: options,
	}

	_, err = completion.GenerateStream(ollamaUrl, query,
		func(answer llm.GenAnswer) error {
			fmt.Print(answer.Response)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

}