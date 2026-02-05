# Build stage - use full Go image for toolchain support
FROM golang:1.22 AS builder

WORKDIR /app

# Allow Go to download newer toolchain if needed
ENV GOTOOLCHAIN=auto

# Copy all source
COPY . .

# Download deps and build in single RUN so GOTOOLCHAIN downloads Go 1.24 for both
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/oracle-universe .

# Runtime stage - minimal
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/oracle-universe .

# Copy entrypoint script
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# Expose port
EXPOSE 8090

# Run entrypoint (creates superuser from env vars, then starts server)
CMD ["./entrypoint.sh"]
