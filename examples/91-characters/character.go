package main
// This package is experimental

import (
	"fmt"
	"time"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

// Struct that interacts with the shared list
type Character struct {
	Id string

	//memory  []llm.Message
	//context []llm.Message
	// TODO: add method to add messages to the memory and context?
	// TODO: add  etc... ?
	Messages         []llm.Message
	OllamaUrl        string
	Model            string
	Options          llm.Options
	MailBroker       *MailBroker
	BeforeCompletion func(self *Character, senderID string, senderMail string) error
	AfterCompletion  func(self *Character, answer string, senderID string, senderMail string) error
	OnError          func(self *Character, err error)
	Verbose          bool
}

// Method for Worker to add a message to the shared list
func (character *Character) SendMail(content string, recipientID string) {
	mail := Mail{
		Content:     content,
		SenderID:    character.Id,
		RecipientID: recipientID,
		Date:        time.Now(),
		Read:        false,
	}
	character.MailBroker.Add(mail)
}

/*
func (character *Character) GetId() string {
	return character.Id
}
*/

// Start method to poll the message list and mark messages as read
func (character *Character) Start() {
	// TODO handel default values
	go func() {
		for {
			list := character.MailBroker.Get()
			for i, msg := range list {
				if msg.RecipientID == character.Id && !msg.Read {

					messagesForQuery := []llm.Message{
						{Role: "user", Content: msg.Content},
					}
					// TODO: add a method to add messages to the memory
					messagesForQuery = append(character.Messages, messagesForQuery...)

					err := character.BeforeCompletion(character, msg.SenderID, msg.Content)

					if err != nil {
						fmt.Println("😡", err)
						character.OnError(character, err)
					}

					//fmt.Println("📝", messagesForQuery)

					query := llm.Query{
						Model:    character.Model,
						Messages: messagesForQuery,
						Options:  character.Options,
						//Stream:  false,
					}

					// TODO: activate or deactivate the print (verbose mode)

					var response llm.Answer

					if character.Verbose  {
						response, err = completion.ChatStream(character.OllamaUrl, query,
							func(answer llm.Answer) error {
								fmt.Print(answer.Message.Content)
								return nil
							})
					} else {
						response, err = completion.Chat(character.OllamaUrl, query)
					}

					if err != nil {
						fmt.Println("😡", err)
						character.OnError(character, err)
					}

					err = character.AfterCompletion(character, response.Message.Content, msg.SenderID, msg.Content)
					if err != nil {
						fmt.Println("😡", err)
						character.OnError(character, err)
					}

					character.MailBroker.MarkAsRead(i)
				}
			}
			time.Sleep(1 * time.Second) // Polling interval
		}
	}()
}
