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

func Generate(url string, query llm.GenQuery) (llm.GenAnswer, error) {
	kindOfCompletion := "generate"

	query.Stream = false

	if query.Options.Verbose {
		fmt.Println("[llm/query]", query.ToJsonString())
		//fmt.Println("[llm/completion]", answer.ToJsonString())
	}

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return llm.GenAnswer{}, err
	}

	req, err := http.NewRequest(http.MethodPost, url+"/api/"+kindOfCompletion, bytes.NewBuffer(jsonQuery))
	if err != nil {
		return llm.GenAnswer{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if query.TokenHeaderName != "" && query.TokenHeaderValue != "" {
		req.Header.Set(query.TokenHeaderName, query.TokenHeaderValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.GenAnswer{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return llm.GenAnswer{}, err
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
				return llm.GenAnswer{}, err
			}
			if strings.HasPrefix(completionError.Error, "model") && strings.HasSuffix(completionError.Error, "not found, try pulling it first") {
				modelNotFound.Code = resp.StatusCode
				modelNotFound.Message = completionError.Error
				modelNotFound.Model = query.Model
			}
			return llm.GenAnswer{}, &modelNotFound
		}

		return llm.GenAnswer{}, errors.New("Error: status code: " + resp.Status + "\n" + string(body))
	}

	var answer llm.GenAnswer
	err = json.Unmarshal(body, &answer)

	if err != nil {
		return llm.GenAnswer{}, err
	}

	if query.Options.Verbose {
		//fmt.Println("[llm/query]", query.ToJsonString())
		fmt.Println()
		fmt.Println("[llm/completion]", answer.ToJsonString())
	}

	return answer, nil

}

func GenerateStream(url string, query llm.GenQuery, onChunk func(llm.GenAnswer) error) (llm.GenAnswer, error) {
	kindOfCompletion := "generate"

	query.Stream = true

	if query.Options.Verbose {
		fmt.Println("[llm/query]", query.ToJsonString())
		//fmt.Println("[llm/completion]", answer.ToJsonString())
	}

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return llm.GenAnswer{}, err
	}
	// -----------------------

	req, err := http.NewRequest(http.MethodPost, url+"/api/"+kindOfCompletion, bytes.NewBuffer(jsonQuery))
	if err != nil {
		return llm.GenAnswer{}, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if query.TokenHeaderName != "" && query.TokenHeaderValue != "" {
		req.Header.Set(query.TokenHeaderName, query.TokenHeaderValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.GenAnswer{}, err
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

	var fullAnswer llm.GenAnswer
	var answer llm.GenAnswer

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
			return llm.GenAnswer{}, errors.New("Error: status code: " + resp.Status)
			//return llm.Answer{}, err
		}

		err = json.Unmarshal(line, &answer)
		if err != nil {
			onChunk(llm.GenAnswer{})
		}
		fullAnswer.Response += answer.Response
		fullAnswer.Context = answer.Context

		// ? 🤔 and if I used answer + error as a parameter?
		err = onChunk(answer)

		// generate an error to stop the stream
		if err != nil {
			return llm.GenAnswer{}, err
		}
	}

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

			return llm.GenAnswer{}, &modelNotFound
		}

		return llm.GenAnswer{}, errors.New("Error: status code: " + resp.Status)
	} else {
		return fullAnswer, nil
	}
}
