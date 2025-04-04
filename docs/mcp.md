# Model Context Protocol

!!! info "🚧 work in progress"

Integration of `github.com/mark3labs/mcp-go/mcp` and `github.com/mark3labs/mcp-go/client`

## STDIO transport


### Overview  

MCP (Model Context Protocol) with STDIO transport allows Parakeet to interact with external tools and services over standard input/output. This is particularly useful for integrating lightweight tools in a portable manner, leveraging command-line processes for execution.  

### Creating an MCP STDIO Client  

#### Initializing the Client  

The following example demonstrates how to initialize an MCP STDIO client:  

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/parakeet-nest/parakeet/mcp-stdio"
)

func main() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create MCP STDIO client
	mcpClient, err := mcpstdio.NewClient(ctx, "docker", []string{}, "run", "--rm", "-i", "mcp-curl")
	if err != nil {
		log.Fatalln("😡", err)
	}

	// Initialize the client
	_, err = mcpClient.Initialize()
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println("🚀 MCP STDIO client initialized successfully.")
}
```

#### Listing Available Tools  

To retrieve a list of tools available on the MCP STDIO server:  

```go
tools, err := mcpClient.ListTools()
if err != nil {
	log.Fatalln("😡 Failed to list tools:", err)
}

fmt.Println("📦 Available Tools:")
for _, tool := range tools {
	fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
}
```

#### Executing a Tool  

This example demonstrates how to call a tool that fetches webpage content:  

```go
messages := []llm.Message{
	{Role: "user", Content: "Fetch this page: https://raw.githubusercontent.com/parakeet-nest/parakeet/main/README.md"},
}

options := llm.SetOptions(map[string]interface{}{
	option.Temperature:   0.0,
	option.RepeatLastN:   2,
	option.RepeatPenalty: 2.0,
})

query := llm.Query{
	Model:    "qwen2.5:0.5b",
	Messages: messages,
	Tools:    tools,
	Options:  options,
	Format:   "json",
}

answer, err := completion.Chat("http://localhost:11434", query)
if err != nil {
	log.Fatalln("😡", err)
}

// Extract the first tool call
toolCall := answer.Message.ToolCalls[0]

fmt.Println("🛠️ Calling:", toolCall.Function.Name, toolCall.Function.Arguments)

// Execute tool
mcpResult, err := mcpClient.CallTool(toolCall.Function.Name, toolCall.Function.Arguments)
if err != nil {
	log.Fatalln("😡", err)
}

fmt.Println("🌍 Content:", mcpResult.Text)
```

### Handling Errors  

MCP STDIO has various error types for debugging:  

- `STDIOClientCreationError`  
- `STDIOClientInitializationError`  
- `STDIOGetToolsError`  
- `STDIOToolCallError`  
- `STDIOResultExtractionError`  
- `STDIOListResourcesError`  
- `STDIOReadResourceError`  

For example, error handling for initialization:  

```go
_, err = mcpClient.Initialize()
if err != nil {
	log.Fatalln("😡 STDIO Initialization Failed:", err)
}
```

### Running MCP STDIO with Docker  

If the MCP server is containerized, you can configure it to run inside a Docker container:  

```bash
docker build -t mcp-curl .
```

To start it:  

```bash
docker run --rm -i mcp-curl
```

Alternatively, using a JSON configuration file:  

```json
{
  "mcpServers": {
    "mcp-curl-with-docker" :{
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "mcp-curl"
      ]
    }
  }
}
```

### Conclusion  

MCP STDIO provides a lightweight and flexible transport mechanism for integrating tools within Parakeet. This is ideal for scenarios where external commands need to be executed without requiring a persistent network server.


## SSE transport

### Overview  

MCP (Model Context Protocol) with SSE (Server-Sent Events) allows Parakeet to interact with an event-driven architecture for fetching resources, executing tools, and processing prompts dynamically. This integration facilitates seamless communication with LLM-based tools while leveraging event streams for efficiency.  

### Using MCP SSE Client with Parakeet  

#### Initializing the Client  

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	mcpsse "github.com/parakeet-nest/parakeet/mcp-sse"
	"github.com/parakeet-nest/parakeet/gear"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡", err)
	}

	mcpSSEServerUrl := gear.GetEnvString("MCP_HOST", "http://0.0.0.0:5001")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create MCP client
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

	mcpClient.Close()
}
```

#### Executing MCP Tools  

##### Fetching Web Page Content  

This example demonstrates how to use the `fetch` tool to retrieve webpage content:  

```go
ollamaTools, err := mcpClient.ListTools()
if err != nil {
	log.Fatalln("😡", err)
}

messages := []llm.Message{
	{Role: "user", Content: "Fetch this page: https://raw.githubusercontent.com/parakeet-nest/parakeet/main/README.md"},
}

options := llm.SetOptions(map[string]interface{}{
	option.Temperature:   0.0,
	option.RepeatLastN:   2,
	option.RepeatPenalty: 2.0,
})

toolsQuery := llm.Query{
	Model:    modelWithToolsSupport,
	Messages: messages,
	Tools:    ollamaTools,
	Options:  options,
	Format:   "json",
}

answer, err := completion.Chat(ollamaUrl, toolsQuery)
if err != nil {
	log.Fatalln("😡", err)
}

toolCall := answer.Message.ToolCalls[0]

fmt.Println("🦙🛠️ 📣 calling:", toolCall.Function.Name, toolCall.Function.Arguments)

mcpResult, err := mcpClient.CallTool(toolCall.Function.Name, toolCall.Function.Arguments)
if err != nil {
	log.Fatalln("😡", err)
}

fmt.Println("🌍 Content:", mcpResult.Text)
```

#### Getting Resources  

The following example lists all available resources from the MCP SSE server:  

```go
resources, err := mcpClient.ListResources()
if err != nil {
	log.Fatalln("😡", err)
}

fmt.Println("📦 Available Resources:")
for _, resource := range resources {
	fmt.Printf("- Name: %s, URI: %s, MIME Type: %s\n", resource.Name, resource.URI, resource.MIMEType)
}
```

### Handling Errors  

MCP SSE has various error types for debugging:  

- `SSEClientCreationError`  
- `SSEClientStartError`  
- `SSEClientInitializationError`  
- `SSEGetToolsError`  
- `SSEToolCallError`  
- `SSEListResourcesError`  
- `SSEReadResourceError`  

For example, error handling for initialization:  

```go
result, err := mcpClient.Initialize()
if err != nil {
	log.Fatalln("😡 SSE Initialization Failed:", err)
}
```

### Conclusion  

MCP SSE provides a structured way to interact with streaming data and tools in an LLM-powered environment using Parakeet. By leveraging this integration, developers can seamlessly manage event-driven AI workflows.

!!! note
	👀 you will find a complete example in:

    - [examples/67-mcp](https://github.com/parakeet-nest/parakeet/tree/main/examples/67-mcp)
    - [examples/75-mcp-sse](https://github.com/parakeet-nest/parakeet/tree/main/examples/75-mcp-sse)
