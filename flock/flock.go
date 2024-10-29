package flock

import (
	"fmt"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

type AgentFunction func(interface{}) (interface{}, error)

// InstructionFunc represents a function that takes context variables and returns a string
type InstructionFunc func(map[string]interface{}) string

// Agent represents an AI agent with its configuration
type Agent struct {
	Name         string          `json:"name"`
	Model        string          `json:"model"`
	OllamaUrl    string          `json:"ollama_url"`
	Instructions interface{}     `json:"instructions"` // Can be string or function returning string
	Functions    map[string]AgentFunction `json:"functions"`
	//ToolChoice   *string         `json:"tool_choice,omitempty"`
	//ParallelToolCalls bool            `json:"parallel_tool_calls"`
	Options llm.Options
}

// SetInstructions sets the instructions for the agent
func (a *Agent) SetInstructions(instructions interface{}) error {
	switch v := instructions.(type) {
	case string:
		a.Instructions = v
	case func() string:
		// Convert simple function to InstructionFunc
		a.Instructions = InstructionFunc(func(map[string]interface{}) string {
			return v()
		})
	case func(map[string]interface{}) string:
		a.Instructions = InstructionFunc(v)
	default:
		return fmt.Errorf("invalid instruction type: must be string, func() string, or func(map[string]interface{}) string")
	}
	return nil
}

// GetInstructions returns the current instructions as a string, using the provided context variables
func (a *Agent) GetInstructions(contextVars map[string]interface{}) string {
	switch v := a.Instructions.(type) {
	case string:
		return v
	case InstructionFunc:
		return v(contextVars)
	case func(map[string]interface{}) string:
		return v(contextVars)
	default:
		return "Invalid instruction type"
	}
}

// Response represents the response from running an agent
type Response struct {
	Messages         []llm.Message          `json:"messages"`
	Agent            Agent                  `json:"agent,omitempty"`
	ContextVariables map[string]interface{} `json:"context_variables"`
}

// Function to get the last message of Response.Messages
func (r *Response) GetLastMessage() llm.Message {
	return r.Messages[len(r.Messages)-1]
}

// Orchestrator represents an API client for running agents
type Orchestrator struct {
}

// Run executes the agent with the given messages and context variables
func (c *Orchestrator) Run(agent Agent, messages []llm.Message, contextVars map[string]interface{}) (Response, error) {
	// Get instructions with context variables
	instructions := agent.GetInstructions(contextVars)
	agentMessages := []llm.Message{}

	agentMessages = append(agentMessages, llm.Message{
		Role:    "system",
		Content: instructions,
	})
	agentMessages = append(agentMessages, messages...)

	queryChat := llm.Query{
		Model:    agent.Model,
		Messages: agentMessages,
		Options:  agent.Options,
		Stream:   false, 
		//Format:  "json",
	}

	// Answer the question
	answer, err := completion.Chat(agent.OllamaUrl, queryChat)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		Messages:         agentMessages,
		Agent:            agent,
		ContextVariables: contextVars,
	}

	// Add a simple response message (for demonstration)
	response.Messages = append(response.Messages, llm.Message{
		//Role:    "assistant",
		Role:    answer.Message.Role,
		Content: answer.Message.Content,
	})

	return response, nil
}

// TODO: Some refactoring ...
func (c *Orchestrator) RunStream(agent Agent, messages []llm.Message, contextVars map[string]interface{}, onChunk func(llm.Answer) error) (Response, error) {
	// Get instructions with context variables
	instructions := agent.GetInstructions(contextVars)
	agentMessages := []llm.Message{}

	agentMessages = append(agentMessages, llm.Message{
		Role:    "system",
		Content: instructions,
	})
	agentMessages = append(agentMessages, messages...)

	queryChat := llm.Query{
		Model:    agent.Model,
		Messages: agentMessages,
		Options:  agent.Options,
		Stream:   true, 
		//Format:  "json",
	}

	// Answer the question
	answer, err := completion.ChatStream(agent.OllamaUrl, queryChat, onChunk)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		Messages:         agentMessages,
		Agent:            agent,
		ContextVariables: contextVars,
	}

	// Add a simple response message (for demonstration)
	response.Messages = append(response.Messages, llm.Message{
		//Role:    "assistant",
		Role:    answer.Message.Role,
		Content: answer.Message.Content,
	})

	return response, nil
}


func (c *Orchestrator) RunWithTools(agent Agent, messages []llm.Message, contextVars map[string]interface{}, tools []llm.Tool, execute bool) (Response, error) {
	// Get instructions with context variables
	instructions := agent.GetInstructions(contextVars)
	agentMessages := []llm.Message{}

	agentMessages = append(agentMessages, llm.Message{
		Role:    "system",
		Content: instructions,
	})
	agentMessages = append(agentMessages, messages...)

	queryChat := llm.Query{
		Model:    agent.Model,
		Messages: agentMessages,
		Options:  agent.Options,
		Stream:   false, 
		Format:  "json",
		Tools:  tools,
	}

	// Answer the question
	answer, err := completion.Chat(agent.OllamaUrl, queryChat)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		Messages:         agentMessages,
		Agent:            agent,
		ContextVariables: contextVars,
	}

	// Execute the tool calls
	var toolsCalls = []struct{Function llm.FunctionTool; Result interface{}; Error error}{}

	for _, toolCall := range answer.Message.ToolCalls {
		newToolCall := struct{Function llm.FunctionTool; Result interface{}; Error error}{}
		newToolCall.Function = toolCall.Function
		
		
		if execute  {
			result, err := agent.Functions[toolCall.Function.Name](toolCall.Function.Arguments)
			newToolCall.Result = result
			newToolCall.Error = err
		}

		
		toolsCalls = append(toolsCalls, newToolCall)
	}
	// Add a simple response message (for demonstration)
	response.Messages = append(response.Messages, llm.Message{
		//Role:    "assistant",
		Role:    answer.Message.Role,
		Content: answer.Message.Content,
		ToolCalls: toolsCalls,
	})


	return response, nil
}