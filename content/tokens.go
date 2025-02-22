package content

import (
	"strings"
	"unicode"
)

// CountTokens counts the number of tokens in a string
// This is a simple implementation that splits on whitespace
func CountTokens(text string) int {
	if len(strings.TrimSpace(text)) == 0 {
		return 0
	}
	return len(strings.Fields(text))
}

// CountTokensAdvanced provides a more sophisticated token counting approach
// that handles punctuation, contractions, and special characters
func CountTokensAdvanced(text string) int {
	if len(strings.TrimSpace(text)) == 0 {
		return 0
	}

	// Preserve contractions but split on other punctuation
	// This regex replaces punctuation with spaces but preserves apostrophes in contractions
	words := strings.FieldsFunc(text, func(r rune) bool {
		// Keep apostrophes that are likely part of contractions
		if r == '\'' {
			return false
		}
		return unicode.IsPunct(r) || unicode.IsSpace(r)
	})

	count := 0
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" {
			count++
		}
	}
	
	return count
}

// EstimateGPTTokens provides a rough estimate of GPT-style tokens
// Note: This is only an approximation as actual GPT tokenization is more complex
func EstimateGPTTokens(text string) int {
	// A very rough approximation: English text averages ~4 characters per token in GPT models
	const avgCharsPerToken = 4
	
	// Count characters excluding whitespace
	charCount := 0
	for _, char := range text {
		if !unicode.IsSpace(char) {
			charCount++
		}
	}
	
	// Add a small constant to account for spaces between words
	wordCount := len(strings.Fields(text))
	
	// Estimate token count
	tokenEstimate := (charCount + wordCount) / avgCharsPerToken
	if tokenEstimate < 1 && len(strings.TrimSpace(text)) > 0 {
		return 1
	}
	
	return tokenEstimate
}