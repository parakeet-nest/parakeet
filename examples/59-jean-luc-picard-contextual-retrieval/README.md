# Chat with Jean-Luc Picard using Contextual Retrieval for RAG


## Generate embeddings

This command will generate embeddings for the documents in the `docs` directory using the specified embeddings model and store them in the `embeddings.db` file.

```go
go run main.go create-embeddings \
--embeddings-model mxbai-embed-large:latest  \
--contextual-model qwen2.5:1.5b \
--docs-path ./docs \
--store-path ./embeddings.db \
--url http://localhost:11434
```

## Chat with Jean-Luc

This command will start a chat session with Jean-Luc Picard using the specified chat model and embeddings model.
You can define a context and instructions for the chat session in the `context.md` and `instructions.md` files.

```go
go run main.go chat \
--embeddings-model mxbai-embed-large:latest \
--chat-model qwen2.5:0.5b \
--store-path ./embeddings.db \
--context-path ./context.md \
--instructions-path ./instructions.md \
--url http://localhost:11434
```