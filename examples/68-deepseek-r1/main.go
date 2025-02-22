package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := os.Getenv("OLLAMA_HOST")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	model := os.Getenv("LLM_CHAT")
	if model == "" {
		model = "deepseek-r1:1.5b"
	}

	fmt.Println("🌍", ollamaUrl, "📕", model)

	systemInstructions, err := os.ReadFile("system-instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	irisInstructions, err := os.ReadFile("iris-instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	irisDatabase, err := os.ReadFile("iris-database.xml")
	if err != nil {
		log.Fatal("😡:", err)
	}

	/*
		// Verginica
		userContent := `Using the above information and the below information,
		Given a specimen with:
		- Petal width: 2,5 cm
		- Petal length: 6 cm
		- Sepal width: 3,3 cm
		- Sepal length: 6,3 cm
		What is the species of the iris?`

		// Versicolor
		userContent := `Using the above information and the below information,
		Given a specimen with:
		- Petal width: 1,5 cm
		- Petal length: 4,5 cm
		- Sepal width: 3,2 cm
		- Sepal length: 6,4 cm
		What is the species of the iris?`

		// Setosa
		userContent := `Using the above information and the below information,
		Given a specimen with:
		- Petal width: 0,2 cm
		- Petal length: 1,4 cm
		- Sepal width: 3,6 cm
		- Sepal length: 5 cm
		What is the species of the iris?`
	*/

	// Verginica
	userContent := `Using the above information and the below information,
	Given a specimen with:
	- Petal width: 1,9 cm
	- Petal length: 5,1 cm
	- Sepal width: 2,7 cm
	- Sepal length: 5,8 cm
	What is the species of the iris?`

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.0,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.2,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	// Prompt construction
	messages := []llm.Message{
		{Role: "system", Content: string(systemInstructions)},
		{Role: "system", Content: string(irisInstructions)},
		{Role: "system", Content: "# Iris Database\n" + string(irisDatabase)},
		{Role: "user", Content: userContent},
	}

	query := llm.Query{
		Model:    model,
		Messages: messages,
		Options:  options,
	}

	_, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		switch e := err.(type) {
		case *completion.ModelNotFoundError:
			fmt.Printf("💥 Got Model Not Found error: %s\n", e.Message)
			fmt.Printf("😡 Error code: %d\n", e.Code)
			fmt.Printf("🧠 Expected Model: %s\n", e.Model)

		case *completion.NoSuchOllamaHostError:
			fmt.Printf("🦙 Got No Such Ollama Host error: %s\n", e.Message)
			fmt.Printf("🌍 Expected Host: %s\n", e.Host)

		default:
			log.Fatal("😡:", err)
		}
	}
}
