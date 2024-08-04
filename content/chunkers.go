package content

import (
	"regexp"
	"strings"
)

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

func SplitTextWithDelimiter(text string, delimiter string) []string {
	return strings.Split(text, delimiter)
}

func SplitTextWithRegex(text string, regexDelimiter string) []string {
	result := regexp.MustCompile(regexDelimiter) 
	return result.Split(text, -1)

}
//TODO: split before or after?