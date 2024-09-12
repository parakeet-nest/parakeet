<!-- TOPIC: WebAssembly plugins for Parakeet SUMMARY: The release of Parakeet's version 0.0.6 brings support for WebAssembly, allowing users to write their own wasm plugins in various languages (Rust, Go, C, etc.) and use them with the "Function Calling" feature. KEYWORDS: Parakeet, WebAssembly, Wasm plugins, Extism project, TinyGo, Function Calling -->
# Wasm plugins

The release `0.0.6` of Parakeet brings the support of **WebAssembly** thanks to the **[Extism project](https://extism.org/)**. That means you can write your own wasm plugins for Parakeet to add new features (for example, a chunking helper for doing RAG) with various languages (Rust, Go, C, ...).

Or you can use the Wasm plugins with the "Function Calling" feature, which is implemented in Parakeet.

!!! note
	You can find an example of "Wasm Function Calling" in [examples/18-call-functions-for-real](https://github.com/parakeet-nest/parakeet/tree/main/examples/18-call-functions-for-real) - the wasm plugin is located in the `wasm` folder and it is built with **[TinyGo](https://tinygo.org/)**.

🚧 more samples to come.
