FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o database-manager cmd/server/main.go

# Create a minimal production image
FROM alpine:3.18

WORKDIR /app

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/database-manager .

# Expose the application port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV DEFAULT_NAMESPACE="ns-uw9b6wey"

# Run the application
CMD ["./database-manager"]