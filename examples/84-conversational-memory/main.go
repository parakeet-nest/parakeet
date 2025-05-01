/*
Topic: Parakeet
Use history.MemoryMessages to handle the conversational history in memory
with Ollama and parakeet
*/
package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/history"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {

	conversation := history.MemoryMessages{
		Messages: make(map[string]llm.MessageRecord),
	}

	// Save the conversation
	_, err := conversation.SaveMessage("", llm.Message{
		Role:    "user",
		Content: "hello, who are you?",
	})

	_, err = conversation.SaveMessage("", llm.Message{
		Role:    "agent",
		Content: "hello, I am Qwen",
	})

	_, err = conversation.SaveMessage("", llm.Message{
		Role:    "user",
		Content: "Nice to meet you",
	})

	_, err = conversation.SaveMessage("", llm.Message{
		Role:    "agent",
		Content: "the same here",
	})

	messages, err := conversation.GetAllMessages()
	// Get the conversation
	fmt.Println("🟢 conversation:", messages)

	conversation.RemoveTopMessage()
	
	messages, err = conversation.GetAllMessages()
	// Get the conversation
	fmt.Println("🟢 conversation:", messages)


	if err != nil {
		log.Fatal("😡:", err)
	}



}
