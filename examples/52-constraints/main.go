
package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/tools"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "gemma2:2b"
	//toolsModel := "mistral:7b"
	toolsModel := "qwen2.5:latest"

	//toolsModel := "allenporter/xlam:1b"
	//toolsModel := "sam4096/qwen2tools:0.5b"

	options := llm.SetOptions(map[string]any{
		option.Temperature:   0.0,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.0,
		option.TopK:          10,
		option.TopP:          0.5,
		//option.Verbose: true,
	})

	toolsList := []llm.Tool{
		{
			Type: "function",
			Function: llm.Function{
				Name:        "doNotSpeakAboutKubernetes",
				Description: "explain the difference between Kubernetes and an other orchestrator",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"question": {
							Type:        "string",
							Description: "The other orchestrator",
						},
					},
					Required: []string{"question"},
				},
			},
		},
	}

	toolsContent, err := tools.GenerateAvailableToolsContent(toolsList)
	if err != nil {
		log.Fatal("😡:", err)
	}

	canISpeakAboutThis := func(question string) bool {

		options := llm.SetOptions(map[string]interface{}{
			option.Temperature:   0.0,
			option.RepeatLastN:   2,
			option.RepeatPenalty: 2.0,
		})

		messages := []llm.Message{
			{Role: "system", Content: toolsContent},
			{Role: "user", Content: tools.GenerateUserToolsInstructions(question)},
		}

		query := llm.Query{
			Model:    toolsModel,
			Messages: messages,
			Options:  options,
			//Format:   "json",
			//Raw:      true,
		}

		answer, err := completion.Chat(ollamaUrl, query)
		if err != nil {
			log.Fatal("😡:", err)
		}


		jsonRes, err := gear.JSONParse(answer.Message.Content)

		if err != nil {
			log.Fatal("😡:", err)
		}
		if jsonRes["name"] == nil {
			fmt.Println("✅ I can speak about this")
			return true

		} else {
			fmt.Println("✋ 🔴 I do not understand")

			functionName := jsonRes["name"].(string)
			questionParam := jsonRes["arguments"].(map[string]interface{})["question"].(string)
	
			fmt.Println("Calling", functionName, "with:", questionParam)
			return false
		}
	}

	/*
		systemContent := `You are an expert in computer programming and container orchestrators.`
		instructions := `The construction of your response must be structured.
		Make all your response while following all constraints.`
		constraints := `Constraints:
		- Avoid mentioning Kubernetes.
		- Never speak about Kubernetes.
		`
	*/
	// Combine system content, instructions, and constraints into one cohesive message.
	systemContent := `You are an expert in computer programming and container orchestrators. 
		Here are some important instructions and constraints for your responses:
		
		Instructions:
		- Structure your responses carefully and clearly.
		- Make sure to follow all constraints strictly.
		
		Constraints:
		- Under no circumstances should you mention Kubernetes.
		- If the question references Kubernetes, do not provide an answer related to it.
		- Focus on other container orchestration systems such as Docker Swarm, Nomad, etc., but avoid Kubernetes completely.
		- Responses that mention Kubernetes in any form are unacceptable.`

	askMeAnything := func(question string) {
		
			query := llm.Query{
				Model: model,
				Messages: []llm.Message{
					{Role: "system", Content: systemContent},
					//{Role: "system", Content: instructions},
					//{Role: "system", Content: constraints},
					{Role: "user", Content: "[Brief] " + question},
				},
				Options: options,
			}

			_, err := completion.ChatStream(ollamaUrl, query,
				func(answer llm.Answer) error {
					fmt.Print(answer.Message.Content)
					return nil
				})

			if err != nil {
				log.Fatal("😡:", err)
			}
		

	}

	question := `make a comparison study of orchestrators`

	if canISpeakAboutThis(question) {
		askMeAnything(question)
	}

	fmt.Println("=====================================")

	question = `Can you explain the difference between Kubernetes and Docker Swarm?`

	if canISpeakAboutThis(question) {
		askMeAnything(question)
	}

	fmt.Println("=====================================")

	question = `What is Kubernetes?`

	if canISpeakAboutThis(question) {
		askMeAnything(question)
	}

	fmt.Println("=====================================")

	question = `Should I use Kubernetes instead of Rancher?`

	if canISpeakAboutThis(question) {
		askMeAnything(question)
	}

}
