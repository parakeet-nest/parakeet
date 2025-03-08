# Build the plugin

**Build**:
```bash
tinygo build -scheduler=none --no-debug \
  -o wasimancer-plugin-fetch.wasm \
  -target wasi main.go
```

**Run**:
```bash
extism call wasimancer-plugin-fetch.wasm fetch \
  --input '{"url": "https://modelcontextprotocol.io/introduction"}' \
  --allow-host "*" \
  --log-level "info" \
  --wasi
```

## Build with Docker

```bash
docker run --rm -v "$PWD":/fetch -w /fetch k33g/wasm-builder:0.0.1 \
  tinygo build -scheduler=none --no-debug \
    -o wasimancer-plugin-fetch.wasm \
    -target wasi main.go
```

## Run with Docker

```bash
docker run --rm -v "$PWD":/fetch -w /fetch k33g/wasm-builder:0.0.1 \
  extism call wasimancer-plugin-fetch.wasm fetch \
  --input '{"url": "https://modelcontextprotocol.io/introduction"}' \
  --allow-host "*" \
  --log-level "info" \
  --wasi
```