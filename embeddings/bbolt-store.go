package embeddings

import (
	"encoding/json"

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

// TODO: if vectorRecord.Id == "" create a uuid
func (bvs *BboltVectorStore) Save(vectorRecord llm.VectorRecord) (llm.VectorRecord, error) {

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

func (bvs *BboltVectorStore) SearchSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64) ([]llm.VectorRecord, error) {

	records, err := bvs.GetAll()
	if err != nil {
		return nil, err
	}
	var recordsFiltered []llm.VectorRecord
	for _, v := range records {
		distance := CosineDistance(embeddingFromQuestion.Embedding, v.Embedding)
		if distance >= limit {
			recordsFiltered = append(recordsFiltered, v)
		}
	}
	return recordsFiltered, nil
}
