package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/enums/provider"
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

	modelRunnerURL := "http://localhost:12434/engines/llama.cpp/v1"
	model := "ai/qwen2.5:latest"
	embeddingsModel := "ai/mxbai-embed-large"

	store := embeddings.MemoryVectorStore{
		Records: make(map[string]llm.VectorRecord),
	}

	// Create embeddings from documents and save them in the store
	for idx, doc := range docs {
		fmt.Println("Creating embedding from document ", idx)
		embedding, err := embeddings.CreateEmbedding(
			modelRunnerURL,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: doc,
			},
			strconv.Itoa(idx),
			provider.DockerModelRunner,
		)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			store.Save(embedding)
		}
	}

	// Question for the Chat system
	userContent := `Who is Philippe Charrière and what spaceship does he work on?`
	//userContent := `What is the nickname of Philippe Charrière?`

	systemContent := `You are an AI assistant. Your name is Seven. 
    Some people are calling you Seven of Nine.
    You are an expert in Star Trek.
    All questions are about Star Trek.
    Using the provided context, answer the user's question
    to the best of your ability using only the resources provided.`

	// Create an embedding from the question
	embeddingFromQuestion, err := embeddings.CreateEmbedding(
		modelRunnerURL,
		llm.Query4Embedding{
			Model:  embeddingsModel,
			Prompt: userContent,
		},
		"question",
		provider.DockerModelRunner,
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}
	fmt.Println("🔎 searching for similarity...")

	similarity, _ := store.SearchMaxSimilarity(embeddingFromQuestion)

	fmt.Println("🎉 similarity", similarity)

	documentsContent := `<context><doc>` + similarity.Prompt + `</doc></context>`

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
	}

	fmt.Println("")
	fmt.Println("🤖 answer:")

	// Answer the question
	_, err = completion.ChatStream(modelRunnerURL, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		}, provider.DockerModelRunner)

	if err != nil {
		log.Fatal("😡:", err)
	}

}
