package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/squawk"
)

func main() {
	// create a `.env` file with the following content:
	// OPENAI_API_KEY=your_openai_api_key
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	sqwk := squawk.New("qwen2.5:0.5b").
		BaseURL(os.Getenv("OLLAMA_BASE_URL")).
		Provider(provider.Ollama).
		Options(llm.SetOptions(map[string]interface{}{
			option.Temperature:   0.0,
			option.RepeatLastN:   2,
			option.RepeatPenalty: 2.2,
		})).
		System("You are a useful AI agent, you are a Star Trek expert.", "instructions").
		User("Who is James T Kirk?", "question-01").
		Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
			if err != nil {
				fmt.Println("😡 Error:", err)
			}
			fmt.Print(answer.Message.Content)
			fmt.Println("\n------------------------------------")
		}).
		SaveAnswer(). // add answer to the history / to the messages list
		User("Who is his best friend?", "question-02").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		}).Exec(func(self *squawk.Squawk) {
			fmt.Println("\n------------------------------------")
		}).
		SaveAnswer().
		User("Who is his worst ennemy?", "question-03").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			//return errors.New("error in the stream")
			return nil
		}).Exec(func(self *squawk.Squawk) {
			fmt.Println("\n------------------------------------")
		})

	fmt.Println(sqwk.LastError())
	fmt.Println("==================")
	fmt.Println("🟢 Answer:", sqwk.LastAnswer().Message.Content)

	//fmt.Println("\n\ndisplay:", display)

	// You can add a label to messages - useful for removing messages

	//fmt.Println("Answer:", sqwk.LastAnswer().Message.Content)

	/*
		answer, err := sqwk.ChatExec()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Answer:", answer.Message.Content)
	*/

	/*
		answer, err = sqwk.
			Model("gpt-4o-mini").
			Provider(provider.OpenAI, os.Getenv("OPENAI_API_KEY")).ChatExec()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Answer:", answer.Message.Content)
	*/

}
