package embeddings

import (
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/sea-monkeys/daphnia"
)

type DaphniaVectoreStore struct {
	store daphnia.VectorStore
}

func (dvs *DaphniaVectoreStore) Initialize(dbPath string) error {
	vectorStore := daphnia.VectorStore{}
	err := vectorStore.Initialize(dbPath)
	if err != nil {
		return err
	}
	dvs.store = vectorStore
	return nil
}

func (dvs *DaphniaVectoreStore) Get(id string) (llm.VectorRecord, error) {
	daphniaRecord, err := dvs.store.Get(id)
	if err != nil {
		return llm.VectorRecord{}, err
	}
	vectorRecord := llm.VectorRecord{
		Id:        id,
		Embedding: daphniaRecord.Embedding,
		Prompt:    daphniaRecord.Prompt,
	}
	return vectorRecord, nil
}

func (dvs *DaphniaVectoreStore) GetAll() ([]llm.VectorRecord, error) {
	daphniaRecords, err := dvs.store.GetAll()
	if err != nil {
		return nil, err
	}
	var vectorRecords []llm.VectorRecord
	for _, daphniaRecord := range daphniaRecords {
		vectorRecord := llm.VectorRecord{
			Id:        daphniaRecord.Id,
			Embedding: daphniaRecord.Embedding,
			Prompt:    daphniaRecord.Prompt,
		}
		vectorRecords = append(vectorRecords, vectorRecord)
	}
	return vectorRecords, nil
}

func (dvs *DaphniaVectoreStore) Save(vectorRecord llm.VectorRecord) (llm.VectorRecord, error) {
	daphniaRecord := daphnia.VectorRecord{
		Id:        vectorRecord.Id,
		Embedding: vectorRecord.Embedding,
		Prompt:    vectorRecord.Prompt,
		Metadata:  vectorRecord.Metadata,
	}
	daphniaRecord, err := dvs.store.Save(daphniaRecord)
	if err != nil {
		return llm.VectorRecord{}, err
	}
	vectorRecord.Id = daphniaRecord.Id
	return vectorRecord, nil
}

func (dvs *DaphniaVectoreStore) SearchMaxSimilarity(embeddingFromQuestion llm.VectorRecord) (llm.VectorRecord, error) {
	records, err := dvs.SearchTopNSimilarities(embeddingFromQuestion, 1.0, 1)
	if err != nil {
		return llm.VectorRecord{}, err
	}
	if len(records) == 0 {
		return llm.VectorRecord{}, nil
	}
	return records[0], nil
}

func (dvs *DaphniaVectoreStore) SearchTopNSimilarities(vectorRecord llm.VectorRecord, threshold float64, topN int) ([]llm.VectorRecord, error) {
	daphniaRecord := daphnia.VectorRecord{
		Id:        vectorRecord.Id,
		Embedding: vectorRecord.Embedding,
		Prompt:    vectorRecord.Prompt,
		Metadata:  vectorRecord.Metadata,
	}
	daphniaRecords, err := dvs.store.SearchTopNSimilarities(daphniaRecord, threshold, topN)
	if err != nil {
		return nil, err
	}
	var vectorRecords []llm.VectorRecord
	for _, daphniaRecord := range daphniaRecords {
		vectorRecord := llm.VectorRecord{
			Id:             daphniaRecord.Id,
			Embedding:      daphniaRecord.Embedding,
			Prompt:         daphniaRecord.Prompt,
			Metadata:       daphniaRecord.Metadata,
			CosineDistance: daphniaRecord.CosineDistance,
		}
		vectorRecords = append(vectorRecords, vectorRecord)
	}
	return vectorRecords, nil
}

func (dvs *DaphniaVectoreStore) SearchSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64) ([]llm.VectorRecord, error) {
	daphniaRecord := daphnia.VectorRecord{
		Id:        embeddingFromQuestion.Id,
		Embedding: embeddingFromQuestion.Embedding,
		Prompt:    embeddingFromQuestion.Prompt,
		Metadata:  embeddingFromQuestion.Metadata,
	}
	daphniaRecords, err := dvs.store.SearchSimilarities(daphniaRecord, limit)
	if err != nil {
		return nil, err
	}
	var vectorRecords []llm.VectorRecord
	for _, daphniaRecord := range daphniaRecords {
		vectorRecord := llm.VectorRecord{
			Id:             daphniaRecord.Id,
			Embedding:      daphniaRecord.Embedding,
			Prompt:         daphniaRecord.Prompt,
			Metadata:       daphniaRecord.Metadata,
			CosineDistance: daphniaRecord.CosineDistance,
		}
		vectorRecords = append(vectorRecords, vectorRecord)
	}
	return vectorRecords, nil
}
