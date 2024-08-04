package embeddings

import (
	"fmt"

	"github.com/parakeet-nest/parakeet/llm"
)

// GenerateContextFromSimilarities generates the context content from a slice of vector records.
//
// Parameters:
// - similarities: a slice of llm.VectorRecord representing the similarities.
//
// Returns:
// - string: the generated context content in XML format.
func GenerateContextFromSimilarities(similarities []llm.VectorRecord) string {
	documentsContent := "<context>\n"
	for _, similarity := range similarities {
		documentsContent += fmt.Sprintf("<doc>%s</doc>\n", similarity.Prompt)
	}
	documentsContent += "</context>"
	return documentsContent
}

// TODO: GenerateContextFromSimilaritiesWithTags

func GenerateContentFromSimilarities(similarities []llm.VectorRecord) string {
	documentsContent := ""
	for _, similarity := range similarities {
		documentsContent += fmt.Sprintf("%s\n", similarity.Prompt)
	}
	return documentsContent
}