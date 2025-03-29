package main

/*
Example from https://docs.anthropic.com/claude/page/code-clarifier
*/
import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	//model := "smollm2:135m"
	model := "smollm:135m"

	/*
	systemContent := `You are a simple assistant that only answers questions using the COMMON QUESTIONS list provided in your context. 
	If the user's question matches one in your list, respond with ONLY the corresponding ANSWER. 
	Do not add any explanations, greetings, or additional text. 
	Do not make up information. 
	If you don't find a matching question, respond with "I don't have information about that."

	*/

	systemContent := `You are a simple lookup system. Match the user's exact question to the COMMON QUESTIONS list. Return ONLY the ANSWER that corresponds to the matched question. Be precise in matching - do not mix up answers between different questions. Never summarize or modify the answers.
	`

	/*
	contextContent := `COMMON QUESTIONS:
	QUESTION: what is the favorite pizza of Philippe?
	ANSWER: Philippe loves Hawaiian pizza.

	QUESTION: what are the ingredients of an Hawaiian pizza?
	ANSWER: An Hawaiian pizza is made with tomato sauce, cheese, ham, and pineapple.
	`
	*/

	contextContent := `COMMON QUESTIONS:
	[Q1] what is the favorite pizza of Philippe?
	[A1] Philippe loves Hawaiian pizza.

	[Q2] what are the ingredients of an Hawaiian pizza?
	[A2] An Hawaiian pizza is made with tomato sauce, cheese, ham, and pineapple.

	[Q3] who is Bob Morane?
	[A3] Bob Morane is a fictional character created by French-speaking Belgian writer Henri Vernes.
	`

	/*
	contextContent = `COMMON QUESTIONS:
	{
	"what is the favorite pizza of Philippe?": "Philippe loves Hawaiian pizza.",
	"what are the ingredients of an Hawaiian pizza?": "An Hawaiian pizza is made with tomato sauce, cheese, ham, and pineapple.",
	"who is Bob Morane?": "Bob Morane is a fictional character created by French-speaking Belgian writer Henri Vernes.",
	}
	`
	*/

	//userContent := `QUESTION: what is the favorite pizza of Philippe?`
	userContent := `QUESTION: give me the list of the ingredients of an Hawaiian pizza?`


	options := llm.SetOptions(map[string]interface{}{
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.5,
		option.Temperature:   0.0,
	})

	//fmt.Println(options)

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: contextContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
		Stream:  false,
	}

	_, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println()
}
