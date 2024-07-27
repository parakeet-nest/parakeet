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
		Temperature:   0.0,
		RepeatLastN:   2,
		RepeatPenalty: 1.5,
	}

	dungeonMaster := &Character{
		Id:         "dungeon-master",
		Model:      model,
		OllamaUrl:  ollamaUrl,
		Options:    options,
		MailBroker: mailBroker,
		Verbose: true,
		Messages: []llm.Message{
			{
				Role: "system",
				Content: `
				You are a Dungeon Master of a role playing game.
				Your duty is to give information about the below characters list, use only the context:
				<characters>
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
					<character>
						Name: Pipkin Took

						Story:
						Pipkin is a relative of Peregrin Took, known as Pippin, and shares the Took family's adventurous spirit. He often finds himself in trouble due to his mischievous nature and penchant for practical jokes. Despite this, Pipkin has a kind heart and a deep loyalty to his friends and family.
						
						Characteristics:
						
						Mischievous: Loves playing pranks and having fun.
						Loyal: Fiercely loyal to his friends and family.
						Brave: Willing to face danger for the sake of his loved ones.
						Clever: Quick-witted and good at thinking on his feet.
					</character>
					<character>
						Name: Primrose Cotton

						Story:
						Primrose, or Prim for short, is from Bywater and is related to Rosie Cotton, Samwise Gamgee’s wife. She is a skilled weaver and dyer, known for creating beautiful tapestries and clothing. Primrose is content with her life in the Shire but harbors a secret desire to see the world beyond its borders.
						
						Characteristics:
						
						Creative: Has a natural talent for crafting and creating.
						Thoughtful: Always considerate of others' feelings and needs.
						Independent: Strong-willed and capable of taking care of herself.
						Dreamer: Often daydreams of adventures and far-off places.
					</character>
					<character>
						Name: Hobson Burrows

						Story:
						Hobson, often called Hoby, comes from a long line of burrow builders in Hobbiton. His family has been responsible for some of the most well-constructed hobbit-holes in the Shire. Hoby, however, has a passion for history and spends his free time studying old maps and manuscripts, hoping to uncover forgotten tales of the Shire and beyond.
						
						Characteristics:
						
						Studious: Loves reading and learning about history.
						Meticulous: Pays great attention to detail, especially in his work.
						Gentle: Soft-spoken and kind-hearted.
						Curious: Always eager to learn more about the past and its secrets.
					</character>
				</characters>
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

			fmt.Println("")

			return nil
		},
		OnError: func(self *Character, err error) {
			fmt.Println("😡", err)
		},
	}
	bob := &Character{
		Id:         "Bob",
		MailBroker: mailBroker,
	}

	dungeonMaster.Start()

	bob.SendMail("[Brief]Give me the characters of the game", dungeonMaster.Id)
	bob.SendMail("[Brief]What is the story of Daisy Bramblefoot", dungeonMaster.Id)
	bob.SendMail("[Brief]What are the characteristics of Primrose Cotton", dungeonMaster.Id)



	<-ctx.Done()

}
