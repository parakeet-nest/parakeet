package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"

	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

type TemplateParameters struct {
	NumberOfVarieties string
	Expertise         string
}

var promptTemplate = `Generate a structured report in markdown format that details {{.NumberOfVarieties}} different {{.Expertise}} varieties. 
Each variety should have a description, followed by a list of specific {{.Expertise}}s belonging to that variety. 
Each {{.Expertise}} should include the following details:

1. **name:** The scientific name of the {{.Expertise}}, bolded.
2. **common name:** The common name of the {{.Expertise}}, also bolded.
3. **description:** A brief description of the {{.Expertise}}, detailing its characteristics and habitat.

The document should be divided into sections for each variety, with the title "Variety: [Variety Name]" followed by a description of the variety. 
Include a '<!-- SPLIT -->' marker between each variety section.

Here's an example structure:

# {{.Expertise}} Varieties Report

## Variety: [name of the variety of {{.Expertise}}]
*Description:*
[description of the variety of {{.Expertise}}]

### {{.Expertise}}s:

#### [name of the {{.Expertise}}]

**name:** **[name of the {{.Expertise}}]**  
**common name:** [common name of the {{.Expertise}}]  
**description:** [description of the {{.Expertise}}]

(Continue for other {{.Expertise}}s in this variety...)

<!-- SPLIT -->

(Continue for other varieties...)
OUTPUT INSTRUCTIONS: do not any remark at the end of the report.
`

func main() {
	// create a `.env` file with the following content:
	// OPENAI_API_KEY=your_openai_api_key
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	// Generate data for a SLM with OpenAI
	openAIURL := "https://api.openai.com/v1"
	model := "gpt-4o-mini"

	systemContent := `You are a useful AI agent.`

	expertise := "dinosaur"

	prompt, err := content.InterpolateString(
		promptTemplate,
		TemplateParameters{
			NumberOfVarieties: "five",
			Expertise:         expertise,
		},
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: prompt},
		},
	}

	/*
		answer, err := completion.ChatWithOpenAI(openAIUrl, query)
		if err != nil {
			log.Fatal("😡:", err)
		}
		fmt.Println(answer.Choices[0].Message.Content)
	*/

	markdownReport, err := completion.ChatStream(openAIURL, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		}, provider.OpenAI, os.Getenv("OPENAI_API_KEY"))

	if err != nil {
		log.Fatal("😡:", err)
	}

	// Create a new file
	file, err := os.Create("report.md")
	if err != nil {
		log.Fatal("😡:", err)
	}
	defer file.Close()

	// Write the content of markdownReport to the file
	_, err = file.WriteString(markdownReport.Message.Content)
	if err != nil {
		log.Fatal("😡:", err)
	}

	// Print a success message
	fmt.Println("\n\n🎉 Report file created successfully.")

	// Make the SLM smarter
	ollamaUrl := "http://localhost:11434"

	smallChatModel := "qwen2.5:1.5b"

	embeddingsModel := "mxbai-embed-large"
	//embeddingsModel := "bge-m3"

	systemContent, err = content.InterpolateString(
		`**Instruction:**
	You are an expert in {{.Expertise}}.
	Please use only the content provided below to answer the question.
	Do not add any external knowledge or assumptions.`,
		TemplateParameters{
			Expertise: expertise,
		},
	)

	if err != nil {
		log.Fatalln("😡:", err)
	}

	documentContent, err := content.ReadTextFile("./report.md")
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
		option.Temperature:   0.0,
		option.RepeatLastN:   10,
		option.RepeatPenalty: 10.0,
		option.TopK:          10,
		option.TopP:          0.5,
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

		similarities, err := store.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 1)

		if err != nil {
			log.Fatalln("😡:", err)
		}

		for _, similarity := range similarities {
			fmt.Println("📝 doc:", similarity.Id, "score:", similarity.CosineSimilarity)
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
