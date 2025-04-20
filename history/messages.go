package history

import "github.com/parakeet-nest/parakeet/llm"

type Messages interface {
	Get(id string) (llm.MessageRecord, error)
	GetMessage(id string) (llm.Message, error)
	GetAll() ([]llm.MessageRecord, error)

	GetAllMessages(patterns ... string) ([]llm.Message, error)
	
	GetAllMessagesOfSession(sessionId string, patterns ... string) ([]llm.Message, error)


	Save(messageRecord llm.MessageRecord) (llm.MessageRecord, error)
	SaveMessage(id string, message llm.Message) (llm.MessageRecord, error)

	SaveMessageWithSessionId(sessionId, messageId string, message llm.Message) (llm.MessageRecord, error)


	RemoveMessage(id string) error
	RemoveAllMessages() error
	RemoveTopMessage() error
	
	RemoveAllMessagesOfSession(sessionId string) error

	SaveMessageWithSession(sessionId string, messagesCounters *map[string]int, message llm.Message) (llm.MessageRecord, error)
	RemoveTopMessageOfSession(sessionId string, messagesCounters *map[string]int, conversationLength int) error

}



