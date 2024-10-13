package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func ChatWithCharacter(instructions, ollamaUrl, model, embeddingsModel string) {

	systemContent := instructions

	//contextContext := description

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.8,
		option.RepeatLastN:   3,
		option.RepeatPenalty: 2.0,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [%s] ask me something> ", model)
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "bye" {
			break
		}

		embeddingFromQuestion, err := CreateEmbeddingsFromQuestion(question, ollamaUrl, embeddingsModel)
		if err != nil {
			log.Fatal("😡:", err)
		}

		fmt.Println("🔎 searching for similarity...")
		similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, 0.3, 1)
		fmt.Println("🎉 number of similarities:", len(similarities))

		for _, similarity := range similarities {
			fmt.Println("🔎 similarity:", similarity.Prompt)
		}

		contextContext := ""
		if len(similarities) == 0 {
			contextContext = "I'm sorry, I don't know the answer to that question."
		} else {
			contextContext = embeddings.GenerateContextFromSimilarities(similarities)
		}

		queryChat := llm.Query{
			Model: model,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: contextContext},
				{Role: "user", Content: question},
			},
			Options:          options,
			TokenHeaderName:  "X-TOKEN",
			TokenHeaderValue: os.Getenv("TOKEN"),
		}

		fmt.Println()
		fmt.Println("🤖 answer:")

		// Answer the question
		_, err = completion.ChatStream(ollamaUrl, queryChat,
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

func PreloadModel(instructions, ollamaUrl, model string) error {
	current_time := time.Now()
	fmt.Println("🟢", current_time.Format("2006-01-02 15:04:05"))

	fmt.Println("🤖 preloading the model ...")
	systemContent := instructions

	_, err := completion.ChatStream(ollamaUrl, llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: "who are you?"},
		},
		Options:          llm.DefaultOptions(),
		TokenHeaderName:  "X-TOKEN",
		TokenHeaderValue: os.Getenv("TOKEN"),
	}, func(answer llm.Answer) error {
		fmt.Print(answer.Message.Content)
		return nil
	})

	if err != nil {
		return err
		//log.Fatal("😡:", err)
	}
	fmt.Println()
	fmt.Println("🤖 model is ready 🎉")
	current_time = time.Now()
	fmt.Println("🟩", current_time.Format("2006-01-02 15:04:05"))
	return nil
}

func CreateEmbeddings(description string, ollamaUrl, model string) error {

	//chunks := content.ParseMarkdown(description)

	chunks := content.SplitTextWithDelimiter(description, "---")

	// Create embeddings from documents and save them in the store
	for idx, doc := range chunks {
		fmt.Println("📝 Creating embedding from document ", idx)

		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:            model,
				Prompt:           doc,
				TokenHeaderName:  "X-TOKEN",
				TokenHeaderValue: os.Getenv("TOKEN"),
			},
			strconv.Itoa(idx),
		)
		if err != nil {
			fmt.Println("😡:", err)
			return err
		} else {
			//embedding.MetaData = "📝 chunk num: " + strconv.Itoa(idx)
			_, err := store.Save(embedding)
			if err != nil {
				fmt.Println("😡:", err)
				return err
			}
		}
		fmt.Println(embedding.Prompt)
	}
	return nil
}

func CreateEmbeddingsFromQuestion(question string, ollamaUrl, model string) (llm.VectorRecord, error) {
	embedding, err := embeddings.CreateEmbedding(
		ollamaUrl,
		llm.Query4Embedding{
			Model:            model,
			Prompt:           question,
			TokenHeaderName:  "X-TOKEN",
			TokenHeaderValue: os.Getenv("TOKEN"),
		},
		"question",
	)
	if err != nil {
		log.Println("😡:", err)
		return llm.VectorRecord{}, err
	}
	return embedding, nil
}

var store = embeddings.MemoryVectorStore{}

func main() {

	// create a `.env` file with the following content:
	// TOKEN=your_token
	// OLLAMA_URL=your_ollama_url
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := os.Getenv("OLLAMA_URL")
	model := "nemotron-mini"
	embeddingsModel := "bge-m3"

	store = embeddings.MemoryVectorStore{
		Records: make(map[string]llm.VectorRecord),
	}

	instructions, err := os.ReadFile("instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	description, err := os.ReadFile("description.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	// Create embeddings from the description
	CreateEmbeddings(string(description), ollamaUrl, embeddingsModel)

	// Preload the model
	err = PreloadModel(string(instructions), ollamaUrl, model)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ChatWithCharacter(
		string(instructions),
		ollamaUrl,
		model,
		embeddingsModel,
	)

}
