# MCP: model context protocol, SSE CLient


## Start Docker distribution of WASImancer

```bash
docker run --rm -p 5001:3001 \
  -e HTTP_PORT=3001 \
  -e PLUGINS_PATH=./plugins \
  -e PLUGINS_DEFINITION_FILE=plugins.yml \
  -v "$(pwd)/plugins":/app/plugins \
  -e RESOURCES_PATH=./resources \
  -e RESOURCES_DEFINITION_FILE=resources.yml \
  -v "$(pwd)/resources":/app/resources \
  k33g/wasimancer:preview 
```

# Build the plugin with the tinygo builder image

**Build**:
```bash
cd plugins/fetch
docker run --rm -v "$PWD":/src -w /src k33g/tinygo-builder:preview \
  tinygo build -scheduler=none --no-debug \
  -o wasimancer-plugin-fetch.wasm \
  -target wasi main.go
```

**Run**:
```bash
docker run --rm -v "$PWD":/app -w /app k33g/tinygo-builder:preview \
  extism call wasimancer-plugin-fetch.wasm fetch \
  --input '{"url":"https://raw.githubusercontent.com/sea-monkeys/WASImancer/main/README.md"}' \
  --log-level "info" \
  --allow-host "*" \
  --wasi
```