FROM --platform=$BUILDPLATFORM golang:1.24.4-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

RUN <<EOF
go mod tidy 
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build
EOF

#FROM debian:bookworm-slim
FROM scratch
WORKDIR /app
COPY --from=builder /app/mcp-http .

EXPOSE 9090

ENTRYPOINT ["./mcp-http"]

# docker build --platform linux/arm64 -t mcp-http:demo .
# docker run --rm -p 9090:9090 mcp-http:demo

# http://host.docker.internal:9090/mcp
# curl http://0.0.0.0:5050/mcp
# curl http://0.0.0.0:9090/mcp