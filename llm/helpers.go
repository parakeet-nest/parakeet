package llm


// SetOfMessages creates or extends a conversation with provided messages
// It can accept either single messages or slices of messages as variadic parameters
func SetOfMessages(messages ...interface{}) []Message {
	setOfMessages := []Message{}

	for _, msg := range messages {
		switch m := msg.(type) {
		case Message:
			setOfMessages = append(setOfMessages, m)
		case []Message:
			setOfMessages = append(setOfMessages, m...)
		}
	}

	return setOfMessages
}
