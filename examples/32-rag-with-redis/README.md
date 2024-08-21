# How to run this example

## Start a Redis server
```bash
docker compose up
```

## Create the embeddings

```bash
cd create-embeddings
go run main.go
```

## Use the embeddings

```bash
cd use-embeddings
go run main.go
```

