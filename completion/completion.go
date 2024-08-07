package completion

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/parakeet-nest/parakeet/llm"
)

func completion(url string, kindOfCompletion string, query llm.Query) (llm.Answer, error) {

	query.Stream = false

	// if tool call is not used
	if query.Tools == nil {
		query.Tools = []llm.Tool{}
	}

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return llm.Answer{}, err
	}

	req, err := http.NewRequest(http.MethodPost, url+"/api/"+kindOfCompletion, bytes.NewBuffer(jsonQuery))
	if err != nil {
		return llm.Answer{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.Answer{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return llm.Answer{}, errors.New("Error: status code: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return llm.Answer{}, err
	}

	var answer llm.Answer
	err = json.Unmarshal(body, &answer)

	//fmt.Println("🔵🟦:", answer.Message.ToolCalls)
	//fmt.Println("🔵🟦:", reflect.TypeOf(answer.Message.ToolCalls))

	if err != nil {
		return llm.Answer{}, err
	}

	return answer, nil

}

func completionStream(url string, kindOfCompletion string, query llm.Query, onChunk func(llm.Answer) error) (llm.Answer, error) {
	query.Stream = true

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return llm.Answer{}, err
	}

	resp, err := http.Post(url+"/api/"+kindOfCompletion, "application/json; charset=utf-8", bytes.NewBuffer(jsonQuery))
	if err != nil {
		return llm.Answer{}, err
	}
	reader := bufio.NewReader(resp.Body)

	var fullAnswer llm.Answer
	var answer llm.Answer
	for {

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return llm.Answer{}, err
		}

		err = json.Unmarshal(line, &answer)
		if err != nil {
			onChunk(llm.Answer{})
		}
		fullAnswer.Message.Content += answer.Message.Content

		// ? 🤔 and if I used answer + error as a parameter?
		err = onChunk(answer)

		// generate an error to stop the stream
		if err != nil {
			return llm.Answer{}, err
		}
	}
	fullAnswer.Message.Role = answer.Message.Role
	return fullAnswer, nil

}
