# MCP Curl server

This is a simple server that creates a tool which can fetch web pages using the `curl` command. Let me break it down:

1. The server part:
   - Creates an MCP server called "mcp-curl" with version "1.0.0"
   - MCP stands for Model Control Protocol - it's a way for AI models to interact with external tools

2. The tool part:
   - Creates a tool called "use_curl" 
   - The tool takes one required parameter: a URL string
   - When called, it will fetch the webpage at that URL

3. The handler part:
   - The `curlHandler` function is what actually does the work
   - It takes the URL from the request
   - Uses the system's `curl` command to fetch the webpage
   - Returns the webpage content as text
   - If anything goes wrong, it returns an error message

4. Server operation:
   - Runs over standard input/output (stdio)

This server essentially acts as a bridge between an AI model and the `curl` command, allowing the AI to fetch web content when needed.


## Build it

```bash
docker build -t mcp-curl .
```

## Use it with a MCP client (like Claude.AI Desktop)

```bash
{
  "mcpServers": {
    "mcp-curl-with-docker" :{
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "mcp-curl"
      ]
    }
  }
}
```
