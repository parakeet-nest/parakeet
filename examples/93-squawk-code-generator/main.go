package main
// ✋ this is a 🚧 WIP, it does not work as expected
import (
	"fmt"

	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/squawk"
)

func main() {
	ollamaBaseUrl := "http://localhost:11434"
	model := "qwen2.5:1.5b"

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.5,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.2,
	})

	squawk.New().
		Model(model).
		BaseURL(ollamaBaseUrl).
		Provider(provider.Ollama).
		Options(options).
		System(`You are a useful AI agent:
		  - You are a Golang expert.
		  - Generate only Golang source code.
		Your job is to transform the JSON string provided by the user into Golang structures.
		`).
		User(`
			{
				"name": "John",
				"age": 30,
				"city": "New York",
				"languages": ["English", "Spanish"],
				"address": {
					"street": "123 Main St",
					"zip": "10001"
				},
				"active": true,
				"balance": 100.50,
				"tags": ["golang", "json", "example"],
				"metadata": {
					"created_at": "2023-10-01T12:00:00Z",
					"updated_at": "2023-10-02T12:00:00Z"
				}
			}
		`).
		Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
			fmt.Println(answer.Message.Content)
		})

}
// https://github.com/traefik/yaegi
// perhaps use the yaegi library to execute the code