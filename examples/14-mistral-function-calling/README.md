# Function Calling with Mistral

https://ollama.com/library/mistral
Mistral 0.3 supports function calling with Ollama’s raw mode.

You need to use the raw mode for function calling.
https://github.com/ollama/ollama/blob/main/docs/api.md#request-raw-mode

In some cases, you may wish to bypass the templating system and provide a full prompt. In this case, you can use the raw parameter to disable templating. Also note that raw mode will not return a context.

Then you can define the list of tools like this and call a tool in the prompt:

```
[AVAILABLE_TOOLS] [{"type": "function", "function": {"name": "get_current_weather", "description": "Get the current weather", "parameters": {"type": "object", "properties": {"location": {"type": "string", "description": "The city and state, e.g. San Francisco, CA"}, "format": {"type": "string", "enum": ["celsius", "fahrenheit"], "description": "The temperature unit to use. Infer this from the users location."}}, "required": ["location", "format"]}}}][/AVAILABLE_TOOLS][INST] What is the weather like today in San Francisco [/INST]
```
