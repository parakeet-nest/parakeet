package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/squawk"
)

func TestEmbeddingsWithModelRunner(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	engineBaseUrl := os.Getenv("MODEL_RUNNER_BASE_URL")
	model := os.Getenv("MODEL_RUNNER_LLM_CHAT")

	toolsList := []llm.Tool{
		{
			Type: "function",
			Function: llm.Function{
				Name:        "hello",
				Description: "Say hello to a given person with his name",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"name": {
							Type:        "string",
							Description: "The name of the person",
						},
					},
					Required: []string{"name"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.Function{
				Name:        "addNumbers",
				Description: "Make an addition of the two given numbers",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"a": {
							Type:        "number",
							Description: "first operand",
						},
						"b": {
							Type:        "number",
							Description: "second operand",
						},
					},
					Required: []string{"a", "b"},
				},
			},
		},
	}

	squawk.New().
		BaseURL(engineBaseUrl).
		Model(model).
		Provider(provider.DockerModelRunner).
		Tools(toolsList).
		Options(llm.SetOptions(map[string]interface{}{
			option.Temperature: 0.0,
		})).
		User(`say "hello" to Bob, say "hello" to Sam`).
		User(`add 2 and 40`).
		FunctionCalling(func(answer llm.Answer, self *squawk.Squawk, err error) {

			var results string
			for _, toolCall := range self.ToolCalls() {

				switch toolCall.Function.Name {
				case "hello":
					results += fmt.Sprintf("Hello %s\n", toolCall.Function.Arguments["name"])
				case "addNumbers":
					a := toolCall.Function.Arguments["a"]
					b := toolCall.Function.Arguments["b"]
					results += fmt.Sprintf("Addition of %v and %v is %v\n", a, b, a.(float64)+b.(float64))
				default:
					results += "" // unknown function
				}

			}
			self.System("RESULTS:\n" + results) // add the result to the system message
		}).
		User("Use the results and format the output with fancy emojis").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

}
