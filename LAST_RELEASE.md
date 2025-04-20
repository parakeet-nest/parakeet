# Parakeet Releases

## Release notes

### v0.2.8 🍩 [doughnut]

- `history.RemoveTopMessage() error`: removes the oldest message from the Messages list.

## v0.2.7 🐳 [spouting whale]

Addition of **Docker Model Runner** support (and OpenAI at the same time) allowing easy development of generative AI applications in Docker containers.

> - All the examples OpenAI related have been udated

## v0.2.6 🍿 [popcorn]

- MCP support progress: 
  - SSE transport Client
  - new STDIO transport Client

- Added a SSE MCP example using the [WASImancer MCP server project](https://github.com/sea-monkeys/WASImancer): `75-mcp-sse`
- Update of the STDIO MCP example: `67-mcp`

## v0.2.5 🥧 [pie] 

### Helpers

#### Estimate the number of tokens in a text

- `content.CountTokens(text string) int`
- `content.CountTokensAdvanced(text string) int`
- `content.EstimateGPTTokens(text string) int`

> this could be useful to estimate the value of `num_ctx`

#### Extract elements from source code:

- `source.ExtractCodeElements(fileContent string, language string) ([]CodeElement, error)`

```golang
// CodeElement represents a code structure element (class, function, method)
type CodeElement struct {
	Type        string // "class", "function", "method"
	Name        string
	Signature   string
	Description string
	LineNumber  int
	ParentClass string // For methods
	Parameters  []string
	Source      string // Source code of the element
}
```

> the `Signature` could be useful to add context to embeddings.

#### Get and cast environment variable value at the same time:

- `gear.GetEnvFloat(key string, defaultValue float64) float64 `
- `gear.GetEnvInt(key string, defaultValue int) int`
- `gear.GetEnvString(key string, defaultValue string) string`

### Conversational history + new samples

#### In Memory

- Added `history.RemoveMessage(id string)`
  - see example: `69-web-chat-bot`
- Added `history.SaveMessageWithSession(sessionId string, messagesCounters *map[string]int, message llm.Message)`
  - see example: `70-web-chat-bot-with-session`
- Added `history.RemoveTopMessageOfSession(sessionId string, messagesCounters *map[string]int, conversationLength int)`
  - see example: `70-web-chat-bot-with-session`

#### Bbolt Memory

- Added `history.RemoveMessage(id string)`
- Added `history.SaveMessageWithSession(sessionId string, messagesCounters *map[string]int, message llm.Message)`
  - see example: `71-web-chat-bot-with-session`
- Added `history.RemoveTopMessageOfSession(sessionId string, messagesCounters *map[string]int, conversationLength int)`
  - see example: `71-web-chat-bot-with-session`

## v0.2.4 🥮 [mooncake]

### RAG

Improving the RAG example with Elasticsearch: `40-rag-with-elastic-markdown` (🙏 Thank you [@codefromthecrypt](https://github.com/codefromthecrypt))

### New examples:
  
- Structured output: `66-structured-outputs`
- Experiments with Hypothetical Document Embeddings (HyDE): `65-hyde` (🚧 this is a work in progress)
- MCP Client: `67-mcp`
- How to use DeepSeek R1 (`1.5b`): `68-deepseek-r1`

### Error management

#### ModelNotFoundError

```golang
// package completion
type ModelNotFoundError struct {
  Code    int
  Message string
  Model   string
}
```

**Usage**:
```golang
answer, err := completion.Chat(ollamaUrl, query)
if err != nil {
  // test if the model is not found
  if modelErr, ok := err.(*completion.ModelNotFoundError); ok {
    fmt.Printf("💥 Got Model Not Found error: %s\n", modelErr.Message)
    fmt.Printf("😡 Error code: %d\n", modelErr.Code)
    fmt.Printf("🧠 Expected Model: %s\n", modelErr.Model)
  } else {
    log.Fatal("😡:", err)
  }
}
```
> See these examples: `04-chat-stream` and `66-structured-outputs`

#### NoSuchOllamaHostError

```golang
// package completion
type NoSuchOllamaHostError struct {
	Host string
	Message string
}
```

**Usage**:
```golang
if noHostErr, ok := err.(*completion.NoSuchOllamaHostError); ok {
  fmt.Printf("🦙 Got No Such Ollama Host error: %s\n", noHostErr.Message)
  fmt.Printf("🌍 Expected Host: %s\n", noHostErr.Host)
}
```

### First MCP support

Integration of `github.com/mark3labs/mcp-go/mcp` and `github.com/mark3labs/mcp-go/client` (this is a work in progress 🚧)

#### Helpers

- `mcphelpers.GetMCPClient(ctx context.Context, command string, env []string, args ...string) (*client.StdioMCPClient, *mcp.InitializeResult, error)`
- `mcphelpers.GetTools(mcpClient *client.StdioMCPClient) ([]llm.Tool, error)`
- `tools.ConvertMCPTools` to convert the MCP tools list to a list compliant with Ollama LLM tools. (used by `GetTools`)
- `mcphelpers.CallTool(ctx context.Context, mcpClient *client.StdioMCPClient, functionName string, arguments map[string]interface{}) (*mcp.CallToolResult, error)`
- `mcphelpers.GetTextFromResult(mcpResult *mcp.CallToolResult) (string, error)`

> See this example: `67-mcp` (an example of a MCP server is provided)

#### Error management (specific type errors)

- `MCPClientCreationError`
- `MCPClientInitializationError`
- `MCPGetToolsError`
- `MCPToolCallError`
- `MCPResultExtractionError`


## v0.2.3 🥧 [pie]

#### What's new in v0.2.3?

Update of the Extism dependency.

## v0.2.2 🧁 [cupcake]

#### What's new in v0.2.2?

##### Flock Agents

> Inspired by: [Swarm by OpenAI](https://github.com/openai/swarm)

Flock is a Parakeet package for creating and managing AI agents using the Ollama backend. It provides a simple way to create conversational agents, orchestrate interactions between them, and implement function calling capabilities.

📝 [Documentation](https://github.com/parakeet-nest/parakeet/blob/main/docs/flock-agents.md)

## v0.2.1 🧇 [waffle]

#### What's new in v0.2.1?

##### Contextual Retrieval

> Inspired by: [Introducing Contextual Retrieval](https://www.anthropic.com/news/contextual-retrieval)

2 new methods are available in the `content` package:

- `CreateChunkContext`
- `CreateChunkContextWithPromptTemplate`

`CreateChunkContext` generates a succinct context for a given chunk within the whole document content.
This context is intended to improve search retrieval of the chunk.


`CreateChunkContextWithPromptTemplate` generates a contextual response based on a given prompt template and document content.
It interpolates the template with the provided document and chunk content, then uses an LLM to generate a response.


##### UI Helpers

2 new methods are available in the `ui` package:


If you use Parakeet to create CLI applications, you can use the `ui` package to create a (very) simple UI.

- `Input`
- `Println`

`Input` displays a prompt with the specified color and waits for user input.

`Println` prints the provided strings with the specified color using the lipgloss styling library.


##### CLI Helpers

8 new methods are available in the `cli` package:

- `Settings` parses command-line arguments and flags.
- `FlagValue` retrieves the value of a flag by its name from a slice of Flag structs.
- `HasArg` checks if an argument with the specified name exists in the provided slice of arguments.
- `HasFlag` checks if a flag with the specified name exists in the provided slice of flags.
- `ArgsTail` extracts the names from a slice of Arg structs and returns them as a slice of strings.
- `FlagsTail` takes a slice of Flag structs and returns a slice of strings containing the names of those flags.
- `FlagsWithNamesTail` takes a slice of Flag structs and returns a slice of strings, where each string is a formatted pair of the flag's name and value in the form "name=value".
- `HasSubsequence` checks if the given subsequence of strings (subSeq) is present in the tail of the provided arguments (args).

**Example**:

```go
// default values
ollamaUrl := "http://localhost:11434"
chatModel := "llama3.1:8b"
embeddingsModel := "bge-m3:latest"

args, flags := cli.Settings()

if cli.HasFlag("url", flags) {
    ollamaUrl = cli.FlagValue("url", flags)
}

if cli.HasFlag("chat-model", flags) {
    chatModel = cli.FlagValue("chat-model", flags)
}

if cli.HasFlag("embeddings-model", flags) {
    embeddingsModel = cli.FlagValue("embeddings-model", flags)
}

switch cmd := cli.ArgsTail(args); cmd[0] {
case "create-embeddings":
    fmt.Println(embeddingsModel)
case "chat":
    fmt.Println(chatModel)
default:
    fmt.Println("Unknown command:", cmd[0])
}
```

##### New samples

- 52-constraints: Preventing an LLM from talking about certain things
- 53-constraints: Preventing an LLM from talking about certain things
- 54-constraints-webapp: Preventing an LLM from talking about certain things
- 55-create-npc: Create a NPC with `nemotron-mini` and chat with him
- 56-jean-luc-picard: Chat with Jean-Luc Picard
- 57-jean-luc-picard-rag: Chat with Jean-Luc Picard + RAG
- 58-michael-burnham: Chat with Michael Burnham
- 59-jean-luc-picard-contextual-retrieval: Chat with Jean-Luc Picard + Contextual Retrieval
- 60-safety-models: Safety Models fine-tuned for content safety classification of LLM inputs and responses

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

