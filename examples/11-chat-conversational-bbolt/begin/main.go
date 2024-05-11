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

	model := "tinydolphin" // fast, and perfect answer (short, brief)

	conversation := history.BboltMessages{}
	conversation.Initialize("../conversation.db")

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


}
