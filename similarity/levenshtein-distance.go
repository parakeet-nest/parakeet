package similarity

// --- Levenshtein distance ---

func min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	} else {
		return c
	}
}

func LevenshteinDistance(str1, str2 string) int {
	m := make([][]int, len(str1)+1)
	for i := range m {
		m[i] = make([]int, len(str2)+1)
	}

	for i := 0; i <= len(str1); i++ {
		for j := 0; j <= len(str2); j++ {
			if i == 0 {
				m[i][j] = j
			} else if j == 0 {
				m[i][j] = i
			} else if str1[i-1] == str2[j-1] {
				m[i][j] = m[i-1][j-1]
			} else {
				m[i][j] = 1 + min(m[i-1][j], m[i][j-1], m[i-1][j-1])
			}
		}
	}

	return m[len(str1)][len(str2)]
}
