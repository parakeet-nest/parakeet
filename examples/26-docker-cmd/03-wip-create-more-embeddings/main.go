package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

type Item struct {
	Input       string `json:"input"`
	Instruction string `json:"instruction"`
	Output      string `json:"output"`
}

var counter = 0

func generateEmbeddingsFromJsonFile(jsonFile string, store embeddings.BboltVectorStore) {
	ollamaUrl := "http://localhost:11434"
	embeddingsModel := "all-minilm:33m"

	// Read the JSON file
	fileContent, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatal("😠 Error reading file:", err)
	}

	// Parse the JSON data
	var items []Item
	err = json.Unmarshal(fileContent, &items)
	if err != nil {
		log.Fatal("😠 Error parsing JSON:", err)
	}

	// Create and save the embeddings
	for _, item := range items {
		fmt.Println("📝 Creating embedding from record ", counter)
		counter += 1

		doc := fmt.Sprintf("Input: %s \n\nOutput:%s", item.Input, item.Output)

		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: doc,
			},
			uuid.NewString(),
		)
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

func main() {

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	generateEmbeddingsFromJsonFile("./NLDockercommands.json", store)
	generateEmbeddingsFromJsonFile("../datasets/docker-build/data.json", store)
	generateEmbeddingsFromJsonFile("../datasets/docker-run/data.json", store)
	generateEmbeddingsFromJsonFile("../datasets/docker-compose/data.json", store)

}
