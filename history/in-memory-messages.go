package history

import "github.com/parakeet-nest/parakeet/llm"

type MemoryMessages struct {
	Messages map[string]llm.MessageRecord
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

func (m *MemoryMessages) GetAllMessages() ([]llm.Message, error) {
	var messages []llm.Message
	for _, messageRecord := range m.Messages {
		messages = append(messages, llm.Message{
			Role:    messageRecord.Role,
			Content: messageRecord.Content,
		})
	}
	return messages, nil
}

func (m *MemoryMessages) Save(messageRecord llm.MessageRecord) (llm.MessageRecord, error) {
	m.Messages[messageRecord.Id] = messageRecord
	return messageRecord, nil
}

func (m *MemoryMessages) SaveMessage(id string, message llm.Message) (llm.MessageRecord, error) {
	messageRecord := llm.MessageRecord{
		Id:      id,
		Role:    message.Role,
		Content: message.Content,
	}
	return m.Save(messageRecord)
}
