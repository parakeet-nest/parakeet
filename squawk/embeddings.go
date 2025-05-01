package squawk

import (
	"log"

	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

// Store sets the vector store for embeddings in the Squawk instance.
// This method configures which storage backend will be used for saving 
// and retrieving embeddings vectors.
//
// Parameters:
//   - store: An implementation of embeddings.VectorStore interface
//   - optionalParameters: Optional slice of strings for additional configuration
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider using memory store:
//   memStore := embeddings.NewMemoryVectorStore()
//   squawk := New().
//     Model("mistral:latest").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(memStore).
//     GenerateEmbeddings([]string{
//       "Go is a programming language",
//       "Python is another language",
//     })
//
// Example with similarity search:
//   squawk := New().
//     Model("codellama:13b").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings(docs).
//     SimilaritySearch("Go programming", 0.8, 5)
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

// generateEmbeddingsFromDocuments creates and stores embeddings for a list of documents.
// This is an internal method used by GenerateEmbeddings to process multiple documents
// and store their vector representations.
//
// Parameters:
//   - docs: A slice of strings containing the documents to be embedded
//   - logs: A boolean flag to enable/disable logging of the embedding process
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// The method performs the following steps for each document:
// 1. Creates an embedding using the configured embeddings model
// 2. Saves the embedding to the configured vector store
// 3. Updates lastError if any errors occur during processing
// 4. Logs progress if logging is enabled
//
// Example usage through GenerateEmbeddings:
//   squawk := New().
//     Model("mistral:latest").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings(
//       []string{
//         "Go is a statically typed language",
//         "Go has built-in concurrency support",
//       },
//       true, // Enable logging
//     )
//
// Error handling:
//   - Errors during embedding creation are logged and stored in lastError
//   - Errors during vector store saving are logged and stored in lastError
//   - Processing continues even if individual documents fail
func (s *Squawk) generateEmbeddingsFromDocuments(docs []string, logs bool) *Squawk {

	for idx, doc := range docs {
		if logs {
			log.Println("📝 Creating embedding from document ", idx)
		}

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

// GenerateEmbeddings creates vector embeddings for a list of documents and stores them
// in the configured vector store. This method provides a high-level interface to the
// embedding generation process.
//
// Parameters:
//   - docs: A slice of strings containing the documents to be embedded
//   - optionalParameters: Optional variadic parameters where the first parameter,
//     if provided and of type bool, controls logging output
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider and logging enabled:
//   squawk := New().
//     Model("mistral:latest").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings(
//       []string{
//         "Go is a statically typed language",
//         "Go has built-in concurrency support",
//       },
//       true, // Enable logging
//     )
//
// Example with simple usage:
//   squawk := New().
//     Model("codellama:13b").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings([]string{
//       "Goroutines are lightweight threads",
//       "Channels are used for communication",
//     })
//
// The method delegates to generateEmbeddingsFromDocuments for the actual
// embedding creation and storage process. Errors during processing can be
// retrieved using the LastError() method.
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

// searchSimilarities performs a similarity search against stored embeddings using
// the provided content. This is an internal method used by SimilaritySearch.
//
// Parameters:
//   - content: The text to search for similar documents
//   - limit: The similarity threshold (0.0 to 1.0) where 1.0 is exact match
//   - max: Maximum number of results to return
//   - logs: Boolean flag to enable/disable logging
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// The method performs these steps:
// 1. Creates an embedding from the search content
// 2. Searches the vector store for similar documents
// 3. Updates the similarities field with results
// 4. Handles and logs any errors that occur
//
// Example usage through SimilaritySearch:
//   squawk := New().
//     Model("mistral:latest").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings([]string{
//       "Go is a concurrent language",
//       "Go uses goroutines for parallelism",
//     }).
//     SimilaritySearch("concurrent programming", 0.7, 3, true)
//
// Error handling:
//   - Embedding creation errors are stored in lastError
//   - Search errors are stored in lastError
//   - Method continues to chain even if errors occur
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

// SimilaritySearch performs a semantic search against stored embeddings to find similar documents.
//
// Parameters:
//   - content: The text to search for similar documents
//   - limit: The similarity threshold (0.0 to 1.0) where 1.0 is exact match
//   - max: Maximum number of results to return
//   - optionalParameters: Optional variadic parameters where the first parameter,
//     if provided and of type bool, controls logging output
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//   squawk := New().
//     Model("mistral:latest").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings([]string{
//       "Go is a statically typed language",
//       "Go supports concurrent programming",
//     }).
//     SimilaritySearch("concurrent Go features", 0.7, 3, true)
//
// The results can be accessed using:
//   similarities := squawk.Similarities()
//   for _, sim := range similarities {
//     fmt.Printf("Score: %f, Content: %s\n", sim.Score, sim.Content)
//   }
//
// Common use cases:
// - Finding relevant documentation
// - Question answering with context
// - Content recommendation
// - Semantic search in documents
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

// SimilaritySearchFromUserMessage performs a semantic search using the content of a labeled user message
// as the search query. This method is useful when you want to find documents similar to
// a previously stored message in the conversation history.
//
// Parameters:
//   - userMessageLabel: The label of the user message to use as search content
//   - limit: The similarity threshold (0.0 to 1.0) where 1.0 is exact match
//   - max: Maximum number of results to return
//   - optionalParameters: Optional variadic parameters where the first parameter,
//     if provided and of type bool, controls logging output
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//   squawk := New().
//     Model("mistral:latest").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings([]string{
//       "Go is a statically typed language",
//       "Go supports concurrent programming",
//     }).
//     User("Tell me about Go concurrency", "question-1").
//     SimilaritySearchFromUserMessage("question-1", 0.7, 3, true)
//
// The results can be accessed using:
//   similarities := squawk.Similarities()
//   for _, sim := range similarities {
//     fmt.Printf("Score: %f, Content: %s\n", sim.Score, sim.Content)
//   }
//
// Error handling:
//   - Returns early if no message is found with the given label
//   - Logs an error message when message is not found
//   - Delegates actual search to searchSimilarities method
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

// AddSimilaritiesToMessages adds the context generated from similarity search results
// to the conversation as a system message. This method is useful for providing
// relevant context to the language model based on vector search results.
//
// Parameters:
//   - optionalParameters: Optional slice of strings where the first element,
//     if provided, is used as a label for the system message
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//   squawk := New().
//     Model("mistral:latest").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings([]string{
//       "Go is a statically typed language",
//       "Go supports concurrent programming",
//     }).
//     SimilaritySearch("concurrent Go features", 0.7, 3).
//     AddSimilaritiesToMessages("context-1").  // Add with label
//     User("Explain these concepts")
//
// Common use cases:
// - Adding search context to conversations
// - Providing relevant documentation to the model
// - Creating context-aware responses
// - Building RAG (Retrieval Augmented Generation) applications
func (s *Squawk) AddSimilaritiesToMessages(optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		//fmt.Println("🔴📝 Adding similarities to messages:", optionalParameters[0])
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "system", Content: embeddings.GenerateContextFromSimilarities(s.similarities), Label: optionalParameters[0]},
		)
		

	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "system", Content: embeddings.GenerateContextFromSimilarities(s.similarities)})
	}

	return s
}

// AddSimilaritiesToMessagesWithPrefix adds the context generated from similarity search results
// to the conversation as a system message, with a custom prefix. This method is useful for
// providing relevant context to the language model with additional instructions or formatting.
//
// Parameters:
//   - prefix: A string that will be prepended to the generated context
//   - optionalParameters: Optional slice of strings where the first element,
//     if provided, is used as a label for the system message
//
// Returns:
//   - *Squawk: A pointer to the same Squawk instance for method chaining
//
// Example with Ollama provider:
//   squawk := New().
//     Model("mistral:latest").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings([]string{
//       "Go is a statically typed language",
//       "Go supports concurrent programming",
//     }).
//     SimilaritySearch("concurrent Go features", 0.7, 3).
//     AddSimilaritiesToMessagesWithPrefix(
//       "Use this documentation as context: \n\n",
//       "context-1",
//     ).
//     User("Explain these concepts")
//
// Common use cases:
// - Adding formatted context to conversations
// - Providing instructions with context
// - Creating structured system messages
// - Customizing RAG output format
func (s *Squawk) AddSimilaritiesToMessagesWithPrefix(prefix string, optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "system", Content: prefix + embeddings.GenerateContextFromSimilarities(s.similarities), Label: optionalParameters[0]},
		)

	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "system", Content: prefix + embeddings.GenerateContextFromSimilarities(s.similarities)})
	}

	return s
}

// Similarities returns the vector records from the most recent similarity search.
// This method provides access to the search results including similarity scores
// and content.
//
// Returns:
//   - []llm.VectorRecord: A slice of vector records containing search results
//     Each record includes:
//     - Content: The text content of the similar document
//     - Score: The similarity score (0.0 to 1.0)
//     - Vector: The embedding vector
//     - Id: Unique identifier for the record
//
// Example with Ollama provider:
//   squawk := New().
//     Model("mistral:latest").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings([]string{
//       "Go is a statically typed language",
//       "Go supports concurrent programming",
//     }).
//     SimilaritySearch("concurrent features", 0.7, 3)
//
//   /* Access search results */
//   similarities := squawk.Similarities()
//   for _, record := range similarities {
//     fmt.Printf("Score: %.2f, Content: %s\n", 
//       record.Score, 
//       record.Content,
//     )
//   }
func (s *Squawk) Similarities() []llm.VectorRecord {
	return s.similarities
}

// ContextFromSimilarities generates a formatted string containing the content from
// similarity search results. This method is useful when you need direct access to
// the context text without adding it to the conversation messages.
//
// Returns:
//   - string: A formatted string containing the content from similar documents
//
// Example with Ollama provider:
//   squawk := New().
//     Model("mistral:latest").
//     EmbeddingsModel("nomic-embed-text").
//     Provider(provider.Ollama).
//     Store(embeddings.NewMemoryVectorStore()).
//     GenerateEmbeddings([]string{
//       "Go is a statically typed language",
//       "Go supports concurrent programming",
//     }).
//     SimilaritySearch("concurrent features", 0.7, 3)
//   
//   /* Get context as string */
//   context := squawk.ContextFromSimilarities()
//   fmt.Println("Retrieved context:", context)
func (s *Squawk) ContextFromSimilarities() string {
	return embeddings.GenerateContextFromSimilarities(s.similarities)
}

