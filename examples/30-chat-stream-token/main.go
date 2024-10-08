/*
Topic: Parakeet
Generate a chat completion with Ollama and parakeet
The output is streamed
*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"

)

func main() {
	// create a `.env` file with the following content:
	// TOKEN=your_token
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := "https://ollama.wasm.ninja"
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
		option.Temperature: 0.5,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 3.0,
		option.Verbose: true,
	})

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options:          options,
		TokenHeaderName:  "X-TOKEN",
		TokenHeaderValue: os.Getenv("TOKEN"),
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
		log.Fatal("😡:", err)
	}
}
