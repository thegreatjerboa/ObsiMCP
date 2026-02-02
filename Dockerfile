# Dockerfile for ObsiMCP
# Go-based MCP server for Obsidian vault CRUD operations

FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install git for go mod download
RUN apk add --no-cache git

# Copy source code
COPY . .

# Build the binary
RUN go build -o obsimcp main.go

# ---

FROM alpine:3.19

WORKDIR /app

# Copy the binary
COPY --from=builder /build/obsimcp /app/obsimcp

# Copy the source config directory structure (needed for runtime.Caller path resolution)
# The binary expects config at the compile-time path, so we recreate it
COPY --from=builder /build/src/config /build/src/config

# Copy entrypoint script and fix line endings (Windows CRLF -> Unix LF)
COPY entrypoint.sh /app/entrypoint.sh
RUN sed -i 's/\r$//' /app/entrypoint.sh && chmod +x /app/entrypoint.sh

# Create default directories
RUN mkdir -p /vault /backup /templates

# Environment variables (can be overridden at runtime)
ENV VAULT_PATH=/vault
ENV BACKUP_PATH=/backup
ENV TEMPLATE_PATH=/templates

ENTRYPOINT ["/app/entrypoint.sh"]
