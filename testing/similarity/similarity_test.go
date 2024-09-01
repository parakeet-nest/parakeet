package testing_similarity

import (
	"fmt"
	"strings"
	"testing"

	"github.com/parakeet-nest/parakeet/similarity"
)

var docs = []string{
	`Michael Burnham is the main character on the Star Trek series, Discovery.  
	She's a human raised on the logical planet Vulcan by Spock's father.  
	Burnham is intelligent and struggles to balance her human emotions with Vulcan logic.  
	She's become a Starfleet captain known for her determination and problem-solving skills.
	Originally played by actress Sonequa Martin-Green`,

	`James T. Kirk, also known as Captain Kirk, is a fictional character from the Star Trek franchise.  
	He's the iconic captain of the starship USS Enterprise, 
	boldly exploring the galaxy with his crew.  
	Originally played by actor William Shatner, 
	Kirk has appeared in TV series, movies, and other media.`,

	`Jean-Luc Picard is a fictional character in the Star Trek franchise.
	He's most famous for being the captain of the USS Enterprise-D,
	a starship exploring the galaxy in the 24th century.
	Picard is known for his diplomacy, intelligence, and strong moral compass.
	He's been portrayed by actor Patrick Stewart.`,

	`Lieutenant Philippe Charrière, known as the **Silent Sentinel** of the USS Discovery, 
	is the enigmatic programming genius whose codes safeguard the ship's secrets and operations. 
	His swift problem-solving skills are as legendary as the mysterious aura that surrounds him. 
	Charrière, a man of few words, speaks the language of machines with unrivaled fluency, 
	making him the crew's unsung guardian in the cosmos. His best friend is Spiderman from the Marvel Cinematic Universe.`,
}

/*
Not entirely interesting with this use case
But If I use a keyword system, perhaps it could be better 🤔
Should work on small pieces of text
*/
func TestJaccardSimilarityCoeff(m *testing.T) {

	userContent := `Who is Jean-Luc Picard?`

	splittedUserContent := strings.Fields(userContent)

	// Calculate Jaccard index for every document
	// the highest index is related to the best similarity
	for idx, doc := range docs {
		jaccardIndex := similarity.JaccardSimilarityCoeff(splittedUserContent, strings.Fields(doc))
		fmt.Println("-", idx, "Jaccard index:", jaccardIndex)
	}

	expected := 0.0377

	jaccardIndex := similarity.JaccardSimilarityCoeff(splittedUserContent, strings.Fields(docs[2]))
	if jaccardIndex < expected {
		m.Fatal("unexpected result")
	}

	fmt.Println("Splitted User Content", splittedUserContent)
	fmt.Println("Splitted Document Content", strings.Fields(docs[2]))
	fmt.Println("Jaccard Index", jaccardIndex)

}
