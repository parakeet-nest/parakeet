package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {

	ollamaUrl := "http://localhost:11434"

	//smallChatModel := "CognitiveComputations/dolphin-gemma2"
	//smallChatModel := "qwen:0.5b"
	smallChatModel := "qwen2:0.5b"


	systemContent := `**Instruction:**
	You are an expert in botanics.
	Please use only the content provided below to answer the question.
	Do not add any external knowledge or assumptions.`

	documentPath := "../data/ferns.2.md"
	//documentPath := "../data/ferns.1.extract.md"
	//documentPath := "../data/ferns.2.extract.md"

	documentContent, err := content.ReadTextFile(documentPath)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	contextContext := "**Content:**: \n\n" + string(documentContent)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 (%s) ask me something> ", smallChatModel)
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "bye" {
			break
		}

		queryChat := llm.Query{
			Model: smallChatModel,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: contextContext},
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

		fmt.Println()
		fmt.Println("🤖 answer:")

		// Answer the question
		_, err = completion.ChatStream(ollamaUrl, queryChat,
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
