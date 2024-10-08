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

	documentationContent, err := content.ReadTextFile("../doc/go1.23.md")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	chunks := content.ParseMarkdown(documentationContent)

	// Create embeddings from documents and save them in the store
	for idx, doc := range chunks {
		fmt.Println("📝 Creating embedding from document ", idx)

		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: fmt.Sprintf("## %s\n\n%s\n\n", doc.Header, doc.Content),
			},
			strconv.Itoa(idx),
		)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			//embedding.MetaData = "📝 chunk num: " + strconv.Itoa(idx)
			_, err := store.Save(embedding)
			if err != nil {
				fmt.Println("😡:", err)
			}
		}

		fmt.Println(embedding.Prompt)
	}

}
