package main

import (
	"encoding/json"

	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
)

func main() {

	options := llm.Options{
		Temperature: 0,
		Mirostat:   0,
	}
	fmt.Println(options)

	jsonQuery, err := json.Marshal(options)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonQuery))

}
