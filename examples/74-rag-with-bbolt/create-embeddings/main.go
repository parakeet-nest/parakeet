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
	//err := store.Initialize("../embeddings.readme.db")
	//err := store.Initialize("../embeddings.db")
	err := store.Initialize("../embeddings.updated.db")


	if err != nil {
		log.Fatalln("😡:", err)
	}

	//textContent, err := content.ReadTextFile("../content.readme.txt")
	//textContent, err := content.ReadTextFile("../content.txt")
	textContent, err := content.ReadTextFile("../updated_content.txt")

	if err != nil {
		log.Fatalln("😡:", err)
	}

	//chunks := content.SplitTextWithRegex(textContent, `## *`)
	//chunks := content.SplitMarkdownBySections(textContent)
	chunks := content.ParseMarkdownWithLineage(textContent)
	//chunks := content.ChunkText(textContent, 2048, 512)

	// Create embeddings from documents and save them in the store
	for idx, doc := range chunks {
		fmt.Println("Creating embedding from document ", idx)
		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: doc.Header + "\n" + doc.Lineage + "\n" + doc.Content,
			},
			strconv.Itoa(idx),
		)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			embedding.SimpleMetaData = "📝 chunk num: " + strconv.Itoa(idx)
			_, err := store.Save(embedding)
			if err != nil {
				fmt.Println("😡:", err)
			}
		}
	}
}
