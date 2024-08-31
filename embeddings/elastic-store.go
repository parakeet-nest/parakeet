package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/parakeet-nest/parakeet/llm"
)

/*


 */

type ElasticsearchStore struct {
	//ctx                 context.Context
	indexName           string
	elasticsearchClient *elasticsearch.Client
}

func (ess *ElasticsearchStore) Initialize(addresses []string, user, pwd string, cert []byte, indexName string) error {
	
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  user,
		Password:  pwd,
		CACert:    cert,
	}
	elasticsearchClient, err := elasticsearch.NewClient(cfg)

	ess.elasticsearchClient = elasticsearchClient
	ess.indexName = indexName

	return err
}

func (ess *ElasticsearchStore) Save(vectorRecord llm.VectorRecord) (llm.VectorRecord, error) {
	// Convert the document/vectorRecord/embedding to JSON
	docJSON, err := json.Marshal(vectorRecord)
	
	if err != nil {
		// Error marshaling document to JSON
		return llm.VectorRecord{}, err
	}
	// Index the document in Elasticsearch
	req := esapi.IndexRequest{
		Index:      ess.indexName,
		DocumentID: vectorRecord.Id,
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), ess.elasticsearchClient)
	if err != nil {
		// Error indexing document
		return llm.VectorRecord{}, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return llm.VectorRecord{}, errors.New("Error response from Elasticsearch: " + res.String())
	} else {
		return vectorRecord, nil
	}
}

func (ess *ElasticsearchStore) SearchTopNSimilarities(embeddingFromQuestion llm.VectorRecord, size int) ([]llm.VectorRecord, error) {
	// Now search for similar embeddings in Elasticsearch
	query := map[string]interface{}{
		"size": size, // Adjust size based on how many results you want
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
	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		// Error encoding query
		return nil, err
	}
	res, err := ess.elasticsearchClient.Search(
		ess.elasticsearchClient.Search.WithContext(context.Background()),
		ess.elasticsearchClient.Search.WithIndex(ess.indexName),
		ess.elasticsearchClient.Search.WithBody(&buf),
		ess.elasticsearchClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		// Error getting response
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		/*
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
		*/
		return nil, errors.New("Error response from Elasticsearch: " + res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		// Error parsing the response body
		return nil, err
	}
	similarities := []llm.VectorRecord{}

	// Process the search results
	// fmt.Printf("Query took %d ms\n", int(r["took"].(float64)))
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		
		docID := hit.(map[string]interface{})["_id"]
		score := hit.(map[string]interface{})["_score"].(float64)
		//fmt.Printf("Document ID: %s, Score: %f\n", docID, score)

		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		prompt := source["prompt"].(string)
		//fmt.Printf("Prompt: %s", prompt)

		similarities = append(similarities, llm.VectorRecord{
			Prompt: prompt,
			Id: docID.(string),
			Score: score,
		})

	}
	return similarities, nil

}
