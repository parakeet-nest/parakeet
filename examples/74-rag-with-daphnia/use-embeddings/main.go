package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := gear.GetEnvString("OLLAMA_BASE_URL", "http://localhost:11434")
	embeddingsModel := gear.GetEnvString("LLM_EMBEDDINGS", "mxbai-embed-large:335m")

	vectorStorePath := gear.GetEnvString("DAPHNIA_STORE_PATH", "../sourcedata.gob")

	store := embeddings.DaphniaVectoreStore{}
	err = store.Initialize(vectorStorePath)

	if err != nil {
		log.Fatalln("😡:", err)
	}

	//var smallChatModel = "qwen2.5:14b"      // This model is for the chat completion
	//var smallChatModel = "qwen2.5:3b" // This model is for the chat completion
	smallChatModel := gear.GetEnvString("LLM_CHAT", "qwen2.5:14b")

	systemContent := `You are a Golang expert and know very well the extism go SDK. Use only the provided content to answer the question.`

	userContent := `How to call a function of a wasm module in golang with the extism-go sdk?`

	//userContent := `how to define a host function?`

	// TODO: add the function name to see if it changes the result
	// Create an embedding from the question
	embeddingFromQuestion, err := embeddings.CreateEmbedding(
		ollamaUrl,
		llm.Query4Embedding{
			Model:  embeddingsModel,
			Prompt: userContent,
		},
		"question",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}
	fmt.Println("🔎 searching for similarity...")

	//similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 5) // qwen2.5:14b
	similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, 0.65, 2) // qwen2.5:7b

	for _, similarity := range similarities {
		fmt.Println("📝 doc:", similarity.Id, "score:", similarity.Score)
		fmt.Println("--------------------------------------------------")
		fmt.Println("📝 metadata:", similarity.Prompt)
		fmt.Println("--------------------------------------------------")
	}


	fmt.Println("🎉 number of similarities:", len(similarities))

	documentsContent := embeddings.GenerateContextFromSimilarities(similarities)

	fmt.Println("📚 documents content:", documentsContent)

	numCtx := gear.GetEnvInt("OPTION_NUM_CTX", 100)

	estimatedTokens:= content.EstimateGPTTokens(systemContent+documentsContent+userContent)
	fmt.Println("================================================")
	fmt.Println("🧩 estimated tokens:", estimatedTokens)
	fmt.Println("================================================")


	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.0,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.2,
		option.TopK:          10,
		option.TopP:          0.5,
		option.NumCtx:        numCtx,
		//option.NumCtx: 40533,

	})

	if numCtx < estimatedTokens {
		fmt.Println("🔥 numCtx is less than estimated tokens")
		options.NumCtx = estimatedTokens + 100

		fmt.Println(options)
	} 

	query := llm.Query{
		Model: smallChatModel,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	fmt.Println()
	fmt.Println("🤖 answer:")

	// Answer the question
	_, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

	fmt.Println()
}
