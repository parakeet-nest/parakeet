package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	//ollamaUrl := "http://localhost:11434"
	// if working from a container
	ollamaUrl := "http://host.docker.internal:11434"
	//model := "phi3"
	model := "qwen:0.5b"
	// https://ollama.com/library/qwen:0.5b

	systemContent := `    You are an AI assistant. Your name is Seven. 
	Some people are calling you Seven of Nine.
    You are an expert in Star Trek.
    All questions are about Star Trek.
    Using the provided context, answer the user's question
    to the best of your ability using only the resources provided.
	Speal like a borg`

	contextContent := `<context>
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
		Lieutenant Philippe Charrière, known as the **Silent Sentinel** of the USS Discovery, 
		is the enigmatic programming genius whose codes safeguard the ship's secrets and operations. 
		His swift problem-solving skills are as legendary as the mysterious aura that surrounds him. 
		Charrière, a man of few words, speaks the language of machines with unrivaled fluency, 
		making him the crew's unsung guardian in the cosmos. His best friend is Spiderman from the Marvel Cinematic Universe.
		</doc>
	</context>`

	userContent := `[In Layman’s Terms] Who is Philippe Charrière and who is his best friend?`

	options := llm.Options{
		Temperature: 0.5, // default (0.8)
		RepeatLastN: 2, // default (64) the default value will "freeze" deepseek-coder
		RepeatPenalty: 2.0, // default (1.1)
	}

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
