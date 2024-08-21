package embeddings

import (
	"encoding/json"

	"github.com/google/uuid"
	bbolt "github.com/parakeet-nest/parakeet/db"
	"github.com/parakeet-nest/parakeet/llm"
	bolt "go.etcd.io/bbolt"
)

const bucketName string = "embeddings-store-bucket"

type BboltVectorStore struct {
	store *bolt.DB
}

func (bvs *BboltVectorStore) Initialize(dbPath string) error {

	db, err := bbolt.Initialize(dbPath, bucketName)
	if err != nil {
		return err
	}
	bvs.store = db
	return nil
}

func (bvs *BboltVectorStore) Get(id string) (llm.VectorRecord, error) {
	jsonStr := bbolt.Get(bvs.store, bucketName, id)
	vectorRecord := llm.VectorRecord{}
	err := json.Unmarshal([]byte(jsonStr), &vectorRecord)
	if err != nil {
		return llm.VectorRecord{}, err
	}
	return vectorRecord, nil
}

func (bvs *BboltVectorStore) GetAll() ([]llm.VectorRecord, error) {
	var records []llm.VectorRecord
	mapStr := bbolt.GetAll(bvs.store, bucketName)
	for _, v := range mapStr {
		vectorRecord := llm.VectorRecord{}
		err := json.Unmarshal([]byte(v), &vectorRecord)
		if err != nil {
			return nil, err
		}
		records = append(records, vectorRecord)
	}
	return records, nil
}

func (bvs *BboltVectorStore) Save(vectorRecord llm.VectorRecord) (llm.VectorRecord, error) {
	if vectorRecord.Id == "" {
		vectorRecord.Id = uuid.New().String()
	}

	jsonStr, err := json.Marshal(vectorRecord)
	if err != nil {
		return llm.VectorRecord{}, err
	}

	err = bbolt.Save(bvs.store, bucketName, vectorRecord.Id, string(jsonStr))
	if err != nil {
		return llm.VectorRecord{}, err
	}
	return vectorRecord, nil

}

// SearchMaxSimilarity searches for the vector record in the BboltVectorStore that has the maximum cosine distance similarity to the given embeddingFromQuestion.
//
// Parameters:
// - embeddingFromQuestion: the vector record to compare similarities with.
//
// Returns:
// - llm.VectorRecord: the vector record with the maximum cosine distance similarity.
// - error: an error if any occurred during the search.
func (bvs *BboltVectorStore) SearchMaxSimilarity(embeddingFromQuestion llm.VectorRecord) (llm.VectorRecord, error) {
	var maxDistance float64 = -1.0
	var selectedKeyRecord string

	records, err := bvs.GetAll()
	if err != nil {
		return llm.VectorRecord{}, err
	}
	for _, v := range records {
		distance := CosineDistance(embeddingFromQuestion.Embedding, v.Embedding)
		if distance > maxDistance {
			maxDistance = distance
			selectedKeyRecord = v.Id
		}
		//fmt.Println("  - ", selectedKeyRecord, v.Id, distance)
	}
	// Return only the nearest vector to the question
	return bvs.Get(selectedKeyRecord)

}

// SearchSimilarities searches for vector records in the BboltVectorStore that have a cosine distance similarity greater than or equal to the given limit.
//
// Parameters:
//   - embeddingFromQuestion: the vector record to compare similarities with.
//   - limit: the minimum cosine distance similarity threshold.
//
// Returns:
//   - []llm.VectorRecord: a slice of vector records that have a cosine distance similarity greater than or equal to the limit.
//   - error: an error if any occurred during the search.
func (bvs *BboltVectorStore) SearchSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64) ([]llm.VectorRecord, error) {
	records, err := bvs.GetAll()

	if err != nil {
		return nil, err
	}
	var recordsFiltered []llm.VectorRecord
	for _, v := range records {

		distance := CosineDistance(embeddingFromQuestion.Embedding, v.Embedding)
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
func (bvs *BboltVectorStore) SearchTopNSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64, max int) ([]llm.VectorRecord, error) {
	records, err := bvs.SearchSimilarities(embeddingFromQuestion, limit)
	if err != nil {
		return nil, err
	}
	return getTopNVectorRecords(records, max), nil
}
