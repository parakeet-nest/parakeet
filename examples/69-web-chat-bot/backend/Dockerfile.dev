FROM golang:1.23.1-alpine AS builder

WORKDIR /app

# Copy the entire parakeet project
COPY . .

# Move to the backend directory
WORKDIR /app/examples/69-web-chat-bot/backend

# Update go.mod to point to the root of the copied parakeet code
RUN go mod edit -replace github.com/parakeet-nest/parakeet=/app

# Build the application
RUN go mod download
RUN go build -o web-chat-bot

FROM scratch
WORKDIR /app
COPY --from=builder /app/examples/69-web-chat-bot/backend/web-chat-bot .

CMD ["./web-chat-bot"]