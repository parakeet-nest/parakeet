package content

import "fmt"

// GenerateContextFromDocs generates the context content from a slice of documents.
//
// Remarks: you can use the generated content to add context to a prompt for an LLM.
//
// Parameters:
// - docs: a slice of strings representing the documents.
//
// Returns:
// - string: the generated context content in XML format.
func GenerateContextFromDocs(docs []string) string {

	documentsContent := "<context>\n"
	for _, doc := range docs {
		documentsContent += fmt.Sprintf("<doc>%s</doc>\n", doc)
	}
	documentsContent += "</context>"
	return documentsContent
}
/*
This is a Go function called GenerateContextFromDocs that takes a slice of strings as input and returns a string in XML format. 
The function generates the context content from a slice of documents by iterating over each document in the slice and appending it to a string in the format <doc>document content</doc>. 
Finally, the function wraps the entire content in <context> tags and returns it.
*/

// TODO: GenerateContextWithTags