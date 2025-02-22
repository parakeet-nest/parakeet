package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/source"

	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := gear.GetEnvString("OLLAMA_BASE_URL", "http://localhost:11434")
	embeddingsModel := gear.GetEnvString("LLM_EMBEDDINGS", "mxbai-embed-large:335m")

	vectorStorePath := gear.GetEnvString("DAPHNIA_STORE_PATH", "../sourcedata.gob")

	store := embeddings.DaphniaVectoreStore{}
	err = store.Initialize(vectorStorePath)

	if err != nil {
		log.Fatalln("😡:", err)
	}

	contentPath := gear.GetEnvString("CONTENT_PATH", "../content.txt")

	allSourceCodes, err := os.ReadFile(contentPath)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	// Split the Gitingest content into chunks
	// Ok, it's not my best idea to use this delimiter
	chunksFromAllSourceCodes := content.SplitTextWithDelimiter(
		string(allSourceCodes),
		`================================================`,
	)

	for idx, chunk := range chunksFromAllSourceCodes {
		//fmt.Println("📝", chunk)
		fmt.Println("================================================")
		// Extract code elements from the chunk
		fmt.Println("🔍 Extracting code elements...")
		elements, err := source.ExtractCodeElements(chunk, "go")
		if err != nil {
			fmt.Println("😡 when extracting element:", err)
			continue
		}

		header := "METADATA:\n"

		for _, element := range elements {
			header += element.Signature + "\n"
			fmt.Println(element.Signature)
		}

		fmt.Println("================================================")

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
			_, err := store.Save(llm.VectorRecord{
				Prompt:    embedding.Prompt,
				Embedding: embedding.Embedding,
				Id:        embedding.Id,
			})

			if err != nil {
				fmt.Println("😡:", err)
			}
		}
	}

}
