# Build stage
FROM golang:1.25.4 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -tags "fts5" -o server src/cmd/webcrawler/main.go

# Final stage
FROM debian:bookworm-slim
WORKDIR /app

# Copy server binary and data folder (schema.sql)
COPY --from=builder /app/server .
COPY data/ ./data/

# Install sqlite3 CLI
RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

# Delete old DB, recreate from schema, then run server
CMD rm -f data/mydb.sqlite && sqlite3 data/mydb.sqlite < data/schema.sql && ./server
