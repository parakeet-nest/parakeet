package embeddings

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/similarity"
)

//const redisKeyPrefix string = "embeddings-store"

type RedisVectorStore struct {
	client         *redis.Client
	ctx            context.Context
	redisKeyPrefix string
}

/*
client := redis.NewClient(&redis.Options{
	Addr:     redisServer, // the name is defined in the compose.yml file
	Password: "",          // no password set
	DB:       0,           // use default DB
})
*/

func (rvs *RedisVectorStore) Initialize(redisAddr string, redisPwd string, storeName string) error {
	rvs.ctx = context.Background()
	rvs.client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPwd,
		DB:       0,
	})
	_, err := rvs.client.Ping(rvs.ctx).Result()
	if err != nil {
		return err
	}
	rvs.redisKeyPrefix = storeName
	return nil
}

func (rvs *RedisVectorStore) Get(id string) (llm.VectorRecord, error) {
	jsonStr, err := rvs.client.Get(rvs.ctx, rvs.redisKeyPrefix+":"+id).Result()
	if err != nil {
		return llm.VectorRecord{}, err
	}
	vectorRecord := llm.VectorRecord{}
	err = json.Unmarshal([]byte(jsonStr), &vectorRecord)
	if err != nil {
		return llm.VectorRecord{}, err
	}
	return vectorRecord, nil
}

func (rvs *RedisVectorStore) GetAll() ([]llm.VectorRecord, error) {
	var records []llm.VectorRecord
	keys, err := rvs.client.Keys(rvs.ctx, rvs.redisKeyPrefix+":*").Result()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		jsonStr, err := rvs.client.Get(rvs.ctx, key).Result()
		if err != nil {
			return nil, err
		}
		vectorRecord := llm.VectorRecord{}
		err = json.Unmarshal([]byte(jsonStr), &vectorRecord)
		if err != nil {
			return nil, err
		}
		records = append(records, vectorRecord)
	}
	return records, nil
}

func (rvs *RedisVectorStore) Save(vectorRecord llm.VectorRecord) (llm.VectorRecord, error) {
	if vectorRecord.Id == "" {
		vectorRecord.Id = uuid.New().String()
	}

	jsonStr, err := json.Marshal(vectorRecord)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	err = rvs.client.Set(rvs.ctx, rvs.redisKeyPrefix+":"+vectorRecord.Id, jsonStr, 0).Err()
	if err != nil {
		return llm.VectorRecord{}, err
	}
	return vectorRecord, nil
}

// SearchMaxSimilarity searches for the vector record in the RedisVectorStore that has the maximum cosine distance similarity to the given embeddingFromQuestion.
//
// Parameters:
// - embeddingFromQuestion: the vector record to compare similarities with.
//
// Returns:
// - llm.VectorRecord: the vector record with the maximum cosine distance similarity.
// - error: an error if any occurred during the search.
func (rvs *RedisVectorStore) SearchMaxSimilarity(embeddingFromQuestion llm.VectorRecord) (llm.VectorRecord, error) {
	var maxDistance float64 = -1.0
	var selectedKeyRecord string

	records, err := rvs.GetAll()
	if err != nil {
		return llm.VectorRecord{}, err
	}
	for _, v := range records {
		distance := similarity.CosineDistance(embeddingFromQuestion.Embedding, v.Embedding)
		if distance > maxDistance {
			maxDistance = distance
			selectedKeyRecord = v.Id
		}
	}
	// Return only the nearest vector to the question
	return rvs.Get(selectedKeyRecord)

}

// SearchSimilarities searches for vector records in the RedisVectorStore that have a cosine distance similarity greater than or equal to the given limit.
//
// Parameters:
//   - embeddingFromQuestion: the vector record to compare similarities with.
//   - limit: the minimum cosine distance similarity threshold.
//
// Returns:
//   - []llm.VectorRecord: a slice of vector records that have a cosine distance similarity greater than or equal to the limit.
//   - error: an error if any occurred during the search.
func (rvs *RedisVectorStore) SearchSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64) ([]llm.VectorRecord, error) {
	records, err := rvs.GetAll()
	if err != nil {
		return nil, err
	}
	var recordsFiltered []llm.VectorRecord
	for _, v := range records {
		distance := similarity.CosineDistance(embeddingFromQuestion.Embedding, v.Embedding)
		if distance >= limit {
			v.CosineDistance = distance
			recordsFiltered = append(recordsFiltered, v)
		}
	}
	return recordsFiltered, nil
}

// SearchTopNSimilarities searches for the top N similar vector records based on the given embedding from a question.
// It returns a slice of vector records and an error if any.
// The limit parameter specifies the minimum similarity score for a record to be considered similar.
// The max parameter specifies the maximum number of vector records to return.
func (rvs *RedisVectorStore) SearchTopNSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64, max int) ([]llm.VectorRecord, error) {
	records, err := rvs.SearchSimilarities(embeddingFromQuestion, limit)
	if err != nil {
		return nil, err
	}
	return similarity.GetTopNVectorRecords(records, max), nil
}
