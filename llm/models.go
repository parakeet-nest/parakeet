package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type ModelList struct {
	Models []ModelInfo `json:"models"`
}

type ModelInfo struct {
	Name       string    `json:"name"`
	ModifiedAt time.Time `json:"modified_at"`
	Size       int64     `json:"size"`
	Digest     string    `json:"digest"`
	Details    Details   `json:"details"`
}

type Details struct {
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}


func getModelsList(url, tokenHeaderName, tokenHeaderValue string) (ModelList, int, error) {

	req, err := http.NewRequest(http.MethodGet, url+"/api/tags", nil)
	if err != nil {
		return ModelList{}, http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if tokenHeaderName != "" && tokenHeaderValue != "" {
		req.Header.Set(tokenHeaderName, tokenHeaderValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ModelList{}, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ModelList{}, resp.StatusCode, errors.New("Error: status code: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ModelList{}, resp.StatusCode, err
	}

	var models ModelList
	err = json.Unmarshal(body, &models)
	if err != nil {
		return ModelList{}, resp.StatusCode, err
	}
	return models, resp.StatusCode, nil
}

func GetModelsList(url string) (ModelList, int, error) {
	return getModelsList(url, "", "")
}

func GetModelsListWithToken(url, tokenHeaderName, tokenHeaderValue string) (ModelList, int, error) {
	return getModelsList(url, tokenHeaderName, tokenHeaderValue) 
}

type ModelInformation struct {
	Modelfile  string `json:"modelfile"`
	Parameters string `json:"parameters"`
	Template   string `json:"template"`
	Details    struct {
		Format            string   `json:"format"`
		Family            string   `json:"family"`
		Families          []string `json:"families"`
		ParameterSize     string   `json:"parameter_size"`
		QuantizationLevel string   `json:"quantization_level"`
	} `json:"details"`
}

// showModelInformation retrieves information about a model from the specified URL.
//
// Parameters:
// - url: the base URL of the API.
// - model: the name of the model to retrieve information for.
//
// Returns:
// - ModelInformation: the information about the model.
// - int: the HTTP status code of the response.
// - error: an error if the request fails.
func showModelInformation(url, model , tokenHeaderName, tokenHeaderValue string) (ModelInformation, int, error) {

	req, err := http.NewRequest(http.MethodPost, url+"/api/show", bytes.NewBuffer([]byte(`{"name":"`+model+`"}`)))
	if err != nil {
		return ModelInformation{}, http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if tokenHeaderName != "" && tokenHeaderValue != "" {
		req.Header.Set(tokenHeaderName, tokenHeaderValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ModelInformation{}, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ModelInformation{}, resp.StatusCode, errors.New("Error: status code: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ModelInformation{}, resp.StatusCode, err
	}

	var info ModelInformation
	err = json.Unmarshal(body, &info)
	if err != nil {
		return ModelInformation{}, resp.StatusCode, err
	}
	return info, resp.StatusCode, nil

}

func ShowModelInformation(url, model string) (ModelInformation, int, error) {
	return showModelInformation(url, model, "", "")
}

func ShowModelInformationWithToken(url, model , tokenHeaderName, tokenHeaderValue string) (ModelInformation, int, error) {
	return showModelInformation(url, model, tokenHeaderName, tokenHeaderValue)
}


type PullResult struct {
	Status string `json:"status"`
}

// pullModel sends a POST request to the specified URL to pull a model with the given name.
//
// Parameters:
// - url: The URL to send the request to.
// - model: The name of the model to pull.
//
// Returns:
// - PullResult: The result of the pull operation.
// - int: The HTTP status code of the response.
// - error: An error if the request fails.
func pullModel(url, model, tokenHeaderName, tokenHeaderValue string) (PullResult, int, error) {

	req, err := http.NewRequest(http.MethodPost, url+"/api/pull", bytes.NewBuffer([]byte(`{"name":"`+model+`","stream":false}`)))
	if err != nil {
		return PullResult{}, http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if tokenHeaderName != "" && tokenHeaderValue != "" {
		req.Header.Set(tokenHeaderName, tokenHeaderValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PullResult{}, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PullResult{}, resp.StatusCode, errors.New("Error: status code: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PullResult{}, resp.StatusCode, err
	}

	var pullResult PullResult
	err = json.Unmarshal(body, &pullResult)
	if err != nil {
		return PullResult{}, resp.StatusCode, err
	}
	return pullResult, resp.StatusCode, nil

}

func PullModel(url, model string) (PullResult, int, error) {
	return pullModel(url, model, "", "")
}

func PullModelWithToken(url, model, tokenHeaderName, tokenHeaderValue string) (PullResult, int, error) {
	return pullModel(url, model, tokenHeaderName, tokenHeaderValue)
}

// TODO:
// - make a stream version of pull
// - make a version with token
