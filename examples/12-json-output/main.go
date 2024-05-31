/*
Topic: Parakeet
Generate a chat completion with Ollama and parakeet
Output in json
no streaming
*/

package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "phi3:mini"
	//model := "tinyllama"
	//model := "qwen:0.5b" // freeze

	dataContent := `<context>
	This is the data related to a chicken
	scientific_name: Gallus gallus domesticus
	main_species: Red Junglefowl
	average_length: 46 cm (male), 33 cm (female)
	average_weight: 5 kg (male), 2.5 kg (female)
	average_lifespan: 6-8 years
	countries: Worldwide
	</context>`

	systemContent := `You are a helpful AI assistant. The user will enter the name of an animal.
	The assistant will then return the following information about the annimal:
	- the scientific name of the animal (the name of json field is: scientific_name)
	- the main species of the animal  (the name of json field is: main_species)
	- the decimal average length of the animal (the name of json field is: average_length)
	- the decimal average weight of the animal (the name of json field is: average_weight)
	- the decimal average lifespan of the animal (the name of json field is: average_lifespan)
	- the countries where the animal lives into json array of strings (the name of json field is: countries)
	Output the results in JSON format and trim the spaces of the sentence.
	Use the provided context to give the data`

	userContent := `chicken`

	options := llm.Options{
		Temperature: 0.0, 
		RepeatLastN: 2, 
		RepeatPenalty: 2.0,
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: dataContent},
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
		Format: "json",
		Raw: true, // Why?
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.Content)
}
