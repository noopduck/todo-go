FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Sørg for at binæren bygges for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Start ny, minimal base
FROM debian:bookworm-slim

# Installer bare det mest nødvendige (f.eks. for TLS/certs)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Kopier over fra build-stage
COPY --from=builder /app/main .
COPY --from=builder /app/entrypoint.sh .
COPY --from=builder /app/todo.db .

# Sørg for at begge er kjørbare
RUN chmod +x ./entrypoint.sh ./main

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
