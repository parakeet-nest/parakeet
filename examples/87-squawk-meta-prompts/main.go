package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/squawk"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	squawk.New().
		Model(os.Getenv("OLLAMA_LLM_CHAT")).
		BaseURL(os.Getenv("OLLAMA_BASE_URL")).
		Provider(provider.Ollama).
		ForKids("Explain Docker").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	fmt.Println("\n====================================")

	squawk.New().
		Model(os.Getenv("OLLAMA_LLM_CHAT")).
		BaseURL(os.Getenv("OLLAMA_BASE_URL")).
		Provider(provider.Ollama).
		Brief("Explain Docker").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	fmt.Println("\n====================================")

	p := squawk.New().
		Model(os.Getenv("MODEL_RUNNER_LLM_CHAT")).
		BaseURL(os.Getenv("MODEL_RUNNER_BASE_URL")).
		Provider(provider.DockerModelRunner).
		InLaymansTerms("Explain Docker").
		Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
			fmt.Print(answer.Message.Content)
			fmt.Println("\n------------------------------------")
		}).
		SummarizeLastAnswer().
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	fmt.Println("\n====================================")

	squawk.New().
		BaseURL(os.Getenv("MODEL_RUNNER_BASE_URL")).
		Provider(provider.DockerModelRunner).
		Summarize(p.LastAnswer().Message.Content).
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		})
	
}
