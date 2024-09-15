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
	"github.com/parakeet-nest/parakeet/enums/option"

	"fmt"
	"log"
)

/*
This Go code snippet defines a main function that performs a chat completion using the Ollama and Parakeet libraries.

Here's a breakdown of what the code does:

1. It sets the ollamaUrl variable to the URL of the Ollama server running on localhost.
2 .It sets the model variable to the name of the model to be used for the chat completion.
3. It defines the dataContent variable, which contains contextual information about a chicken.
4 .It defines the systemContent variable, which describes the role of the AI assistant and the expected inputs and outputs.
5 .It defines the userContent variable, which represents the user's input.
6. It creates an options variable with various settings for the chat completion, such as temperature, repeat penalties, and formatting options.
7. It creates a query variable, which contains the model name, the messages to be exchanged during the chat completion, and the options.
8. It calls the completion.Chat function with the Ollama URL and the query to perform the chat completion.
9. If the chat completion is successful, it prints the content of the last message received from the chat.
*/
func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	
	model := "tinydolphin"
	//model := "qwen:0.5b" // freeze

	/*
	dataContent := `<context>
	This is the data related to a chicken
	scientific_name: Gallus gallus domesticus
	main_species: Red Junglefowl
	average_length: 46 cm (male), 33 cm (female)
	average_weight: 5 kg (male), 2.5 kg (female)
	average_lifespan: 6-8 years
	countries: Worldwide
	</context>`
	*/

	systemContent := `You are a helpful AI assistant. The user will enter the name of an animal.
	The assistant will then return the following information about the animal:
	- the scientific name of the animal (the name of json field is: scientific_name)
	- the main species of the animal  (the name of json field is: main_species)
	- the decimal average length of the animal (the name of json field is: average_length)
	- the decimal average weight of the animal (the name of json field is: average_weight)
	- the decimal average lifespan of the animal (the name of json field is: average_lifespan)
	- the countries where the animal lives into json array of strings (the name of json field is: countries)
	Output the results in JSON format and trim the spaces of the sentence.
	Use the provided context to give the data`

	userContent := `chicken`

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 2.0,
	})

	
	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			//{Role: "system", Content: dataContent},
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
