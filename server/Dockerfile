FROM golang:1.25.4 AS builder
WORKDIR /app

# dl dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# compile server to binaries
RUN go build -tags "fts5" -o server src/cmd/webcrawler/main.go


FROM debian:bookworm-slim
WORKDIR /app

# sqlite dependency
RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*

# copy binary
COPY --from=builder /app/server .

# copy data directory 
COPY --from=builder /app/data ./data

# copy entry script
COPY entry.sh .

# make entrypoint executable
RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
