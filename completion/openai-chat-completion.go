package completion

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/parakeet-nest/parakeet/llm"
)

func ChatWithOpenAI(url string, query llm.OpenAIQuery) (llm.OpenAIAnswer, error) {
	if url == "" {
		url = "https://api.openai.com"
	}
	kindOfCompletion := "/v1/chat/completions"

	query.Stream = false

	//query.Verbose = nil

	if query.Verbose {
		fmt.Println("[llm/query]", query.ToJsonString())
		//fmt.Println("[llm/completion]", answer.ToJsonString())
	}

	// if tool call is not used
	if query.Tools == nil {
		query.Tools = []llm.Tool{}
	}

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return llm.OpenAIAnswer{}, err
	}

	req, err := http.NewRequest(http.MethodPost, url+kindOfCompletion, bytes.NewBuffer(jsonQuery))
	if err != nil {
		return llm.OpenAIAnswer{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+query.OpenAIAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.OpenAIAnswer{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return llm.OpenAIAnswer{}, err
	}

	if resp.StatusCode != http.StatusOK {
		// we need to create a new error because
		// because, even if the status is not ok (ex 401 Unauthorized)
		// the error == nil
		return llm.OpenAIAnswer{}, errors.New("Error: status code: " + resp.Status + "\n" + string(body))
	}

	var answer llm.OpenAIAnswer
	err = json.Unmarshal(body, &answer)

	if err != nil {
		return llm.OpenAIAnswer{}, err
	}

	if query.Verbose {
		//fmt.Println("[llm/query]", query.ToJsonString())
		fmt.Println()
		fmt.Println("[llm/completion]", answer.ToJsonString())
	}

	return answer, nil

}

func ChatWithOpenAIStream(url string, query llm.OpenAIQuery, onChunk func(llm.OpenAIAnswer) error) (llm.OpenAIAnswer, error) {
	return llm.OpenAIAnswer{}, nil
}
