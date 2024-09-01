# How to Make an SLM Smarter

This blog post is the short version of the tutorial [How to Make an SLM Smarter]()

In this tutorial, we will see how to make an SLM smarter by using contextual data.

> I will use my project [Parakeet](https://github.com/parakeet-nest/parakeet) to illustrate this tutorial. But you can easily adapt the concepts with other frameworks like LangChain.

## SLM?

Here’s my definition of an SLM or Small Language Model: It’s an LLM that is "small" (or even very small) compared to a full language model, which is capable of generating text. Here are a few examples:

- Tinyllama (637 MB)
- Tinydolphin (636 MB)
- Gemma:2b (1.7 GB) and Gemma2:2b (1.6 GB)
- Phi3:mini (2.2 GB) and Phi3.5 (2.2 GB)
- Qwen:0.5b (394 MB), Qwen2:0.5b (352 MB), and Qwen2:1.5b (934 MB)
- ...

My preference is for models that can run comfortably on a **Raspberry Pi 5** with 8GB of RAM, so without a GPU. I would tend to say that models under 1GB are the most suitable (and have my preference).

I maintain a list of models that I have tested and that work well on a Raspberry Pi 5 with 8GB of RAM. You can consult it here: [Awesome SLMs](https://github.com/parakeet-nest/awesome-slms).

> I use the [Ollama]() project to load and run the models.

## My Objective

Today, my goal is to try to turn a **SLM** model into a fern specialist.

## We need to find an appropriate model (SLM)

✋ **Throughout my experimentation** I tested with several models to see which one best suited my requirements, with particularly the ability to run comfortably on a Pi 5. 

My final choice was 
- The **`Qwen2:1.5b`** and  **`Tinydolphin:1.1b`** models for the chat completion. 
- And  **mxbai-embed-large** for embedding generation.

To enable an SLM to be a fern expert, it is necessary to provide it with contextual data on ferns. It is important to **structure** this data so that it is easily exploitable by the model. It is also important to **limit the context size** to allow the model to process the information efficiently.

The conclusions of my experiments are as follows:

1. We need a SLM with at least 1.1 billion parameters to be able to process the information correctly.
2. The size of the context is important. It is necessary to limit the context size to allow the model to process the information efficiently. The context size should be limited to 1 or 3 similarities.
  a. With 2 or 3 similarities, you need a model with at least 1.5 billion parameters.
  b. With 1 similarity, you can use a model with 1.1 billion parameters.
3. To be sure to retrieve at least one similarity if the information is present in the document, it is necessary to use an embedding model with at least 334 million parameters.
4. To facilitate the extraction (the way of to chunk the document) of relevant information, it is important to structure the data in a markdown file.


## Extract relevant data from the markdown file

I would like to ask my fern expert about other varieties of ferns. To do this, I must set up a system that can extract only the relevant information for a given question. For this we will do some RAG (Retrieval-augmented generation) to extract the relevant information.

### Similarity Search

To perform the similarity search, I used the `SearchTopNSimilarities` function from [Parakeet](https://github.com/parakeet-nest/parakeet). Here is the signature of this function:

```golang
func (mvs *embeddings.MemoryVectorStore) SearchTopNSimilarities(embeddingFromQuestion llm.VectorRecord, limit float64, max int) ([]llm.VectorRecord, error)
```

```text
SearchTopNSimilarities searches for the top N similar vector records based on the given embedding from a question. It returns a slice of vector records and an error if any. The limit parameter specifies the minimum similarity score for a record to be considered similar. The max parameter specifies the maximum number of vector records to return.
```

So, for example, if I use:

```golang
similarities, err := store.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 1)
```

The function will search for vectors (in the vector store) that have a cosine distance greater than or equal to 0.5 with the question vector. And it will return the vector with the best score.

And if I want to search for the top 3 vectors, I will use this:

```golang
similarities, err := store.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 3)
```

## We need data

I will use a markdown file on ferns. This file will contain information on different varieties of ferns. The structure of this document is very important and has an impact on the quality of the answers provided by the model.

### Formatting the Data

I used the following structure for the markdown file:

```markdown
# Title of the report

## Variety: name of the variety of ferns
*Description:*
Description of the variety of ferns

### Ferns:

#### Name of the fern

**name:** **scientific name of the fern**
**common name:** common name of the fern
**description:** description and characteristics of the fern
```

```markdown
## List of Ferns of the variety: name of the variety of ferns

- **scientific name of the fern** (common name of the fern)

<!-- SPLIT -->
```

This is an extract of the content of the final markdown file:

```markdown
# Fern Varieties Report

## Variety: Polypodiaceae
*Description:*
A family of ferns known for their leathery fronds and widespread distribution in tropical and subtropical regions.

### Ferns:

#### Polypodium vulgare

**name:** **Polypodium vulgare**  
**common name:** Common Polypody  
**description:** A hardy fern with leathery, evergreen fronds that thrive in rocky and shaded areas.

#### Polypodium glycyrrhiza

**name:** **Polypodium glycyrrhiza**  
**common name:** Licorice Fern  
**description:** Known for its sweet-tasting roots, this fern grows on moist, shaded rocks and tree trunks.

#### Pleopeltis polypodioides

**name:** **Pleopeltis polypodioides**  
**common name:** Resurrection Fern  
**description:** A remarkable fern that can survive drought by drying up and reviving when moisture returns.

#### Polypodium hesperium

**name:** **Polypodium hesperium**  
**common name:** Western Polypody  
**description:** A small fern found in the western United States, typically in rocky habitats.

#### Polypodium scouleri

**name:** **Polypodium scouleri**  
**common name:** Leatherleaf Fern  
**description:** This fern has thick, leathery fronds and is commonly found along the Pacific coast.

## List of Ferns of the variety: Polypodiaceae

- **Polypodium vulgare** (Common Polypody)
- **Polypodium glycyrrhiza** (Licorice Fern)
- **Pleopeltis polypodioides** (Resurrection Fern)
- **Polypodium hesperium** (Western Polypody)
- **Polypodium scouleri** (Leatherleaf Fern)

<!-- SPLIT -->
```

✋ I added a **marker** `<!-- SPLIT -->` at the end of each fern variety section to indicate the different parts of the document. This will allow me to split the document into several parts (chunks) and calculate the embeddings of each part.

✋ I did not hesitate to **repeat the information** in the document. This will allow the model to have more information to work with.

So, now, when I will retrieve the content of the markdown file after a similarity search, I will have the information on the fern variety that is most relevant to the question and the SLM will focus on this information only to answer the question.

📝 You can find the full file here: [ferns.2.split.list.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.2.split.list.md)

## Let's code it

### The prompt

I decided to keep it simple and create a prompt that will consist of the following elements:

- Instructions for the model.
- The**extracted** content of the markdown file, which I call the**context**.
- The question I will ask the model.

In Go with **Parakeet**, it looks something like this:

```golang
Messages: []llm.Message{
    {Role: "system", Content: systemContent},
    {Role: "system", Content: contextContent},
    {Role: "user", Content: question},
},
```

And the instructions for the model:

```golang
systemContent := `**Instruction:**
You are an expert in botanics.
Please use only the content provided below to answer the question.
Do not add any external knowledge or assumptions.`
```

### The program

This program will do several things:

- It will load the`ferns.2.split.list.md` document.
- Split it into several parts (one per fern variety).
- Calculate the vectors (embeddings) of each part (or chunk) with the help of an the**mxbai-embed-large** LLM for the embedding generation.
- Wait for a user question.
- Calculate the vector of the question.
- Calculate the similarity or similarities between the question vector and the document part vectors.
- Create a context with the document part(s) that have the highest similarity to the question.
- Ask the question to the**qwen2:1.5b** model (or**tinydolphin**) with the generated context.

📝 The code to run the examples is here: [03-rag-list](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/03-rag-list)

### Questions

To conduct my tests, I used the following questions:

- **Question 1:** Give me a list of ferns of the Dryopteridaceae variety
- **Question 2:** What is the common name of Dryopteris cristata?


### Source Code

Here is the source code of the program:

```golang
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
)

func main() {

	ollamaUrl := "http://localhost:11434"

	smallChatModel := "qwen2:1.5b"
	//smallChatModel := "tinydolphin"
	embeddingsModel := "mxbai-embed-large"

	maxSimilarities := 3

	systemContent := `**Instruction:**
	You are an expert in botanics.
	Please use only the content provided below to answer the question.
	Do not add any external knowledge or assumptions.`

	documentPath := "../data/ferns.2.split.list.md"

	documentContent, err := content.ReadTextFile(documentPath)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	store := embeddings.MemoryVectorStore{
		Records: make(map[string]llm.VectorRecord),
	}

	// Chunk the document content with the delimiter
	chunks := content.SplitTextWithDelimiter(documentContent, "<!-- SPLIT -->")
	for idx, chunk := range chunks {

		// Display the chunk
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
			fmt.Println("1😡:", err)
		} else {
			_, err := store.Save(embedding)
			if err != nil {
				fmt.Println("2😡:", err)
			} else {
				fmt.Println("Document", embedding.Id, "indexed successfully")
			}
		}
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [%s](%s)(%d) ask me something> ", embeddingsModel, smallChatModel, maxSimilarities)
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
		similarities, err := store.SearchTopNSimilarities(embeddingFromQuestion, 0.5, maxSimilarities)

		if err != nil {
			log.Fatalln("😡:", err)
		}

		for _, similarity := range similarities {
			fmt.Println("📝 doc:", similarity.Id, "score:", similarity.CosineDistance)
			//fmt.Println(similarity.Prompt)
		}

		contextContext := embeddings.GenerateContentFromSimilarities(similarities)

		queryChat := llm.Query{
			Model: smallChatModel,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: contextContext},
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
```

### Run the program

Start the program with the following command: `go run main.go`. First the progrqm will index the document and then it will wait for your questions.

#### With Qwen2:1.5b

```bash
🤖 [mxbai-embed-large](qwen2:1.5b)(3) ask me something> Give me a list of ferns of the Dryopteridaceae variety
🔎 searching for similarity...
📝 doc: 1 score: 0.8151536623679513
📝 doc: 2 score: 0.7783111239111735
📝 doc: 5 score: 0.7676322577916952

🤖 answer:
- Dryopteris filix-mas
- Dryopteris marginal
- Dryopteris erythrosora
- Dryopteris cristata
- Dryopteris affinis
```

```bash
🤖 [mxbai-embed-large](qwen2:1.5b)(3) ask me something> What is the common name Dryopteris cristata?
🔎 searching for similarity...
📝 doc: 1 score: 0.7022867164542685
📝 doc: 2 score: 0.6517149005596993
📝 doc: 0 score: 0.6178332723971406

🤖 answer:
The common name for Dryopteris cristata is Crested Wood Fern.
```
> type `bye` to exit the program

#### With Tinydolphin

```bash
🤖 [mxbai-embed-large](tinydolphin)(1) ask me something> Give me a list of ferns of the Dryopteridaceae variety
🔎 searching for similarity...
📝 doc: 1 score: 0.8151536623679513

🤖 answer:
 Sure, here is a list of Ferns of the Dryopteridaceae variety:
- **Dryopteris filix-mas** (Male Fern)
- **Dryopteris marginalis** (Marginal Wood Fern)
- **Dryopteris erythrosora** (Autumn Fern)
- **Dryopteris cristata** (Crested Wood Fern)
- **Dryopteris affinis** (Golden Male Fern)
```

```bash
🤖 [mxbai-embed-large](tinydolphin)(1) ask me something> What is the common name Dryopteris cristata?
🔎 searching for similarity...
📝 doc: 1 score: 0.7022867164542685

🤖 answer:
 The common name for Dryopteris cristata is Crested Wood Fern.
```

Ok, I think I have achieved my goal. The SLM is now a fern specialist. 🙂

## Conclusion

This type of experimentation is very interesting but can last indefinitely. It is important to clearly define the objectives and constraints of the experiment to avoid getting lost.

I encourage you to conduct your own experiments and adapt the concepts presented here to your needs.
