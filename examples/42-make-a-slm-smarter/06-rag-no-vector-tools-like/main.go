package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/tools"
)

type SearchTool struct {
	Name      string `json:"name"`
	Arguments struct {
		Name string `json:"name"`
	} `json:"arguments"`
}

func IsSearchTool(ollamaUrl string, model string, nativeFunction bool, question string) (bool, SearchTool, error) {

	toolsList := []llm.Tool{
		{
			Type: "function",
			Function: llm.Function{
				Name:        "variety_list",
				Description: "I want the list of a given fern variety. Use only the list command.",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"name": {
							Type:        "string",
							Description: "name of the fern variety",
						},
					},
					Required: []string{"name"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.Function{
				Name:        "commom_name",
				Description: "I want the common name of a given fern. Use only the common command.",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"name": {
							Type:        "string",
							Description: "name of the fern",
						},
					},
					Required: []string{"name"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.Function{
				Name:        "variety_information",
				Description: "I want the information of a given fern variety. Use only the variety information command.",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"name": {
							Type:        "string",
							Description: "name of the fern variety",
						},
					},
					Required: []string{"name"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.Function{
				Name:        "fern_details",
				Description: "I want the details of a given fern. Use only the details command.",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"name": {
							Type:        "string",
							Description: "name of the fern",
						},
					},
					Required: []string{"name"},
				},
			},
		},
	}

	toolsContent, err := tools.GenerateAvailableToolsContent(toolsList)

	fmt.Println("🟡 Tools content", toolsContent)

	if err != nil {
		return false, SearchTool{}, err
	}

	// only if implemented

	options := llm.Options{
		Temperature:   0.0,
		RepeatLastN:   2,
		RepeatPenalty: 2.0,
		Seed:          123,
	}

	query := llm.Query{}

	if nativeFunction {
		question = tools.GenerateUserToolsInstructions(question)

		query = llm.Query{
			Model: model,
			Messages: []llm.Message{
				{Role: "system", Content: toolsContent},
				{Role: "user", Content: question},
			},
			Options: options,
			Format:  "json",
			Raw:     true, // try with false
		}
	} else { // The LLM do not implement the tools natively
		systemContentIntroduction := `You have access to the following tools:`
		systemContentInstructions := tools.GenerateSystemToolsInstructions()

		query = llm.Query{
			Model: model,
			Messages: []llm.Message{
				{Role: "system", Content: systemContentIntroduction},
				{Role: "system", Content: toolsContent},
				{Role: "system", Content: systemContentInstructions},
				{Role: "user", Content: question},
			},
			Options: options,
			Format:  "json",
			Raw:     true, // try with false
		}
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		return false, SearchTool{}, err
	}

	result, err := gear.PrettyString(answer.Message.Content)
	if err != nil {
		return false, SearchTool{}, err
	}

	fmt.Println("🟢 Search tool", result)

	searchTool := SearchTool{}

	err = json.Unmarshal([]byte(result), &searchTool)
	if err != nil {
		// that means the JSON string is not a SearchTool or it is not a valid JSON string
		// return false / empty MoveTool but not an error
		// then the chat model will handle the question
		return false, SearchTool{}, err
	}

	// ✋ the tool will always try to return a tool, even if it is not a valid tool

	if searchTool.Name == "variety_list" ||
		searchTool.Name == "variety_information" ||
		searchTool.Name == "fern_details" ||
		searchTool.Name == "common_name" {
		return true, searchTool, nil
	} else {
		return false, SearchTool{}, errors.New("the tool is not a valid tool")
	}

}

/*
SLM with tools: the format are different than the one used with Mistral

https://ollama.com/sam4096/qwen2tools:0.5b
ollama pull sam4096/qwen2tools:0.5b

https://ollama.com/allenporter/xlam:1b
ollama pull allenporter/xlam:1b
*/

// TODO: try something simpler (in another sample?) with 13-simple-fake-function-calling
// TODO: for the current example, use only native tools (function calls)

func main() {

	ollamaUrl := "http://localhost:11434"

	//smallChatModel := "qwen2:1.5b"
	//toolModel := "dolphin-phi:2.7b"
	//toolModel := "phi3:mini"
	//toolModel := "allenporter/xlam:1b"
	toolModel := "mistral:7b"
	nativeFunction := true

	//toolModel := "qwen2:1.5b"
	// nativeFunction := false
	smallChatModel := "tinydolphin"

	systemContent := `**Instruction:**
	You are an expert in botanics.
	Please use only the content provided below to answer the question.
	Do not add any external knowledge or assumptions.`

	fmt.Println(systemContent)
	fmt.Println()

	/*
		- variety_list     **Question:** Give me a list of ferns of the Dryopteridaceae variety
		- commom_name      **Question:** What is the common name of Dryopteris cristata?
		- commom_name      **Question:** What is the common name of this fern: Dryopteris cristata?
		- fern_details     **Question:** I want details about of the fern Dryopteris cristata
		- variety_information **Question:** I want information about of the fern variety Dryopteridaceae

	*/

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [%s](%s) ask me something> ", toolModel, smallChatModel)
		question, _ := reader.ReadString('\n')
		question = strings.TrimSpace(question)

		if question == "bye" {
			break
		}

		fmt.Println("🔎 searching for the best similarity...")

		isSearch, searchTool, err := IsSearchTool(ollamaUrl, toolModel, nativeFunction, question)

		//contextContext := ""

		if err != nil {
			fmt.Println("😡:", err)
		}
		if isSearch {
			fmt.Println("🎉🔎 search tool found: name", searchTool.Name, "argument:", searchTool.Arguments.Name)

			// Call the tool + context build from query in the json file

		} else {
			question = "Explain that you do not understand the question."
			fmt.Println("🤔", question)
		}

		//contextContext := embeddings.GenerateContentFromSimilarities(similarities)

		/*
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
		*/

	}

}
