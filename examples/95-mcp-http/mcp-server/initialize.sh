#!/bin/bash
: <<'COMMENT'
# Initialize

MCP (Model Context Protocol) requires a proper initialization handshake 
before any other operations can be performed. 
This is by design - it ensures that both client and server agree 
on protocol version and capabilities before exchanging data.
COMMENT

MCP_SERVER=${MCP_SERVER:-"http://localhost:9090"}


#MCP_SERVER=http://host.docker.internal:9090/mcp


# STEP 1: Initialize the server
read -r -d '' INIT_DATA <<- EOM
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "id": "init-uuid",
  "params": {
    "protocolVersion": "2024-11-05"
  }
}
EOM

echo "Sending initialization request..."

# STEP 2: Get the session ID

SESSION_ID=$(curl -i -s -X POST \
  -H "Content-Type: application/json" \
  -d "${INIT_DATA}" \
  ${MCP_SERVER}/mcp | grep -i "mcp-session-id:" | cut -d' ' -f2 | tr -d '\r\n')

echo "🌍 Session ID: '$SESSION_ID'"
echo ""

# STEP 3: Get the server's JSON response (without headers)

JSON_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "${INIT_DATA}" \
  ${MCP_SERVER}/mcp)

echo "📝 Response:"
echo "$JSON_RESPONSE" | jq '.'

# STEP 4: Store the session ID in an environment variable
export MCP_SESSION_ID="$SESSION_ID"
echo ""
echo "SESSION_ID=$SESSION_ID" >| mcp.env
echo "✅ Session ID stored in mcp.env"
