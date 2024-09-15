package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"

)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "gemma2:2b"

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

	systemContent := `As an expert in Go and Docker, generate a Dockerfile and Docker Compose file for a typical Go project.
        
	Regarding the Dockerfile:
	1. Use an official Go base image (I use go version '1.22.1')
	2. Use the /app working directory
	3. Copy the project source files
	4. Compile the Go application
	5. Create a lightweight final image for execution 
		- Use the /app working directory
		- Copy the binary in the /app directory
		- Copy the static assets (and directory structure) from /public to /app/public
	6. Start the Golang application
	
	Regarding the Compose file:
	The Golang application is using a Redis database
	1. Add a redis service with "redis-server" as name to the compose file.
		- The redis service is listening on the default port
	2. Add a webapp service 
		- The webapp service uses an environment variable (REDIS_URL)  to connect to the "redis-server" service
		- To set the value of REDIS_URL, use only the name of the redis service and the default redis port
		- The webapp service is listening on the 8080 HTTP port
		- the webapp service depends on the redis-server service
			
	Ensure the Dockerfile and the Compose file are well-commented and follows best practices. 
	Briefly explain each step after the Dockerfile and the Compose file.
	`

	userContent := `
	Here's the project structure:
	.
	├── git.sh
	├── go.mod
	├── go.sum
	├── init.sh
	├── load-data
	│  ├── bulk_loading.sh
	│  └── data.txt
	├── main.go
	├── public
	│  ├── components
	│  │  ├── App.js
	│  │  └── Title.js
	│  ├── css
	│  │  ├── install-pico.md
	│  │  └── pico.min.css
	│  ├── index.html
	│  ├── info.txt
	│  └── js
	│     ├── install-preact.md
	│     ├── preact-htm.js
	│     └── update.js
	└── README.md
	`

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 2.0,
		option.TopK: 10,
		option.TopP: 0.5,
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
