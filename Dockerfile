FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

#  Build binary for linux os
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

# New minimal base
FROM debian:bookworm-slim

WORKDIR /app

# Kopier over fra build-stage
COPY --from=builder /app/main .
COPY --from=builder /app/entrypoint.sh .
COPY --from=builder /app/todo.db .

# Make sure the files are executable
RUN chmod +x ./entrypoint.sh ./main

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
