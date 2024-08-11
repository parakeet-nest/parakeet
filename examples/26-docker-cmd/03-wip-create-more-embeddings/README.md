# Make a Small Language Model Smarter: Teach Docker Commands

Unfortunately, we have made our SLM smarter for only three examples. We will need to give it the ability to build appropriate contexts based on user requests.

## Let's do some RAG and create Embeddings

RAG (Retrieval-Augmented Generation) and embeddings are key concepts in Generative AI:

**Embeddings**: These are numerical representations of words, phrases, or entire documents in a high-dimensional vector space. They capture semantic meaning, allowing similar concepts to be close to each other in this space.

**RAG**: This is a technique that combines retrieval of relevant information with language generation. It works by:

1. Creating embeddings of a knowledge base
2. For a given query, finding the most relevant information using embedding similarity
3. Providing this retrieved context to a LLM to generate a response

RAG enhances the accuracy and relevance of AI responses by grounding them in specific, retrieved information, rather than relying solely on the model's pre-trained knowledge.

### Create Embeddings for Docker Commands

To create the embeddings database, we need a knowledge base. We will use the Docker commands dataset from the Hugging Face Datasets library: https://huggingface.co/datasets/adeocybersecurity/DockerCommand

To download the dataset, run the following command:

```bash
wget https://huggingface.co/datasets/adeocybersecurity/DockerCommand/resolve/main/NLDockercommands.json
```

You will get a JSON file containing the Docker commands dataset with the following structure:

```json
[
  {
    "input": "Give me a list of containers that have the Ubuntu image as their ancestor.",
    "instruction": "translate this sentence in docker command",
    "output": "docker ps --filter 'ancestor=ubuntu'"
  },
]
```

We will use the **all-minilm:33m** LLM to generate embeddings from every record in the dataset. The **all-minilm:33m** model is a tiny language model of 67MB. This model is only used for embeddings and not for chatting.

We need to create a new Golang project and add the **Parakeet** package to it. This project will be a simple command-line application that interacts with **all-minilm:33m** to generate the embeddings and store them in a vector database.

The provided Go code reads a JSON file containing a list of items, processes each item to create embeddings using **all-minilm:33m**, and then saves these embeddings to a vector store.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

// 1️⃣ Item 
type Item struct {
	Input       string `json:"input"`
	Instruction string `json:"instruction"`
	Output      string `json:"output"`
}

func main() {
	ollamaUrl := "http://localhost:11434"
	embeddingsModel := "all-minilm:33m" 

	// 2️⃣ Initialize the vector store
	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	// 3️⃣ Read the JSON file
	fileContent, err := os.ReadFile("./NLDockercommands.json")
	if err != nil {
		log.Fatal("😠 Error reading file:", err)
	}

	// 4️⃣ Parse the JSON data
	var items []Item
	err = json.Unmarshal(fileContent, &items)
	if err != nil {
		log.Fatal("😠 Error parsing JSON:", err)
	}

	// 5️⃣ Create and save the embeddings
	for i, item := range items {
		fmt.Println("📝 Creating embedding from record ", i+1)

		doc := fmt.Sprintf("Input: %s \n\nOutput:%s", item.Input, item.Output)

		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: doc,
			},
			strconv.Itoa(i+1),
		)
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
```

1. **Item Struct**: Defines a struct `Item` with three fields: `Input`, `Instruction`, and `Output`, all of which are strings and will be parsed from the JSON file.
2. **Initialization**:
  - Initializes a `BboltVectorStore` to store embeddings, pointing to a database file ../embeddings.db. (*Parakeet provides [two kinds of vector store](https://github.com/parakeet-nest/parakeet?tab=readme-ov-file#vector-stores): a in memory vector store and a bbolt vector store. The last one use [Bbolt](https://github.com/etcd-io/bbolt), an embedded key/value database for Go to persist the vectors and the related data.*)
3. Reads the content of `NLDockercommands.json` into `fileContent`.
4. Parse JSON Data: unmarshals the JSON content into a slice of Item structs.
5. Create and Save Embeddings:
  - Iterates over each Item in the parsed data.
  - For each item, it constructs a document string doc combining the `Input` and `Output` fields.
  - Calls `embeddings.CreateEmbedding` with the constructed document to create an embedding.
  - Saves the embedding to the vector store.


Ok, now, let's run the program to create the embeddings and store them in the vector database: `embeddings.db`.

```bash
go run main.go
```

Wait for the program to finish processing all the items in the dataset (*we run it only once, except if we want to update the dataset*). Once completed, you will have a database of embeddings for the Docker commands dataset of `2415` records.

Now, we are ready to update the first program to use the embeddings database to provide context to the LLM.
