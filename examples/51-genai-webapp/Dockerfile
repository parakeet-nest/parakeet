FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY main.go .
COPY go.mod .

RUN <<EOF
go mod tidy 
go build
EOF


FROM scratch
WORKDIR /app
COPY public ./public
#COPY parakeet-webapp .
COPY --from=builder /app/parakeet-webapp .
CMD ["./parakeet-webapp"]
