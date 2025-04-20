package llm


// Conversation creates or extends a conversation with provided messages
// It can accept either single messages or slices of messages as variadic parameters
func Conversation(messages ...interface{}) []Message {
	conversation := []Message{}

	for _, msg := range messages {
		switch m := msg.(type) {
		case Message:
			conversation = append(conversation, m)
		case []Message:
			conversation = append(conversation, m...)
		}
	}

	return conversation
}
