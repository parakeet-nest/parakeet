package embeddings

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/parakeet-nest/parakeet/embeddings/typesprovider/openai"
	"github.com/parakeet-nest/parakeet/llm"
)

func modelRunnerCreateEmbedding(modelRunnerURL string, query llm.Query4Embedding, id string) (llm.VectorRecord, error) {

	var openAIQuery4Embedding = openai.Query4Embedding{
		Input:        query.Prompt,
		Model:        query.Model,
		OpenAIAPIKey: "no-key",
	}

	jsonData, err := json.Marshal(openAIQuery4Embedding)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	// curl https://api.openai.com/v1/embeddings \
	// https://platform.openai.com/docs/guides/embeddings/what-are-embeddings
	req, err := http.NewRequest(http.MethodPost, modelRunnerURL+"/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return llm.VectorRecord{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+openAIQuery4Embedding.OpenAIAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//return llm.VectorRecord{}, err

		// we need to create a new error because
		// because, even if the status is not ok (ex 401 Unauthorized)
		// the error == nil
		return llm.VectorRecord{}, errors.New("Error: status code: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	var answer openai.EmbeddingResponse
	err = json.Unmarshal([]byte(string(body)), &answer)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	vectorRecord := llm.VectorRecord{
		Prompt:    openAIQuery4Embedding.Input,
		Embedding: answer.Data[0].Embedding,
		Id:        id,
	}

	// Sometime vectorRecord.Embedding is empty
	if len(vectorRecord.Embedding) == 0 {
		return llm.VectorRecord{}, errors.New("embedding is empty")
	}

	return vectorRecord, nil

}
