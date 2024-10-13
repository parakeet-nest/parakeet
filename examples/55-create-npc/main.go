package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

type Character struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

func GenerateName(species, ollamaUrl, model string) (Character, error) {


	instructionsContent, err := os.ReadFile("name.instructions.md")
	if err != nil {
		return Character{}, err
		//log.Fatal("😡:", err)
	}

	userContentTpl, err := os.ReadFile("name.usertpl.md")
	if err != nil {
		return Character{}, err
		//log.Fatal("😡:", err)
	}
	userContent := fmt.Sprintf(string(userContentTpl), species)

	fmt.Println("🤖 user message for name generation>", userContent)

	options := llm.DefaultOptions()

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: string(instructionsContent)},
			{Role: "user", Content: userContent},
		},
		Options: options,
		Format:  "json",
		Raw:     true,
	}

	// Answer the question
	answer, err := completion.Chat(ollamaUrl, query)

	if err != nil {
		return Character{}, err
		//log.Fatal("😡:", err)
	}

	fmt.Println("🤖 generated name>", answer.Message.Content)

	// Unmarshal the JSON data into a struct
	var character Character
	err = json.Unmarshal([]byte(answer.Message.Content), &character)
	if err != nil {
		log.Fatal("😡:", err)
	}

	return Character{
		Name: character.Name,
		Kind: character.Kind,
	}, nil
}

func GenerateDescription(character Character, ollamaUrl, model string) (string, error) {

	instructionsContent, err := os.ReadFile("description.instructions.md")
	if err != nil {
		return "", err
		//log.Fatal("😡:", err)
	}

	stepsContent, err := os.ReadFile("description.steps.md")
	if err != nil {
		return "", err
		//log.Fatal("😡:", err)
	}

	userContent := fmt.Sprintf("Create a %s with this name:%s", character.Kind, character.Name)

	fmt.Println("🤖 user message for description generation>", userContent)

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.9,
		option.RepeatLastN:   3,
		option.RepeatPenalty: 2.0,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: string(instructionsContent)},
			{Role: "system", Content: string(stepsContent)},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	// Answer the question
	answer, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println()
	return answer.Message.Content, nil

}

func ChatWithCharacter(character Character, description, ollamaUrl, model string) {

	systemContentTpl , err := os.ReadFile("chat.instructiontpl.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	systemContent := fmt.Sprintf(string(systemContentTpl), character.Kind, character.Name)

	contextContext := description

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.0,
		option.RepeatLastN:   3,
		option.RepeatPenalty: 2.0,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [%s] ask me something> ", character.Name)
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "bye" {
			break
		}

		queryChat := llm.Query{
			Model: model,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: contextContext},
				{Role: "user", Content: question},
			},
			Options: options,
		}

		fmt.Println()
		fmt.Println("🤖 answer:")

		// Answer the question
		_, err := completion.ChatStream(ollamaUrl, queryChat,
			func(answer llm.Answer) error {
				fmt.Print(answer.Message.Content)
				return nil
			})

		if err != nil {
			log.Fatal("😡:", err)
		}

		fmt.Println()
	}

}

func main() {

	ollamaUrl := "http://localhost:11434"
	model := "nemotron-mini"

	npc, err := GenerateName("elf", ollamaUrl, model)
	if err != nil {
		log.Fatalln("😡:", err)
	}
	fmt.Println(npc)

	// Character sheet
	description, err := GenerateDescription(npc, ollamaUrl, model)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	// Save the description to a file
	err = os.WriteFile("description.md", []byte("# CHARACTER SHEET\n\n"+description), 0644)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ChatWithCharacter(npc, description, ollamaUrl, model)
}
