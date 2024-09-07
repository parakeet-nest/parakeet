<!-- TOPIC: Function Calling SUMMARY: A feature in LLMs that allows them to provide a specific output with the same format (predictable output format). KEYWORDS: function calling, predictable output format, LLMs. -->
# Function Calling (before tool support)

!!! info "almost depecrated"

What is **"Function Calling"**? First, it's not a feature where a LLM can call and execute a function. "Function Calling" is the ability for certain LLMs to provide a specific output with the same format (we could say: "a predictable output format").

So, the principle is simple:

- You (or your GenAI application) will create a prompt with a delimited list of tools (the functions) composed by name, descriptions, and parameters: `SayHello`, `AddNumbers`, etc.
- Then, you will add your question ("Hey, say 'hello' to Bob!") to the prompt and send all of this to the LLM.
- If the LLM "understand" that the `SayHello` function can be used to say "hello" to Bob, then the LLM will answer with only the name of the function with the parameter(s). For example: `{"name":"SayHello","arguments":{"name":"Bob"}}`.

Then, it will be up to you to implement the call of the function.

The [latest version (v0.3) of Mistral 7b](https://ollama.com/library/mistral:7b) supports function calling and is available for Ollama.

## Define a list of tools

First, you have to provide the LLM with a list of tools with the following format:

```golang
toolsList := []llm.Tool{
    {
        Type: "function",
        Function: llm.Function{
            Name:        "hello",
            Description: "Say hello to a given person with his name",
            Parameters: llm.Parameters{
                Type: "object",
                Properties: map[string]llm.Property{
                    "name": {
                        Type:        "string",
                        Description: "The name of the person",
                    },
                },
                Required: []string{"name"},
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
```

## Generate a prompt from the tools list and the user instructions

The `tools.GenerateContent` method generates a string with the tools in JSON format surrounded by `[AVAILABLE_TOOLS]` and `[/AVAILABLE_TOOLS]`:
```golang
toolsContent, err := tools.GenerateContent(toolsList)
if err != nil {
    log.Fatal("😡:", err)
}
```


The `tools.GenerateInstructions` method generates a string with the user instructions surrounded by `[INST]` and `[/INST]`:
```golang
userContent := tools.GenerateInstructions(`say "hello" to Bob`)
```

Then, you can add these two strings to the messages list:
```golang
messages := []llm.Message{
    {Role: "system", Content: toolsContent},
    {Role: "user", Content: userContent},
}
```

## Send the prompt (messages) to the LLM

It's important to set the `Temperature` to `0.0`:
```golang
options := llm.Options{
    Temperature:   0.0,
    RepeatLastN:   2,
    RepeatPenalty: 2.0,
}

You must set the `Format` to `json` and `Raw` to `true`:
query := llm.Query{
    Model: model,
    Messages: messages,
    Options: options,
    Format:  "json",
    Raw:     true,
}
```
> When building the payload to be sent to Ollama, we need to set the `Raw` field to true, thanks to that, no formatting will be applied to the prompt (we override the prompt template of Mistral), and we need to set the `Format` field to `"json"`.

No you can call the `Chat` method. The answer of the LLM will be in JSON format:
```golang
answer, err := completion.Chat(ollamaUrl, query)
if err != nil {
    log.Fatal("😡:", err)
}
// PrettyString is a helper that prettyfies the JSON string
result, err := gear.PrettyString(answer.Message.Content)
if err != nil {
    log.Fatal("😡:", err)
}
fmt.Println(result)
```

You should get this answer:
```json
{
  "name": "hello",
  "arguments": {
    "name": "Bob"
  }
}
```

You can try with the other tool (or function):
```golang
userContent := tools.GenerateInstructions(`add 2 and 40`)
```

You should get this answer:
```json
{
  "name": "addNumbers",
  "arguments": {
    "a": 2,
    "b": 40
  }
}
```

> **Remark**: always test the format of the output, even if Mistral is trained for "function calling", the result are not entirely predictable.

Look at this sample for a complete sample: [examples/15-mistral-function-calling](https://github.com/parakeet-nest/parakeet/tree/main/examples/15-mistral-function-calling)
