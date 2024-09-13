<!-- TOPIC: Function Calling with LLMs that do not implement Function Calling SUMMARY: A technique to reproduce function calling feature in LLMs without native support by adding specific messages at the beginning and end of the conversation. KEYWORDS: phi3, mini, golang, JSON, tool calling, argument passing -->
# Function Calling with LLMs that do not implement tools support

It is possible to reproduce this feature with some LLMs that do not implement the "Function Calling" feature natively, but we need to supervise them and explain precisely what we need. The result (the output) will be less predictable, so you will need to add some tests before using the output, but with some "clever" LLMs, you will obtain correct results. I did my experiments with **[phi3:mini](https://ollama.com/library/phi3:mini)**.

The trick is simple:

Add this message at the begining of the list of messages:
```golang
systemContentIntroduction := `You have access to the following tools:`
```

Add this message at the end of the list of messages, just before the user message:
```golang
systemContentInstructions := `If the question of the user matched the description of a tool, the tool will be called.
To call a tool, respond with a JSON object with the following structure: 
{
    "name": <name of the called tool>,
    "arguments": {
    <name of the argument>: <value of the argument>
    }
}

search the name of the tool in the list of tools with the Name field
`
```

At the end, you will have this:
```golang
messages := []llm.Message{
    {Role: "system", Content: systemContentIntroduction},
    {Role: "system", Content: toolsContent},
    {Role: "system", Content: systemContentInstructions},
    {Role: "user", Content: `say "hello" to Bob`},
}
```

!!! note
	Look at this sample for a complete sample: [examples/17-fake-function-calling](https://github.com/parakeet-nest/parakeet/tree/main/examples/17-fake-function-calling)
<!-- split -->

