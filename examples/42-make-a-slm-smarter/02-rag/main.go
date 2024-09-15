package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"

)

func main() {

	ollamaUrl := "http://localhost:11434"

	//smallChatModel := "qwen2:0.5b"
	//smallChatModel := "qwen2:1.5b"
	smallChatModel := "tinydolphin"

	//embeddingsModel := "all-minilm:33m"
	//embeddingsModel := "nomic-embed-text"
	embeddingsModel := "mxbai-embed-large"

	systemContent := `**Instruction:**
	You are an expert in botanics.
	Please use only the content provided below to answer the question.
	Do not add any external knowledge or assumptions.`

	documentPath := "../data/ferns.2.split.md"

	documentContent, err := content.ReadTextFile(documentPath)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	store := embeddings.MemoryVectorStore{
		Records: make(map[string]llm.VectorRecord),
	}
	// Chunk the document content

	chunks := content.SplitTextWithDelimiter(documentContent, "<!-- SPLIT -->")
	for idx, chunk := range chunks {

		fmt.Println("---------------------------------------------")
		fmt.Println(chunk)
		fmt.Println("---------------------------------------------")

		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: chunk,
			},
			strconv.Itoa(idx),
		)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			_, err := store.Save(embedding)
			if err != nil {
				fmt.Println("😡:", err)
			} else {
				fmt.Println("Document", embedding.Id, "indexed successfully")
			}
		}
	}

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 3.0,
		option.TopK: 10,
		option.TopP: 0.5,
	})


	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [%s](%s) ask me something> ", embeddingsModel, smallChatModel)
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "bye" {
			break
		}

		// Create an embedding from the question
		embeddingFromQuestion, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: question,
			},
			"question",
		)
		if err != nil {
			log.Fatalln("😡:", err)
		}
		fmt.Println("🔎 searching for similarity...")

		/*
			- **Question:** Give me a list of ferns of the Dryopteridaceae variety
			- **Question:** What is the common name Dryopteris cristata?
		*/

		similarities, err := store.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 1)

		if err != nil {
			log.Fatalln("😡:", err)
		}

		for _, similarity := range similarities {
			fmt.Println("📝 doc:", similarity.Id, "score:", similarity.CosineDistance)
			fmt.Println(similarity.Prompt)
		}

		contextContext := embeddings.GenerateContentFromSimilarities(similarities)

		queryChat := llm.Query{
			Model: smallChatModel,
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
