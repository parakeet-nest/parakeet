package main

import (
	"fmt"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

var ollamaUrl = "http://localhost:11434"

// if working from a container
// use: "http://host.docker.internal:11434"

var systemContent = `You are an expert of the Star Trek franchise.
Using the provided context, answer the user's question to the best of your ability using only the resources provided.
`

var facts = map[string]string{
	"burnham" : `Michael Burnham is the main character on the Star Trek series Discovery.
	Michael Burnham's best friend is Sylvia Tilly`,
	"kirk": `James T. Kirk, also known as Captain Kirk, is the iconic captain of the starship USS Enterprise.
	Kirk's best friend is Spock.`,
	"picard": `Jean-Luc Picard is the captain of the USS Enterprise-D
	Jean-Luc Picard's best friend is Dr. Beverly Crusher.`,
	"spock": `Spock is most famous for being the half-Vulcan, half-human science officer
	and first officer on the starship USS Enterprise.
	Spock's best friend is Kirk.`,
	"keegorg": `Lieutenant KeegOrg is the enigmatic programming genius whose codes safeguard the ship's secrets and operations.
	KeegOrg's best friend is Spiderman from the Marvel Cinematic Universe.`,
}


var model = "qwen:0.5b"

func Question(userContent string) (string, error) {
	options := llm.Options{
		Temperature: 0.0, // default (0.8)
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	answer, err := completion.Chat(ollamaUrl, query)

	fmt.Println("🙂 >", userContent)
	fmt.Println("🤖 >", answer.Message.Content)
	fmt.Println()

	if err != nil {
		return "", err
	}

	return answer.Message.Content, nil

}

func QuestionWithContext(userContent string, contextContent string) (string, error) {
	options := llm.Options{
		Temperature: 0.0, // default (0.8)
		RepeatLastN: 2,   // default (64) the default value will "freeze" deepseek-coder
		//MirostatTau: 0.0, // default (5.0) with 5 the test will fail

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

	answer, err := completion.Chat(ollamaUrl, query)

	fmt.Println("🙂 >", userContent)
	fmt.Println("🤖 >", answer.Message.Content)
	fmt.Println()

	if err != nil {
		return "", err
	}

	return answer.Message.Content, nil

}

func SetOfQuestions() {
	userContent := "Who is James T Kirk and who is his best friend?"
	_, err := Question(userContent)

	if err != nil {
		fmt.Println("😡:", err)
	}

	userContent = "Who is Jean-Luc Picard and who is his best friend?"
	_, err = Question(userContent)

	if err != nil {
		fmt.Println("😡:", err)
	}

	userContent = "Who is Michael Burnham and who is his best friend?"
	_, err = Question(userContent)

	if err != nil {
		fmt.Println("😡:", err)
	}

	userContent = "Who is Spock and who is his best friend?"
	_, err = Question(userContent)

	if err != nil {
		fmt.Println("😡:", err)
	}

	/*
		KeegOrg does not exist in the Star Trek franchise
		but the agent knows that Star Trek is the main topic
		so, the agent try to answer as if the character exists in the Star Trek franchise
	*/
	userContent = "Who is KeegOrg and who is his best friend?"
	_, err = Question(userContent)

	if err != nil {
		fmt.Println("😡:", err)
	}

}

func SetOfQuestionsWithContext() {

	contextContent := `<context>`

	for _, v := range facts {
		contextContent += `<doc>` + v + `</doc>`
	}

	contextContent += `</context>`

	fmt.Println(contextContent)



	// Be careful, the whole context is used, so if you don't precise
	// who is the best friend of Kirk, the agent will try to answer
	// with another name
	userContent := "Who is James T Kirk and who is his best friend?"
	_, err := QuestionWithContext(userContent, contextContent)

	if err != nil {
		fmt.Println("😡:", err)
	}

	userContent = "Who is Jean-Luc Picard and who is his best friend?"
	_, err = QuestionWithContext(userContent, contextContent)

	if err != nil {
		fmt.Println("😡:", err)
	}

	userContent = "Who is Michael Burnham and who is his best friend?"
	_, err = QuestionWithContext(userContent, contextContent)

	if err != nil {
		fmt.Println("😡:", err)
	}

	userContent = "Who is Spock and who is his best friend?"
	_, err = QuestionWithContext(userContent, contextContent)

	if err != nil {
		fmt.Println("😡:", err)
	}

	/*
		KeegOrg does not exist in the Star Trek franchise
		but the agent knows that Star Trek is the main topic
		so, the agent try to answer as if the character exists in the Star Trek franchise
	*/
	userContent = "Who is KeegOrg and who is his best friend?"
	_, err = QuestionWithContext(userContent, contextContent)

	if err != nil {
		fmt.Println("😡:", err)
	}

}

func SetOfQuestionsWithTargetedContext() {
	// Be careful, the whole context is used, so if you don't precise
	// who is the best friend of Kirk, the agent will try to answer
	// with another name
	userContent := "Who is James T Kirk and who is his best friend?"
	_, err := QuestionWithContext(userContent, "<context><doc>" + facts["kirk"] + "</doc></context>")

	if err != nil {
		fmt.Println("😡:", err)
	}

	userContent = "Who is Jean-Luc Picard and who is his best friend?"
	_, err = QuestionWithContext(userContent, "<context><doc>" + facts["picard"] + "</doc></context>")

	if err != nil {
		fmt.Println("😡:", err)
	}

	userContent = "Who is Michael Burnham and who is his best friend?"
	_, err = QuestionWithContext(userContent, "<context><doc>" +facts["burnham"] + "</doc></context>")

	if err != nil {
		fmt.Println("😡:", err)
	}

	userContent = "Who is Spock and who is his best friend?"
	_, err = QuestionWithContext(userContent, "<context><doc>" +facts["spock"] + "</doc></context>")

	if err != nil {
		fmt.Println("😡:", err)
	}

	/*
		KeegOrg does not exist in the Star Trek franchise
		but the agent knows that Star Trek is the main topic
		so, the agent try to answer as if the character exists in the Star Trek franchise
	*/
	userContent = "Who is KeegOrg and who is his best friend?"
	_, err = QuestionWithContext(userContent, "<context><doc>" + facts["keegorg"] + "</doc></context>")

	if err != nil {
		fmt.Println("😡:", err)
	}

}




func main() {
	//SetOfQuestions()
	//fmt.Println("-------------------------------------------")
	//SetOfQuestionsWithContext()
	SetOfQuestionsWithTargetedContext()
}

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
