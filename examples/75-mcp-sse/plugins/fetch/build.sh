#!/bin/bash
tinygo build -scheduler=none --no-debug \
  -o wasimancer-plugin-fetch.wasm \
  -target wasi main.go

ls -lh *.wasm
