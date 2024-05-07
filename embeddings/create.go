package embeddings

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/parakeet-nest/parakeet/llm"
)

type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

func CreateEmbedding(ollamaUrl string, query llm.Query4Embedding, id string) (llm.VectorRecord, error) {
	//log.Println("⏳ Creating embedding... ", id)
	jsonData, err := json.Marshal(query)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	req, err := http.NewRequest(http.MethodPost, ollamaUrl+"/api/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return llm.VectorRecord{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return llm.VectorRecord{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	var answer EmbeddingResponse
	err = json.Unmarshal([]byte(string(body)), &answer)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	vectorRecord := llm.VectorRecord{
		Prompt:    query.Prompt,
		Embedding: answer.Embedding,
		Id:        id,
	}

	return vectorRecord, nil
}
