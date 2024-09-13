/*
Topic: Parakeet
Generate a simple completion with Ollama and parakeet
no streaming
*/
package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
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
	model := "tinydolphin"

	info, statusCode, err := llm.ShowModelInformationWithToken(ollamaUrl, model, "X-TOKEN", os.Getenv("TOKEN"))

	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(statusCode, info)

	list, statusCode, err := llm.GetModelsListWithToken(ollamaUrl, "X-TOKEN", os.Getenv("TOKEN"))

	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(statusCode, list.Models)

	fmt.Println()

}
