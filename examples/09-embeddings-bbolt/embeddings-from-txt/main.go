package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	var embeddingsModel = "all-minilm" // This model is for the embeddings of the documents

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	textFiles, err := os.ReadDir("./data")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	docs := []string{}
	for _, textFile := range textFiles {
		fmt.Println(textFile.Name())
		pathfile := "./data/" + textFile.Name()
		extension := filepath.Ext(pathfile)
		if extension == ".txt" {
			data, err := os.ReadFile(pathfile)
			if err != nil {
				log.Fatal("😡 when reading file:", err)
				return
			}
			docs = append(docs, string(data))
		}
	}

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
