package similarity
// 🚧 experimental
// --- Levenshtein distance ---
// LevenshteinDistance calculates the Levenshtein distance between two strings.
// For comparisons, the lowest distance is related to the best similarity
// ✋ This method is not precise for large strings.
func LevenshteinDistance(a, b string) int {
	// Convert strings to slices of runes to handle multibyte characters
	runeA := []rune(a)
	runeB := []rune(b)

	lenA := len(runeA)
	lenB := len(runeB)

	// Initialize a 2D slice to store distances
	dp := make([][]int, lenA+1)
	for i := range dp {
		dp[i] = make([]int, lenB+1)
	}

	// Initialize the distance matrix
	for i := 0; i <= lenA; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= lenB; j++ {
		dp[0][j] = j
	}

	// Compute the distances
	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			cost := 0
			if runeA[i-1] != runeB[j-1] {
				cost = 1
			}
			dp[i][j] = min(dp[i-1][j]+1, dp[i][j-1]+1, dp[i-1][j-1]+cost)
		}
	}

	return dp[lenA][lenB]
}

// Helper function to find the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
