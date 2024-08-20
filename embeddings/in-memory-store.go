package embeddings

import (
	"github.com/google/uuid"
	"github.com/parakeet-nest/parakeet/llm"
)

type MemoryVectorStore struct {
	Records map[string]llm.VectorRecord
}

func (mvs *MemoryVectorStore) Get(id string) (llm.VectorRecord, error) {
	return mvs.Records[id], nil
}

func (mvs *MemoryVectorStore) GetAll() ([]llm.VectorRecord, error) {
	var records []llm.VectorRecord
	for _, record := range mvs.Records {
		records = append(records, record)
	}
	return records, nil
}

func (mvs *MemoryVectorStore) Save(vectorRecord llm.VectorRecord) (llm.VectorRecord, error) {
	if vectorRecord.Id == "" {
		vectorRecord.Id = uuid.New().String()
	}
	mvs.Records[vectorRecord.Id] = vectorRecord
	return vectorRecord, nil
}

// SearchMaxSimilarity finds the vector record in MemoryVectorStore with the maximum cosine distance similarity to the provided vector record.
//
// Parameters:
//   - embeddingFromQuestion: llm.VectorRecord - the vector record to compare similarities with.
//
// Returns:
//   - llm.VectorRecord: The vector record with the maximum similarity.
//   - error: Error if any.
func (mvs *MemoryVectorStore) SearchMaxSimilarity(embeddingFromQuestion llm.VectorRecord) (llm.VectorRecord, error) {

	var maxDistance float64 = -1.0
	var selectedKeyRecord string
	for k, v := range mvs.Records {
		distance := CosineDistance(embeddingFromQuestion.Embedding, v.Embedding)
		if distance > maxDistance {
			maxDistance = distance
			selectedKeyRecord = k
		}
	}
	// Return only the nearest vector to the question
	return mvs.Records[selectedKeyRecord], nil
}

// SearchSimilarities searches for vector records in the MemoryVectorStore that have a cosine distance similarity greater than or equal to the given limit.
//
// Parameters:
//   - embeddingFromQuestion: the vector record to compare similarities with.
//   - limit: the minimum cosine distance similarity threshold.
//
// Returns:
//   - []llm.VectorRecord: a slice of vector records that have a cosine distance similarity greater than or equal to the limit.
//   - error: an error if any occurred during the search.
func (mvs *MemoryVectorStore) SearchSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64) ([]llm.VectorRecord, error) {

	var records []llm.VectorRecord

	for _, v := range mvs.Records {
		distance := CosineDistance(embeddingFromQuestion.Embedding, v.Embedding)
		if distance >= limit {
			v.CosineDistance = distance
			records = append(records, v)
		}
	}
	return records, nil
}

// SearchTopNSimilarities searches for the top N similar vector records based on the given embedding from a question.
// It returns a slice of vector records and an error if any.
// The limit parameter specifies the minimum similarity score for a record to be considered similar.
// The max parameter specifies the maximum number of vector records to return.
func (mvs *MemoryVectorStore) SearchTopNSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64, max int) ([]llm.VectorRecord, error) {
	records, err := mvs.SearchSimilarities(embeddingFromQuestion, limit)
	if err != nil {
		return nil, err
	}
	return getTopNVectorRecords(records, max), nil
}
