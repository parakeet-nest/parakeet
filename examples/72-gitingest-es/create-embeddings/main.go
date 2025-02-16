package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
)

type RepositoryChunk struct {
	Header string
	Code   string
}

func createChunk(piecesOfContent []string) []RepositoryChunk {
	// Calculate how many items we can create (chunks length divided by 2)
	numItems := len(piecesOfContent) / 2

	// Create the items slice with the calculated capacity
	items := make([]RepositoryChunk, numItems)

	// Fill the items slice
	for i := 0; i < numItems; i++ {
		items[i] = RepositoryChunk{
			Code:   piecesOfContent[i*2],                      // Odd indices (1, 3, 5, ...) for codes
			Header: strings.TrimSpace(piecesOfContent[i*2+1]), // Even indices (0, 2, 4, ...) for titles
		}
	}

	return items
}

func main() {
	fmt.Println("Hello, World!")

	contentPath := gear.GetEnvString("CONTENT_PATH", "../data/content.txt")

	ollamaUrl := gear.GetEnvString("OLLAMA_BASE_URL", "http://localhost:11434")

	embeddingsModel := gear.GetEnvString("LLM_EMBEDDINGS", "mxbai-embed-large")


	elasticStore := embeddings.ElasticsearchStore{}
	err := elasticStore.Initialize(
		[]string{
			os.Getenv("ELASTICSEARCH_HOSTS"),
		},
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
		nil,
		os.Getenv("ELASTICSEARCH_INDEX"),
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	// open ../data/content.txt
	// read the content
	allSourceCodes, err := os.ReadFile(contentPath)
	if err != nil {
		log.Fatal(err)
	}

	// Ok, it's not my best idea to use this delimiter
	chunksFromAllSourceCodes := content.SplitTextWithDelimiter(
		string(allSourceCodes),
		`================================================`,
	)

	/*
		for _, chunk := range chunksFromAllSourceCodes {
			fmt.Println("📝", chunk)
		}
	*/

	bigChunks := createChunk(chunksFromAllSourceCodes)

	for idxFirsLevel, bigChunk := range bigChunks {
		fmt.Println("✋", idxFirsLevel, bigChunk.Header)
		fmt.Println("📝", bigChunk.Code)

		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: "## " + bigChunk.Header + ":\n" + bigChunk.Code,
			},
			strconv.Itoa(idxFirsLevel)+"-"+strconv.Itoa(idxFirsLevel),
		)
		if err != nil {
			log.Fatalln("😡:", err)
		}

		embedding.SimpleMetaData = bigChunk.Header

		if _, err = elasticStore.Save(embedding); err != nil {
			log.Fatalln("😡:", err)
		}

		fmt.Println("🎉 Document", embedding.Id, "indexed successfully")

		/*
			smallerChunks := content.ChunkText(bigChunk.Code, 2048, 512)

			for idxSecondLevel, smallChunk := range smallerChunks {
				fmt.Println("🪚", idxSecondLevel, smallChunk)

				embedding, err := embeddings.CreateEmbedding(
					ollamaUrl,
					llm.Query4Embedding{
						Model:  embeddingsModel,
						Prompt: "## "+bigChunk.Header + ":\n" + smallChunk,
					},
					strconv.Itoa(idxFirsLevel)+"-"+strconv.Itoa(idxSecondLevel),
				)
				if err != nil {
					log.Fatalln("😡:", err)
				}

				embedding.SimpleMetaData = bigChunk.Header

				if _, err = elasticStore.Save(embedding); err != nil {
					log.Fatalln("😡:", err)
				}

				fmt.Println("🎉 Document", embedding.Id, "indexed successfully")

			}
		*/

	}

}
