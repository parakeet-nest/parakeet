package history


import (
	"encoding/json"

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

func (b *BboltMessages) GetAllMessages() ([]llm.Message, error) {
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
