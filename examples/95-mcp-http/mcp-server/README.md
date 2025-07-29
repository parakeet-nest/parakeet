# MCP HTTP Server Example

Build and start the MCP HTTP server with the following commands:
```bash
docker build --platform linux/arm64 -t mcp-http:demo .
docker run --rm -p 9090:9090 mcp-http:demo
```