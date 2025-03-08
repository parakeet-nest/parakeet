#!/bin/bash
extism call wasimancer-plugin-fetch.wasm fetch \
  --input '{"url": "https://modelcontextprotocol.io/introduction"}' \
  --allow-host "*" \
  --log-level "info" \
  --wasi
echo ""
