# Dockerfile for Backend
FROM golang:1.20 as builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -o main .

# Final stage
FROM alpine:3.18

WORKDIR /root/

# Copy compiled binary from builder
COPY --from=builder /app/main .

# Expose the app port
EXPOSE 8084

# Command to run the executable
CMD ["./main"]
