# Contextual Retrieval

!!! info "📦 `content` package"

!!! note "Inspired by: [Introducing Contextual Retrieval](https://www.anthropic.com/news/contextual-retrieval)"

## `CreateChunkContext`

`CreateChunkContext` generates a succinct context for a given chunk within the whole document content.
This context is intended to improve search retrieval of the chunk.

**Parameters**:
  - wholeDocumentContent: The entire content of the document as a string.
  - chunk: The specific chunk of the document for which context is to be generated.
  - ollamaUrl: The URL for the Ollama service.
  - contextualModel: The model used for generating the context.
  - options: Additional options for the LLM (Language Model).

**Returns**:
  - A string containing the succinct context for the chunk.
  - An error if the context generation fails.

> `CreateChunkContext` use a default template for the prompt. If you want to use a custom template, you can use `CreateChunkContextWithPromptTemplate`.

## `CreateChunkContextWithPromptTemplate`

`CreateChunkContextWithPromptTemplate` generates a contextual response based on a given prompt template and document content.
It interpolates the template with the provided document and chunk content, then uses an LLM to generate a response.

**Parameters**:
  - promptTpl: A string template for the prompt.
  - wholeDocumentContent: The content of the entire document.
  - chunk: A Chunk struct containing a portion of the document content.
  - ollamaUrl: The URL of the LLM service.
  - contextualModel: The model to be used for generating the response.
  - options: Options for the LLM query.

**Returns**:
  - A string containing the generated response.
  - An error if the process fails at any step.

### Template example

```golang
promptTemplateForContext := `
<chunk> 
{{.chunkContent}} 
</chunk> 
Generate a brief context of the above chunk to situate this chunk within the below document. 
<document> 
{{.wholeDocument}} 
</document> 
`
```


!!! note
	👀 you will find a complete example in:

    - [examples/59-jean-luc-picard-contextual-retrieval](https://github.com/parakeet-nest/parakeet/tree/main/examples/59-jean-luc-picard-contextual-retrieval)

