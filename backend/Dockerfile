# Build stage
FROM golang:1.21 AS builder

WORKDIR /app

# Copy and download dependency definitions
COPY go.mod go.sum ./
RUN go mod download

# Copy the application code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend main.go

# Runtime stage
FROM alpine:latest

# Install CA certificates for HTTPS support
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/backend .

# Expose the application port
EXPOSE 8081

# Command to run the application
CMD ["./backend"]