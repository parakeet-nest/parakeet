/*
Topic: Parakeet
Use history.BboltMessages{} to handle the conversational history in a Bbolt bucket
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

	conversation := history.BboltMessages{}
	conversation.Initialize("./conversation.db")

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
