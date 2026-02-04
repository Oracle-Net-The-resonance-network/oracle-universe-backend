# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/oracle-universe .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/oracle-universe .

# Expose port
EXPOSE 8090

# Run command - creates admin if env vars set
CMD ["./oracle-universe", "serve", "--http=0.0.0.0:8090"]
