package squawk

import "github.com/parakeet-nest/parakeet/llm"

// System adds a system message to the conversation context.
// System messages help define the behavior and context for the AI assistant.
//
// Parameters:
//   - message: A string containing the system instructions or context
//   - optionalParameters: Optional slice of strings where the first element,
//     if provided, is used as a label for the message
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with basic system message:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  System("You are a helpful AI assistant specializing in Go programming")
//
// Example with labeled system message:
//
//	squawk := New().
//	  Model("llama2:13b").
//	  Provider(provider.Ollama).
//	  System("You are a code review expert", "code-reviewer").
//	  User("Review this code...", "review-request")
//
// The label parameter is useful for:
//   - Message filtering
//   - Message replacement
//   - Context management
//   - Message tracking
func (s *Squawk) System(message string, optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "system", Content: message, Label: optionalParameters[0]},
		)

	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "system", Content: message})
	}

	return s
}

// User adds a user message to the conversation context.
//
// Parameters:
//   - message: A string containing the user's input or question
//   - optionalParameters: Optional slice of strings where the first element,
//     if provided, is used as a label for the message
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with basic user message:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  System("You are a helpful AI assistant").
//	  User("What is the capital of France?")
//
// Example with labeled user message:
//
//	squawk := New().
//	  Model("codellama:13b").
//	  Provider(provider.Ollama).
//	  System("You are a code review expert", "role").
//	  User("Review this Go code...", "code-review-1").
//	  User("What about this pattern?", "code-review-2")
//
// Usage with streaming response:
//
//	squawk := New().
//	  Model("llama2:13b").
//	  Provider(provider.Ollama).
//	  System("You are a helpful assistant").
//	  User("Explain quantum computing", "quantum-basics").
//	  ChatStream(func(answer llm.Answer, self *Squawk) error {
//	      fmt.Print(answer.Message.Content)
//	      return nil
//	  })
func (s *Squawk) User(message string, optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "user", Content: message, Label: optionalParameters[0]},
		)
	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "user", Content: message})
	}
	return s
}

// Assistant adds an assistant message to the conversation context.
// Assistant messages represent AI responses that you want to include manually
// in the conversation history, separate from automated responses.
//
// Parameters:
//   - message: A string containing the assistant's response or message
//   - optionalParameters: Optional slice of strings where the first element,
//     if provided, is used as a label for the message
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with basic assistant message:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  System("You are a Go expert").
//	  User("What is a goroutine?").
//	  Assistant("A goroutine is a lightweight thread of execution in Go")
//
// Example with labeled assistant message:
//
//	squawk := New().
//	  Model("codellama:13b").
//	  Provider(provider.Ollama).
//	  System("You are a code reviewer").
//	  User("Review this function", "request").
//	  Assistant("The function looks good but needs error handling", "review-1")
//
// Common use cases:
//   - Manual injection of AI responses into conversation
//   - Creating training examples
//   - Testing conversation flows
//   - Building conversation history
func (s *Squawk) Assistant(message string, optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "assistant", Content: message, Label: optionalParameters[0]},
		)
	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "assistant", Content: message})
	}
	return s
}

// AddSetOfMessages creates or extends a conversation with provided messages.
// It can accept either single messages or slices of messages as variadic parameters.
//
// Parameters:
//   - messages: A variadic parameter that accepts either llm.Message or []llm.Message
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with single messages:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  AddSetOfMessages(
//	    llm.Message{Role: "system", Content: "You are a Go expert"},
//	    llm.Message{Role: "user", Content: "Explain channels"},
//	  )
//
// Example with message slices:
//
//	previousMessages := []llm.Message{
//	  {Role: "system", Content: "You are a Go expert"},
//	  {Role: "user", Content: "What is concurrency?"},
//	  {Role: "assistant", Content: "Concurrency is..."},
//	}
//
//	squawk := New().
//	  Model("codellama:13b").
//	  Provider(provider.Ollama).
//	  AddSetOfMessages(previousMessages).
//	  User("How does this relate to goroutines?")
//
// Common use cases:
//   - Restoring previous conversations
//   - Batch adding multiple messages
//   - Combining message sets from different sources
//   - Building complex conversation contexts
func (s *Squawk) AddSetOfMessages(messages ...interface{}) *Squawk {
	// AddSetOfMessages creates or extends a conversation with provided messages
	// It can accept either single messages or slices of messages as variadic parameters
	for _, msg := range messages {
		switch m := msg.(type) {
		case llm.Message:
			s.setOfMessages = append(s.setOfMessages, m)
		case []llm.Message:
			s.setOfMessages = append(s.setOfMessages, m...)
		}
	}
	return s
}

// Messages returns the current conversation context as a slice of messages.
// This method provides access to all messages in the conversation, including
// system instructions, user inputs, and assistant responses.
//
// Returns:
//   - []llm.Message: A slice containing all messages in the conversation
//
// Example usage:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  System("You are a Go expert").
//	  User("What is a slice?")
//
//	messages := squawk.Messages()
//	for _, msg := range messages {
//	  fmt.Printf("Role: %s, Content: %s\n", msg.Role, msg.Content)
//	}
//
// Common use cases:
//   - Inspecting the current conversation state
//   - Saving conversations for later use
//   - Debugging conversation flow
//   - Processing message history
func (s *Squawk) Messages() []llm.Message {
	return s.setOfMessages
}

// SaveAssistantAnswer adds the last model response to the conversation history as an assistant message.
// This method is similar to SaveAnswer and is used for maintaining conversation context by
// preserving model responses.
//
// Parameters:
//   - optionalParameters: Optional slice of strings where the first element,
//     if provided, is used as a label for the saved message
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Chat completion:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  System("You are a Go expert").
//	  User("What is error handling in Go?").
//	  Chat(func(answer llm.Answer, self *Squawk, err error) {
//	      if err != nil {
//	          return
//	      }
//	  }).
//	  SaveAssistantAnswer("error-explanation").  // Save with label
//	  User("Can you provide an example?")
//
// Example with labeled streaming response:
//
//	squawk := New().
//	  Model("codellama:13b").
//	  Provider(provider.Ollama).
//	  System("You are a code reviewer").
//	  User("Review this function").
//	  ChatStream(func(answer llm.Answer, self *Squawk) error {
//	      if !answer.Done {
//	          fmt.Print(answer.Message.Content)
//	          return nil
//	      }
//	      return nil
//	  }).
//	  SaveAssistantAnswer("code-review-1")  // Save complete response
func (s *Squawk) SaveAssistantAnswer(optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "assistant", Content: s.lastAnswer.Message.Content, Label: optionalParameters[0]},
		)
	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "assistant", Content: s.lastAnswer.Message.Content})
	}
	return s

}

// LastAnswer gets or sets the most recent answer from the language model.
// When called without parameters, it returns the last answer. When called with
// an answer parameter, it sets the last answer and returns it.
//
// Parameters:
//   - optionalAnswer: Optional llm.Answer to set as the last answer
//
// Returns:
//   - llm.Answer: The current last answer or the newly set answer
//
// Example retrieving the last answer:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  System("You are a Go expert").
//	  User("What is a channel?").
//	  Chat(func(answer llm.Answer, self *Squawk, err error) {
//	      if err != nil {
//	          return
//	      }
//	  })
//
//	lastAnswer := squawk.LastAnswer()
//	fmt.Println("Last response:", lastAnswer.Message.Content)
//
// Example setting a custom answer:
//
//	customAnswer := llm.Answer{
//	  Message: llm.Message{
//	    Role: "assistant",
//	    Content: "Custom response",
//	  },
//	}
//	squawk.LastAnswer(customAnswer)
func (s *Squawk) LastAnswer(optionalAnswer ...llm.Answer) llm.Answer {
	if len(optionalAnswer) > 0 {
		s.lastAnswer = optionalAnswer[0]
	}
	return s.lastAnswer
}

// RemoveMessageByLabel removes messages with the specified label from the conversation context.
// This method is useful for managing conversation flow by selectively removing messages.
//
// Parameters:
//   - label: A string representing the label of messages to remove
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//
//	squawk := New().
//	  Model("mistral:latest").
//	  Provider(provider.Ollama).
//	  System("You are a Go expert", "role").
//	  User("What is a channel?", "question-1").
//	  Chat(func(answer llm.Answer, self *Squawk, err error) {
//	      if err != nil {
//	          return
//	      }
//	  }).
//	  SaveAssistantAnswer("answer-1").
//	  RemoveMessageByLabel("question-1").  // Remove specific question
//	  User("How do goroutines work?", "question-2")
//
// Common use cases:
//   - Cleaning up conversation history
//   - Managing context window size
//   - Removing outdated or irrelevant messages
//   - Implementing conversation history pruning
func (s *Squawk) RemoveMessageByLabel(label string) *Squawk {
	// Remove messages by label
	var newMessages []llm.Message
	for _, message := range s.setOfMessages {
		if message.Label != label {
			newMessages = append(newMessages, message)
		}
	}
	s.setOfMessages = newMessages
	return s
}
