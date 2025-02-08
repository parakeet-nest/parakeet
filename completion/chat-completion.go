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

	"github.com/parakeet-nest/parakeet/llm"
)

type ModelNotFoundError struct {
	Code    int
	Message string
	Model   string
}

func (e *ModelNotFoundError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Model: %s", e.Code, e.Message, e.Model)
}

type CompletionError struct {
	Error string
}

func Chat(url string, query llm.Query) (llm.Answer, error) {
	kindOfCompletion := "chat"

	/*
		var outputSchema = ""
		if query.Format != "" && query.Format != "json" {
			outputSchema = query.Format
		}
	*/

	query.Stream = false

	if query.Options.Verbose {
		fmt.Println("[llm/query]", query.ToJsonString())
		//fmt.Println("[llm/completion]", answer.ToJsonString())
	}

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

	if query.TokenHeaderName != "" && query.TokenHeaderValue != "" {
		req.Header.Set(query.TokenHeaderName, query.TokenHeaderValue)
	}

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
		// because, even if the status is not ok (ex 401 Unauthorized)
		// the error == nil

		if resp.StatusCode == http.StatusNotFound {
			var completionError CompletionError
			var modelNotFound ModelNotFoundError
			err = json.Unmarshal(body, &completionError)
			if err != nil {
				return llm.Answer{}, err
			}
			if strings.HasPrefix(completionError.Error, "model") && strings.HasSuffix(completionError.Error, "not found, try pulling it first") {
				modelNotFound.Code = resp.StatusCode
				modelNotFound.Message = completionError.Error
				modelNotFound.Model = query.Model
			}
			return llm.Answer{}, &modelNotFound
		}

		return llm.Answer{}, errors.New("Error: status code: " + resp.Status + "\n" + string(body))

	}

	var answer llm.Answer
	err = json.Unmarshal(body, &answer)

	if err != nil {
		return llm.Answer{}, err
	}

	if query.Options.Verbose {
		//fmt.Println("[llm/query]", query.ToJsonString())
		fmt.Println()
		fmt.Println("[llm/completion]", answer.ToJsonString())
	}

	return answer, nil

}

func ChatStream(url string, query llm.Query, onChunk func(llm.Answer) error) (llm.Answer, error) {
	kindOfCompletion := "chat"

	query.Stream = true

	if query.Options.Verbose {
		fmt.Println("[llm/query]", query.ToJsonString())
		//fmt.Println("[llm/completion]", answer.ToJsonString())
	}

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return llm.Answer{}, err
	}
	// -----------------------

	req, err := http.NewRequest(http.MethodPost, url+"/api/"+kindOfCompletion, bytes.NewBuffer(jsonQuery))
	if err != nil {
		return llm.Answer{}, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if query.TokenHeaderName != "" && query.TokenHeaderValue != "" {
		req.Header.Set(query.TokenHeaderName, query.TokenHeaderValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.Answer{}, err
	}

	defer resp.Body.Close()

	// -----------------------
	/*
		resp, err := http.Post(url+"/api/"+kindOfCompletion, "application/json; charset=utf-8", bytes.NewBuffer(jsonQuery))
		if err != nil {
			return llm.Answer{}, err
		}
	*/
	// -----------------------

	reader := bufio.NewReader(resp.Body)

	var fullAnswer llm.Answer
	var answer llm.Answer

	for {

		line, err := reader.ReadBytes('\n')
		if err != nil {

			if err == io.EOF {
				//&& resp.Status == "200"
				break
			}
			// we need to create a new error because
			// because, even if the status is not ok (ex 401 Unauthorized)
			// the error == nil
			return llm.Answer{}, errors.New("Error: status code: " + resp.Status)
			//return llm.Answer{}, err
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

	if query.Options.Verbose {
		//fmt.Println("[llm/query]", query.ToJsonString())
		fmt.Println()
		fmt.Println("[llm/completion]", answer.ToJsonString())
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			var modelNotFound ModelNotFoundError

			modelNotFound.Code = resp.StatusCode
			modelNotFound.Message = "model " + query.Model + " not found, try pulling it first"
			modelNotFound.Model = query.Model

			return llm.Answer{}, &modelNotFound
		}

		return llm.Answer{}, errors.New("Error: status code: " + resp.Status)
	} else {
		return fullAnswer, nil
	}

}
