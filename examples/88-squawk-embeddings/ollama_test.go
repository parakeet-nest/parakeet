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

func TestEmbeddingsWithOllama(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	ollamaBaseUrl := os.Getenv("OLLAMA_BASE_URL")
	model := os.Getenv("OLLAMA_LLM_CHAT")
	embeddingsModel := os.Getenv("OLLAMA_LLM_EMBEDDINGS")

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

	ollamaParrot := squawk.New().
		EmbeddingsModel(embeddingsModel).
		Model(model).
		BaseURL(ollamaBaseUrl).Provider(provider.Ollama)

	ollamaParrot.Store(&store).Embeddings(docs, true)

	ollamaParrot.
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


		ollamaParrot.
		User("Who is Jean-Luc Picard?","label:picard").
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

		ollamaParrot.
		SimilaritySearch("Who is Philippe Charriere?", 0.6, 1, true).
		RemoveMessageByLabel("label:picard").
		RemoveMessageByLabel("label:similarities").
		AddSimilaritiesToMessages("label:similarities").
		User("Who is Philippe Charriere?","label:charriere").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		})



}
