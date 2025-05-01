package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/squawk"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡 When loading env vars:", err)
	}
	engineBaseUrl := os.Getenv("OLLAMA_BASE_URL")
	model := os.Getenv("OLLAMA_LLM_CHAT")
	embeddingsModel := os.Getenv("OLLAMA_LLM_EMBEDDINGS")

	store := embeddings.BboltVectorStore{}
	err = store.Initialize("./embeddings.db")
	if err != nil {
		log.Fatalln("😡 When creating store:", err)
	}

	var starWarsChars = []string{
		`Luke Skywalker is the main protagonist of the original Star Wars trilogy.
		Born on Tatooine, he discovers his Force sensitivity and trains to become a Jedi Knight.
		Son of Anakin Skywalker (Darth Vader), Luke plays a crucial role in defeating the Empire.
		He later attempts to rebuild the Jedi Order and briefly trains Rey before becoming one with the Force.
		Originally portrayed by actor Mark Hamill.`,

		`Princess Leia Organa is a leader of the Rebel Alliance and a key figure in the fight against the Empire.
		As Luke's twin sister and daughter of Anakin Skywalker, she is Force-sensitive but chooses a political path.
		Known for her tactical brilliance, determination, and distinctive hairstyles.
		She later becomes a general in the Resistance against the First Order.
		Originally portrayed by actress Carrie Fisher.`,

		`Darth Vader, formerly Anakin Skywalker, is the infamous Sith Lord serving Emperor Palpatine.
		Seduced by the dark side of the Force, he hunts Jedi and enforces the Emperor's will across the galaxy.
		Known for his black armor, intimidating breathing, and powerful Force abilities.
		He ultimately finds redemption by saving his son Luke and destroying the Emperor.
		Originally portrayed by actors David Prowse and voiced by James Earl Jones.`,

		`Ahsoka Tano is the former Padawan of Anakin Skywalker who left the Jedi Order.
		A skilled Togruta warrior with distinctive white lightsabers and montrals (head-tails).
		After Order 66, she helps the early Rebellion and later assists the Mandalorian.
		She becomes a symbol of balance, operating outside traditional Jedi constraints.
		Originally voiced by Ashley Eckstein in animated series and portrayed by Rosario Dawson in live action.`,
	}

	squawk.New().
		EmbeddingsModel(embeddingsModel).
		Model(model).
		BaseURL(engineBaseUrl).
		Provider(provider.Ollama).
		Options(llm.SetOptions(map[string]interface{}{
			option.Temperature:   0.0,
			option.RepeatLastN:   2,
			option.RepeatPenalty: 2.2,
		})).
		Store(&store).
		GenerateEmbeddings(starWarsChars, true).
		System("You are a useful AI agent, you are a Star Wars expert.").
		User("Who is Luke Skywalker?", "q1").
		SimilaritySearchFromUserMessage("q1", 0.6, 1, true).
		AddSimilaritiesToMessages("sim1").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		}).
		SaveAssistantAnswer("a1"). // save the answe to the history / to the messages list
		Cmd(func(self *squawk.Squawk) {
			// List of messages
			fmt.Println("\n====================================")
			fmt.Println("List of messages:")
			for _, m := range self.Messages() {
				fmt.Println("  - ", m.Role, m.Content)
			}
			fmt.Println("====================================")

		}).
		RemoveMessageByLabel("q1").
		RemoveMessageByLabel("sim1").
		User("Who is Leia?", "q2").
		SimilaritySearchFromUserMessage("q2", 0.6, 1, true).
		AddSimilaritiesToMessages().
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		})
}
