package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func getFirstLine(s string) string {
	// Split the string by the newline character
	lines := strings.SplitN(s, "\n", 2)
	// Return the first line
	return lines[0]
}

func main() {
	ollamaUrl := "http://localhost:11434"
	var embeddingsModel = "all-minilm:33m" // This model is for the embeddings of the documents

	store := embeddings.BboltVectorStore{}
	err := store.Initialize("../embeddings.db")

	if err != nil {
		log.Fatalln("😡:", err)
	}

	documentationContent, err := content.ReadTextFile("../doc/go1.23.md")
	if err != nil {
		log.Fatalln("😡:", err)
	}

	chunks := content.SplitMarkdownBySections(documentationContent)

	medatData := ""
	
	// Create embeddings from documents and save them in the store
	for idx, doc := range chunks {
		fmt.Println("Creating embedding from document ", idx)

		var prompt string
		if strings.HasPrefix(doc, "## ") {
			prompt = doc
		} else {
			// Add the first line of level 2 sections to the metadata (and then to the next level == sub sections)
			prompt = "<!-- "+medatData+" -->\n"+doc
		}

		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: prompt,
			},
			strconv.Itoa(idx),
		)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			//embedding.MetaData = "📝 chunk num: " + strconv.Itoa(idx)
			_, err := store.Save(embedding)
			if err != nil {
				fmt.Println("😡:", err)
			}
		}
		// Add the first line of level 2 sections to the metadata (and then to the next level == sub sections)
		if strings.HasPrefix(doc, "## ") {
			medatData = getFirstLine(doc)
		}
		fmt.Println(embedding.Prompt)
	}

}
