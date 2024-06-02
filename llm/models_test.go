package llm

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetInformationModel(t *testing.T) {
	ollamaUrl := "http://localhost:11434"

	info, status, err := ShowModelInformation(ollamaUrl, "tinyllama:latest")

	if err != nil {
		t.Fatal(err, status)
	}
	fmt.Println(info)

}

func TestGetInformationUnknownModel(t *testing.T) {
	ollamaUrl := "http://localhost:11434"

	_, status, err := ShowModelInformation(ollamaUrl, "tiny_llama")

	if err != nil {
		t.Log(status)
		if status != http.StatusNotFound {
			t.Fatal(err, status)
		}
		fmt.Println("Model not found")
	}
}
