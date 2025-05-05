package squawk

import (
	"encoding/json"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

// Schema sets a schema for structured output from the language model.
// This method configures how the model should format its response,
// particularly useful when you need responses in specific JSON structures.
//
// Parameters:
//   - schema: A map[string]any defining the expected structure of the output
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  Schema(map[string]any{
//	      "type": "object",
//	      "properties": map[string]any{
//	          "name": map[string]any{
//	              "type": "string",
//	              "description": "The name of the person",
//	          },
//	          "age": map[string]any{
//	              "type": "integer",
//	              "description": "The age of the person",
//	          },
//	      },
//	  }).
//	  User("Extract name and age from: John is 25 years old").
//	  StructuredOutput(func(answer llm.Answer, self *Squawk, err error) {
//	      fmt.Printf("%+v\n", answer.StructuredOutput)
//	  })
func (s *Squawk) Schema(schema map[string]any) *Squawk {
	s.schema = schema
	return s
}

func (s *Squawk) SchemaJSONS(schemaJSONString string) *Squawk {
	// transform the JSON string into a map with go methods
	s.schema = make(map[string]any)
	err := json.Unmarshal([]byte(schemaJSONString), &s.schema)
	if err != nil {
		s.lastError = err
	}
	return s
}

// structuredOutputExec executes a chat completion request and returns a structured response.
// This is an internal method used by StructuredOutput to generate responses in a specific format.
//
// The method performs the following steps:
// 1. Constructs a query with the current model, messages, options, and schema
// 2. Sets Raw to false to ensure proper JSON parsing
// 3. Calls the completion.Chat function with the configured provider
// 4. Updates the lastAnswer field with the structured response
//
// Returns:
//   - llm.Answer: The model's response containing the structured output
//   - error: Any error that occurred during the completion request
//
// Example usage with Ollama provider:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  Schema(map[string]any{
//	      "type": "object",
//	      "properties": map[string]any{
//	          "summary": map[string]any{
//	              "type": "string",
//	              "description": "Brief summary of the code",
//	          },
//	          "issues": map[string]any{
//	              "type": "array",
//	              "items": map[string]any{
//	                  "type": "string",
//	              },
//	          },
//	      },
//	  }).
//	  User("Review this code: func main() {}").
//	  StructuredOutput(func(answer llm.Answer, self *Squawk, err error) {
//	      if err != nil {
//	          fmt.Println("Error:", err)
//	          return
//	      }
//	      fmt.Printf("%+v\n", answer.StructuredOutput)
//	  })
func (s *Squawk) structuredOutputExec() (llm.Answer, error) {
	query := llm.Query{
		Model:    s.chatModel,
		Messages: s.setOfMessages,
		Options:  s.options,
		Format:   s.schema,
		Raw:      false,
	}
	answer, err := completion.Chat(s.apiUrl, query, s.provider, s.openAPIKey)
	if err != nil {
		return llm.Answer{}, err
	}
	s.lastAnswer = answer
	return answer, nil
}

// StructuredOutput processes a structured chat completion request and handles the response
// through a callback function. This method is used to get responses in a specific JSON format
// defined by the schema set using the Schema method.
//
// Parameters:
//   - callBack: A function that receives the structured response (llm.Answer),
//     a pointer to the Squawk instance, and any error that occurred
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  Schema(map[string]any{
//	      "type": "object",
//	      "properties": map[string]any{
//	          "functionName": map[string]any{
//	              "type": "string",
//	              "description": "Name of the function",
//	          },
//	          "complexity": map[string]any{
//	              "type": "string",
//	              "enum": ["O(1)", "O(n)", "O(n^2)"],
//	              "description": "Time complexity of the function",
//	          },
//	      },
//	  }).
//	  User("Analyze this function: func BubbleSort(arr []int) []int { ... }").
//	  StructuredOutput(func(answer llm.Answer, self *Squawk, err error) {
//	      if err != nil {
//	          fmt.Println("Error:", err)
//	          return
//	      }
//	      result := answer.StructuredOutput
//	      fmt.Printf("Function: %s\nComplexity: %s\n",
//	          result["functionName"],
//	          result["complexity"])
//	  })
//
// TODO: rewrite the examples
func (s *Squawk) StructuredOutput(callBack func(answer llm.Answer, self *Squawk, err error)) *Squawk {
	answer, err := s.structuredOutputExec()
	if err != nil {
		callBack(answer, s, err)
		s.lastError = err
		return s
	}
	callBack(answer, s, nil)
	return s
}
