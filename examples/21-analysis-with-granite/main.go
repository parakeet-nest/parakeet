package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/prompt"
	"github.com/parakeet-nest/parakeet/enums/option"

)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "granite-code:3b"

	info, status, err := llm.ShowModelInformation(ollamaUrl, model)
	if err != nil {
		if status == 404 {
			fmt.Println("✋ we need to pull the model")
			result, status, errPull := llm.PullModel(ollamaUrl, model)
			if errPull != nil {
				log.Fatal(errPull, status)
			}
			fmt.Println(result)
		}
		log.Fatal(err)
	}

	fmt.Println("📝 Model information:", info.Details.Family)

	systemContent := `You are an expert in computer programming.`

	userContent := prompt.StepByStep(`can you explain the following code:
	function generateUUID() {
		// Parts of the UUID
		const parts = [];
		const hex = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'];
	  
		// Random function (replace if needed)
		function random(max) {
		  return Math.floor(Math.random() * max);
		}
	  
		// Generate random hex digits
		for (let i = 0; i < 36; i++) {
		  if (i === 8 || i === 13 || i === 18 || i === 23) {
			parts.push('-');
		  } else {
			parts.push(hex[random(16)]);
		  }
		}
	  
		return parts.join('');
	  }
	  
	  const uuid = generateUUID();
	  console.log(uuid); // Example output: 12b3e4f5-1a2b-3c4d-5e6f-7890a1b2c3d4
	  
	`)

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.5,
		option.RepeatLastN: 2,
	})

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	_, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}
}
