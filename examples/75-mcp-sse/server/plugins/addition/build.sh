#!/bin/bash
#set -e

cargo clean
cargo install cargo-cache
cargo cache -a
cargo build --release
cp target/wasm32-unknown-unknown/release/wasimancer_plugin_addition.wasm ./
