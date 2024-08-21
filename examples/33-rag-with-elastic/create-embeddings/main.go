package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

/*
type MyDocument struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Embedding []float64 `json:"embedding"`
}
*/

func main() {

	cert, _ := os.ReadFile(os.Getenv("ELASTIC_CERT_PATH"))

	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTIC_ADDRESS"),
		},
		Username: os.Getenv("ELASTIC_USERNAME"),
		Password: os.Getenv("ELASTIC_PASSWORD"),
		CACert:   cert,
	}
	elasticsearchClient, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	ollamaUrl := "http://localhost:11434"
	var embeddingsModel = "all-minilm:33m" // This model is for the embeddings of the documents

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

			//embedding.MetaData = "📝 chunk num: " + strconv.Itoa(idx)
			//redisStore.Save(embedding)
			// Create a document with an embedding

			/*
			doc := MyDocument{
				ID:        embedding.Id,
				Text:      embedding.Prompt,
				Embedding: embedding.Embedding,
			}
			*/

			// Convert the document to JSON
			docJSON, err := json.Marshal(embedding)
			if err != nil {
				log.Fatalf("Error marshaling document to JSON: %s", err)
			}

			// Index the document in Elasticsearch
			req := esapi.IndexRequest{
				Index:      "my-index", // Your index name
				DocumentID: embedding.Id,
				Body:       bytes.NewReader(docJSON),
				Refresh:    "true",
			}

			//log.Println(es.SSL.Certificates)

			res, err := req.Do(context.Background(), elasticsearchClient)
			log.Println(res.String())

			if err != nil {

				log.Fatalf("Error indexing document: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("Error response from Elasticsearch: %s", res.String())
			} else {
				log.Printf("Document indexed successfully")
			}
		}
	}
}
