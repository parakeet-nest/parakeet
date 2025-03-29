package completion

import (
	"fmt"
	"testing"

	"github.com/parakeet-nest/parakeet/llm"
)

func TestMetricsOfAnswers(t *testing.T) {
	ollamaUrl := "http://localhost:11434"
	model := "qwen2.5:3b"
	userContent := "Who is James T Kirk?"

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
		RepeatLastN: 2,
	}
	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			//{Role: "system", Content: systemContent},
			//{Role: "system", Content: contextContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	answer, err := Chat(ollamaUrl, query)
	if err != nil {
		t.Fatal("😡:", err)
	}
	fmt.Println("created: ", answer.CreatedAt)
	fmt.Println("total duration: ", answer.TotalDuration)
	fmt.Println("load duration: ", answer.LoadDuration)
	fmt.Println("prompt eval count: ", answer.PromptEvalCount)
	fmt.Println("prompt eval duration: ", answer.PromptEvalDuration)
	fmt.Println("eval count: ", answer.EvalCount)
	fmt.Println("eval duration: ", answer.EvalDuration)
	

}
