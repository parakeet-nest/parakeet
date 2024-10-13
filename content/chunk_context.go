package content

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

// CreateChunkContextWithPromptTemplate generates a contextual response based on a given prompt template and document content.
// It interpolates the template with the provided document and chunk content, then uses an LLM to generate a response.
//
// Parameters:
//   - promptTpl: A string template for the prompt.
//   - wholeDocumentContent: The content of the entire document.
//   - chunk: A Chunk struct containing a portion of the document content.
//   - ollamaUrl: The URL of the LLM service.
//   - contextualModel: The model to be used for generating the response.
//   - options: Options for the LLM query.
//
// Returns:
//   - A string containing the generated response.
//   - An error if the process fails at any step.
func CreateChunkContextWithPromptTemplate(promptTpl, wholeDocumentContent string, chunk Chunk, ollamaUrl, contextualModel string, options llm.Options) (string, error) {

	// Contextual retrieval
	data := map[string]interface{}{
		"wholeDocument": wholeDocumentContent,
		"chunkContent":  chunk.Content,
	}
	contextualPrompt, err := InterpolateString(promptTpl, data)
	if err != nil {
		return "", err
	}

	question := llm.GenQuery{
		Model:   contextualModel,
		Prompt:  contextualPrompt,
		Options: options,
	}
	answer, err := completion.Generate(ollamaUrl, question)
	if err != nil {
		return "", err
	}
	return answer.Response, nil
}


// CreateChunkContext generates a succinct context for a given chunk within the whole document content.
// This context is intended to improve search retrieval of the chunk.
//
// Parameters:
//   - wholeDocumentContent: The entire content of the document as a string.
//   - chunk: The specific chunk of the document for which context is to be generated.
//   - ollamaUrl: The URL for the Ollama service.
//   - contextualModel: The model used for generating the context.
//   - options: Additional options for the LLM (Language Model).
//
// Returns:
//   - A string containing the succinct context for the chunk.
//   - An error if the context generation fails.
func CreateChunkContext(wholeDocumentContent string, chunk Chunk, ollamaUrl, contextualModel string, options llm.Options) (string, error) {
	promptTemplateForContext := `<document> 
	{{.wholeDocument}} 
	</document> 
	Here is the below chunk we want to situate within the above whole document 
	<chunk> 
	{{.chunkContent}} 
	</chunk> 
	Please give a short succinct context to situate this chunk within the overall document for the purposes of improving search retrieval of the chunk. 
	Answer only with the succinct context and nothing else. 
	`
	return CreateChunkContextWithPromptTemplate(promptTemplateForContext, wholeDocumentContent, chunk, ollamaUrl, contextualModel, options)

}
