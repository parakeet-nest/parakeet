package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

type Item struct {
	Input       string `json:"input"`
	Instruction string `json:"instruction"`
	Output      string `json:"output"`
}

func main() {
	ollamaUrl := "http://localhost:11434"
	embeddingsModel := "all-minilm:33m" 

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	// Read the JSON file
	fileContent, err := os.ReadFile("./NLDockercommands.json")
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
	for i, item := range items {
		fmt.Println("📝 Creating embedding from record ", i+1)

		doc := fmt.Sprintf("Input: %s \n\nOutput:%s", item.Input, item.Output)

		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: doc,
			},
			strconv.Itoa(i+1),
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
