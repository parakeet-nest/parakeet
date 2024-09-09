package embeddings

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/parakeet-nest/parakeet/llm"
)

func CreateEmbeddingWithOpenAI(url string, query llm.OpenAIQuery4Embedding, id string) (llm.VectorRecord, error) {
	//log.Println("⏳ Creating embedding... ", id)
	jsonData, err := json.Marshal(query)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	// curl https://api.openai.com/v1/embeddings \
	// https://platform.openai.com/docs/guides/embeddings/what-are-embeddings
	req, err := http.NewRequest(http.MethodPost, url+"/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return llm.VectorRecord{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+query.OpenAIAPIKey)

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

	var answer llm.OpenAIEmbeddingResponse
	err = json.Unmarshal([]byte(string(body)), &answer)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	vectorRecord := llm.VectorRecord{
		Prompt:    query.Input,
		Embedding: answer.Data[0].Embedding,
		Id:        id,
	}

	// Sometime vectorRecord.Embedding is empty
	if len(vectorRecord.Embedding) == 0 {
		return llm.VectorRecord{}, errors.New("embedding is empty")
	}

	return vectorRecord, nil
}
