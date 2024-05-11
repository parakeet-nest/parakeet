package history

import "github.com/parakeet-nest/parakeet/llm"

type Messages interface {
	Get(id string) (llm.MessageRecord, error)
	GetMessage(id string) (llm.Message, error)
	GetAll() ([]llm.MessageRecord, error)
	GetAllMessages() ([]llm.Message, error)
	Save(messageRecord llm.MessageRecord) (llm.MessageRecord, error)
	SaveMessage(id string, message llm.Message) (llm.MessageRecord, error)
}
// TODO: implement Delete
// TODO: get with prefix


