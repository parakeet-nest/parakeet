<!-- TOPIC: Parakeet - A Go Library for Creating GenAI Apps SUMMARY: Parakeet is a simple Go library used to create text-based GenAI apps, allowing users to generate new content based on training data. KEYWORDS: Parakeet, GenAI, Go, Library, Text Generation, AI -->
# 🦜🪺 Parakeet

Parakeet is the simplest Go library to create **GenAI apps** with **[Ollama](https://ollama.com/)**.

> A GenAI app is an application that uses generative AI technology. Generative AI can create new text, images, or other content based on what it's been trained on. So a GenAI app could help you write a poem, design a logo, or even compose a song! These are still under development, but they have the potential to be creative tools for many purposes. - [Gemini](https://gemini.google.com)

> ✋ Parakeet is only for creating GenAI apps generating **text** (not image, music,...).

## Install

!!! note
	current release: `v0.2.2 🧁 [cupcake]`

```bash
go get github.com/parakeet-nest/parakeet
```
<!-- split -->

<!-- TOPIC: Simple Completion in Golang using Parakeet and LLaMA SUMMARY: This code snippet demonstrates the use of simple completion in Golang to generate a response for a given prompt with a provided model, specifically using Parakeet and LLaMA. KEYWORDS: Golang, Parakeet, LLaMA, Simple Completion, AI-powered Text Generation -->
## 🚀 Getting Started - First completion
> `generate`

The simple completion can be used to generate a response for a given prompt with a provided model.

```golang
package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"

	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "tinydolphin"

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.5,
	})

	question := llm.GenQuery{
		Model: model,
		Prompt: "Who is James T Kirk?",
		Options: options,
	}

	answer, err := completion.Generate(ollamaUrl, question)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)
}
```
<!-- split -->

