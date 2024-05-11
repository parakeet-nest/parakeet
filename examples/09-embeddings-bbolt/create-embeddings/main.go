package main

import (
	"fmt"
	"strconv"

	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

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

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	var embeddingsModel = "all-minilm" // This model is for the embeddings of the documents

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	// Create embeddings from documents and save them in the store
	for idx, doc := range docs {
		fmt.Println("📝 Creating embedding from document ", idx)
		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: doc,
			},
			strconv.Itoa(idx), // don't forget the id (unique identifier)
		)
		fmt.Println("📦", embedding.Id, embedding.Prompt)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			_, err := store.Save(embedding)
			if err != nil {
				fmt.Println("😡:", err)
			}
		}
	}

}
