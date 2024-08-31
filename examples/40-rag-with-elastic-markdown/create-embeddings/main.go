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
	chunks := content.ParseMarkdownWithLineage(documentContent)
	// Prepare the pieces of markdown for the embeddings
	for idx, chunk := range chunks {

		// you can add meta data to the chunk: chunk.MetaData
		// you can add keywords to the chunk: chunk.KeyWords
		// these metadata and keywords could be added to the embedding using the templates

		pieceOfMarkdown := ""

		if chunk.ParentHeader == "" {
			// Create the markdown section
			mdTemplate := `{{.Prefix}} {{.Header}}

{{.Content}}`

			pieceOfMarkdown, err = content.InterpolateString(mdTemplate, chunk)
			if err != nil {
				log.Println("😡:", err)
				pieceOfMarkdown = ""
			}
		} else {

			// Add metadata to the markdown section
			mdTemplate := `{{.Prefix}} {{.Header}}
<!-- Parent Section: {{.ParentPrefix}} {{.ParentHeader}} -->
<!-- Lineage: {{.Lineage}} -->

{{.Content}}`

			pieceOfMarkdown, err = content.InterpolateString(mdTemplate, chunk)
			if err != nil {
				log.Println("😡:", err)
				pieceOfMarkdown = ""
			}
		}

		if len(pieceOfMarkdown) > 0 {
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
				// You can add metadata to the embedding
				// It could be useful for debugging and filtering with Elasticsearch
				// TODO: see how to use this metadata in the search
				embedding.MetaData = "👋 hello from Parakeet 🦜🪺"


				_, err := elasticStore.Save(embedding)
				if err != nil {
					fmt.Println("😡:", err)
				} else {
					fmt.Println("Document", embedding.Id, "indexed successfully")
				}
			}
		}
	}
}
