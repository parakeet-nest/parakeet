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
	err := godotenv.Load(".env", "openai.key.env")	
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaParrot := squawk.New().Model(os.Getenv("OLLAMA_LLM_CHAT")).
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
		}).Cmd(func(self *squawk.Squawk) {
			fmt.Println("\n------------------------------------")
		}).
		SaveAnswer().
		User("Who is his worst ennemy?", "question-03").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		}).Cmd(func(self *squawk.Squawk) {
			fmt.Println("\n------------------------------------")
		})

	fmt.Println(ollamaParrot.LastError())
	fmt.Println("==================")
	fmt.Println("🟢 Answer:", ollamaParrot.LastAnswer().Message.Content)
	fmt.Println("==================")



	dockerModelRunnerParrot := squawk.New().Model(os.Getenv("MODEL_RUNNER_LLM_CHAT")).
		BaseURL(os.Getenv("MODEL_RUNNER_BASE_URL")).
		Provider(provider.DockerModelRunner).
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
		}).Cmd(func(self *squawk.Squawk) {
			fmt.Println("\n------------------------------------")
		}).
		SaveAnswer().
		User("Who is his worst ennemy?", "question-03").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		}).Cmd(func(self *squawk.Squawk) {
			fmt.Println("\n------------------------------------")
		})

	fmt.Println(dockerModelRunnerParrot.LastError())
	fmt.Println("==================")
	fmt.Println("🟩 Answer:", dockerModelRunnerParrot.LastAnswer().Message.Content)
	fmt.Println("==================")

	openAIParrot := squawk.New().Model(os.Getenv("OPENAI_LLM_CHAT")).
		BaseURL(os.Getenv("OPENAI_BASE_URL")).
		Provider(provider.OpenAI, os.Getenv("OPENAI_API_KEY")).
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
		}).Cmd(func(self *squawk.Squawk) {
			fmt.Println("\n------------------------------------")
		}).
		SaveAnswer().
		User("Who is his worst ennemy?", "question-03").
		ChatStream(func(answer llm.Answer, self *squawk.Squawk) error {
			fmt.Print(answer.Message.Content)
			return nil
		}).Cmd(func(self *squawk.Squawk) {
			fmt.Println("\n------------------------------------")
		})

		fmt.Println(openAIParrot.LastError())
		fmt.Println("==================")
		fmt.Println("💚 Answer:", openAIParrot.LastAnswer().Message.Content)
		fmt.Println("==================")

}
