package llm

import (
	"encoding/json"
	"fmt"
	"testing"
	//"github.com/parakeet-nest/parakeet/tools"
)

func TestTools(t *testing.T) {

	helloTool := Tool{
		Type: "function",
		Function: Function{
			Name:        "hello",
			Description: "Say hello to a given person with his name",
			Parameters: Parameters{
				Type: "object",
				Properties: map[string]Property{
					"name": {
						Type:        "string",
						Description: "The name of the person",
					},
				},
				Required: []string{"name"},
			},
		},
	}

	addNumbersTool := Tool{
		Type: "function",
		Function: Function{
			Name:        "addNumbers",
			Description: "Make an addition of the two given numbers",
			Parameters: Parameters{
				Type: "object",
				Properties: map[string]Property{
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
	}

	tools := []Tool{
		helloTool, addNumbersTool,
	}

	// marshall tools to JSON string
	toolsJSON, err := json.Marshal(&tools)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("📝 toolsJSON", string(toolsJSON))

}
