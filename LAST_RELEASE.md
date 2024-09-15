# Parakeet Releases

## Release notes

### v0.2.0 🍕 [pizza]

#### What's new in v0.2.0?

##### New way to set the options

**Problem**: 
The `omitempty` tag prevents a field from being serialised if its value is the zero value for the field's type (e.g., 0.0 for float64).

That means when `Temperature` equals `0.0`, the field is not serialised (then Ollama will use the `Temperature` default value, which equals `0.8`).

The problem will happen for every value equal to `0` or `0.0`

**Solution(s)**:

###### Set all the fields:

```golang
options := Options{
		NumPredict: -1,

		NumKeep:          4,
		Temperature:      0.8,
		TopK:             40,
		TopP:             0.9,
		TFSZ:             1.0,
		TypicalP:         1.0,
		RepeatLastN:      64,
		RepeatPenalty:    1.1,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
		Mirostat:         0,
		MirostatTau:      5.0,
		MirostatEta:      0.1,
		PenalizeNewline:  true,
		Seed:             -1,
}
```


###### Default Options + overriding:

```golang
options := llm.DefaultOptions()
// override the default value
options.Temperature = 0.5
```

###### Use the `SetOptions` helper:

Define only the fields you want to override:
```golang
options := llm.SetOptions(map[string]interface{}{
  "Temperature": 0.5,
})
```
The `SetOptions` helper will set the default values for the fields not defined in the map.


Or use the `SetOptions` helper with the `option` enums:
```golang
options := llm.SetOptions(map[string]interface{}{
  option.Temperature: 0.5,
  option.RepeatLastN: 2,
})
```

**Note**: the results should be more accurate.

##### New sample.

- `51-genai-webapp`: GenAI web application demo


### v0.1.9 🌭 [hot dog]

#### What's new in v0.1.9?

- `llm.GetModelsList(url string) (ModelList, int, error)`
- `llm.GetModelsListWithToken(url, tokenHeaderName, tokenHeaderValue string) (ModelList, int, error)`
- `llm.ShowModelInformationWithToken(url, model , tokenHeaderName, tokenHeaderValue string) (ModelInformation, int, error)`
- `llm.PullModelWithToken(url, model, tokenHeaderName, tokenHeaderValue string) (PullResult, int, error)`


### v0.1.8 🍔 [hamburger]

#### What's new in v0.1.8?

- Return a string of the answer at the end of the stream: `func ChatWithOpenAIStream(url string, query llm.OpenAIQuery, onChunk func(llm.OpenAIAnswer) error) (string, error) {}`
- Helper to create Embedding objects with the OpenAI API: `CreateEmbeddingWithOpenAI(url string, query llm.OpenAIQuery4Embedding, id string) (llm.VectorRecord, error) {}`


### v0.1.7 🥯 [bagel]

#### What's new in v0.1.7?

- A website is available at [https://parakeet-nest.github.io/parakeet/](https://parakeet-nest.github.io/parakeet/)
- (Ollama) Generate completion: replace `llm.Query` with `llm.GenQuery` and `llm.Answer` with `llm.GenAnswer` + 🐛 fix
- Add `Suffix` field to `llm.GenQuery`
- OpenAI API Chat completion support (only tested with the `gpt-4o-mini` model):
  - `func ChatWithOpenAI(url string, query llm.OpenAIQuery) (llm.OpenAIAnswer, error) {}`
  - `func ChatWithOpenAIStream(url string, query llm.OpenAIQuery, onChunk func(llm.OpenAIAnswer) error) error {}`
  - Tools: planned.

> Ollama provides **experimental** [compatibility](https://github.com/ollama/ollama/blob/main/docs/openai.md ) with parts of the [OpenAI API](https://platform.openai.com/docs/api-reference). As it's experimental, I prefer to keep the completion methods of Ollama and OpenAI "separated."

- New samples in the `examples` directory:
  - `44-chat-openai`
  - `45-chat-stream-openai`
  - `47-function-calling-xp`: call several tools in the same prompt
  - `48-testing-models`: test models with different prompts
    - `yi-coder/01-completion`: write an algorithm
    - `yi-coder/02-insertion`: find a problem in the code (and fix it)
    - `yi-coder/03-qa`: ask a question about the code
    - `yi-coder/04-gitlab-ci`: explain a CI/CD pipeline
    - `mathstral/01-completion`: solve a math problem

### v0.1.6 🥨 [pretzel]

#### What's new in v0.1.6?

- Move the cosine similarity function to the `similarity` package
- Implement the Jaccard index calculation for text similarity (🚧 experimental)
- Implement the Levenshtein Distance calculation for set similarity (🚧 experimental)
- Renaming of methods in the `tools` package:
  - `tools.GenerateSystemInstructions()` to `tools.tools.GenerateSystemToolsInstructions()`
  - `tools.GenerateContent()` to `tools.GenerateAvailableToolsContent()`
  - `tools.GenerateInstructions()` to `tools.GenerateUserToolsInstructions()`
- Improve function calling in the `tools` package

### v0.1.5 🥖 [baguette]

#### What's new in v0.1.5?

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

