#!/bin/bash
extism call wasimancer_plugin_addition.wasm add \
  --input '{"left":30, "right":12}' \
  --log-level "info" \
  --wasi
echo ""
