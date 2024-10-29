# Flock Agents

Flock is a Parakeet package for creating and managing AI agents using the Ollama backend. It provides a simple way to create conversational agents, orchestrate interactions between them, and implement function calling capabilities.

## Basic Concepts

### Agent
An Agent represents an AI entity with specific configurations, instructions, and capabilities. Each agent can have:
- Name
- Model (Ollama model to use)
- Instructions (static or dynamic)
- Functions (for tool calling)
- Options (model parameters)

### Orchestrator
The Orchestrator manages agent execution and interactions. It provides methods to:
- Run single agent interactions
- Stream responses
- Execute function calls

## Creating Agents

### Basic Agent

```go
agent := flock.Agent{
    Name: "Bob",
    Model: "qwen2.5:3b",
    OllamaUrl: "http://localhost:11434",
    Options: llm.SetOptions(map[string]interface{}{
        option.Temperature: 0.7,
        option.TopK: 40,
        option.TopP: 0.9,
    }),
}

// Setting static instructions
agent.SetInstructions("Help the user with their queries.")

// Setting dynamic instructions with context
agent.SetInstructions(func(contextVars map[string]interface{}) string {
    userName := contextVars["userName"].(string)
    return fmt.Sprintf("Help %s with their queries.", userName)
})
```

## Using the Orchestrator

### Basic Interaction

```go
orchestrator := flock.Orchestrator{}

response, err := orchestrator.Run(
    agent,
    []llm.Message{
        {Role: "user", Content: "Hello, what's the best pizza?"},
    },
    map[string]interface{}{"userName": "Sam"},
)

if err != nil {
    log.Fatal(err)
}

fmt.Println(response.Messages[len(response.Messages)-1].Content)

// Or:
fmt.Println(response.GetLastMessage().Content)

```

### Streaming Responses

```go
response, err := orchestrator.RunStream(
    agent,
    messages,
    contextVars,
    func(answer llm.Answer) error {
        fmt.Print(answer.Message.Content)
        return nil
    },
)
```

### Function Calling

```go
options := llm.SetOptions(map[string]interface{}{
    option.Temperature: 0.0,
})

// Create agent with function capability
calculator := flock.Agent{
    Name:    "Calculator",
    Model:   "allenporter/xlam:1b",
    OllamaUrl: "http://localhost:11434",
    Options: options,
    Functions: map[string]flock.AgentFunction{
        "multiply": func(args interface{}) (interface{}, error) {
            argsMap := args.(map[string]interface{})
            return argsMap["a"].(float64) * argsMap["b"].(float64), nil
        },
    },
}

// Define tool
tools := []llm.Tool{
    {
        Type: "function",
        Function: llm.Function{
            Name:        "multiply",
            Description: "Make a multiplication of the two given numbers",
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

// Run calculation
orchestrator := flock.Orchestrator{}

response, _ := orchestrator.RunWithTools(
    calculator,
    []llm.Message{{Role: "user", Content: "multiply 5 and 3"}},
    map[string]interface{}{},
    tools,
    true, // execute the tool(s)    
)

// Access result
fmt.Println("Result:", response.GetLastMessage().ToolCalls[0].Result)

```

## Examples

### Basic Agent

!!! note
    Look at this sample:

    - [examples/61-agents](https://github.com/parakeet-nest/parakeet/tree/main/examples/61-agents)

### Conversational Agents

!!! note
	Look at these samples:

    - [examples/62-agents-chat](https://github.com/parakeet-nest/parakeet/tree/main/examples/62-agents-chat)
    - [examples/63-agents-chat-stream](https://github.com/parakeet-nest/parakeet/tree/main/examples/63-agents-chat-stream)


### Function Calling Example

!!! note
    Look at this sample:

    - [examples/64-agents-with-tools](https://github.com/parakeet-nest/parakeet/tree/main/examples/64-agents-with-tools)
