package main

import (
	"fmt"
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

	embeddingsModel := "all-minilm:33m"

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	// Parse all source code of the examples
	// Create embeddings from documents and save them in the store

	readme, _ := content.ReadTextFile("../../README.md")

	//chunks := strings.Split(readme, "<!-- split -->")
	chunks := content.SplitTextWithRegex(readme, `## *`)

	counter := 0
	for _, chunk := range chunks {
		fmt.Println(chunk)
		fmt.Println("📝 Creating embedding for:", chunk)
		counter++
		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: string(chunk),
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
	}

}
