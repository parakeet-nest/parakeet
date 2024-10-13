package main

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func ChatWithCharacter(instructions, description, ollamaUrl, model string) {

	systemContent := instructions

	contextContext := description

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.5,
		option.RepeatLastN:   3,
		option.RepeatPenalty: 2.0,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [%s] ask me something> ", model)
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "bye" {
			break
		}

		queryChat := llm.Query{
			Model: model,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: contextContext},
				{Role: "user", Content: question},
			},
			Options:          options,
			TokenHeaderName:  "X-TOKEN",
			TokenHeaderValue: os.Getenv("TOKEN"),
		}

		fmt.Println()
		fmt.Println("🤖 answer:")

		// Answer the question
		_, err := completion.ChatStream(ollamaUrl, queryChat,
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

func PreloadModel(instructions, description, ollamaUrl, model string) error {
	current_time := time.Now()
	fmt.Println("🟢", current_time.Format("2006-01-02 15:04:05"))

	fmt.Println("🤖 preloading the model (with data)...")
	systemContent := instructions

	contextContext := description

	_, err := completion.ChatStream(ollamaUrl, llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: contextContext},
			{Role: "user", Content: "who are you?"},
		},
		Options:          llm.DefaultOptions(),
		TokenHeaderName:  "X-TOKEN",
		TokenHeaderValue: os.Getenv("TOKEN"),
	}, func(answer llm.Answer) error {
		fmt.Print(answer.Message.Content)
		return nil
	})

	if err != nil {
		return err
		//log.Fatal("😡:", err)
	}
	fmt.Println()
	fmt.Println("🤖 model is ready 🎉")
	current_time = time.Now()
	fmt.Println("🟩", current_time.Format("2006-01-02 15:04:05"))
	return nil
}


func main() {
	// create a `.env` file with the following content:
	// TOKEN=your_token (let it empty if you don't have a token)
	// OLLAMA_URL=your_ollama_url

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := os.Getenv("OLLAMA_URL")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	model := os.Getenv("MODEL")
	if model == "" {
		model = "nemotron-mini"
	}
	// nemotron-mini 🤩
	// tinydolphin:latest 😡
	// tinyllama:latest 😡
	// phi3:mini 😡
	// qwen2:0.5b 😡
	// qwen2:1.5b 🙂
	// gemma2:2b 🙂
	// phi3.5:latest 😡
	// dolphin-phi:2.7b 🙂

	// some questions to ask:
	// what is your name?
	// give me the list without detail of your qualities
	// where are you from?
	// where are you located?
	// where are you living?
	// where did you grow up?
	// who is your best friend?
	// who is your worst enemy?
	// give me the list without detail of all your friends
	// give me the list without detail of all your enemies


	instructions, err := os.ReadFile("instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}
	description, err := os.ReadFile("description.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	_ = PreloadModel(string(instructions),
		string(description),
		ollamaUrl,
		model,
	)

	ChatWithCharacter(
		string(instructions),
		string(description),
		ollamaUrl,
		model,
	)
}
