FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY main.go .
COPY go.mod .

RUN <<EOF
go mod edit -dropreplace github.com/parakeet-nest/parakeet
go mod download
go mod tidy 
go build -o web-chat-bot
EOF

FROM scratch
WORKDIR /app
COPY --from=builder /app/web-chat-bot .

CMD ["./web-chat-bot"]