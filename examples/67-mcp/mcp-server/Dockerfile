FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .

RUN <<EOF
go mod tidy 
go build
EOF

FROM curlimages/curl:8.6.0
WORKDIR /app
COPY --from=builder /app/mcp-curl .
ENTRYPOINT ["./mcp-curl"]



