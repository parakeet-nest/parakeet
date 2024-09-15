# Parakeet GenAI Webapp demo

```bash
go run main.go
```

## Start the project with Docker Compose

This 🐳 Compose GenAI Stack is an example of the usage of the [Parakeet](https://github.com/parakeet-nest/parakeet) 🦜🪺 GoLang library, made to simplify the development of small generative AI applications.

If you want to use the stack with your local install of Ollama:

```bash
HTTP_PORT=9999 OLLAMA_BASE_URL=http://host.docker.internal:11434 docker compose --profile webapp up
```
> Make sure that Ollama is started

If you want to use the stack with Ollama running in a container:
```bash
HTTP_PORT=7777 OLLAMA_BASE_URL=http://ollama:11434 docker compose --profile container up
```

> ✋ Compose will start the pull of the LLM, so be patient (my advice would be to use "small LLM" as tinyllama, tinydolphin or gemma:2b)

## Remarks

If you have a GPU on your workstation, you can turn on GPU access with Docker Compose: https://docs.docker.com/compose/gpu-support. With Mac M1, M2 and M3 use the local install of Ollama.

## Install dependencies (front)

**JavaScript**:
```bash
cd public/js
wget https://cdn.jsdelivr.net/npm/markdown-it@14.1.0/dist/markdown-it.min.js
wget https://cdn.jsdelivr.net/npm/beercss@3.7.8/dist/cdn/beer.min.js
wget https://cdn.jsdelivr.net/npm/material-dynamic-colors@1.1.2/dist/cdn/material-dynamic-colors.min.js
wget https://unpkg.com/htmx.org@2.0.2/dist/htmx.min.js

wget https://unpkg.com/htmx-ext-client-side-templates@2.0.0/client-side-templates.js
wget https://unpkg.com/mustache@latest -O mustache.js
```

**Css**:
```bash
cd public/css
wget https://cdn.jsdelivr.net/npm/beercss@3.7.8/dist/cdn/beer.min.css
wget https://github.com/marella/material-symbols/raw/main/material-symbols/material-symbols-outlined.woff2
```

