package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/squawk"
)

func TestEmbeddingsWithOpenAI(t *testing.T) {
	// create a `openai.key.env` file with the following content:
	// OPENAI_API_KEY=your_openai_api_key
	err := godotenv.Load(".env", "openai.key.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	openAIBaseUrl := os.Getenv("OPENAI_BASE_URL")
	model := os.Getenv("OPENAI_LLM_CHAT")
	embeddingsModel := os.Getenv("OPENAI_LLM_EMBEDDINGS")

	fmt.Println("OpenAI Base URL:", openAIBaseUrl)
	fmt.Println("OpenAI LLM Chat:", model)
	fmt.Println("OpenAI LLM Embeddings:", embeddingsModel)
	fmt.Println("====================================")

	store := embeddings.MemoryVectorStore{
		Records: make(map[string]llm.VectorRecord),
	}

	var docs = []string{
		`Michael Burnham is the main character on the Star Trek series, Discovery.  
		She's a human raised on the logical planet Vulcan by Spock's father.  
		Burnham is intelligent and struggles to balance her human emotions with Vulcan logic.  
		She's become a Starfleet captain known for her determination and problem-solving skills.
		Originally played by actress Sonequa Martin-Green`,

		`James T. Kirk, also known as Captain Kirk, is a fictional character from the Star Trek franchise.  
		He's the iconic captain of the starship USS Enterprise, 
		boldly exploring the galaxy with his crew.  
		Originally played by actor William Shatner, 
		Kirk has appeared in TV series, movies, and other media.`,

		`Jean-Luc Picard is a fictional character in the Star Trek franchise.
		He's most famous for being the captain of the USS Enterprise-D,
		a starship exploring the galaxy in the 24th century.
		Picard is known for his diplomacy, intelligence, and strong moral compass.
		He's been portrayed by actor Patrick Stewart.`,

		`Lieutenant Philippe Charrière, known as the **Silent Sentinel** of the USS Discovery, 
		is the enigmatic programming genius whose codes safeguard the ship's secrets and operations. 
		His swift problem-solving skills are as legendary as the mysterious aura that surrounds him. 
		Charrière, a man of few words, speaks the language of machines with unrivaled fluency, 
		making him the crew's unsung guardian in the cosmos. His best friend is Spiderman from the Marvel Cinematic Universe.`,
	}

	openAIParrot := squawk.New().
		EmbeddingsModel(embeddingsModel).
		Model(model).
		BaseURL(openAIBaseUrl).Provider(provider.OpenAI, os.Getenv("OPENAI_API_KEY"))

	openAIParrot.Store(&store).GenerateEmbeddings(docs, true)

	openAIParrot.
		SimilaritySearch("Who is James T Kirk?", 0.6, 1, true).
		Cmd(func(self *squawk.Squawk) {
			if self.LastError() != nil {
				t.Log("😡 Error:", self.LastError())
			}

			fmt.Println(embeddings.GenerateContextFromSimilarities(self.Similarities()))
			fmt.Println("====================================")
			//Simpler
			fmt.Println(self.ContextFromSimilarities())

		})

	fmt.Println("++++++++++++++++++++++++++++++++++")

	openAIParrot.
		User("Who is Jean-Luc Picard?", "label:picard").
		SimilaritySearchFromUserMessage("label:picard", 0.6, 1, true).
		Cmd(func(self *squawk.Squawk) {
			fmt.Println(self.ContextFromSimilarities())
		}).
		AddSimilaritiesToMessages("label:similarities").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	fmt.Println("++++++++++++++++++++++++++++++++++")
	
	// TODO: find why this doesn't work
	openAIParrot.
		SimilaritySearch("Who is Philippe Charriere?", 0.5, 2, true).
		RemoveMessageByLabel("label:picard").
		RemoveMessageByLabel("label:similarities").
		AddSimilaritiesToMessages("label:similarities").
		User("Who is Philippe Charriere?", "label:charriere").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

}
