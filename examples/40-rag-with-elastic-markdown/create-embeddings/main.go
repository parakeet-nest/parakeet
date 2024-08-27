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

	//embeddingsModel := "all-minilm:33m" // This model is for the embeddings of the documents
	//embeddingsModel := "nomic-embed-text"

	embeddingsModel := "mxbai-embed-large"

	// Create an Elasticsearch client
	cert, _ := os.ReadFile(os.Getenv("ELASTIC_CERT_PATH"))

	elasticStore := embeddings.ElasticsearchStore{}
	err = elasticStore.Initialize(
		[]string{
			os.Getenv("ELASTIC_ADDRESS"),
		},
		os.Getenv("ELASTIC_USERNAME"),
		os.Getenv("ELASTIC_PASSWORD"),
		cert,
		"hierarchy-mxbai-golang-index",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	documentContent, err := content.ReadTextFile("./go1.23.md")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	// Chunk the document content
	//chunks := content.ParseMarkdownWithHierarchy(documentContent)
	chunks := content.ParseMarkdownWithLineage(documentContent)
	// Prepare the pieces of markdown for the embeddings
	for idx, chunk := range chunks {
		tpl := ""
		pieceOfMarkdown := ""
		if chunk.ParentHeader == "" {
			tpl = "%s %s \n\n %s"

			pieceOfMarkdown = fmt.Sprintf(
				tpl,
				chunk.Prefix,
				chunk.Header,
				chunk.Content,
			)
		} else {
			// Add parent section information to the markdown section
			tpl = "%s %s \n\n <!-- Parent Section: %s %s --> \n\n <!-- Parent Lineage: %s --> \n\n %s"

			pieceOfMarkdown = fmt.Sprintf(
				tpl,
				chunk.Prefix,
				chunk.Header,
				chunk.ParentPrefix,
				chunk.ParentHeader,
				chunk.Lineage,
				chunk.Content,
			)
		}

		fmt.Println("---------------------------------------------")
		fmt.Println(pieceOfMarkdown)

		fmt.Println("📝 Creating embedding from document ", idx)
		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: pieceOfMarkdown,
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
