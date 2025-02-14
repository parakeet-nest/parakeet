package history

import (
	"encoding/json"
	"strconv"

	bbolt "github.com/parakeet-nest/parakeet/db"
	"github.com/parakeet-nest/parakeet/llm"
	bolt "go.etcd.io/bbolt"
)

const bucketName string = "messages-bucket"

type BboltMessages struct {
	messages *bolt.DB
}

func (b *BboltMessages) Initialize(dbPath string) error {

	db, err := bbolt.Initialize(dbPath, bucketName)
	if err != nil {
		return err
	}
	b.messages = db
	return nil
}

func (b *BboltMessages) Get(id string) (llm.MessageRecord, error) {
	jsonStr := bbolt.Get(b.messages, bucketName, id)
	messageRecord := llm.MessageRecord{}
	err := json.Unmarshal([]byte(jsonStr), &messageRecord)
	if err != nil {
		return llm.MessageRecord{}, err
	}
	return messageRecord, nil
}

func (b *BboltMessages) GetMessage(id string) (llm.Message, error) {
	messageRecord, err := b.Get(id)
	if err != nil {
		return llm.Message{}, err
	}
	return llm.Message{
		Role:    messageRecord.Role,
		Content: messageRecord.Content,
	}, nil
}

func (b *BboltMessages) GetAll() ([]llm.MessageRecord, error) {
	var records []llm.MessageRecord
	mapStr := bbolt.GetAll(b.messages, bucketName)
	for _, v := range mapStr {
		messageRecord := llm.MessageRecord{}
		err := json.Unmarshal([]byte(v), &messageRecord)
		if err != nil {
			return nil, err
		}
		records = append(records, messageRecord)
	}
	return records, nil
}

// TODO: implement the filter by pattern()
func (b *BboltMessages) GetAllMessages(patterns ...string) ([]llm.Message, error) {
	var messages []llm.Message
	records, err := b.GetAll()
	if err != nil {
		return nil, err
	}
	for _, messageRecord := range records {
		messages = append(messages, llm.Message{
			Role:    messageRecord.Role,
			Content: messageRecord.Content,
		})
	}
	return messages, nil
}

// TODO: to test
// TODO: implement the filter by pattern()
func (b *BboltMessages) GetAllMessagesOfSession(sessionId string, patterns ...string) ([]llm.Message, error) {
	var messages []llm.Message
	records, err := b.GetAll()
	if err != nil {
		return nil, err
	}
	for _, messageRecord := range records {
		if messageRecord.SessionId == sessionId {
			messages = append(messages, llm.Message{
				Role:    messageRecord.Role,
				Content: messageRecord.Content,
			})
		}
	}
	return messages, nil
}

func (b *BboltMessages) Save(messageRecord llm.MessageRecord) (llm.MessageRecord, error) {

	jsonStr, err := json.Marshal(messageRecord)
	if err != nil {
		return llm.MessageRecord{}, err
	}

	err = bbolt.Save(b.messages, bucketName, messageRecord.Id, string(jsonStr))
	if err != nil {
		return llm.MessageRecord{}, err
	}
	return messageRecord, nil
}

func (b *BboltMessages) SaveMessage(id string, message llm.Message) (llm.MessageRecord, error) {
	messageRecord := llm.MessageRecord{
		Id:      id,
		Role:    message.Role,
		Content: message.Content,
	}
	return b.Save(messageRecord)
}

func (b *BboltMessages) SaveMessageWithSessionId(sessionId, messageId string, message llm.Message) (llm.MessageRecord, error) {
	messageRecord := llm.MessageRecord{
		Id:        messageId,
		Role:      message.Role,
		Content:   message.Content,
		SessionId: sessionId,
	}
	return b.Save(messageRecord)
}

func (b *BboltMessages) RemoveMessage(id string) error {
	return bbolt.Delete(b.messages, bucketName, id)
}

// TODO: to test
func (b *BboltMessages) RemoveAllMessages() error {
	records, err := b.GetAll()
	if err != nil {
		return err
	}
	for _, messageRecord := range records {
		err := b.RemoveMessage(messageRecord.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: to test
func (b *BboltMessages) RemoveAllMessagesOfSession(sessionId string) error {
	records, err := b.GetAll()
	if err != nil {
		return err
	}
	for _, messageRecord := range records {
		if messageRecord.SessionId == sessionId {
			err := b.RemoveMessage(messageRecord.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *BboltMessages) SaveMessageWithSession(sessionId string, messagesCounters *map[string]int, message llm.Message) (llm.MessageRecord, error) {
	//? generate an id for the message
	generateId := func(counter int, sessionId string) string {
		return strconv.Itoa(counter) + "-" + sessionId
	}

	//* increment the counter and save the user message
	(*messagesCounters)[sessionId]++
	messageRecord, err := b.SaveMessageWithSessionId(sessionId, generateId((*messagesCounters)[sessionId], sessionId), llm.Message{
		Role:    "user",
		Content: message.Content,
	})
	if err != nil {
		return llm.MessageRecord{}, err
	}

	return messageRecord, nil
}

func (b *BboltMessages) RemoveTopMessageOfSession(sessionId string, messagesCounters *map[string]int, conversationLength int) error {
	//? generate an id for the message
	generateId := func(counter int, sessionId string) string {
		return strconv.Itoa(counter) + "-" + sessionId
	}

	//? get the top message id of a conversation of maxMessages messages for a given sessionId
	getTopMessageId := func(conversationLength, counter int, sessionId string) string {
		return generateId(counter-(conversationLength-1), sessionId)
	}

	if (*messagesCounters)[sessionId] >= conversationLength {
		//fmt.Println("🔵 counter:", (*messagesCounters)[sessionId])

		topMessageId := getTopMessageId(conversationLength, (*messagesCounters)[sessionId], sessionId)

		//m, _ := b.Get(topMessageId)
		//fmt.Println("🟦 message:", m.Id, m.Role, m.Content)

		err := b.RemoveMessage(topMessageId)
		if err != nil {
			return err
		}
	}
	return nil
}
