/*
Topic: Parakeet
Use history.BboltMessages{} to handle the conversational history in a Bbolt bucket
with Ollama and parakeet
*/

package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/history"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	//ollamaUrl := "http://bob.local:11434" // Pi5

	model := "tinydolphin" // fast, and perfect answer (short, brief)

	conversation := history.BboltMessages{}
	conversation.Initialize("conversation.db")

	systemContent := `You are an expert with the Star Trek series. use the history of the conversation to answer the question`


	options := llm.Options{
		Temperature: 0.5, // default (0.8)
		RepeatLastN: 2,   // default (64) the default value will "freeze" deepseek-coder
	}

	// New question
	userContent := `Who is his best friend ?`

	previousMessages, _ := conversation.GetAllMessages()

	// (Re)Create the conversation
	conversationMessages := []llm.Message{}
	// instruction
	conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: systemContent})
	// history
	conversationMessages = append(conversationMessages, previousMessages...)
	// last question
	conversationMessages = append(conversationMessages, llm.Message{Role: "user", Content: userContent})

	query := llm.Query{
		Model:    model,
		Messages: conversationMessages,
		Options:  options,
	}

	fmt.Println()
	fmt.Println()

	_, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		},
	)
	fmt.Println()
	if err != nil {
		log.Fatal("😡:", err)
	}

}
