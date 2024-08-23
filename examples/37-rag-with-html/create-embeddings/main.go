package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	var embeddingsModel = "all-minilm:33m" // This model is for the embeddings of the documents

	store := embeddings.BboltVectorStore{}
	err := store.Initialize("../embeddings.db")

	if err != nil {
		log.Fatalln("😡:", err)
	}

	rulesContent, err := content.ReadTextFile("./chronicles.html")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	chunks, _ := content.SplitHTMLBySections(rulesContent)

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
			embedding.MetaData = "📝 chunk num: " + strconv.Itoa(idx)
			store.Save(embedding)
		}
	}
}
