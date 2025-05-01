package history

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/parakeet-nest/parakeet/llm"
)

type MemoryMessages struct {
	Messages map[string]llm.MessageRecord
	Keys     []string // Add this field to maintain order

}

// Add a new constructor
func NewMemoryMessages() *MemoryMessages {
	return &MemoryMessages{
		Messages: make(map[string]llm.MessageRecord),
		Keys:     make([]string, 0),
	}
}

func (m *MemoryMessages) Get(id string) (llm.MessageRecord, error) {
	return m.Messages[id], nil
}

func (m *MemoryMessages) GetMessage(id string) (llm.Message, error) {
	messageRecord, err := m.Get(id)
	if err != nil {
		return llm.Message{}, err
	}
	return llm.Message{
		Role:    messageRecord.Role,
		Content: messageRecord.Content,
	}, nil
}

func (m *MemoryMessages) GetAll() ([]llm.MessageRecord, error) {
	var records []llm.MessageRecord
	for _, messageRecord := range m.Messages {
		records = append(records, messageRecord)
	}
	return records, nil
}

// GetLastNMessages returns the last n messages from the Messages list.
// If n <= 0, returns an error.
// If n > total messages, returns all messages.
func (m *MemoryMessages) GetLastNMessages(n int) ([]llm.Message, error) {
    if n <= 0 {
        return nil, fmt.Errorf("n must be positive, got %d", n)
    }

    totalMessages := len(m.Keys)
    if totalMessages == 0 {
        return []llm.Message{}, nil
    }

    // Calculate starting index
    startIndex := totalMessages - n
    if startIndex < 0 {
        startIndex = 0
    }

    // Get last n messages using Keys slice for order
    var messages []llm.Message
    for _, key := range m.Keys[startIndex:] {
        messages = append(messages, llm.Message{
            Role:    m.Messages[key].Role,
            Content: m.Messages[key].Content,
        })
    }

    return messages, nil
}

// TODO: implement the filter by pattern()
func (m *MemoryMessages) GetAllMessages(patterns ...string) ([]llm.Message, error) {
	var messages []llm.Message

	for _, key := range m.Keys {
		messages = append(messages, llm.Message{
			Role:    m.Messages[key].Role,
			Content: m.Messages[key].Content,
		})
	}

	return messages, nil
}

// TODO: implement the filter by pattern()
func (m *MemoryMessages) GetAllMessagesOfSession(sessionId string, patterns ...string) ([]llm.Message, error) {

	keys := make([]string, 0, len(m.Messages))
	for k := range m.Messages {
		keys = append(keys, k)
	}

	// Sort the keys
	sort.Strings(keys)

	// Create ordered slice of messages
	var messages []llm.Message

	for _, key := range keys {
		if m.Messages[key].SessionId == sessionId {
			messages = append(messages, llm.Message{
				Role:    m.Messages[key].Role,
				Content: m.Messages[key].Content,
			})
		}
	}

	return messages, nil

}

// TODO: private or public?
func (m *MemoryMessages) Save(messageRecord llm.MessageRecord) (llm.MessageRecord, error) {
	if _, exists := m.Messages[messageRecord.Id]; !exists {
		m.Keys = append(m.Keys, messageRecord.Id)
	}
	m.Messages[messageRecord.Id] = messageRecord
	return messageRecord, nil
}

func (m *MemoryMessages) SaveMessage(id string, message llm.Message) (llm.MessageRecord, error) {
	if id == "" {
		// generate a unique for the message
		id = uuid.New().String()
	}
	messageRecord := llm.MessageRecord{
		Id:      id,
		Role:    message.Role,
		Content: message.Content,
	}
	return m.Save(messageRecord)
}

func (m *MemoryMessages) SaveMessageWithSession(sessionId, messageId string, message llm.Message) (llm.MessageRecord, error) {
	if messageId == "" {
		// generate a unique for the message
		messageId = uuid.New().String()
	}
	messageRecord := llm.MessageRecord{
		Id:        messageId,
		Role:      message.Role,
		Content:   message.Content,
		SessionId: sessionId,
	}
	return m.Save(messageRecord)
}

func (m *MemoryMessages) RemoveMessage(id string) error {
	delete(m.Messages, id)
	// Remove from Keys
	for i, key := range m.Keys {
		if key == id {
			m.Keys = append(m.Keys[:i], m.Keys[i+1:]...)
			break
		}
	}
	return nil
}

// TODO: to test
func (m *MemoryMessages) RemoveAllMessages() error {
	m.Messages = make(map[string]llm.MessageRecord)
	m.Keys = make([]string, 0)
	return nil
}

// TODO: to test
func (m *MemoryMessages) RemoveAllMessagesOfSession(sessionId string) error {
	for id, messageRecord := range m.Messages {
		if messageRecord.SessionId == sessionId {
			delete(m.Messages, id)
		}
	}
	return nil
}

func (m *MemoryMessages) RemoveTopMessageOfSession(sessionId string) error {
	if len(m.Keys) == 0 {
		return nil // No messages to remove
	}

	// Find the first message of the session
	var topMessageId string
	for _, key := range m.Keys {
		if m.Messages[key].SessionId == sessionId {
			topMessageId = key
			break
		}
	}

	// No messages found for this session
	if topMessageId == "" {
		return nil
	}

	// Remove from Messages map
	delete(m.Messages, topMessageId)

	// Remove from Keys slice
	for i, key := range m.Keys {
		if key == topMessageId {
			m.Keys = append(m.Keys[:i], m.Keys[i+1:]...)
			break
		}
	}

	return nil
}

// RemoveTopMessage removes the oldest message from the Messages map and updates the Keys slice
func (m *MemoryMessages) RemoveTopMessage() error {
	if len(m.Keys) == 0 {
		return nil // No messages to remove
	}

	// Get the oldest message ID (first in Keys slice)
	topMessageId := m.Keys[0]

	// Remove from Messages map
	delete(m.Messages, topMessageId)

	// Remove from Keys slice
	m.Keys = m.Keys[1:] // Remove first element

	return nil
}

// KeepLastN removes all messages except the last n messages
func (m *MemoryMessages) KeepLastN(n int) error {
	if n < 0 {
		return fmt.Errorf("n must be positive, got %d", n)
	}

	if len(m.Keys) <= n {
		return nil // Nothing to remove
	}

	// Calculate how many messages to remove
	removeCount := len(m.Keys) - n

	// Remove oldest messages first
	for i := 0; i < removeCount; i++ {
		// Get oldest message ID
		oldestId := m.Keys[0]

		// Remove from Messages map
		delete(m.Messages, oldestId)

		// Remove from Keys slice
		m.Keys = m.Keys[1:]
	}

	return nil
}

// KeepLastNOfSession removes all messages of a session except the last n messages
// KeepLastNOfSession removes all messages of a session except the last n messages
func (m *MemoryMessages) KeepLastNOfSession(sessionId string, n int) error {
	if n < 0 {
		return fmt.Errorf("n must be positive, got %d", n)
	}

	// Get all keys for the session in order
	var sessionKeys []string
	for _, key := range m.Keys {
		if m.Messages[key].SessionId == sessionId {
			sessionKeys = append(sessionKeys, key)
		}
	}

	if len(sessionKeys) <= n {
		return nil // Nothing to remove
	}

	// Calculate how many messages to remove
	removeCount := len(sessionKeys) - n

	// Remove oldest messages first
	for i := 0; i < removeCount; i++ {
		oldestId := sessionKeys[i]

		// Remove from Messages map
		delete(m.Messages, oldestId)

		// Remove from Keys slice
		for j, key := range m.Keys {
			if key == oldestId {
				m.Keys = append(m.Keys[:j], m.Keys[j+1:]...)
				break
			}
		}
	}

	return nil
}
