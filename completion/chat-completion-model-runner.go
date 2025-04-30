package completion

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/parakeet-nest/parakeet/completion/typesprovider/openai"
	"github.com/parakeet-nest/parakeet/llm"
)

func modelRunnerChat(url string, query llm.Query) (llm.Answer, error) {
	
	openAIQuery := openai.Query{
		Model:    query.Model,
		Messages: query.Messages,

		Stop:        query.Options.Stop,
		Seed:        query.Options.Seed,
		Temperature: query.Options.Temperature,
		TopP:        query.Options.TopP,

		PresencePenalty:  query.Options.PresencePenalty,
		FrequencyPenalty: query.Options.FrequencyPenalty,

		// *** OpenAI specific options ***
		// TODO
		// *** End of OpenAI specific options ***

		Verbose: query.Options.Verbose,

		OpenAIAPIKey: "no-key",

	}

	//TODO: test if the query.Format is empty
	openAIQuery.Responseformat = map[string]interface{}{
		"type": "json_schema",
		"json_schema": map[string]interface{}{
			"name": "my_schema",
			"schema": query.Format,
		},
	}
	//TODO: make it with OpenAI too
	
	
	openAIQuery.Tools = query.Tools

	// if tool call is not used
	if openAIQuery.Tools == nil {
		openAIQuery.Tools = []llm.Tool{}
	}

	if len(openAIQuery.Tools) > 0 {
		openAIQuery.ToolChoice = "auto"
	}

	kindOfCompletion := "/chat/completions"

	openAIQuery.Stream = false

	if openAIQuery.Verbose {
		fmt.Println("[llm/query]", openAIQuery.ToJsonString())
	}

	jsonQuery, err := json.Marshal(openAIQuery)
	if err != nil {
		return llm.Answer{}, err
	}

	//fmt.Println("[llm/query]", string(jsonQuery))

	req, err := http.NewRequest(http.MethodPost, url+kindOfCompletion, bytes.NewBuffer(jsonQuery))
	if err != nil {
		return llm.Answer{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+openAIQuery.OpenAIAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.Answer{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return llm.Answer{}, err
	}

	if resp.StatusCode != http.StatusOK {
		// we need to create a new error because
		// even if the status is not ok (ex 401 Unauthorized)
		// the error == nil
		return llm.Answer{}, errors.New("Error: status code: " + resp.Status + "\n" + string(body))
	}

	var openAIAnswer openai.Answer
	err = json.Unmarshal(body, &openAIAnswer)

	if err != nil {
		return llm.Answer{}, err
	}

	if openAIQuery.Verbose {
		fmt.Println()
		fmt.Println("[llm/completion]", openAIAnswer.ToJsonString())
	}

	return convertOpenAIAnswerToAnswer(openAIAnswer)


}

func modelRunnerChatStream(url string, query llm.Query, onChunk func(llm.Answer) error) (llm.Answer, error) {

	openAIQuery := openai.Query{
		Model:    query.Model,
		Messages: query.Messages,

		Stop:        query.Options.Stop,
		Seed:        query.Options.Seed,
		Temperature: query.Options.Temperature,
		TopP:        query.Options.TopP,

		PresencePenalty:  query.Options.PresencePenalty,
		FrequencyPenalty: query.Options.FrequencyPenalty,

		// *** OpenAI specific options ***
		// TODO
		// *** End of OpenAI specific options ***

		//Stream:            query.Stream,
		Tools: query.Tools,

		Verbose: query.Options.Verbose,

		OpenAIAPIKey: "no-key",
	}

	kindOfCompletion := "/chat/completions"

	openAIQuery.Stream = true

	if openAIQuery.Verbose {
		fmt.Println("[llm/query]", openAIQuery.ToJsonString())
	}

	jsonQuery, err := json.Marshal(openAIQuery)
	if err != nil {
		//return err
		return llm.Answer{}, err
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, url+kindOfCompletion, strings.NewReader(string(jsonQuery)))
	//req, err := http.NewRequest(http.MethodPost, url+kindOfCompletion, bytes.NewBuffer(jsonQuery))
	if err != nil {
		//return err
		return llm.Answer{}, err
	}
	// Set headers
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+openAIQuery.OpenAIAPIKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.Answer{}, err
		//return err
	}

	defer resp.Body.Close()

	fullAnswer := ""
	var returnAnswer = llm.Answer{}

	fullResponse := openai.Answer{}
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
			var response openai.Answer
			err := json.Unmarshal([]byte(data), &response)
			if err != nil {
				return llm.Answer{}, err
			}

			fullAnswer += response.Choices[0].Delta.Content

			fullResponse = response
			// Conversion of openAIAnswer  to llm.Answer
			answer := llm.Answer{
				Model: response.Model,
				Message: llm.Message{
					Role:    response.Choices[0].Message.Role,
					Content: response.Choices[0].Delta.Content,
				},
			}

			err = onChunk(answer)
			if err != nil {
				return llm.Answer{}, err
			}

			returnAnswer = llm.Answer{
				Model: response.Model,
				Message: llm.Message{
					Role:    response.Choices[0].Message.Role,
					Content: fullAnswer,
				},
			}

		}
	}

	if err := scanner.Err(); err != nil {
		return llm.Answer{}, err
	}

	if openAIQuery.Verbose {
		fmt.Println()
		fmt.Println("[llm/completion]", fullResponse.ToJsonString())
	}

	if resp.StatusCode != http.StatusOK {
		return llm.Answer{}, errors.New("Error: status code: " + resp.Status)
	} else {
		return returnAnswer, nil
	}

}
