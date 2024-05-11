package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func createDocFromCode(path string, fileExtension string) ([]string, error) {
	textFiles, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	docs := []string{}
	for _, textFile := range textFiles {
		fmt.Println(textFile.Name())
		pathfile := path + "/" + textFile.Name()
		extension := filepath.Ext(pathfile)
		if extension == fileExtension {
			data, err := os.ReadFile(pathfile)
			if err != nil {
				return nil, err
			}
			docs = append(docs, string(data))
		}
	}
	return docs, nil
}

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	var embeddingsModel = "magicoder:latest" 

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	docs := []string{}

	first, err := createDocFromCode("../../examples/01-generate", ".go")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	docs = append(docs, first...)
	second, err := createDocFromCode("../../examples/02-generate-stream", ".go")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	docs = append(docs, second...)
	third, err := createDocFromCode("../../examples/03-chat", ".go")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	docs = append(docs, third...)
	fourth, err := createDocFromCode("../../examples/04-chat-stream", ".go")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	docs = append(docs, fourth...)
	fifth, err := createDocFromCode("../../examples/05-context", ".go")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	docs = append(docs, fifth...)

	sixth, err := createDocFromCode("../../examples/10-chat-conversational-memory", ".go")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	docs = append(docs, sixth...)

	seventh, err := createDocFromCode("../../examples/11-chat-conversational-bbolt/begin", ".go")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	docs = append(docs, seventh...)

	eight, err := createDocFromCode("../../examples/11-chat-conversational-bbolt/resume", ".go")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	docs = append(docs, eight...)

	fmt.Println(docs)

	// Create embeddings from documents and save them in the store
	for idx, doc := range docs {
		fmt.Println("📝 Creating embedding from document ", idx)
		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: doc,
			},
			strconv.Itoa(idx), // don't forget the id (unique identifier)
		)
		fmt.Println("📦", embedding.Id, embedding.Prompt)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			_, err := store.Save(embedding)
			if err != nil {
				fmt.Println("😡:", err)
			}
		}
	}

}
