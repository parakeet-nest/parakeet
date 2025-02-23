package source

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

//TODO: language detection

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

// LanguagePatterns contains regex patterns for different languages
type LanguagePatterns struct {
	ClassPattern    *regexp.Regexp
	FunctionPattern *regexp.Regexp
	MethodPattern   *regexp.Regexp
	CommentPattern  *regexp.Regexp
}

// GetLanguagePatterns returns regex patterns for a specific language
func GetLanguagePatterns(language string) LanguagePatterns {
	switch strings.ToLower(language) {
	case "rust":
		return LanguagePatterns{
			// Matches both struct and trait declarations
			ClassPattern: regexp.MustCompile(`^\s*(?:struct|trait)\s+([A-Za-z0-9_]+)(?:<[^>]+>)?(?:\s*(?:where\s+[^{]+)?)\s*(?:\{|;)`),
			// Matches standalone functions
			FunctionPattern: regexp.MustCompile(`^\s*(?:pub\s+)?(?:async\s+)?fn\s+([A-Za-z0-9_]+)\s*(?:<[^>]+>)?\s*\((.*)\)(?:\s*->\s*[^{]+)?\s*(?:where\s+[^{]+)?\s*\{`),
			// Matches impl block methods
			MethodPattern: regexp.MustCompile(`^\s*(?:pub\s+)?(?:async\s+)?fn\s+([A-Za-z0-9_]+)\s*(?:<[^>]+>)?\s*\(&(?:mut\s+)?self(?:\s*,\s*(.*))?\)(?:\s*->\s*[^{]+)?\s*(?:where\s+[^{]+)?\s*\{`),
			// Matches both /// and /** ... */ style comments
			CommentPattern: regexp.MustCompile(`(?s)^\s*(?:///|/\*\*)(.*?)(?:\*/)?$`),
		}
	case "python":
		return LanguagePatterns{
			ClassPattern:    regexp.MustCompile(`^\s*class\s+([A-Za-z0-9_]+)(?:\(([^)]*)\))?:`),
			FunctionPattern: regexp.MustCompile(`^\s*def\s+([A-Za-z0-9_]+)\s*\((.*)\)\s*(?:->.*?)?:`),
			MethodPattern:   regexp.MustCompile(`^\s+def\s+([A-Za-z0-9_]+)\s*\((self|cls)(.*)\)\s*(?:->.*?)?:`),
			CommentPattern:  regexp.MustCompile(`(?s)^\s*(?:"""|''')(.*?)(?:"""|''')`),
		}
	case "javascript":
		return LanguagePatterns{
			ClassPattern:    regexp.MustCompile(`^\s*class\s+([A-Za-z0-9_]+)(?:\s+extends\s+([A-Za-z0-9_]+))?\s*\{`),
			FunctionPattern: regexp.MustCompile(`^\s*(?:function|const|let|var)\s+([A-Za-z0-9_]+)\s*=?\s*(?:function)?\s*\((.*)\)`),
			MethodPattern:   regexp.MustCompile(`^\s+(?:async\s+)?([A-Za-z0-9_]+)\s*\((.*)\)\s*\{`),
			CommentPattern:  regexp.MustCompile(`(?s)^\s*\/\*\*(.*?)\*\/`),
		}
	case "go":
		return LanguagePatterns{
			ClassPattern:    regexp.MustCompile(`^\s*type\s+([A-Za-z0-9_]+)\s+struct\s*\{`),
			FunctionPattern: regexp.MustCompile(`^\s*func\s+([A-Za-z0-9_]+)\s*\((.*)\)\s*(?:\(?(.*?)\)?)?\s*\{`),
			MethodPattern:   regexp.MustCompile(`^\s*func\s+\(\s*([A-Za-z0-9_]+)\s+\*?([A-Za-z0-9_]+)\s*\)\s+([A-Za-z0-9_]+)\s*\((.*)\)\s*(?:\(?(.*?)\)?)?\s*\{`),
			CommentPattern:  regexp.MustCompile(`^\s*\/\/\s*(.*)`),
		}
	case "java":
		return LanguagePatterns{
			ClassPattern:    regexp.MustCompile(`^\s*(?:public|private|protected)?\s*(?:final|abstract)?\s*class\s+([A-Za-z0-9_]+)(?:\s+extends\s+([A-Za-z0-9_]+))?(?:\s+implements\s+([A-Za-z0-9_,\s]+))?\s*\{`),
			FunctionPattern: regexp.MustCompile(`^\s*(?:public|private|protected)?\s*(?:static)?\s*(?:[A-Za-z0-9_<>[\],\s]+)\s+([A-Za-z0-9_]+)\s*\((.*)\)\s*(?:throws\s+([A-Za-z0-9_,\s]+))?\s*\{`),
			MethodPattern:   regexp.MustCompile(`^\s*(?:public|private|protected)?\s*(?:static)?\s*(?:[A-Za-z0-9_<>[\],\s]+)\s+([A-Za-z0-9_]+)\s*\((.*)\)\s*(?:throws\s+([A-Za-z0-9_,\s]+))?\s*\{`),
			CommentPattern:  regexp.MustCompile(`(?s)^\s*\/\*\*(.*?)\*\/`),
		}
	// Add more languages as needed
	default:
		// Generic patterns
		return LanguagePatterns{
			ClassPattern:    regexp.MustCompile(`^\s*(?:class|struct)\s+([A-Za-z0-9_]+)`),
			FunctionPattern: regexp.MustCompile(`^\s*(?:function|func|def)\s+([A-Za-z0-9_]+)\s*\((.*)\)`),
			MethodPattern:   regexp.MustCompile(`^\s*(?:function|func|def)\s+([A-Za-z0-9_]+)\s*\((.*)\)`),
			CommentPattern:  regexp.MustCompile(`^\s*(?:\/\/|#)\s*(.*)`),
		}
	}
}

// ExtractCodeElements extracts code elements from a string containing source code
func ExtractCodeElements(fileContent string, language string) ([]CodeElement, error) {
	// Split content into lines for processing
	allLines := strings.Split(fileContent, "\n")

	patterns := GetLanguagePatterns(language)
	var elements []CodeElement
	var currentClass string
	var lineNumber int
	var commentBuffer string
	var currentElementStartLine int
	var inElement bool

	// Track element scope
	var bracketLevel int

	scanner := bufio.NewScanner(strings.NewReader(fileContent))
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Track indentation and scope
		trimmedLine := strings.TrimSpace(line)

		// Track bracket level for languages that use them
		if strings.Contains(trimmedLine, "{") {
			bracketLevel += strings.Count(trimmedLine, "{")
		}
		if strings.Contains(trimmedLine, "}") {
			bracketLevel -= strings.Count(trimmedLine, "}")
		}

		// Check for comments
		if match := patterns.CommentPattern.FindStringSubmatch(line); match != nil {
			commentBuffer += strings.TrimSpace(match[1]) + " "
			continue
		}

		// Check for classes
		if match := patterns.ClassPattern.FindStringSubmatch(line); match != nil {
			className := match[1]
			currentClass = className
			// Start tracking a new element
			if inElement {
				// Extract source for the previous element
				extractSourceForLastElement(&elements, allLines, currentElementStartLine, lineNumber-1, language)
			}

			currentElementStartLine = lineNumber
			inElement = true

			elements = append(elements, CodeElement{
				Type:        "class",
				Name:        className,
				Signature:   strings.TrimSpace(line),
				Description: strings.TrimSpace(commentBuffer),
				LineNumber:  lineNumber,
			})
			commentBuffer = ""
			continue
		}

		// Check for methods (if inside a class)
		if currentClass != "" {
			if match := patterns.MethodPattern.FindStringSubmatch(line); match != nil {
				methodName := match[1]
				var params []string
				if len(match) > 2 {
					paramStr := strings.TrimSpace(match[2])
					if paramStr != "" {
						params = splitParameters(paramStr)
					}
				}
				elements = append(elements, CodeElement{
					Type:        "method",
					Name:        methodName,
					Signature:   strings.TrimSpace(line),
					Description: strings.TrimSpace(commentBuffer),
					LineNumber:  lineNumber,
					ParentClass: currentClass,
					Parameters:  params,
				})
				commentBuffer = ""
				continue
			}
		}

		// Check for functions
		if match := patterns.FunctionPattern.FindStringSubmatch(line); match != nil {
			functionName := match[1]
			var params []string
			if len(match) > 2 {
				paramStr := strings.TrimSpace(match[2])
				if paramStr != "" {
					params = splitParameters(paramStr)
				}
			}
			// Start tracking a new element
			if inElement {
				// Extract source for the previous element
				extractSourceForLastElement(&elements, allLines, currentElementStartLine, lineNumber-1, language)
			}

			currentElementStartLine = lineNumber
			inElement = true

			elements = append(elements, CodeElement{
				Type:        "function",
				Name:        functionName,
				Signature:   strings.TrimSpace(line),
				Description: strings.TrimSpace(commentBuffer),
				LineNumber:  lineNumber,
				Parameters:  params,
			})
			commentBuffer = ""
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Extract source for the last element if any
	if inElement {
		extractSourceForLastElement(&elements, allLines, currentElementStartLine, len(allLines), language)
	}

	return elements, nil
}

// splitParameters splits a parameter string into individual parameters
func splitParameters(paramStr string) []string {
	var params []string
	var buffer string
	parenthesesCount := 0
	bracketCount := 0

	for i := 0; i < len(paramStr); i++ {
		char := paramStr[i]
		switch char {
		case '(':
			parenthesesCount++
			buffer += string(char)
		case ')':
			parenthesesCount--
			buffer += string(char)
		case '[':
			bracketCount++
			buffer += string(char)
		case ']':
			bracketCount--
			buffer += string(char)
		case ',':
			if parenthesesCount == 0 && bracketCount == 0 {
				params = append(params, strings.TrimSpace(buffer))
				buffer = ""
			} else {
				buffer += string(char)
			}
		default:
			buffer += string(char)
		}
	}

	if buffer != "" {
		params = append(params, strings.TrimSpace(buffer))
	}

	return params
}

// PrintElements prints extracted code elements
func PrintElements(elements []CodeElement) {
	for _, elem := range elements {
		fmt.Printf("Type: %s\n", elem.Type)
		fmt.Printf("Name: %s\n", elem.Name)
		fmt.Printf("Line: %d\n", elem.LineNumber)
		fmt.Printf("Signature: %s\n", elem.Signature)

		if elem.ParentClass != "" {
			fmt.Printf("Parent Class: %s\n", elem.ParentClass)
		}

		if len(elem.Parameters) > 0 {
			fmt.Printf("Parameters: %v\n", elem.Parameters)
		}

		if elem.Description != "" {
			fmt.Printf("Description: %s\n", elem.Description)
		}

		sourcePreview := elem.Source
		if len(sourcePreview) > 100 {
			sourcePreview = sourcePreview[:100] + "..."
		}
		fmt.Printf("Source: %s\n", sourcePreview)

		fmt.Println(strings.Repeat("-", 50))
	}
}

// extractSourceForLastElement extracts the source code for the last element in the list
func extractSourceForLastElement(elements *[]CodeElement, allLines []string, startLine, endLine int, language string) {
	if len(*elements) == 0 {
		return
	}

	// Determine the actual end line based on language and structure
	actualEndLine := determineElementEndLine(allLines, startLine, endLine, language, (*elements)[len(*elements)-1].Type)

	// Extract source lines
	sourceLines := allLines[startLine-1 : actualEndLine]
	source := strings.Join(sourceLines, "\n")

	// Set the source for the last element
	lastIdx := len(*elements) - 1
	(*elements)[lastIdx].Source = source
}

// determineElementEndLine finds the actual end line of a code element
func determineElementEndLine(lines []string, startLine, maxEndLine int, language string, elementType string) int {
	// Default to maxEndLine if we can't determine better
	endLine := maxEndLine
	if endLine > len(lines) {
		endLine = len(lines)
	}

	// Different strategies based on language
	switch strings.ToLower(language) {
	case "python":
		// Python uses indentation
		baseIndent := getIndentationLevel(lines[startLine-1])

		// Start from the line after the definition
		for i := startLine; i < endLine; i++ {
			// Skip empty lines
			if strings.TrimSpace(lines[i]) == "" {
				continue
			}

			currentIndent := getIndentationLevel(lines[i])
			// If we find a line with same or less indentation, we've exited the element
			if currentIndent <= baseIndent && i > startLine {
				return i
			}
		}

	case "go", "java", "javascript", "c", "cpp", "rust":
		// Bracket-based languages
		bracketLevel := 0
		foundOpeningBracket := false

		for i := startLine - 1; i < endLine; i++ {
			line := lines[i]

			// Count brackets
			bracketLevel += strings.Count(line, "{")
			bracketLevel -= strings.Count(line, "}")

			// Once we find the opening bracket, start tracking
			if strings.Contains(line, "{") && !foundOpeningBracket {
				foundOpeningBracket = true
			}

			// If brackets are balanced and we found at least one, we're done
			if foundOpeningBracket && bracketLevel == 0 && i > startLine-1 {
				return i + 1
			}
		}
	}

	return endLine
}

// getIndentationLevel counts the leading whitespace to determine indentation
func getIndentationLevel(line string) int {
	return len(line) - len(strings.TrimLeft(line, " \t"))
}

// SaveToJSON saves extracted elements to a JSON file
func SaveToJSON(elements []CodeElement, outputPath string) error {
	jsonData, err := json.MarshalIndent(elements, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}

// ExtractCodeElementsFromFile is a utility function that extracts code elements from a file
// This preserves compatibility with code that expects to work with files
func ExtractCodeElementsFromFile(filePath string, language string) ([]CodeElement, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return ExtractCodeElements(string(fileContent), language)
}
