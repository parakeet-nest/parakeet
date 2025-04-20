package history

import (
	"encoding/json"
	"fmt"

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

func (b *BboltMessages) SaveMessageWithSession(sessionId, messageId string, message llm.Message) (llm.MessageRecord, error) {
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

func (b *BboltMessages) RemoveTopMessageOfSession(sessionId string) error {
	return b.messages.Update(func(tx *bolt.Tx) error {
		orderBucket := tx.Bucket([]byte(orderBucketName))
		messagesBucket := tx.Bucket([]byte(messagesBucketName))

		if orderBucket == nil || messagesBucket == nil {
			return fmt.Errorf("buckets not found")
		}

		// Find the first message of the session
		var topOrderKey, topMessageId []byte
		cursor := orderBucket.Cursor()
		for orderKey, messageId := cursor.First(); orderKey != nil; orderKey, messageId = cursor.Next() {
			messageData := messagesBucket.Get(messageId)
			if messageData == nil {
				continue
			}

			var messageRecord llm.MessageRecord
			if err := json.Unmarshal(messageData, &messageRecord); err != nil {
				return err
			}

			if messageRecord.SessionId == sessionId {
				topOrderKey = orderKey
				topMessageId = messageId
				break
			}
		}

		// No messages found for this session
		if topMessageId == nil {
			return nil
		}

		// Remove from messages bucket
		if err := messagesBucket.Delete(topMessageId); err != nil {
			return err
		}

		// Remove from order bucket
		return orderBucket.Delete(topOrderKey)
	})
}

// RemoveTopMessage removes the oldest message from the database
func (b *BboltMessages) RemoveTopMessage() error {
	return b.messages.Update(func(tx *bolt.Tx) error {
		orderBucket := tx.Bucket([]byte(orderBucketName))
		messagesBucket := tx.Bucket([]byte(messagesBucketName))

		if orderBucket == nil || messagesBucket == nil {
			return fmt.Errorf("buckets not found")
		}

		// Get the first (oldest) message
		cursor := orderBucket.Cursor()
		orderKey, messageId := cursor.First()
		if orderKey == nil {
			return nil // No messages to remove
		}

		// Remove from messages bucket
		err := messagesBucket.Delete(messageId)
		if err != nil {
			return err
		}

		// Remove from order bucket
		return orderBucket.Delete(orderKey)
	})
}

// KeepLastN removes all messages except the last n messages
func (b *BboltMessages) KeepLastN(n int) error {
	if n < 0 {
		return fmt.Errorf("n must be positive, got %d", n)
	}

	return b.messages.Update(func(tx *bolt.Tx) error {
		orderBucket := tx.Bucket([]byte(orderBucketName))
		messagesBucket := tx.Bucket([]byte(messagesBucketName))

		if orderBucket == nil || messagesBucket == nil {
			return fmt.Errorf("buckets not found")
		}

		// Count total messages
		total := 0
		cursor := orderBucket.Cursor()
		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			total++
		}

		if total <= n {
			return nil // Nothing to remove
		}

		// Remove oldest messages until we have n messages left
		removeCount := total - n
		cursor = orderBucket.Cursor()
		for k, v := cursor.First(); k != nil && removeCount > 0; k, v = cursor.Next() {
			if err := messagesBucket.Delete(v); err != nil {
				return err
			}
			if err := orderBucket.Delete(k); err != nil {
				return err
			}
			removeCount--
		}

		return nil
	})
}

// KeepLastNOfSession removes all messages of a session except the last n messages
func (b *BboltMessages) KeepLastNOfSession(sessionId string, n int) error {
	if n < 0 {
		return fmt.Errorf("n must be positive, got %d", n)
	}

	return b.messages.Update(func(tx *bolt.Tx) error {
		orderBucket := tx.Bucket([]byte(orderBucketName))
		messagesBucket := tx.Bucket([]byte(messagesBucketName))

		if orderBucket == nil || messagesBucket == nil {
			return fmt.Errorf("buckets not found")
		}

		// First, collect all message IDs for this session
		var sessionMessages []struct {
			orderKey  []byte
			messageId []byte
		}

		cursor := orderBucket.Cursor()
		for orderKey, messageId := cursor.First(); orderKey != nil; orderKey, messageId = cursor.Next() {
			var messageRecord llm.MessageRecord
			messageData := messagesBucket.Get(messageId)
			if messageData == nil {
				continue
			}

			if err := json.Unmarshal(messageData, &messageRecord); err != nil {
				return err
			}

			if messageRecord.SessionId == sessionId {
				sessionMessages = append(sessionMessages, struct {
					orderKey  []byte
					messageId []byte
				}{orderKey, messageId})
			}
		}

		if len(sessionMessages) <= n {
			return nil // Nothing to remove
		}

		// Remove oldest messages until we have n messages left
		removeCount := len(sessionMessages) - n
		for i := 0; i < removeCount; i++ {
			if err := messagesBucket.Delete(sessionMessages[i].messageId); err != nil {
				return err
			}
			if err := orderBucket.Delete(sessionMessages[i].orderKey); err != nil {
				return err
			}
		}

		return nil
	})
}
