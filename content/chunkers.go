package content

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
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

// used by SplitMarkdownBySections and SplitAsciiDocBySections
func splitContentBySectionWithRegex(content string, regexDelimiter string) []string {
	var sections []string
	var currentSection []string

	// Regex to detect Markdown/AsciiDoc titles
	re := regexp.MustCompile(regexDelimiter)

	// use a scanner to read the content line by line
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is a title
		if matches := re.FindStringSubmatch(line); matches != nil {
			// If a new title is found, ends the current section and starts a new section
			if len(currentSection) > 0 {
				sections = append(sections, strings.Join(currentSection, "\n"))
			}
			currentSection = []string{line}
		} else {
			// Adds the content lines to the current section
			currentSection = append(currentSection, line)
		}
	}

	// Add the last section if it exists
	if len(currentSection) > 0 {
		sections = append(sections, strings.Join(currentSection, "\n"))
	}

	return sections
}

func SplitMarkdownBySections(content string) []string {
	return splitContentBySectionWithRegex(content, `^(#+)\s+(.*)`)
}

/*
sections := SplitMarkdownByLevelSections(content, 1) // Split by top-level sections

	for i, section := range sections {
		fmt.Printf("Section %d:\n%s\n\n", i+1, section)
	}

subsections := SplitMarkdownByLevelSections(content, 2) // Split by second-level sections

	for i, subsection := range subsections {
		fmt.Printf("Subsection %d:\n%s\n\n", i+1, subsection)
	}
*/
func SplitMarkdownByLevelSections(content string, level int) []string {
	// Create a regex that matches the specified level of header
	regexDelimiter := fmt.Sprintf(`^(%s)\s+(.*)`, strings.Repeat("#", level))
	return splitContentBySectionWithRegex(content, regexDelimiter)
}

func SplitAsciiDocBySections(content string) []string {
	return splitContentBySectionWithRegex(content, `^(=+|\#+)\s+(.*)`)
}

func SplitHTMLBySections(content string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	var sections []string
	var currentSection []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "h1" || n.Data == "h2" || n.Data == "h3" || n.Data == "h4" || n.Data == "h5" || n.Data == "h6") {
			if len(currentSection) > 0 {
				sections = append(sections, strings.Join(currentSection, ""))
				currentSection = []string{}
			}
		}

		htmlContent := renderNode(n)
		currentSection = append(currentSection, htmlContent)

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	if len(currentSection) > 0 {
		sections = append(sections, strings.Join(currentSection, ""))
	}

	return sections, nil
}

// renderNode convertit un noeud HTML en chaîne de caractères
func renderNode(n *html.Node) string {
	var sb strings.Builder
	html.Render(&sb, n)
	return sb.String()
}
