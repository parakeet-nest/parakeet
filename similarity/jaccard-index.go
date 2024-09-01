package similarity

// --- Jaccard index ---

// check if an item is a part of a set
func contains(set []string, element string) bool {
	for _, s := range set {
		if s == element {
			return true
		}
	}
	return false
}

// https://en.wikipedia.org/wiki/Jaccard_index
func JaccardSimilarityCoeff(set1, set2 []string) float64 {
	intersection := 0
	union := len(set1) + len(set2) - intersection

	for _, element := range set1 {
		if contains(set2, element) {
			intersection++
		}
	}

	return float64(intersection) / float64(union)
}
