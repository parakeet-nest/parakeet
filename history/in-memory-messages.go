package history

import (
	"sort"
	"strconv"
	"github.com/google/uuid"
	"github.com/parakeet-nest/parakeet/llm"
)

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

// TODO: implement the filter by pattern()
func (m *MemoryMessages) GetAllMessages(patterns ...string) ([]llm.Message, error) {

	keys := make([]string, 0, len(m.Messages))
	for k := range m.Messages {
		keys = append(keys, k)
	}

	// Sort the keys
	sort.Strings(keys)

	// Create ordered slice of messages
	var messages []llm.Message

	for _, key := range keys {

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

func (m *MemoryMessages) SaveMessageWithSessionId(sessionId, messageId string, message llm.Message) (llm.MessageRecord, error) {
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
	return nil
}

// TODO: to test
func (m *MemoryMessages) RemoveAllMessages() error {
	m.Messages = map[string]llm.MessageRecord{}
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

func (m *MemoryMessages) SaveMessageWithSession(sessionId string, messagesCounters *map[string]int, message llm.Message) (llm.MessageRecord, error) {
	//? generate an id for the message
	generateId := func(counter int, sessionId string) string {
		return strconv.Itoa(counter) + "-" + sessionId
	}

	//* increment the counter and save the user message
	(*messagesCounters)[sessionId]++
	messageRecord, err := m.SaveMessageWithSessionId(sessionId, generateId((*messagesCounters)[sessionId], sessionId), llm.Message{
		Role:    "user",
		Content: message.Content,
	})
	if err != nil {
		return llm.MessageRecord{}, err
	}

	return messageRecord, nil

}

func (m *MemoryMessages) RemoveTopMessageOfSession(sessionId string, messagesCounters *map[string]int, conversationLength int) error {
	//? generate an id for the message
	generateId := func(counter int, sessionId string) string {
		return strconv.Itoa(counter) + "-" + sessionId
	}

	//? get the top message id of a conversation of maxMessages messages for a given sessionId
	getTopMessageId := func(conversationLength, counter int, sessionId string) string {
		return generateId(counter-(conversationLength-1), sessionId)
	}

	if (*messagesCounters)[sessionId] >= conversationLength {
		//fmt.Println("🟢 counter:", (*messagesCounters)[sessionId])

		topMessageId := getTopMessageId(conversationLength, (*messagesCounters)[sessionId], sessionId)

		//msg, _ := m.Get(topMessageId)
		//fmt.Println("🟩 message:", msg.Id, msg.Role, msg.Content)

		err := m.RemoveMessage(topMessageId)
		if err != nil {
			return err
		}
	}
	return nil
}
