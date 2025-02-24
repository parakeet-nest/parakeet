/*
Topic: Parakeet
Generate a chat completion with Ollama and parakeet
The output is streamed
*/

package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "deepseek-coder"

	systemContent := `You are an expert in computer programming.
	Please make friendly answer for the noobs.
	Add source code examples if you can.`

	userContent := `I need a clear explanation regarding the following question:
	Can you create a "hello world" program in Golang?
	And, please, be structured with bullet points`

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.5,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 3.0,
		option.Verbose:       false,
	})

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	fullAnswer, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	fmt.Println("📝 Full answer:")
	fmt.Println(fullAnswer.Message.Role)
	fmt.Println(fullAnswer.Message.Content)

	if err != nil {
		// test if the model is not found
		if modelErr, ok := err.(*completion.ModelNotFoundError); ok {
			fmt.Printf("💥 Got Model Not Found error: %s\n", modelErr.Message)
			fmt.Printf("😡 Error code: %d\n", modelErr.Code)
			fmt.Printf("🧠 Expected Model: %s\n", modelErr.Model)
		} 

		if noHostErr, ok := err.(*completion.NoSuchOllamaHostError); ok {
			fmt.Printf("🦙 Got No Such Ollama Host error: %s\n", noHostErr.Message)
			fmt.Printf("🌍 Expected Host: %s\n", noHostErr.Host)
		}
		
		log.Fatal("😡:", err)
		
	}
}
