package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := "http://localhost:11434"
	embeddingsModel := "all-minilm:33m" // This model is for the embeddings of the documents
	//smallChatModel := "llama3.1:8b"

	cert, _ := os.ReadFile(os.Getenv("ELASTIC_CERT_PATH"))

	elasticStore := embeddings.ElasticSearchStore{}
	err = elasticStore.Initialize(
		[]string{
			os.Getenv("ELASTIC_ADDRESS"),
		},
		os.Getenv("ELASTIC_USERNAME"),
		os.Getenv("ELASTIC_PASSWORD"),
		cert,
		"new-golang-index",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	documentContent, err := content.ReadTextFile("./go1.23.md")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	chunks := content.ParseMarkdown(documentContent)

	/*
	newMarkdown := ""
	for _, chunk := range chunks {
		newMarkdown += fmt.Sprintf("## %s\n\n%s\n\n", chunk.Header, chunk.Content)
	}
	err = os.WriteFile("./newMarkdown.md", []byte(newMarkdown), 0644)
	if err != nil {
		log.Fatal("😡:", err)
	}
	*/

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
			_, err := elasticStore.Save(embedding)
			if err != nil {
				fmt.Println("😡:", err)
			} else {
				fmt.Println("Document", embedding.Id, "indexed successfully")
			}
		}
	}
}
