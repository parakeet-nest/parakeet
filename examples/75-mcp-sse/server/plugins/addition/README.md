# Extism Rust PDK Plugin

## Create Extism plugin

```bash
mkdir addition && cd addition
extism generate plugin 
```
> see: https://github.com/extism/cli?tab=readme-ov-file#generate-a-plugin


**Build**:
```bash
cd addition
cargo build --release 
cp target/wasm32-unknown-unknown/release/wasimancer_plugin_addition.wasm ./
```

**Run**:
```
extism call wasimancer_plugin_addition.wasm add \
  --input '{"left":30, "right":12}' \
  --log-level "info" \
  --wasi
```

## Build with Docker

```bash
docker run --rm -v "$PWD":/addition -w /addition k33g/wasm-builder:0.0.1 ./build.sh
```

Or:
```bash
docker run --rm -v "$PWD":/addition -w /addition k33g/wasm-builder:0.0.1 \
  bash -c "
    cargo clean && \
    cargo install cargo-cache && \
    cargo cache -a && \
    cargo build --release && \
    cp target/wasm32-unknown-unknown/release/wasimancer_plugin_addition.wasm ./
  "
```