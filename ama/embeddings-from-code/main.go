package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	//var embeddingsModel = "magicoder:latest"
	
	embeddingsModel := "all-minilm"

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	// Parse all source code of the examples
	// Create embeddings from documents and save them in the store

	counter := 0
	_, err := content.ForEachFile("../../examples", ".go", func(path string) error {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		fmt.Println("📝 Creating embedding from:", path)
		counter++
		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: string(data),
			},
			strconv.Itoa(counter), // don't forget the id (unique identifier)
		)
		fmt.Println("📦 Created: ", len(embedding.Embedding))

		if err != nil {
			fmt.Println("😡:", err)
		} else {
			_, err := store.Save(embedding)
			if err != nil {
				fmt.Println("😡:", err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalln("😡:", err)
	}


}
