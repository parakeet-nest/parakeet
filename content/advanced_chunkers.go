package content

// ✋ this is an experimental package, and it is subject to change in the future.
import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Chunk struct {
	Header         string
	Content        string
	Level          int
	Prefix         string
	ParentLevel    int
	ParentHeader   string
	ParentPrefix   string
	Lineage        string
	SimpleMetaData string                 // Additional metadata if needed
	Metadata       map[string]interface{} // additional metadata
	KeyWords       []string               // Keywords that could be extracted from the content
}

func ParseMarkdown(content string) []*Chunk {
	var chunks []*Chunk
	var currentHeader string
	var currentContent bytes.Buffer

	md := goldmark.New()
	reader := text.NewReader([]byte(content))
	doc := md.Parser().Parse(reader)

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch n.Kind() {
		case ast.KindHeading:
			if entering {
				if currentContent.Len() > 0 {
					chunks = append(chunks, &Chunk{
						Header:  currentHeader,
						Content: currentContent.String(),
					})
					currentContent.Reset()
				}
				heading := n.(*ast.Heading)
				currentHeader = string(heading.Text(reader.Source()))
			}
		case ast.KindParagraph:
			if entering {
				if currentHeader != "" && currentContent.Len() > 0 {
					currentContent.WriteString("\n\n") // Separate paragraphs
				}
			}
		case ast.KindText:
			if entering {
				if currentContent.Len() > 0 && currentContent.Bytes()[currentContent.Len()-1] != ' ' {
					currentContent.WriteString(" ") // Add space before appending new text
				}
				currentContent.Write(n.Text(reader.Source()))
			}
		}

		return ast.WalkContinue, nil
	})

	// Handle any remaining content
	if currentContent.Len() > 0 {
		chunks = append(chunks, &Chunk{
			Header:  currentHeader,
			Content: currentContent.String(),
		})
	}

	return chunks

}

/*
chunks := parseMarkdown(markdown)
for _, chunk := range chunks {
	fmt.Printf("Header: %s\nContent: %s\n\n", chunk.Header, chunk.Content)
}
*/

/*
Steps to chunk a markdown document while maintaining semantic meaning
and preserving the relationship between sections:

### Step 1: Parse the Markdown Document
Use a Markdown parser to convert the document into a structured format (e.g., an Abstract Syntax Tree, or AST).

### Step 2: Traverse the AST
Traverse the parsed AST to identify and extract different sections based on headers (`#`, `##`, etc.).
During this traversal, you should maintain a stack or a similar data structure to keep track
of the current context and hierarchical relationships between sections.

### Step 3: Chunk the Document
As you traverse the AST, you can group content under each header as a chunk. To maintain semantic meaning:
- **Main Sections**: Treat each top-level header (e.g., `# Title`) as a separate chunk.
- **Subsections**: Combine subsections with their parent sections.
  For instance, if you have a `##` header under a `#` header, the content under `##` should be included
  in the same chunk as the `#` header.
- **Size Constraints**: If you're splitting based on size (e.g., token count or character length),
  ensure that sections are not split in a way that disrupts the semantic flow.

### Step 4: Preserve Relationships
Ensure that each chunk maintains a clear relationship with its parent sections.
This can be done by including metadata in each chunk, such as the header level, parent section,
or even the path of headers leading to that chunk.

### Step 5: Output the Chunks
Once the chunks are identified, format them according to your needs, which might be as a JSON array,
each element containing the chunked content and its metadata (e.g., section titles, header levels).

### Explanation:
- **Parsing**: The `parseMarkdown` function uses `goldmark` to parse the Markdown content.
- **Chunking**: It walks through the AST, creating chunks whenever it encounters a header or paragraph.
  The content under each header is grouped into a chunk.
- **Preserving Relations**: The example groups content under the current header, but you could modify it
  to include metadata for more complex relations.

### Final Output:
The output will give you the content grouped under each header, preserving the structure and
meaning of the document. You can further customize this logic to handle different header levels,
include metadata, or chunk based on size constraints.

*/

func ParseMarkdownWithHierarchy(content string) []Chunk {
	lines := strings.Split(content, "\n")
	var chunks []Chunk
	var stack []Chunk

	headerRegex := regexp.MustCompile(`^(#+)\s+(.*)$`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if matches := headerRegex.FindStringSubmatch(line); matches != nil {
			level := len(matches[1])
			header := matches[2]
			prefix := matches[1]

			// Find content for this header
			contentLines := []string{}
			for j := i + 1; j < len(lines); j++ {
				if headerRegex.MatchString(lines[j]) {
					break
				}
				contentLines = append(contentLines, lines[j])
			}
			content := strings.Join(contentLines, "\n")

			// Determine parent header
			var parent Chunk
			for len(stack) > 0 && stack[len(stack)-1].Level >= level {
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				parent = stack[len(stack)-1]
			}

			chunk := Chunk{
				Level:        level,
				Prefix:       prefix,
				Header:       header,
				Content:      strings.TrimSpace(content),
				ParentPrefix: parent.Prefix,
				ParentLevel:  parent.Level,
				ParentHeader: parent.Header,
			}

			chunks = append(chunks, chunk)
			stack = append(stack, chunk)
		}
	}

	return chunks
}

// ParseMarkdownWithLineage parses the given markdown content and returns a slice of Chunk structs.
// Each Chunk represents a header and its associated content, along with its hierarchical lineage.
//
// The function processes the markdown content line by line, identifying headers and their levels
// using a regular expression. It then collects the content associated with each header and
// determines the parent header to build the hierarchical structure.
//
// Parameters:
//   - content: A string containing the markdown content to be parsed.
//
// Returns:
//   - A slice of Chunk structs, each representing a header and its associated content, along with
//     its hierarchical lineage.
func ParseMarkdownWithLineage(content string) []Chunk {
	lines := strings.Split(content, "\n")
	var chunks []Chunk
	var stack []Chunk

	headerRegex := regexp.MustCompile(`^(#+)\s+(.*)$`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if matches := headerRegex.FindStringSubmatch(line); matches != nil {
			level := len(matches[1])
			header := matches[2]
			prefix := matches[1]

			// Find content for this header
			contentLines := []string{}
			for j := i + 1; j < len(lines); j++ {
				if headerRegex.MatchString(lines[j]) {
					break
				}
				contentLines = append(contentLines, lines[j])
			}
			content := strings.Join(contentLines, "\n")

			// Determine parent header
			var parent Chunk
			for len(stack) > 0 && stack[len(stack)-1].Level >= level {
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				parent = stack[len(stack)-1]
			}

			// Build lineage
			lineage := buildLineage(stack, header)

			chunk := Chunk{
				Level:        level,
				Prefix:       prefix,
				Header:       header,
				Content:      strings.TrimSpace(content),
				ParentPrefix: parent.Prefix,
				ParentLevel:  parent.Level,
				ParentHeader: parent.Header,
				Lineage:      lineage,
			}
			//if chunk.Content != "" {
				chunks = append(chunks, chunk)
				stack = append(stack, chunk)
			//}
		}
	}

	return chunks
}

func buildLineage(stack []Chunk, currentHeader string) string {
	var lineage []string
	for _, chunk := range stack {
		lineage = append(lineage, chunk.Header)
	}
	lineage = append(lineage, currentHeader)
	return strings.Join(lineage, " > ")
}

/*
func ChunkToText(chunk Chunk, template string) (string, error) {
    res, err := InterpolateString(template, chunk)
    return res, err
}

`{{.Prefix}} {{.Header}} {{.Level}} -> {{.ParentPrefix}} {{.ParentHeader}} {{.ParentLevel}} -> {{.Lineage}} {{.Content}}`
*/
