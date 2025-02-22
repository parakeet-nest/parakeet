package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/source"

	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	fmt.Println("Hello, World!")

	contentPath := gear.GetEnvString("CONTENT_PATH", "../data/content.txt")

	ollamaUrl := gear.GetEnvString("OLLAMA_BASE_URL", "http://localhost:11434")

	embeddingsModel := gear.GetEnvString("LLM_EMBEDDINGS", "mxbai-embed-large")

	// Initialize the vector store
	var vectorStorePath = gear.GetEnvString("DAPHNIA_STORE_PATH", "../store/sourcedata.gob")

	vectorStore := embeddings.DaphniaVectoreStore{}
	err := vectorStore.Initialize(vectorStorePath)

	if err != nil {
		log.Fatalln("😡:", err)
	}

	

	// open ../data/content.txt
	// read the content
	allSourceCodes, err := os.ReadFile(contentPath)
	if err != nil {
		log.Fatal(err)
	}

	// First pass: Split the Gitingest content into chunks
	// Ok, it's not my best idea to use this delimiter
	chunksFromAllSourceCodes := content.SplitTextWithDelimiter(
		string(allSourceCodes),
		`================================================`,
	)

	// Second pass: Extract code elements from the chunk and create embeddings
	for idx, chunk := range chunksFromAllSourceCodes {

		fmt.Println("🔍 Extracting code elements...")
		// For example, extract the function signatures
		elements, err := source.ExtractCodeElements(chunk, "go")
		if err != nil {
			fmt.Println("😡 when extracting element:", err)
			continue
		}
		header := "METADATA:\n"
		// use the function signatures as metadata
		for _, element := range elements {
			header += element.Signature + "\n"
			fmt.Println("📝", element.Signature)
		}
		fmt.Println("================================================")
		// Create the embeddings
		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: header + chunk,
			},
			strconv.Itoa(idx),
		)

		if err != nil {
			fmt.Println("😡 when generating embedding:", err)
		} else {

			_, err := vectorStore.Save(llm.VectorRecord{
				Prompt:    embedding.Prompt,
				Embedding: embedding.Embedding,
				Id:        embedding.Id,
			})

			if err != nil {
				fmt.Println("😡:", err)
			}


			fmt.Println("🎉 Document", embedding.Id, "indexed successfully")
		}

	}

}
