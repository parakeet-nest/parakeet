#!/bin/bash
tinygo build -scheduler=none --no-debug \
  -o plugin.wasm \
  -target wasi main.go