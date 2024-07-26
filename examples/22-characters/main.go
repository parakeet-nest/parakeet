package main

// This demo uses an experimental package

import (
	"context"
	"fmt"

	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ctx := context.Background()

	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "qwen2:0.5b"

	mailBroker := &MailBroker{}

	options := llm.Options{
		Temperature:   0.6,
		RepeatLastN:   2,
		RepeatPenalty: 1.5,
	}

	daisyBramblefoot := &Character{
		Id:         "Daisy_Bramblefoot",
		Model:      model,
		OllamaUrl:  ollamaUrl,
		Options:    options,
		MailBroker: mailBroker,
		Verbose:    true,
		Messages: []llm.Message{
			{
				Role: "system",
				Content: `
					You are this character:
					<character>
						Name: Daisy Bramblefoot

						Story:
						Daisy is from Buckland, near the borders of the Old Forest. She grew up with stories of the mysterious forest and its denizens, and unlike many hobbits, she has a fascination with the unknown. Daisy is known for her exceptional skills in herbology and healing, often sought out by her fellow hobbits for remedies and advice.
						
						Characteristics:
						
						Knowledgeable: Has extensive knowledge of herbs and healing.
						Adventurous: Enjoys exploring and learning about the world beyond the Shire.
						Compassionate: Always ready to help those in need.
						Determined: Once she sets her mind to something, she sees it through.
					</character>
					If someone speaks to you, be nice and have a friendly conversation.
					Your main topic is the Old Forest
					Be creative, speak about what you know about the Old Forest
					`,
			},
		},
		BeforeCompletion: func(self *Character, senderID string, senderMail string) error {
			fmt.Println("📝", senderID, "sent a message:", senderMail)
			fmt.Println("🤖:")
			return nil
		},
		AfterCompletion: func(self *Character, answer string, senderID string, senderMail string) error {
			// keep the memory of the conversation
			self.Messages = append(self.Messages, llm.Message{Role: "user", Content: senderMail})
			self.Messages = append(self.Messages, llm.Message{Role: "system", Content: answer})

			self.SendMail("What do you think?", senderID)
			fmt.Println()

			return nil
		},
		OnError: func(self *Character, err error) {
			fmt.Println("😡", err)
		},
	}

	miloGoodbody := &Character{
		Id:         "Milo_Goodbody",
		Model:      model,
		OllamaUrl:  ollamaUrl,
		Options:    options,
		MailBroker: mailBroker,
		Verbose:    true,
		Messages: []llm.Message{
			{
				Role: "system",
				Content: `
					You are this character:
					<character>
						Name: Milo Goodbody

						Story:
						Milo hails from Hobbiton and is a distant cousin of the famous Samwise Gamgee. He grew up hearing tales of the great adventures of Bilbo and Frodo Baggins, which fueled his own desire for adventure. Despite the peaceful nature of the Shire, Milo's curiosity often leads him to explore the forests and hills around his home.
						
						Characteristics:
						
						Curious: Always seeking out new experiences and knowledge.
						Brave: Though small, he has a heart as stout as any adventurer.
						Resourceful: Good at finding solutions to unexpected problems.
						Friendly: Easily makes friends, both hobbit and non-hobbit alike.
					</character>
					If someone speaks to you, be nice and have a friendly conversation.
					Your main topic is the Old Forest
					Be creative, speak about what you know about the Old Forest
					`,
			},
		},
		BeforeCompletion: func(self *Character, senderID string, senderMail string) error {
			fmt.Println("📝", senderID, "sent a message:", senderMail)
			fmt.Println("🤖:")
			return nil
		},
		AfterCompletion: func(self *Character, answer string, senderID string, senderMail string) error {
			// keep the memory of the conversation
			self.Messages = append(self.Messages, llm.Message{Role: "user", Content: senderMail})
			self.Messages = append(self.Messages, llm.Message{Role: "system", Content: answer})

			self.SendMail("please resume", senderID)
			fmt.Println()

			return nil
		},
		OnError: func(self *Character, err error) {
			fmt.Println("😡", err)
		},
	}

	daisyBramblefoot.Start()
	miloGoodbody.Start()

	//daisyBramblefoot.SendMail("Hey! Who are you? Let's have a chat!", miloGoodbody.Id)
	miloGoodbody.SendMail("Do you know the Old Forest?", daisyBramblefoot.Id)

	<-ctx.Done()

}
