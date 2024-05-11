/*
Topic: Parakeet
Use history.MemoryMessages to handle the conversational history in memory
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
	//ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	ollamaUrl := "http://bob.local:11434" // Pi5

	//model := "qwen:0.5b" // speed is fast, but the model seems not recognize the history
	model := "tinydolphin" // fast, and perfect answer (short, brief)
	//model := "tinyllama" // fast, and perfect answer
	//model := "gemma:2b" // speed is correct, but the model seems not recognize the history
	//model := "phi3" // slow, but goos answer

	conversation := history.MemoryMessages{
		Messages: make(map[string]llm.MessageRecord),
	}

	systemContent := `You are an expert with the Star Trek series. use the history of the conversation to answer the question`

	userContent := `Who is James T Kirk?`

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
		RepeatLastN: 2,   // default (64) the default value will "freeze" deepseek-coder
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	// Ask the question
	answer, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		},
	)
	if err != nil {
		log.Fatal("😡:", err)
	}

	// Save the conversation
	_, err = conversation.SaveMessage("1", llm.Message{
		Role:    "user",
		Content: userContent,
	})
	if err != nil {
		log.Fatal("😡:", err)
	}

	_, err = conversation.SaveMessage("2", llm.Message{
		Role:    "system",
		Content: answer.Message.Content,
	})

	if err != nil {
		log.Fatal("😡:", err)
	}

	// New question
	//userContent = `Who is his best friend in the star trek series?` // qwen has an answer, but false
	userContent = `Who is his best friend ?` // qwen does not know the answer
	//userContent = `What is his first name?` // qwen knows the answer


	previousMessages, _ := conversation.GetAllMessages()

	// (Re)Create the conversation
	conversationMessages := []llm.Message{}
	// instruction
	conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: systemContent})
	// history
	conversationMessages = append(conversationMessages, previousMessages...)
	// last question
	conversationMessages = append(conversationMessages, llm.Message{Role: "user", Content: userContent})

	query = llm.Query{
		Model:    model,
		Messages: conversationMessages,
		Options:  options,
	}

	fmt.Println()
	fmt.Println()

	answer, err = completion.ChatStream(ollamaUrl, query,
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
