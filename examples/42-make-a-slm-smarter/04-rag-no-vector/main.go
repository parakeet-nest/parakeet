package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/similarity"
)

func main() {

	ollamaUrl := "http://localhost:11434"

	//smallChatModel := "qwen2:1.5b"
	smallChatModel := "tinydolphin"

	systemContent := `**Instruction:**
	You are an expert in botanics.
	Please use only the content provided below to answer the question.
	Do not add any external knowledge or assumptions.`

	documentPath := "../data/ferns.2.split.list.md"

	documentContent, err := content.ReadTextFile(documentPath)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	// Chunk the document content with the delimiter
	chunks := content.SplitTextWithDelimiter(documentContent, "<!-- SPLIT -->")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [%s] ask me something> ", smallChatModel)
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "bye" {
			break
		}

		fmt.Println("🔎 searching for the best similarity...")

		currentJaccardIndex := 0.0
		nearestContent := ""
		splittedQuestion := strings.Fields(question)

		for idx, chunk := range chunks {

			sixtyFirstCharactersOfChunk := chunk[:60]

			jaccardIndex := similarity.JaccardSimilarityCoeff(splittedQuestion, strings.Fields(sixtyFirstCharactersOfChunk))
			fmt.Println("-", idx, "Jaccard index:", jaccardIndex)
			if jaccardIndex >= currentJaccardIndex {
				nearestContent = chunk
				currentJaccardIndex = jaccardIndex
			}
		}
		// handle the case when no content is found (jaccardIndex == 0.0 or nearestContent == "")

		/*
			- **Question:** Give me a list of ferns of the Dryopteridaceae variety
			- **Question:** What is the common name Dryopteris cristata?
		*/

		fmt.Println("🔢 score:", currentJaccardIndex)
		fmt.Println("📝 doc:", nearestContent)

		queryChat := llm.Query{
			Model: smallChatModel,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: nearestContent},
				{Role: "user", Content: question},
			},
			Options: llm.Options{
				Temperature:   0.0,
				RepeatLastN:   2,
				RepeatPenalty: 3.0,
				TopK:          10,
				TopP:          0.5,
			},
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
