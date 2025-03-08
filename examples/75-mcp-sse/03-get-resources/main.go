package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/gear"
	mcpsse "github.com/parakeet-nest/parakeet/mcp-sse"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡", err)
	}

	//ollamaUrl := gear.GetEnvString("OLLAMA_HOST", "http://localhost:11434")
	//modelWithToolsSupport := gear.GetEnvString("LLM_WITH_TOOLS_SUPPORT", "qwen2.5:0.5b")
	//chatModel := gear.GetEnvString("LLM_CHAT", "qwen2.5:0.5b")
	mcpSSEServerUrl := gear.GetEnvString("MCP_HOST", "http://0.0.0.0:5001")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new mcp client
	mcpClient, err := mcpsse.NewClient(ctx, mcpSSEServerUrl)

	if err != nil {
		log.Fatalln("😡", err)
	}
	err = mcpClient.Start()
	if err != nil {
		log.Fatalln("😡", err)
	}
	result, err := mcpClient.Initialize()
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println("🚀 Initialized with server:", result.ServerInfo.Name, result.ServerInfo.Version)

	resources, err := mcpClient.ListResources()
	if err != nil {
		log.Fatalln("😡", err)
	}
	fmt.Println("📦", resources)

	// Print the list of available resources
	fmt.Println("🌍 Available Static Resources:")
	for _, resource := range resources {
		fmt.Printf("- Name: %s, URI: %s, MIME Type: %s\n",
			resource.Name, resource.URI, resource.MIMEType)
	}


	resourceResult, err := mcpClient.ReadResource("system://instructions")
	if err != nil {
		log.Fatalf("😡 Failed to read resource: %v", err)
	}
	fmt.Println("📖", resourceResult.Contents)

	for _, content := range resourceResult.Contents {
		fmt.Println("- 📝 [", content["kind"],"]:", content["text"])

	}
	
	mcpClient.Close()
	fmt.Println("👋 Bye!")
}
