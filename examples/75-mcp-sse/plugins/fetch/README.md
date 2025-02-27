# Build the plugin with the tinygo builder image

**Build**:
```bash
docker run --rm -v "$PWD":/src -w /src k33g/tinygo-builder:preview \
  tinygo build -scheduler=none --no-debug \
  -o wasimancer-plugin-fetch.wasm \
  -target wasi main.go
```

**Run**:
```bash
docker run --rm -v "$PWD":/app -w /app k33g/tinygo-builder:preview \
  extism call wasimancer-plugin-fetch.wasm fetch \
  --input '{"url":"https://modelcontextprotocol.io/introduction"}' \
  --allow-host "*" \
  --log-level "info" \
  --wasi
```