# Step 1: Build stage
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
# If you have a go.sum, uncomment the line below
# COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# -o main: defines the output binary name
RUN go build -o main .

# Step 2: Final stage (Runtime)
FROM alpine:latest

# Install ca-certificates (needed for HTTPS requests to CoinGecko/FRED)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]