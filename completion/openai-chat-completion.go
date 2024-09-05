package completion

// OpenAI API support
// https://beta.openai.com/docs/api-reference/chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/parakeet-nest/parakeet/llm"
)

func ChatWithOpenAI(url string, query llm.OpenAIQuery) (llm.OpenAIAnswer, error) {

	/*
		if url == "" {
			url = "https://api.openai.com/v1"
		}
	*/
	kindOfCompletion := "/chat/completions"

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

// func ChatWithOpenAIStream(url string, query llm.OpenAIQuery, onChunk func(llm.OpenAIAnswer) error) (llm.OpenAIAnswer, error) {
func ChatWithOpenAIStream(url string, query llm.OpenAIQuery, onChunk func(llm.OpenAIAnswer) error) error {

	kindOfCompletion := "/chat/completions"

	query.Stream = true

	if query.Verbose {
		fmt.Println("[llm/query]", query.ToJsonString())
		//fmt.Println("[llm/completion]", answer.ToJsonString())
	}

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return err
		//return llm.OpenAIAnswer{}, err
	}
	// -----------------------

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, url+kindOfCompletion, strings.NewReader(string(jsonQuery)))
	//req, err := http.NewRequest(http.MethodPost, url+kindOfCompletion, bytes.NewBuffer(jsonQuery))
	if err != nil {
		return err
		//return llm.OpenAIAnswer{}, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+query.OpenAIAPIKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//return llm.OpenAIAnswer{}, err
		return err
	}

	defer resp.Body.Close()

	var fullAnswer llm.OpenAIAnswer
	// Read and stream the response
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			// Parse the response JSON
			var response llm.OpenAIAnswer
			err := json.Unmarshal([]byte(data), &response)
			if err != nil {
				return err
				//return llm.OpenAIAnswer{}, err
				//fmt.Println("Error parsing JSON:", err)
				//continue
			}
			fullAnswer = response // for the verbose mode
			//fullAnswer.Choices[0].Delta.Content += response.Choices[0].Delta.Content
			err = onChunk(response)
			if err != nil {
				return err
			}
			// Print the streamed text
			//fmt.Print("", response.Choices[0].Text)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
		//return llm.OpenAIAnswer{}, err
		//fmt.Println("Error reading response:", err)
	}

	if query.Verbose {
		//fmt.Println("[llm/query]", query.ToJsonString())
		fmt.Println()
		fmt.Println("[llm/completion]", fullAnswer.ToJsonString())
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Error: status code: " + resp.Status)
		//return llm.OpenAIAnswer{}, errors.New("Error: status code: " + resp.Status)
	} else {
		return nil
		//return fullAnswer, nil
	}

}
