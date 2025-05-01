package squawk

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

// Tools sets the list of available tools for function calling capabilities.
// This method configures which tools the language model can use during conversations.
//
// Parameters:
//   - toolsList: A slice of llm.Tool representing the available tools
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//   squawk := New().
//     Model("mistral:latest").
//     Provider(provider.Ollama).
//     Tools([]llm.Tool{
//       {
//         Name: "get_weather",
//         Description: "Get current weather for a location",
//         Parameters: map[string]interface{}{
//           "location": map[string]string{
//             "type": "string",
//             "description": "City name",
//           },
//         },
//       },
//     }).
//     User("What's the weather in Paris?").
//     FunctionCalling(func(answer llm.Answer, s *Squawk, err error) {
//       if err != nil {
//         fmt.Println("Error:", err)
//         return
//       }
//       /* Handle tool calls */
//       for _, call := range s.ToolCalls() {
//         fmt.Printf("Tool called: %s with args: %v\n", 
//           call.Name, 
//           call.Arguments,
//         )
//       }
//     })
//
// Common use cases:
// - Adding custom functionality to conversations
// - Enabling API integrations
// - Creating interactive chatbots
// - Building tool-augmented assistants
func (s *Squawk) Tools(toolsList []llm.Tool) *Squawk {
	s.tools = toolsList
	return s
}

// functionCallingExec executes a function calling request with the configured tools
// and returns the model's response. This is an internal method used by FunctionCalling
// to process tool-enabled interactions.
//
// Returns:
//   - llm.Answer: The model's response including any tool calls
//   - error: Any error that occurred during the request
//
// The method performs these steps:
// 1. Creates a query with the current conversation context and tools
// 2. Sends the query to the language model
// 3. Updates the instance with the response and tool calls
// 4. Returns the answer or any error that occurred
//
// Example usage through FunctionCalling:
//   squawk := New().
//     Model("mistral:latest").
//     Provider(provider.Ollama).
//     Tools([]llm.Tool{
//       {
//         Name: "calculate",
//         Description: "Perform basic math operations",
//         Parameters: map[string]interface{}{
//           "operation": map[string]string{
//             "type": "string",
//             "description": "Math operation to perform",
//           },
//         },
//       },
//     }).
//     User("What is 2 + 2?").
//     FunctionCalling(func(answer llm.Answer, s *Squawk, err error) {
//       if err != nil {
//         fmt.Println("Error:", err)
//         return
//       }
//       /* Process tool calls */
//       for _, call := range s.ToolCalls() {
//         fmt.Printf("Tool: %s, Args: %v\n", call.Name, call.Arguments)
//       }
//     })
//
// Internal use:
// - Called by FunctionCalling method
// - Updates lastAnswer and toolCalls fields
// - Handles error propagation
func (s *Squawk) functionCallingExec() (llm.Answer, error) {
	query := llm.Query{
		Model:    s.chatModel,
		Messages: s.setOfMessages,
		Tools:    s.tools,
		Options:  s.options,
	}

	answer, err := completion.Chat(s.apiUrl, query, s.provider, s.openAPIKey)
	if err != nil {
		return llm.Answer{}, err
	}
	s.lastAnswer = answer
	s.toolCalls = answer.Message.ToolCalls.Tools()
	return answer, nil
}

// FunctionCalling executes a function calling request and processes the results through
// a callback function. This method enables tool-augmented conversations where the model
// can call predefined functions.
//
// Parameters:
//   - callBack: A function that receives the model's answer, the Squawk instance,
//     and any error that occurred during processing
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//   squawk := New().
//     Model("mistral:latest").
//     Provider(provider.Ollama).
//     Tools([]llm.Tool{
//       {
//         Name: "search_database",
//         Description: "Search for records in database",
//         Parameters: map[string]interface{}{
//           "query": map[string]string{
//             "type": "string",
//             "description": "Search query",
//           },
//         },
//       },
//     }).
//     User("Find all records about Go programming").
//     FunctionCalling(func(answer llm.Answer, s *Squawk, err error) {
//       if err != nil {
//         fmt.Println("Error:", err)
//         return
//       }
//       
//       /* Process tool calls */
//       for _, call := range s.ToolCalls() {
//         fmt.Printf("Function: %s\nArguments: %v\n", 
//           call.Name, 
//           call.Arguments,
//         )
//       }
//     })
//
// Error handling:
// - Callback receives any error that occurred during execution
// - Errors are stored in lastError field
// - Processing continues with error passed to callback
//
// Common use cases:
// - Implementing custom tool handlers
// - API integrations
// - Database queries
// - External service calls
func (s *Squawk) FunctionCalling(callBack func(answer llm.Answer, self *Squawk, err error)) *Squawk {

	answer, err := s.functionCallingExec()
	if err != nil {
		callBack(answer, s, err)
		s.lastError = err
		return s
	}
	callBack(answer, s, nil)
	return s
}

// ToolCalls returns the tool calls from the most recent function calling interaction.
// This method provides access to the tool calls made by the language model during
// the last FunctionCalling execution.
//
// Returns:
//   - []llm.ToolCall: A slice of tool calls containing the function names and arguments
//
// Example with Ollama provider:
//   squawk := New().
//     Model("mistral:latest").
//     Provider(provider.Ollama).
//     Tools([]llm.Tool{
//       {
//         Name: "fetch_data",
//         Description: "Fetch data from a source",
//         Parameters: map[string]interface{}{
//           "source": map[string]string{
//             "type": "string",
//             "description": "Data source identifier",
//           },
//         },
//       },
//     }).
//     User("Get data from source A").
//     FunctionCalling(func(answer llm.Answer, s *Squawk, err error) {
//       if err != nil {
//         return
//       }
//       /* Access tool calls */
//       calls := s.ToolCalls()
//       for _, call := range calls {
//         fmt.Printf("Called: %s with %v\n", call.Name, call.Arguments)
//       }
//     })
//
// Common use cases:
// - Processing tool call results
// - Implementing tool handlers
// - Debugging function calls
// - Auditing model interactions
func (s *Squawk) ToolCalls() []llm.ToolCall {
	return s.toolCalls
}
