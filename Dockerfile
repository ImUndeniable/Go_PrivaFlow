# Stage 1: Build the code
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build BOTH binaries
RUN go build -o api_binary cmd/api/main.go
RUN go build -o worker_binary cmd/worker/main.go

# Stage 2: Run the code (Tiny image)
FROM alpine:latest
WORKDIR /root/

# Copy the binaries from the builder stage
COPY --from=builder /app/api_binary .
COPY --from=builder /app/worker_binary .

# We don't specify CMD here, we will do it in docker-compose!