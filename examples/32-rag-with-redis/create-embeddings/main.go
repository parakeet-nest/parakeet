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

	redisStore := embeddings.RedisVectorStore{}
	err := redisStore.Initialize("localhost:6379", "", "chronicles-bucket")

	if err != nil {
		log.Fatalln("😡:", err)
	}

	rulesContent, err := content.ReadTextFile("./chronicles.md")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	chunks := content.SplitTextWithRegex(rulesContent, `## *`)

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
			redisStore.Save(embedding)
		}
	}
}
