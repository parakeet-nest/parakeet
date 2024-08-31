# How to Make an SLM Smarter

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

Today, my goal is to try to turn the **`Qwen:0.5b`** model into a fern specialist. To do this, I will provide it with contextual data about ferns, and I will ask it various questions.

## Prerequisites:

### Verify that the model or models I will use know nothing about ferns

To check that the **`Qwen2:0.5b`** model knows nothing about ferns, I will ask it a question about ferns. For this, I use the following script:

```bash
curl http://localhost:11434/api/chat \
-H "Content-Type: application/json" \
-d '
{
  "model": "qwen2:0.5b",
  "messages": [
    {
      "role": "system",
      "content": "You are an expert in botanics and ferns."
    },
    {
      "role": "user",
      "content": "Give me a list of ferns of the Dryopteridaceae variety"
    }
  ],
  "stream": false
}' | jq '.message.content'
```

✋ **Throughout my experimentation** I tested with several models to see which one best suited my requirements, particularly the ability to run comfortably on a Pi 5. My final choice was the **`Qwen2:0.5b`** model. So I will only present the results obtained with this model. But nothing stops you from testing with other models. I even encourage you to do so. I will also offer some hypotheses and conclusions to explain the choices and results obtained. This should then allow you to conduct your own experiments and adapt them to your needs.

### Preparing the Data

To start, I need to generate a data file on ferns that will be my only source of truth. For this, I created a `ferns.json` file by asking **ChatGPT4o** for help. I used the following prompt:

```
You are an expert in botanics and your main topic is about the ferns.

I want a list of 10 varieties of ferns with their characteristics and 5 ferns per variety.

the output format is in JSON:
```json
[
    {
        "variety": "name of the variety of ferns",
        "description": "Description of the variety of ferns",
        "ferns": [
            {
                "name": "scientific name of the fern",
                "common_name": "common name of the fern",
                "description": "description and characteristics of the fern"
            },
        ]
    },
]
```

Here is an excerpt from this file:

```json
[
    {
        "variety": "Polypodiaceae",
        "description": "A family of ferns known for their leathery fronds and widespread distribution in tropical and subtropical regions.",
        "ferns": [
            {
                "name": "Polypodium vulgare",
                "common_name": "Common Polypody",
                "description": "A hardy fern with leathery, evergreen fronds that thrive in rocky and shaded areas."
            },
            {
                "name": "Polypodium glycyrrhiza",
                "common_name": "Licorice Fern",
                "description": "Known for its sweet-tasting roots, this fern grows on moist, shaded rocks and tree trunks."
            },
```

📝 You can find the full file here: [ferns.json](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.json)

### Formatting the Data

Next, I transformed the JSON file into two markdown files with the same data but with different structures. Here are the structures of these files:

#### `ferns.1.md`

```markdown
# Title of the report

## Variety: name of the variety of ferns
*Description:*
Description of the variety of ferns

### Ferns:
- **name:** **scientific name of the fern**
  **common name:** common name of the fern
  **description:** description and characteristics of the fern
```

#### `ferns.2.md`

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

> The second file is a bit more structured than the first, thanks to paragraph titles. This will allow me to test if the structure of the data impacts the results obtained.

📝 You can find the full files here: [ferns.1.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.1.md) and [ferns.2.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.2.md)

## First Experiments

I decided to keep it simple and create a prompt that will consist of the following elements:

- Instructions for the model.
- The **complete** content of the `ferns.1.md` or `ferns.2.md` file, which I call the **context**.
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

📝 You can find the complete code here: [01-context](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/01-context)

### Questions

To conduct my tests, I used the following questions:

- **Question 1:** Give me a list of ferns of the Dryopteridaceae variety
- **Question 2:** What is the common name of Dryopteris cristata?

For this first series of tests, I used three models: **`Qwen:0.5b`**, **`Qwen2:0.5b`**, and **`CognitiveComputations/dolphin-gemma2`** (1.6 GB).

The results obtained are as follows:

#### With `ferns.1.md`

| LLM + ferns.1.md | Question 1 | Question 2 |
| ---------------- | ---------- | ---------- |
| qwen:0.5b        | 😡         | 😡         |
| qwen2:0.5b       | 😡         | 😡         |
| dolphin-gemma2   | 🙂         | 🙂         |

> - 😡: incorrect or incomplete result
> - 🙂: satisfactory result

#### With `ferns.2.md`

| LLM + ferns.2.md | Question 1 | Question 2 |
| ---------------- | ---------- | ---------- |
| qwen:0.5b        | 😡         | 😡         |
| qwen2:0.5b       | 😡         | 😡         |
| dolphin-gemma2   | 😡         | 🙂         |

> - 😡: incorrect or incomplete result
> - 🙂: satisfactory result

### Hypotheses and Observations

- **dolphin-gemma2** already has knowledge of ferns, so it was able to answer the question correctly by relying on its knowledge and the information provided (even though theoretically the instructions were to use only the provided context).
- **dolphin-gemma2** seems to be able to handle a large context.
- **qwen:0.5b** and **qwen2:0.5b** could not answer the questions correctly. This may be due to their ability to handle larger contexts.
- The data structure seems to have an impact on the results obtained with **dolphin-gemma2** - but I would have thought that the second document structure would be easier to use. We will see if this is confirmed later.

My first hypothesis or observation is that **qwen:0.5b** and **qwen2:0.5b** (very small LLMs) are not able to exploit a very large context, even if it is structured.

Therefore, I need to help **qwen:0.5b** and **qwen2:0.5b** "focus" on specific and "closer" data to the questions. So, I repeated the same experiments but reduced the context size.

## Second Series of Experiments: Reducing the Context Size

For this second series of experiments, I reduced the context size by keeping only the information on a single variety of ferns. I created a `ferns.1.extract.md` and `ferns.2.extract.md` file that contains information on the **Dryopteridaceae** variety.

📝 I use the same code to run the examples: [01-context](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/01-context)

**For `ferns.1.extract.md`:**

```markdown
# Fern Varieties Report

## Variety: Dryopteridaceae
*Description:*
A large family of ferns with robust and often leathery fronds, commonly found in woodlands.

### Ferns:
- **name:** **Dryopteris filix-mas**
  **common name:** Male Fern
  **description:** A sturdy fern with pinnate fronds, commonly found in temperate forests.

- **name:** **Dryopteris marginalis**
  **common name:** Marginal Wood Fern
  **description:** Known for its evergreen fronds, this fern thrives in rocky, shaded environments.

- **name:** **Dryopteris erythrosora**
  **common name:** Autumn Fern
  **description:** This fern features striking copper-red fronds that mature to green.

- **name:** **Dryopteris cristata**
  **common name:** Crested Wood Fern
  **description:** A fern with uniquely crested fronds, typically found in wetland areas.

- **name:** **Dryopteris affinis**
  **common name:** Golden Male Fern
  **description:** A robust fern with yellowish fronds and a preference for moist, shaded habitats.
```

**For `ferns.2.extract.md`:**

```markdown
# Fern Varieties Report

## Variety: Dryopteridaceae
*Description:*
A large family of ferns with robust and often leathery fronds, commonly found in woodlands.

### Ferns:

#### Dryopteris filix-mas

**name:** **Dryopteris filix-mas**  
**common name:** Male Fern  
**description:** A sturdy fern with pinnate fronds, commonly found in temperate forests.

#### Dryopteris marginalis

**name:** **Dryopteris marginalis**  
**common name:** Marginal Wood Fern  
**description:** Known for its evergreen fronds, this fern thrives in rocky, shaded environments.

#### Dryopteris erythrosora

**name:** **Dryopteris erythrosora**  
**common name:** Autumn Fern  
**description:** This fern features striking copper-red fronds that mature to green.

#### Dryopteris cristata

**name:** **Dryopteris cristata**  
**common name:** Crested Wood Fern  
**description:** A fern with uniquely crested fronds, typically found in wetland areas.

#### Dryopteris affinis

**name:** **Dryopteris affinis**  
**common name:** Golden Male Fern  
**description:** A robust fern with yellowish fronds and a preference for moist, shaded habitats.
```

📝 You can find the full files here: [ferns.1.extract.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.1.extract.md) and [ferns.2.extract.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.2.extract.md)

I then repeated the same experiments as before (same questions).

### Results

The results obtained are as follows:

| LLM + ferns.1.extract.md | Question 1 | Question 2 |
| ------------------------ | ---------- | ---------- |
| qwen:0.5b                | 😡         | 😡         |
| qwen2:0.5b               | 🙂         | 🙂         |
| dolphin-gemma2           | 🙂         | 😡         |

| LLM + ferns.2.extract.md | Question 1 | Question 2 |
| ------------------------ | ---------- | ---------- |
| qwen:0.5b                | 🙂         | 😡         |
| qwen2:0.5b               | 🙂         | 🙂         |
| dolphin-gemma2           | 🙂         | 🙂         |

Clearly, reducing the context size allowed **qwen2:0.5b** to better exploit the information and answer the questions correctly. **dolphin-gemma2** had mixed results, but it was still able to answer some questions correctly.

Moreover, the data structure seems to impact the results obtained with all three models. The second document structure seems to be more easily exploitable.

### Conclusion

✋ **I realize that my hypotheses and conclusions are based on a limited number of tests. It would be interesting to repeat these experiments with a larger number of models and questions.**

Nevertheless, for the continuation of my experiments, I will continue to use **qwen2:0.5b** and the second document structure.

Of course, I would like to be able to ask my fern expert about other varieties of ferns. To do this, I will need to set up a system that will be able to extract only the relevant information for a given question. We will therefore move on to the next phase of our experiments and do some **RAG** (Retrieval-augmented generation) to extract the relevant information.

## Third Series of Experiments: Similarity Search to Provide a More Relevant but Smaller Context

📝 This time the code to run the examples is here: [02-rag](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/02-rag)

This program will do several things:

- It will load the `ferns.2.split.md` document.
- Split it into several parts (one per fern variety).
- Calculate the vectors (embeddings) of each part (or chunk) with the help of an appropriate LLM for embedding generation. I will try **all-minilm:33m**, **nomic-embed-text**, and **mxbai-embed-large**.
- Wait for a user question.
- Calculate the vector of the question.
- Calculate the similarity or similarities between the question vector and the document part vectors.
- Create a context with the document part(s) that have the highest similarity to the question.
- Ask the question to the **qwen2:0.5b** model with the generated context.

> I will also try **qwen2:1.5b**.

✋ `ferns.2.split.md` is a markdown file that contains the same information as `ferns.2.md` (It is also structured the same way), but I added a **marker** `<!-- SPLIT -->` at the end of each fern variety section to indicate the different parts of the document. This will allow me to split the document into several parts and calculate the embeddings of each part.

📝 You can find the full file here: [ferns.2.split.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.2.split.md)

✋ For this experiment, I used an "in-memory vector store" to store the vectors of the document parts, and the "cosine distance" to calculate the similarity between the document part vectors and the question vector. These features are available in the [Parakeet](https://github.com/parakeet-nest/parakeet) project.

> Parakeet provides other features for RAG, notably with **Elasticsearch**. You can consult the documentation and examples for more information.

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

The

 function will search for vectors (in the vector store) that have a cosine distance greater than or equal to 0.5 with the question vector. And it will return the vector with the best score.

And if I want to search for the top 3 vectors:

```golang
similarities, err := store.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 3)
```

### Questions

For this new experiment, I used the same questions as in the previous experiment:

- **Question 1:** Give me a list of ferns of the Dryopteridaceae variety
- **Question 2:** What is the common name of Dryopteris cristata?

### Results

The results obtained are as follows:

|   | LLM + ferns.2.split.md | Question 1 | Question 2 | TopNSimilarities |
| - | ---------------------- | ---------- | ---------- | ---------------- |
| 1 | qwen2:0.5b             | 😡         | 🙂😡       | 3                |
| 2 | qwen2:0.5b             | 🙂         | 🙂😡       | 1                |

1. **Returns up to 3 similarities**: In the 1st case, **qwen2:0.5b** cannot answer the question correctly; it loops. Regarding question 2, it answers correctly if I only have one similarity. But it cannot answer correctly if I have more than one similarity. And sometimes it finds no similarity even though the information exists.
2. **Returns 1 similarity**: In the 2nd case, **qwen2:0.5b** correctly answers questions 1 and 2. But sometimes it finds no similarity even though the information exists.

#### First Hypotheses and Observations

- Only one similarity should be returned for **qwen2:0.5b** to answer the question correctly (to stay focused).
- No similarity should be returned for **qwen2:0.5b** to answer the question correctly (to have the information).
- I must therefore find a way to improve similarity search to ensure the information is found.

So I will use another model to generate embeddings: **nomic-embed-text**.

### Using nomic-embed-text

|   | LLM + ferns.2.split.md | Question 1 | Question 2 | TopNSimilarities |
| - | ---------------------- | ---------- | ---------- | ---------------- |
| 1 | qwen2:0.5b             | 😡         | 🙂         | 3                |
| 2 | qwen2:0.5b             | 🙂         | 😡         | 1                |

Even though there is some improvement, I still end up with unsatisfactory similarity search results (0 similarity even though the information exists). So I will try with **mxbai-embed-large**.

### Using mxbai-embed-large

|   | LLM + ferns.2.split.md | Question 1 | Question 2 | TopNSimilarities |
| - | ---------------------- | ---------- | ---------- | ---------------- |
| 1 | qwen2:0.5b             | 😡         | 🙂         | 3                |
| 2 | qwen2:0.5b             | 🙂         | 🙂😡       | 1                |

Clearly, I must limit it to one similarity for **qwen2:0.5b** to answer the question correctly.

The use of **mxbai-embed-large** brought a significant improvement. I no longer have unsatisfactory similarity search results (the information is found correctly).

However, **qwen2:0.5b** does not always answer question 2 correctly and uses information from another fern of the same variety.

So I will do the same test with other models to see if I can get better results.

## Fourth Series of Experiments: Trying Other Models

I keep **mxbai-embed-large** for embedding generation, and I will try other models for completion.

|   | LLM + ferns.2.split.md | Question 1 | Question 2 | TopNSimilarities |
| - | ---------------------- | ---------- | ---------- | ---------------- |
| 1 | qwen2:1.5b             | 🙂         | 🙂         | 3                |
| 2 | qwen2:1.5b             | 🙂         | 🙂         | 1                |
| 3 | tinydolphin            | 😡         | 😡         | 3                |
| 4 | tinydolphin            | 😡🙂       | 🙂         | 1                |

The model parameter sizes are as follows:

- **qwen2:0.5b** (352 MB) `0.5b`
- **qwen2:1.5b** (934 MB) `1.5b`
- **tinydolphin** (636 MB) `1.1b`

Unfortunately, **qwen2:0.5b** does not fully satisfy my fern expert use case.

**qwen2:1.5b** is much better than **qwen2:0.5b**. It answers both questions 1 and 2 correctly. It can answer correctly even if I return 3 similarities.

**tinydolphin** can also answer questions 1 and 2 correctly. But for question 1, it sometimes returns duplicate results. However, it is necessary to limit similarity search to one result to get a satisfactory answer.

I wonder if I could help **tinydolphin** a bit by providing it with more precise and structured context to allow it to answer question 1 correctly (getting the list of ferns of a given variety).

## Fifth Series of Experiments: Trying with a More Precise and Structured Context

I created a new file, `ferns.2.split.list.md`, which contains the same information as `ferns.2.split.md`, but at the end of each fern variety section, I added a list of the fern names of the variety. Like this, for example:

```markdown
## List of Ferns of the variety: Dryopteridaceae

- **Dryopteris filix-mas** (Male Fern)
- **Dryopteris marginalis** (Marginal Wood Fern)
- **Dryopteris erythrosora** (Autumn Fern)
- **Dryopteris cristata** (Crested Wood Fern)
- **Dryopteris affinis** (Golden Male Fern)
```

In fact, I present the same information several times in the document. But in a different form.

📝 You can find the full file here: [ferns.2.split.list.md](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/data/ferns.2.split.list.md)

I keep **mxbai-embed-large** for embedding generation, and I will perform the same tests as in the previous experiment but with the `ferns.2.split.list.md` file.

📝 This time the code to run the examples is here: [03-rag-list](https://github.com/parakeet-nest/parakeet/tree/main/examples/42-make-a-slm-smarter/03-rag-list)

### Results

|   | LLM + ferns.2.split.list.md | Question 1 | Question 2 | TopNSimilarities |
| - | --------------------------- | ---------- | ---------- | ---------------- |
| 1 | qwen2:1.5b                  | 🙂         | 🙂         | 3                |
| 2 | qwen2:1.5b                  | 🙂         | 🙂         | 1                |
| 3 | tinydolphin                 | 😡         | 😡🙂       | 3                |
| 4 | tinydolphin                 | 🙂         | 🙂         | 1                |

This time, by adding additional information to the context, I was able to get better results with **tinydolphin**. It answers questions 1 and 2 correctly if I keep only one similarity.

## Conclusion

This type of experimentation is very interesting but can last indefinitely. It is important to clearly define the objectives and constraints of the experiment to avoid getting lost.

For my use case, my conclusions are as follows:

To enable an SLM to be a fern expert, it is necessary to provide it with contextual data on ferns. It is important to structure this data so that it is easily exploitable by the model. It is also important to limit the context size to allow the model to process the information efficiently.

And I would retain the following candidates to create a good fern expert:

|   | LLM + ferns.2.split.list.md | ferns.2.list.md | ferns.2.split.list.md | TopNSimilarities |
| - | --------------------------- | --------------- | --------------------- | ---------------- |
| 1 | qwen2:1.5b                  | ✅              | ✅                    | 3                |
| 2 | qwen2:1.5b                  | ✅              | ✅                    | 1                |
| 4 | tinydolphin                 |                 | ✅                    | 1                |

I encourage you to conduct your own experiments and adapt the concepts presented here to your needs. 
