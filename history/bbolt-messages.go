package history

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	bbolt "github.com/parakeet-nest/parakeet/db"
	"github.com/parakeet-nest/parakeet/llm"
	bolt "go.etcd.io/bbolt"
)

const (
	messagesBucketName string = "messages-bucket"
	orderBucketName    string = "messages-order-bucket" // New bucket for ordering
)

type BboltMessages struct {
	messages *bolt.DB
	counter  int // Track message count for ordering

}

func (b *BboltMessages) Initialize(dbPath string) error {
	db, err := bbolt.Initialize(dbPath, messagesBucketName)
	if err != nil {
		return err
	}

	// Create order bucket in the same db
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(orderBucketName))
		return err
	})
	if err != nil {
		return err
	}

	b.messages = db
	b.counter = 0
	return nil
}

func (b *BboltMessages) Get(id string) (llm.MessageRecord, error) {
    var messageRecord llm.MessageRecord
    
    err := b.messages.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(messagesBucketName))
        if bucket == nil {
            return fmt.Errorf("bucket %s not found", messagesBucketName)
        }
        
        data := bucket.Get([]byte(id))
        if data == nil {
            return fmt.Errorf("message with id %s not found", id)
        }
        
        return json.Unmarshal(data, &messageRecord)
    })
    
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
    
    err := b.messages.View(func(tx *bolt.Tx) error {
        orderBucket := tx.Bucket([]byte(orderBucketName))
        messagesBucket := tx.Bucket([]byte(messagesBucketName))
        
        if orderBucket == nil || messagesBucket == nil {
            return fmt.Errorf("buckets not found")
        }
        
        return orderBucket.ForEach(func(orderKey, messageId []byte) error {
            messageData := messagesBucket.Get(messageId)
            if messageData == nil {
                return nil // Skip if message was deleted
            }
            
            var messageRecord llm.MessageRecord
            if err := json.Unmarshal(messageData, &messageRecord); err != nil {
                return err
            }
            
            records = append(records, messageRecord)
            return nil
        })
    })
    
    if err != nil {
        return nil, err
    }
    
    return records, nil
}

// TODO: implement the filter by pattern()
func (b *BboltMessages) GetAllMessages(patterns ...string) ([]llm.Message, error) {
    var messages []llm.Message
    
    err := b.messages.View(func(tx *bolt.Tx) error {
        orderBucket := tx.Bucket([]byte(orderBucketName))
        messagesBucket := tx.Bucket([]byte(messagesBucketName))
        
        return orderBucket.ForEach(func(orderKey, messageId []byte) error {
            messageData := messagesBucket.Get(messageId)
            if messageData == nil {
                return nil // Skip if message was deleted
            }
            
            var messageRecord llm.MessageRecord
            if err := json.Unmarshal(messageData, &messageRecord); err != nil {
                return err
            }
            
            messages = append(messages, llm.Message{
                Role:    messageRecord.Role,
                Content: messageRecord.Content,
            })
            return nil
        })
    })
    
    if err != nil {
        return nil, err
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
	// Increment counter for ordering
	b.counter++
	orderKey := fmt.Sprintf("%019d", b.counter) // Pad with zeros for proper sorting

	// Start transaction to save both message and order
	err := b.messages.Update(func(tx *bolt.Tx) error {
		// Save message
		messagesBucket := tx.Bucket([]byte(messagesBucketName))
		jsonStr, err := json.Marshal(messageRecord)
		if err != nil {
			return err
		}
		err = messagesBucket.Put([]byte(messageRecord.Id), jsonStr)
		if err != nil {
			return err
		}

		// Save order
		orderBucket := tx.Bucket([]byte(orderBucketName))
		return orderBucket.Put([]byte(orderKey), []byte(messageRecord.Id))
	})

	if err != nil {
		return llm.MessageRecord{}, err
	}
	return messageRecord, nil
}

func (b *BboltMessages) SaveMessage(id string, message llm.Message) (llm.MessageRecord, error) {
	if id == "" {
		// generate a unique for the message
		id = uuid.New().String()
	}

	messageRecord := llm.MessageRecord{
		Id:      id,
		Role:    message.Role,
		Content: message.Content,
	}
	return b.Save(messageRecord)
}

func (b *BboltMessages) SaveMessageWithSessionId(sessionId, messageId string, message llm.Message) (llm.MessageRecord, error) {
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
	return b.Save(messageRecord)
}

func (b *BboltMessages) RemoveMessage(id string) error {
    return b.messages.Update(func(tx *bolt.Tx) error {
        messagesBucket := tx.Bucket([]byte(messagesBucketName))
        orderBucket := tx.Bucket([]byte(orderBucketName))
        
        // Find and remove from order bucket first
        c := orderBucket.Cursor()
        for k, v := c.First(); k != nil; k, v = c.Next() {
            if string(v) == id {
                if err := orderBucket.Delete(k); err != nil {
                    return err
                }
                break
            }
        }
        
        // Remove from messages bucket
        return messagesBucket.Delete([]byte(id))
    })
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
