package squawk

import (
	"log"

	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func (s *Squawk) Store(store embeddings.VectorStore, optionalParameters ...string) *Squawk {

	/*
		switch v := store.(type) {
		case *embeddings.MemoryVectorStore:
			s.vectorStore = v
		case *embeddings.RedisVectorStore:
			s.vectorStore = v
		case *embeddings.BboltVectorStore:
			s.vectorStore = v
		case *embeddings.DaphniaVectoreStore:
			s.vectorStore = v
		case *embeddings.ElasticsearchStore:
			s.vectorStore = v
		default:
			// Handle unknown store type or set a default
			log.Printf("Warning: Unknown vector store type: %T", store)
		}
	*/
	s.vectorStore = store
	if len(optionalParameters) > 0 {
	}
	return s

}

// TODO: reset the vector : set an empty one

// Generate embeddings for the given documents
// and store them in the vector store.
// The function returns the Squawk instance for method chaining.
func (s *Squawk) generateEmbeddingsFromDocuments(docs []string, logs bool) *Squawk {

	for idx, doc := range docs {
		if logs {
			log.Println("📝 Creating embedding from document ", idx)
		}

		//TODO:query


		embedding, err := embeddings.CreateEmbedding(
			s.apiUrl,
			llm.Query4Embedding{
				Model:  s.embeddingsModel,
				Prompt: doc,
			},
			"",
			s.provider, s.openAPIKey,
		)
		if err != nil {
			if logs {
				log.Println("😡 When generating embedding:", err)
			}
			s.lastError = err
		} else {
			record, err := s.vectorStore.Save(embedding)
			if err != nil {
				if logs {
					log.Println("😡 When saving embedding:", err)
				}
				s.lastError = err
			} else {
				if logs {
					log.Println("📦 Embedding saved:", record.Id)
				}
			}

		}
	}

	return s
}

func (s *Squawk) GenerateEmbeddings(docs []string, optionalParameters ...any) *Squawk {
	if len(optionalParameters) > 0 {
		if logs, ok := optionalParameters[0].(bool); ok {
			s.generateEmbeddingsFromDocuments(docs, logs)
		} else {
			s.generateEmbeddingsFromDocuments(docs, false)
		}
	}
	return s
}

func (s *Squawk) searchSimilarities(content string, limit float64, max int, logs bool) *Squawk {

	// Create an embedding from the content
	embeddingFromContent, err := embeddings.CreateEmbedding(
		s.apiUrl,
		llm.Query4Embedding{
			Model:  s.embeddingsModel,
			Prompt: content,
		},
		"content",
		s.provider, s.openAPIKey,

	)
	if err != nil {
		if logs {
			log.Println("😡 [Similarity Search] When creating embedding:", err)
			s.lastError = err
		}
		return s
	}
	if logs {
		log.Println("🔎 searching for similarity...")
	}
	similarities, err := s.vectorStore.SearchTopNSimilarities(embeddingFromContent, limit, max)
	if err != nil {
		if logs {
			log.Println("😡 When searching similarities:", err)
			s.lastError = err
		}
		return s
	}
	if logs {
		log.Println("🎉 similarities", len(similarities))
	}
	s.similarities = similarities

	return s
}

func (s *Squawk) SimilaritySearch(content string, limit float64, max int, optionalParameters ...any) *Squawk {
	if len(optionalParameters) > 0 {
		if logs, ok := optionalParameters[0].(bool); ok {
			s.searchSimilarities(content, limit, max, logs)
		} else {
			s.searchSimilarities(content, limit, max, false)
		}
	} else {
		s.searchSimilarities(content, limit, max, false)
	}
	return s
}

func (s *Squawk) SimilaritySearchFromUserMessage(userMessageLabel string, limit float64, max int, optionalParameters ...any) *Squawk {
	
	// search the message with the label
	var content string
	for _, message := range s.setOfMessages {
		if message.Label == userMessageLabel {
			content = message.Content
			break
		}
	}
	// TODO: handle the case where the message is not found
	if content == "" {
		log.Println("😡 No message found with label:", userMessageLabel)
		return s
	}
	
	if len(optionalParameters) > 0 {
		if logs, ok := optionalParameters[0].(bool); ok {
			s.searchSimilarities(content, limit, max, logs)
		} else {
			s.searchSimilarities(content, limit, max, false)
		}
	} else {
		s.searchSimilarities(content, limit, max, false)
	}
	return s
}

func (s *Squawk) AddSimilaritiesToMessages(optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "system", Content: embeddings.GenerateContextFromSimilarities(s.similarities), Label: optionalParameters[0]},
		)

	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "system", Content: embeddings.GenerateContextFromSimilarities(s.similarities)})
	}

	return s
}

func (s *Squawk) AddSimilaritiesToMessagesWithPrefix(prefix string, optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "system", Content: prefix+embeddings.GenerateContextFromSimilarities(s.similarities), Label: optionalParameters[0]},
		)

	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "system", Content: prefix+embeddings.GenerateContextFromSimilarities(s.similarities)})
	}

	return s
}



func (s *Squawk) Similarities() []llm.VectorRecord {
	return s.similarities
}

func (s *Squawk) ContextFromSimilarities() string {
	return embeddings.GenerateContextFromSimilarities(s.similarities)
}

// TODO: create a function to generate the context from the similarities
// and add it to the system message (as a system message)

// 	documentsContent := embeddings.GenerateContextFromSimilarities(similarities)

/*

func (s *Squawk) embeddingsFromChunks() *Squawk {

	return s
}

func (s *Squawk) embedding(doc string) *Squawk {

	return s
}
*/
