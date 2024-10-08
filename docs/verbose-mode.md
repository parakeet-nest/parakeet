# Verbose mode

You can activate the "verbose mode" with all kinds of completions.


```golang
options := llm.SetOptions(map[string]interface{}{
  option.Temperature: 0.5,
  option.RepeatLastN: 2,
  option.RepeatPenalty: 2.0,
  option.Verbose: true,
})
```
You will get an output like this (with the query and the completion):

```json
[llm/query] {
  "model": "deepseek-coder",
  "messages": [
    {
      "role": "system",
      "content": "You are an expert in computer programming.\n\tPlease make friendly answer for the noobs.\n\tAdd source code examples if you can."
    },
    {
      "role": "user",
      "content": "I need a clear explanation regarding the following question:\n\tCan you create a \"hello world\" program in Golang?\n\tAnd, please, be structured with bullet points"
    }
  ],
  "options": {
    "repeat_last_n": 2,
    "temperature": 0.5,
    "repeat_penalty": 2,
    "Verbose": true
  },
  "stream": false,
  "prompt": "",
  "context": null,
  "tools": null,
  "TokenHeaderName": "",
  "TokenHeaderValue": ""
}

[llm/completion] {
  "model": "deepseek-coder",
  "message": {
    "role": "assistant",
    "content": "Sure, here's a simple \"Hello, World!\" program in Golang.\n\t1. First, you need to have Golang installed on your machine.\n\t2. Open your text editor, and write the following code:\n\t```go\n\tpackage main\n\timport \"fmt\"\n\tfunc main() {\n\t    fmt.Println(\"Hello, World!\")\n\t} \n\t```\n\t3. Save the file with a `.go` extension (like `hello.go`).\n\t4. In your terminal, navigate to the directory containing the `.go` file.\n\t5. Run the program with the command:\n\t```\n\tgo run hello.go\n\t```\n\t6. If everything goes well, you should see \"Hello, World!\" printed in your terminal.\n\t7. If there's an error, you will see the error message.\n\t8. If everything is correct, you'll see \"Hello, World!\" printed in your terminal.\n"
  },
  "done": true,
  "response": "",
  "context": null,
  "created_at": "2024-08-19T05:57:23.979361Z",
  "total_duration": 3361191958,
  "load_duration": 2044932125,
  "prompt_eval_count": 79,
  "prompt_eval_duration": 95034000,
  "eval_count": 222,
  "eval_duration": 1216689000
}
```
