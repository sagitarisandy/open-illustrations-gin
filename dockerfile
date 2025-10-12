# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy module files dan download dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary
RUN go build -o main ./main.go

# Stage 2: Run (secure)
FROM alpine:latest
RUN adduser -D appuser
USER appuser
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]