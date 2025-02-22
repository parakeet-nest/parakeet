<!-- TOPIC: Embeddings and Vector Stores SUMMARY: This document discusses creating embeddings, storing them in an efficient way using a vector store, and searching for similar embeddings. KEYWORDS: Embeddings, Vector Store, MemoryVectorStore, BBolt -->
# Embeddings

!!! info "📦 `embeddings` package"

## Create embeddings

```golang
embedding, err := embeddings.CreateEmbedding(
	ollamaUrl,
	llm.Query4Embedding{
		Model:  "all-minilm",
		Prompt: "Jean-Luc Picard is a fictional character in the Star Trek franchise.",
	},
	"Picard", // identifier
)
```

## Vector stores

A vector store allows to store and search for embeddings in an efficient way.

### In memory vector store

**Create a store**:
```golang
store := embeddings.MemoryVectorStore{
	Records: make(map[string]llm.VectorRecord),
}
```

**Save embeddings**:
```golang
store.Save(embedding)
```

**Search embeddings**:
```golang
embeddingFromQuestion, err := embeddings.CreateEmbedding(
	ollamaUrl,
	llm.Query4Embedding{
		Model:  "all-minilm",
		Prompt: "Who is Jean-Luc Picard?",
	},
	"question",
)
// find the nearest vector
similarity, _ := store.SearchMaxSimilarity(embeddingFromQuestion)

documentsContent := `<context><doc>` + similarity.Prompt + `</doc></context>`
```

> 👀 you will find a complete example in `examples/08-embeddings`

### Bbolt vector store

**[Bbolt](https://github.com/etcd-io/bbolt)** is an embedded key/value database for Go.

**Create a store, and open an existing store**:
```golang
store := embeddings.BboltVectorStore{}
store.Initialize("../embeddings.db")
```

!!! note
	👀 you will find a complete example in:

    - [examples/09-embeddings-bbolt](https://github.com/parakeet-nest/parakeet/tree/main/examples/09-embeddings-bbolt)
    - [examples/09-embeddings-bbolt/create-embeddings](https://github.com/parakeet-nest/parakeet/tree/main/examples/09-embeddings-bbolt/create-embeddings): create and populate the vector store
    - [examples/09-embeddings-bbolt/use-embeddings](https://github.com/parakeet-nest/parakeet/tree/main/examples/09-embeddings-bbolt/use-embeddings): search similarities in the vector store

<!-- split -->

### Redis vector store

**Create a store, and open an existing store**:
```golang
redisStore := embeddings.RedisVectorStore{}
err := redisStore.Initialize("localhost:6379", "", "chronicles-bucket")

if err != nil {
	log.Fatalln("😡:", err)
}
```

!!! note
	👀 you will find a complete example in:

    - [examples/32-rag-with-redis](https://github.com/parakeet-nest/parakeet/tree/main/examples/32-rag-with-redis)
    - [examples/32-rag-with-redis/create-embeddings](https://github.com/parakeet-nest/parakeet/tree/main/examples/32-rag-with-redis/create-embeddings): create and populate the vector store
    - [examples/32-rag-with-redis/use-embeddings](https://github.com/parakeet-nest/parakeet/tree/main/examples/32-rag-with-redis/use-embeddings): search similarities in the vector store


### Elasticsearch vector store

**Create a store, and open an existing store**:
```golang
cert, _ := os.ReadFile(os.Getenv("ELASTIC_CERT_PATH"))

elasticStore := embeddings.ElasticSearchStore{}
err := elasticStore.Initialize(
	[]string{
		os.Getenv("ELASTIC_ADDRESS"),
	},
	os.Getenv("ELASTIC_USERNAME"),
	os.Getenv("ELASTIC_PASSWORD"),
	cert,
	"chronicles-index",
)
```

!!! note
	👀 you will find a complete example in:

    - [examples/33-rag-with-elastic](https://github.com/parakeet-nest/parakeet/tree/main/examples/33-rag-with-elastic)
    - [examples/33-rag-with-elastic/create-embeddings](https://github.com/parakeet-nest/parakeet/tree/main/examples/33-rag-with-elastic/create-embeddings): create and populate the vector store
    - [examples/33-rag-with-elastic/use-embeddings](https://github.com/parakeet-nest/parakeet/tree/main/examples/33-rag-with-elastic/use-embeddings): search similarities in the vector store

### Daphnia vector store

[Daphnia](https://github.com/sea-monkeys/daphnia) is another one of my projects to create an embedded vector database (useful to experiment).

```golang
// Initialize the vector store
vectorStore := embeddings.DaphniaVectoreStore{}
vectorStore.Initialize("my-data.gob")
```

!!! note
	👀 you will find a complete example in:

    - [examples/65-hyde](https://github.com/parakeet-nest/parakeet/tree/main/examples/65-hyde)
    - [examples/74-rag-with-daphnia](https://github.com/parakeet-nest/parakeet/tree/main/examples/4-rag-with-daphnia)

### Additional data

you can add additional data to a vector record (embedding):

```golang
embedding.Text()
embedding.Reference()
embedding.MetaData()
```

<!-- TOPIC: Natural Language Processing, Computer Programming, Golang SUMMARY: This document discusses the process of creating embeddings from text files and performing similarity searches using a Bolt Vector Store. The document also demonstrates how to use these embeddings to generate context for a chat system. KEYWORDS: Embeddings, Similarity Search, Bolt Vector Store, Natural Language Processing, Computer Programming, Golang -->
## Create embeddings from text files and Similarity search

### Create embeddings
```golang
ollamaUrl := "http://localhost:11434"
embeddingsModel := "all-minilm"

store := embeddings.BboltVectorStore{}
store.Initialize("../embeddings.db")

// Parse all golang source code of the examples
// Create embeddings from documents and save them in the store
counter := 0
_, err := content.ForEachFile("../../examples", ".go", func(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fmt.Println("📝 Creating embedding from:", path)
	counter++
	embedding, err := embeddings.CreateEmbedding(
		ollamaUrl,
		llm.Query4Embedding{
			Model:  embeddingsModel,
			Prompt: string(data),
		},
		strconv.Itoa(counter), // don't forget the id (unique identifier)
	)
	fmt.Println("📦 Created: ", len(embedding.Embedding))

	if err != nil {
		fmt.Println("😡:", err)
	} else {
		_, err := store.Save(embedding)
		if err != nil {
			fmt.Println("😡:", err)
		}
	}
	return nil
})
if err != nil {
	log.Fatalln("😡:", err)
}
```

### Similarity search

```golang
ollamaUrl := "http://localhost:11434"
embeddingsModel := "all-minilm"
chatModel := "magicoder:latest"

store := embeddings.BboltVectorStore{}
store.Initialize("../embeddings.db")

systemContent := `You are a Golang developer and an expert in computer programming.
Please make friendly answer for the noobs. Use the provided context and doc to answer.
Add source code examples if you can.`

// Question for the Chat system
userContent := `How to create a stream chat completion with Parakeet?`

// Create an embedding from the user question
embeddingFromQuestion, err := embeddings.CreateEmbedding(
	ollamaUrl,
	llm.Query4Embedding{
		Model:  embeddingsModel,
		Prompt: userContent,
	},
	"question",
)
if err != nil {
	log.Fatalln("😡:", err)
}
fmt.Println("🔎 searching for similarity...")

similarities, _ := store.SearchSimilarities(embeddingFromQuestion, 0.3)

// Generate the context from the similarities
// This will generate a string with a content like this one:
// `<context><doc>...<doc><doc>...<doc></context>`
documentsContent := embeddings.GenerateContextFromSimilarities(similarities)

fmt.Println("🎉 similarities", len(similarities))

options := llm.SetOptions(map[string]interface{}{
	option.Temperature: 0.4,
	option.RepeatLastN: 2,
})

query := llm.Query{
	Model: chatModel,
	Messages: []llm.Message{
		{Role: "system", Content: systemContent},
		{Role: "system", Content: documentsContent},
		{Role: "user", Content: userContent},
	},
	Options: options,
	Stream: false,
}

fmt.Println("")
fmt.Println("🤖 answer:")

// Answer the question
_, err = completion.ChatStream(ollamaUrl, query,
	func(answer llm.Answer) error {
		fmt.Print(answer.Message.Content)
		return nil
	})

if err != nil {
	log.Fatal("😡:", err)
}
```

### Other similarity search methods

`SearchMaxSimilarity` searches for the vector record in the `BboltVectorStore` that has the maximum **cosine distance similarity** to the given `embeddingFromQuestion`:
```golang
similarity, _ := store.SearchMaxSimilarity(embeddingFromQuestion)
```

`SearchTopNSimilarities` searches for vector records in the `BboltVectorStore` that have a **cosine distance similarity** greater than or equal to the given `limit` and returns the top `n` records:
```golang
similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, limit, n)
```
