# # Используем официальный образ Go
# FROM golang:1.24-alpine

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

# COPY . .

# # RUN go build -o excel-service .
# RUN go build -o excel-service ./cmd/main.go

# CMD ["./excel-service"]


# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Make sure it's executable
RUN chmod +x ./main

CMD ["./main"]