package content

import (
	"regexp"
	"strings"
)


// ChunkText takes a text string and divides it into chunks of a specified size with a given overlap.
// It returns a slice of strings, where each string represents a chunk of the original text.
//
// Parameters:
//   - text: The input text to be chunked.
//   - chunkSize: The size of each chunk.
//   - overlap: The amount of overlap between consecutive chunks.
//
// Returns:
//   - []string: A slice of strings representing the chunks of the original text.
func ChunkText(text string, chunkSize, overlap int) []string {
	chunks := []string{}
	for start := 0; start < len(text); start += chunkSize - overlap {
		end := start + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[start:end])
	}
	return chunks
}


// SplitTextWithDelimiter splits the given text using the specified delimiter and returns a slice of strings.
//
// Parameters:
//   - text: The text to be split.
//   - delimiter: The delimiter used to split the text.
//
// Returns:
//   - []string: A slice of strings containing the split parts of the text.
func SplitTextWithDelimiter(text string, delimiter string) []string {
	return strings.Split(text, delimiter)
}


// SplitTextWithRegex splits the given text using the provided regular expression delimiter.
// It returns a slice of strings containing the split parts of the text.
func SplitTextWithRegex(text string, regexDelimiter string) []string {
	result := regexp.MustCompile(regexDelimiter) 
	return result.Split(text, -1)
}
//TODO: split before or after?