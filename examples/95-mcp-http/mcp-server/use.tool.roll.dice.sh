#!/bin/bash
: <<'COMMENT'
# Use tool "add"
COMMENT

# STEP 1: Load the session ID from the environment file
source mcp.env

MCP_SERVER=${MCP_SERVER:-"http://localhost:9090"}

read -r -d '' DATA <<- EOM
{
  "jsonrpc": "2.0",
  "id": "test",
  "method": "tools/call",
  "params": {
    "name": "rool_dices",
    "arguments": {
      "nb_dices": 5,
      "nb_sides": 8
    }
  }
}
EOM

curl ${MCP_SERVER}/mcp \
  -H "Content-Type: application/json" \
  -H "Mcp-Session-Id: $SESSION_ID" \
  -d "${DATA}" | jq 


