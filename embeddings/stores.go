package embeddings

import "github.com/parakeet-nest/parakeet/llm"

type VectorStore interface {
	Get(id string) (llm.VectorRecord, error)
	GetAll() ([]llm.VectorRecord, error)
	Save(vectorRecord llm.VectorRecord) (llm.VectorRecord, error)
	SearchMaxSimilarity(embeddingFromQuestion llm.VectorRecord) (llm.VectorRecord, error)
	SearchSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64) ([]llm.VectorRecord, error)
	SearchTopNSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64, max int) ([]llm.VectorRecord, error)
}
// TODO: implement Delete

