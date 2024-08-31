# Parakeet Releases

## Release notes

### ### v0.1.5 🥖 [baguette]

#### What's new in the next release?

- `tools.GenerateSystemInstructions() string` generates a string containing the system content instructions for using "function calling". (✋ Use it only if the LLM does not implement function calling).
- `content.SplitMarkdownByLevelSections(content string, level int) []string` allows choosing the level of the section you want to split
- `content.ParseMarkdown(content string) []*Chunk` chunk a markdown document. (🚧 experimental)
- `content.ParseMarkdownWithLineage(content string) []Chunk` chunk a markdown document while maintaining semantic meaning and preserving the relationship between sections.
- New types: `QA`, `IO` and `Card` (🚧 experimental, used to create prompt, context, datasets...)
- Unit tests in progress

### v0.1.4 🥐 [croissant]

#### What's new in v0.1.4?

New split methods (to create document chunks) are available in the `content` package:

- `content.SplitMarkdownBySections()`
- `content.SplitAsciiDocBySections()`
- `content.SplitHTMLBySections()`

### v0.1.3 📚 [books]

#### What's new in v0.1.3?

##### Elastic vector store

The ElasticSearch and Kibana services are now started with Docker Compose. The certificates are generated and stored in the `certs` directory.

**Start Elasticsearch and Kibana**
```bash
docker compose up -d
```

> 👀 you will find a complete example in `examples/33-rag-with-elastic`


### v0.1.2 📕 [red-textbook]

#### What's new in v0.1.2?

##### Elastic vector store

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

> 👀 you will find a complete example in `examples/33-rag-with-elastic`
> - `examples/33-rag-with-elastic/create-embeddings`: create and populate the vector store
> - `examples/33-rag-with-elastic/use-embeddings`: search similarities in the vector store


### v0.1.1 📗 [green-textbook]

#### What's new in v0.1.1?

##### Redis vector store

**Create a store, and open an existing store**:
```golang
redisStore := embeddings.RedisVectorStore{}
err := redisStore.Initialize("localhost:6379", "", "chronicles-bucket")

if err != nil {
	log.Fatalln("😡:", err)
}
```

> 👀 you will find a complete example in `examples/32-rag-with-redis`
> - `examples/32-rag-with-redis/create-embeddings`: create and populate the vector store
> - `examples/32-rag-with-redis/use-embeddings`: search similarities in the vector store

___

### v0.1.0 📘 [blue-book]

#### What's new in v0.1.0?

##### Completion

**Verbose mode**:

```golang
options := llm.Options{
    Temperature: 0.5,
    RepeatLastN: 2,
    RepeatPenalty: 2.0,
    Verbose: true,
}
```
You will get an output like this (with the query and the completion):

```json
[llm/query] {
  "model": "deepseek-coder",
  "messages": [
    {
      "role": "system",
      "content": "You are an expert in computer programming.\n\tPlease make friendly answer for the noobs.\n\tAdd source code examples if you can."
    },
    {
      "role": "user",
      "content": "I need a clear explanation regarding the following question:\n\tCan you create a \"hello world\" program in Golang?\n\tAnd, please, be structured with bullet points"
    }
  ],
  "options": {
    "repeat_last_n": 2,
    "temperature": 0.5,
    "repeat_penalty": 2,
    "Verbose": true
  },
  "stream": false,
  "prompt": "",
  "context": null,
  "tools": null,
  "TokenHeaderName": "",
  "TokenHeaderValue": ""
}

[llm/completion] {
  "model": "deepseek-coder",
  "message": {
    "role": "assistant",
    "content": "Sure, here's a simple \"Hello, World!\" program in Golang.\n\t1. First, you need to have Golang installed on your machine.\n\t2. Open your text editor, and write the following code:\n\t```go\n\tpackage main\n\timport \"fmt\"\n\tfunc main() {\n\t    fmt.Println(\"Hello, World!\")\n\t} \n\t```\n\t3. Save the file with a `.go` extension (like `hello.go`).\n\t4. In your terminal, navigate to the directory containing the `.go` file.\n\t5. Run the program with the command:\n\t```\n\tgo run hello.go\n\t```\n\t6. If everything goes well, you should see \"Hello, World!\" printed in your terminal.\n\t7. If there's an error, you will see the error message.\n\t8. If everything is correct, you'll see \"Hello, World!\" printed in your terminal.\n"
  },
  "done": true,
  "response": "",
  "context": null,
  "created_at": "2024-08-19T05:57:23.979361Z",
  "total_duration": 3361191958,
  "load_duration": 2044932125,
  "prompt_eval_count": 79,
  "prompt_eval_duration": 95034000,
  "eval_count": 222,
  "eval_duration": 1216689000
}
```

##### Vector store

Add additional data to a vector record (embedding):

```golang
embedding.Text()
embedding.Reference()
embedding.MetaData()
```

##### Protected endpoint
If your Ollama endpoint is protected with a header token, you can specify the token like this:

```golang
query := llm.Query{
    Model: model,
    Messages: []llm.Message{
        {Role: "system", Content: systemContent},
        {Role: "user", Content: userContent},
    },
    Options: options,
    TokenHeaderName: "X-TOKEN",
    TokenHeaderValue: "john_doe",
}
```

