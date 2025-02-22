FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY main.go .
COPY go.mod .

RUN <<EOF
go mod edit -dropreplace github.com/parakeet-nest/parakeet
go mod download
go mod tidy 
go build -o create-embeddings
EOF

FROM scratch
WORKDIR /app
COPY --from=builder /app/create-embeddings .

CMD ["./create-embeddings"]