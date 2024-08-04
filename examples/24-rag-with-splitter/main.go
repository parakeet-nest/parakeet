package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/content"

)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	var embeddingsModel = "all-minilm:33m" // This model is for the embeddings of the documents
	//var embeddingsModel = "all-minilm" // This model is for the embeddings of the documents
	var smallChatModel = "qwen2:0.5b"  // This model is for the chat completion
	//var smallChatModel = "tinyllama"  // This model is for the chat completion


	store := embeddings.MemoryVectorStore{
		Records: make(map[string]llm.VectorRecord),
	}


	rulesContent, err := content.ReadTextFile("./chronicles.md")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	chunks := content.SplitTextWithDelimiter(rulesContent, "<!-- SPLIT -->")
	//chunks := content.SplitText(rulesContent, "##")


	// Create embeddings from documents and save them in the store
	for idx, doc := range chunks {
		fmt.Println("Creating embedding from document ", idx)
		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: doc,
			},
			strconv.Itoa(idx),
		)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			store.Save(embedding)
		}
	}

	// Question for the Chat system
	//userContent := `Who are the monsters of Chronicles of Aethelgard? Give details for every monster.`
	userContent := `Tell me more about Keegorg`
	//userContent := `Who are the monsters of Chronicles of Aethelgard?`


	systemContent := `You are the dungeon master,
	expert at interpreting and answering questions based on provided sources.
	Using only the provided context, answer the user's question 
	to the best of your ability using only the resources provided. 
	Be verbose!`

	// Create an embedding from the question
	embeddingFromQuestion, err := embeddings.CreateEmbedding(
		ollamaUrl,
		llm.Query4Embedding{
			Model:  embeddingsModel,
			Prompt: userContent,
		},
		"question",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}
	fmt.Println("🔎 searching for similarity...")

	similarities, _ := store.SearchSimilarities(embeddingFromQuestion, 0.3)


	fmt.Println("🎉 similarities:", len(similarities))

	documentsContent := embeddings.GenerateContentFromSimilarities(similarities)

	//fmt.Println("📙", documentsContent)


	query := llm.Query{
		Model: smallChatModel,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
		Options: llm.Options{
			Temperature: 0.0,
			RepeatLastN: 2,
			RepeatPenalty: 3.0,
			TopK: 10,
			TopP: 0.5,
		},
	}

	fmt.Println("")
	fmt.Println("🤖 answer:")

	// Answer the question
	_, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

}

