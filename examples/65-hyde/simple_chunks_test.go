package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func TestGenerateSimpleChunks(t *testing.T) {

	ollamaUrl := os.Getenv("OLLAMA_URL")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}
	embeddingsModel := "mxbai-embed-large:latest"
	/*
		options := llm.SetOptions(map[string]interface{}{
			option.Temperature: 0.0,
		})
	*/

	// Initialize the vector store
	vectorStore := embeddings.DaphniaVectoreStore{}
	vectorStore.Initialize("no-context.gob")

	content.ForEachFile("./docs", ".md", func(documentPath string) error {
		fmt.Println("📝 Creating embedding from document ", documentPath)

		// Read the content of the file
		wholeDocumentContent, err := content.ReadTextFile(documentPath)
		if err != nil {
			log.Fatalln("😡:", err)
		}

		chunks := content.ParseMarkdownWithLineage(wholeDocumentContent)

		//chunks := content.ParseMarkdownWithHierarchy(wholeDocumentContent)

		fmt.Println("👋 Found", len(chunks), "chunks")

		// Create embeddings from documents and save them in the store
		for idx, doc := range chunks {
			fmt.Println("📝 Creating embedding from document ", idx)
			/*
				fmt.Println("Level:", doc.Level)
				fmt.Println("Prefix:", doc.Prefix)
				fmt.Println("ParentPrefix:", doc.ParentPrefix)
				fmt.Println("ParentHeader:", doc.ParentHeader)
			*/
			fmt.Println("🖼️", doc.Header)
			fmt.Println("Lineage:", doc.Lineage)

			embedding, err := embeddings.CreateEmbedding(
				ollamaUrl,
				llm.Query4Embedding{
					Model:  embeddingsModel,
					Prompt: fmt.Sprintf("METADATA: %s\n\n ## %s\n\n%s\n\n", doc.Lineage, doc.Header, doc.Content),
				},
				documentPath+"-"+strconv.Itoa(idx),
			)
			if err != nil {
				fmt.Println("😡:", err)
			} else {

				record, err := vectorStore.Save(llm.VectorRecord{
					Prompt:    embedding.Prompt,
					Embedding: embedding.Embedding,
					Id:        embedding.Id,
				})
				//fmt.Println("📝 Embedding:", record.Embedding)
				fmt.Println("📝 Embedding:", record.Id)


				if err != nil {
					fmt.Println("😡:", err)
				}

			}

			fmt.Println("---[Improved chunk]--------------------------------")
			fmt.Println(embedding.Prompt)
			fmt.Println("Lineage:", doc.Lineage)
			fmt.Println("---------------------------------------------------")

		}

		return nil
	})

}
