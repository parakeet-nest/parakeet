<!-- TOPIC: Parakeet - A Go Library for Creating GenAI Apps SUMMARY: Parakeet is a simple Go library used to create text-based GenAI apps, allowing users to generate new content based on training data. KEYWORDS: Parakeet, GenAI, Go, Library, Text Generation, AI -->

# 🦜🪺 Parakeet

Parakeet is the simplest Go library to create **GenAI apps** with **[Ollama](https://ollama.com/)**.

> A GenAI app is an application that uses generative AI technology. Generative AI can create new text, images, or other content based on what it's been trained on. So a GenAI app could help you write a poem, design a logo, or even compose a song! These are still under development, but they have the potential to be creative tools for many purposes. - [Gemini](https://gemini.google.com)

> ✋ Parakeet is only for creating GenAI apps generating **text** (not image, music,...).

## Install

```bash
go get github.com/parakeet-nest/parakeet
```

## Some examples

## Chat with streaming

```golang
ollamaUrl := "http://localhost:11434"
model := "deepseek-coder"

systemContent := `You are an expert in computer programming.
Please make friendly answer for the noobs.
Add source code examples if you can.`

userContent := `Ccreate a "hello world" program in Golang.`

options := llm.SetOptions(map[string]interface{}{
    option.Temperature:   0.5,
    option.RepeatLastN:   2,
    option.RepeatPenalty: 2.2,
})

query := llm.Query{
    Model: model,
    Messages: []llm.Message{
        {Role: "system", Content: systemContent},
        {Role: "user", Content: userContent},
    },
    Options: options,
}

_, err := completion.ChatStream(ollamaUrl, query,
    func(answer llm.Answer) error {
        fmt.Print(answer.Message.Content)
        return nil
    })
```

## Tools (function calling)

```golang

ollamaUrl := "http://localhost:11434"
model := "allenporter/xlam:1b"

toolsList := []llm.Tool{
    {
        Type: "function",
        Function: llm.Function{
            Name:        "multiplyNumbers",
            Description: "Make a multiplication of the two given numbers",
            Parameters: llm.Parameters{
                Type: "object",
                Properties: map[string]llm.Property{
                    "a": {
                        Type:        "number",
                        Description: "first operand",
                    },
                    "b": {
                        Type:        "number",
                        Description: "second operand",
                    },
                },
                Required: []string{"a", "b"},
            },
        },
    },
    {
        Type: "function",
        Function: llm.Function{
            Name:        "addNumbers",
            Description: "Make an addition of the two given numbers",
            Parameters: llm.Parameters{
                Type: "object",
                Properties: map[string]llm.Property{
                    "a": {
                        Type:        "number",
                        Description: "first operand",
                    },
                    "b": {
                        Type:        "number",
                        Description: "second operand",
                    },
                },
                Required: []string{"a", "b"},
            },
        },
    },
}

messages := []llm.Message{
    {Role: "user", Content: `add 2 and 40`},
    {Role: "user", Content: `multiply 2 and 21`},
}

options := llm.SetOptions(map[string]interface{}{
    option.Temperature: 0.0,
    option.RepeatLastN: 2,
    option.RepeatPenalty: 2.0,
})

query := llm.Query{
    Model:    model,
    Messages: messages,
    Tools:    toolsList,
    Options:  options,
    Format:   "json",
}

answer, err := completion.Chat(ollamaUrl, query)
if err != nil {
    log.Fatal("😡:", err)
}

for idx, toolCall := range answer.Message.ToolCalls {
    result, err := toolCall.Function.ToJSONString()
    if err != nil {
        log.Fatal("😡:", err)
    }
    // displqy the tool to call
    fmt.Println("ToolCall", idx, ":", result)

    /* Results:
    ToolCall 0 : {"name":"addNumbers","arguments":{"a":2,"b":40}}
    ToolCall 1 : {"name":"multiplyNumbers","arguments":{"a":2,"b":21}}
    */
}
```

## Structured Output

```golang
ollamaUrl := "http://localhost:11434"
model := "qwen2.5:0.5b"

options := llm.SetOptions(map[string]interface{}{
    option.Temperature: 1.5,
})

// define schema for a structured output
schema := map[string]any{
    "type": "object",
    "properties": map[string]any{
        "name": map[string]any{
            "type": "string",
        },
        "capital": map[string]any{
            "type": "string",
        },
        "languages": map[string]any{
            "type": "array",
            "items": map[string]any{
                "type": "string",
            },
        },
    },
    "required": []string{"name", "capital", "languages"},
}

query := llm.Query{
    Model: model,
    Messages: []llm.Message{
        {Role: "user", Content: "Tell me about Canada."},
    },
    Options: options,
    Format:  schema,
    Raw:     false,
}

answer, err := completion.Chat(ollamaUrl, query)

fmt.Println(answer.Message.Content)

/* Results:
{
  "capital": "Ottawa",
  "languages": ["English", "French"],
  "name": "Canada of the West: Land of Ice and Rainbows"
}
*/
```

## Quick RAG

```golang
docs := []string{
    `Michael Burnham is the main character on the Star Trek series, Discovery.  
    She's a human raised on the logical planet Vulcan by Spock's father.  
    Burnham is intelligent and struggles to balance her human emotions with Vulcan logic.  
    She's become a Starfleet captain known for her determination and problem-solving skills.
    Originally played by actress Sonequa Martin-Green`,

    `James T. Kirk, also known as Captain Kirk, is a fictional character from the Star Trek franchise.  
    He's the iconic captain of the starship USS Enterprise, 
    boldly exploring the galaxy with his crew.  
    Originally played by actor William Shatner, 
    Kirk has appeared in TV series, movies, and other media.`,

    `Jean-Luc Picard is a fictional character in the Star Trek franchise.
    He's most famous for being the captain of the USS Enterprise-D,
    a starship exploring the galaxy in the 24th century.
    Picard is known for his diplomacy, intelligence, and strong moral compass.
    He's been portrayed by actor Patrick Stewart.`,

    `Lieutenant Philippe Charrière, known as the **Silent Sentinel** of the USS Discovery, 
    is the enigmatic programming genius whose codes safeguard the ship's secrets and operations. 
    His swift problem-solving skills are as legendary as the mysterious aura that surrounds him. 
    Charrière, a man of few words, speaks the language of machines with unrivaled fluency, 
    making him the crew's unsung guardian in the cosmos. His best friend is Spiderman from the Marvel Cinematic Universe.`,
}

ollamaUrl := "http://localhost:11434"
embeddingsModel := "mxbai-embed-large:latest" // This model is for the embeddings of the documents
smallChatModel := "qwen2.5:1.5b"   // This model is for the chat completion

store := embeddings.MemoryVectorStore{
    Records: make(map[string]llm.VectorRecord),
}

// Create embeddings from documents and save them in the store
for idx, doc := range docs {
    fmt.Println("Creating embedding from document ", idx)
    embedding, err := embeddings.CreateEmbedding(
        ollamaUrl,
        llm.Query4Embedding{
            Model:  embeddingsModel,
            Prompt: doc,
        },
        strconv.Itoa(idx),
    )
    if err != nil {
        fmt.Println("😡:", err)
    } else {
        store.Save(embedding)
    }
}

// Question for the Chat system
userContent := `Who is Philippe Charrière and what spaceship does he work on?`

systemContent := `You are an AI assistant. Your name is Seven. 
Some people are calling you Seven of Nine.
You are an expert in Star Trek.
All questions are about Star Trek.
Using the provided context, answer the user's question
to the best of your ability using only the resources provided.`

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

//🔎 searching for similarity...
similarity, _ := store.SearchMaxSimilarity(embeddingFromQuestion)

documentsContent := `<context><doc>` + similarity.Prompt + `</doc></context>`

query := llm.Query{
    Model: smallChatModel,
    Messages: []llm.Message{
        {Role: "system", Content: systemContent},
        {Role: "system", Content: documentsContent},
        {Role: "user", Content: userContent},
    },
    Options: llm.SetOptions(map[string]interface{}{
        option.Temperature: 0.4,
        option.RepeatLastN: 2,
    }),
}

fmt.Println("🤖 answer:")

// Answer the question
_, err = completion.ChatStream(ollamaUrl, query,
    func(answer llm.Answer) error {
        fmt.Print(answer.Message.Content)
        return nil
    })
```