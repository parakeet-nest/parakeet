package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"

	"github.com/elastic/go-elasticsearch/v8"
	//"github.com/elastic/go-elasticsearch/v8/esapi"
)

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
	var smallChatModel = "qwen2:0.5b"      // This model is for the chat completion

	if err != nil {
		log.Fatalln("😡:", err)
	}

	userContent := `Who are the monsters of Chronicles of Aethelgard?`
	//userContent := `Tell me more about Keegorg`

	// Create an embedding from the question
	embeddingFromQuestion, err := embeddings.CreateEmbedding(
		ollamaUrl,
		llm.Query4Embedding{
			Model:  embeddingsModel,
			Prompt: userContent,
		},
		"question",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}
	fmt.Println("🔎 searching for similarity...")

	// Now search for similar embeddings in Elasticsearch
	query := map[string]interface{}{
		"size": 3, // Adjust size based on how many results you want
		"query": map[string]interface{}{
			"script_score": map[string]interface{}{
				"query": map[string]interface{}{
					"match_all": map[string]interface{}{},
				},
				"script": map[string]interface{}{
					"source": "cosineSimilarity(params.query_vector, 'embedding') + 1.0",
					"params": map[string]interface{}{
						"query_vector": embeddingFromQuestion.Embedding,
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := elasticsearchClient.Search(
		elasticsearchClient.Search.WithContext(context.Background()),
		elasticsearchClient.Search.WithIndex("my-index"), // Replace with your actual index name
		elasticsearchClient.Search.WithBody(&buf),
		elasticsearchClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	similarities := []llm.VectorRecord{}

	// Process the search results
	fmt.Printf("Query took %d ms\n", int(r["took"].(float64)))
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		docID := hit.(map[string]interface{})["_id"]
		score := hit.(map[string]interface{})["_score"].(float64)
		fmt.Printf("Document ID: %s, Score: %f\n", docID, score)

		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		prompt := source["prompt"].(string)
		
		//fmt.Printf("Prompt: %s", prompt)

		similarities = append(similarities, llm.VectorRecord{
			Prompt: prompt,
		})

	}

	documentsContent := embeddings.GenerateContentFromSimilarities(similarities)

	systemContent := `You are the dungeon master,
	expert at interpreting and answering questions based on provided sources.
	Using only the provided context, answer the user's question 
	to the best of your ability using only the resources provided. 
	Be verbose!`

	queryChat := llm.Query{
		Model: smallChatModel,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
		Options: llm.Options{
			Temperature:   0.0,
			RepeatLastN:   2,
			RepeatPenalty: 3.0,
			TopK:          10,
			TopP:          0.5,
		},
	}

	fmt.Println()
	fmt.Println("🤖 answer:")

	// Answer the question
	_, err = completion.ChatStream(ollamaUrl, queryChat,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

	fmt.Println()

}
