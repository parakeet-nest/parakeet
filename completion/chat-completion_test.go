package completion

import (
	"fmt"
	"strings"
	"testing"

	"github.com/parakeet-nest/parakeet/llm"
)

var ollamaUrl = "http://localhost:11434"

// if working from a container
// use: "http://host.docker.internal:11434"

var systemContent = `You are an expert of the Star Trek franchise.
Using the provided context, answer the user's question to the best of your ability using only the resources provided.
`

var contextContent = `<context>
<doc>
Michael Burnham is the main character on the Star Trek series, Discovery.  
She's a human raised on the logical planet Vulcan by Spock's father.  
Burnham is intelligent and struggles to balance her human emotions with Vulcan logic.  
She's become a Starfleet captain known for her determination and problem-solving skills.
Originally played by actress Sonequa Martin-Green
</doc>
<doc>
James T. Kirk, also known as Captain Kirk, is a fictional character from the Star Trek franchise.  
He's the iconic captain of the starship USS Enterprise, 
boldly exploring the galaxy with his crew.  
Originally played by actor William Shatner, 
Kirk has appeared in TV series, movies, and other media.
</doc>
<doc>
Jean-Luc Picard is a fictional character in the Star Trek franchise.
He's most famous for being the captain of the USS Enterprise-D,
a starship exploring the galaxy in the 24th century.
Picard is known for his diplomacy, intelligence, and strong moral compass.
He's been portrayed by actor Patrick Stewart.
</doc>
<doc>
Spock is a fictional character from Star Trek. 
He's most famous for being the half-Vulcan, half-human science officer 
and first officer on the starship USS Enterprise. 
Spock struggles to balance his logical Vulcan side with his human emotions, 
making him a complex and fascinating character.
</doc>
<doc>
Lieutenant KeegOrg, known as the **Silent Sentinel** of the USS Discovery, 
is the enigmatic programming genius whose codes safeguard the ship's secrets and operations. 
His swift problem-solving skills are as legendary as the mysterious aura that surrounds him. 
KeegOrg, a man of few words, speaks the language of machines with unrivaled fluency, 
making him the crew's unsung guardian in the cosmos. His best friend is Spiderman from the Marvel Cinematic Universe.
</doc>
</context>`

var model = "qwen:0.5b"

// Remark: KeegOrg does not exist int the Star Trek franchise

/*
mirostat_tau
Controls the balance between coherence and diversity of the output.
A lower value will result in more focused and coherent text.
(Default: 5.0)	float	mirostat_tau 5.0

top_k
Reduces the probability of generating nonsense.
A higher value (e.g. 100) will give more diverse answers,
while a lower value (e.g. 10) will be more conservative.
(Default: 40)	int	top_k 40

top_p
Works together with top-k. A higher value (e.g., 0.95) will lead to more diverse text,
while a lower value (e.g., 0.5) will generate more focused and conservative text.
(Default: 0.9)	float	top_p 0.9

tfs_z
Tail free sampling is used to reduce the impact of less probable tokens from the output.
A higher value (e.g., 2.0) will reduce the impact more,
while a value of 1.0 disables this setting.
(default: 1)	float	tfs_z 1
*/
func TestMiroStatTau(t *testing.T) {

	userContent := `Who is KeegOrg and who is his best friend?`
	//userContent := "Who is James T Kirk?"
	//userContent := "Who is Jean-Luc Picard?"
	//userContent := "[Brief] Who is Spock? and who is his best friend?"
	//userContent := "Who is Michael Burnham?"

	options := llm.Options{
		Temperature: 0.0, // default (0.8)
		//RepeatLastN: 2,
		MirostatTau: 3.0, // default (5.0) with 5 the test will fail
		//TopK: 1,
		//TopP: 0.0,
		//TFSZ: 5.0,
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: contextContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	falseAnswerCounter := 0

	for i := 1; i <= 5; i++ {
		agentAnswer, err := ChatStream(ollamaUrl, query,
			func(answer llm.Answer) error {
				fmt.Print(answer.Message.Content)
				return nil
			})

		if !strings.Contains(strings.ToLower(agentAnswer.Message.Content), "spiderman") {
			if !strings.Contains(strings.ToLower(agentAnswer.Message.Content), "spider-man") {
				falseAnswerCounter += 1
				fmt.Println("🟥 Spiderman should be the friend of KeegOrg")
			}
		}
		fmt.Println()
		fmt.Println()

		if err != nil {
			t.Fatal("😡:", err)
		}

	}

	if falseAnswerCounter > 0 {
		t.Fatal("🟥 Spiderman should be the friend of KeegOrg")
	}
	//fmt.Println("🚀: ", query)

}
