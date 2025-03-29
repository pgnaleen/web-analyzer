# Start from the official Golang image as a build stage
FROM golang:1.24.1 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./

# this is to download dependencies form golang
RUN go env -w GOPROXY=http://proxy.golang.org,direct

RUN go mod download

# Copy the application source code
COPY . .

RUN chmod +x ./scripts/start.sh

# Build the Go application
RUN go build -o web-analyzer .

# Use a minimal image for deployment
FROM alpine:latest

# Install necessary certificates for HTTPS requests
# is used in Alpine-based Docker images to install CA (Certificate Authority) certificates,
# which are required to make secure HTTPS requests from within the container.
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built application from the builder stage
COPY --from=builder /app/web-analyzer .

# Expose the application port
EXPOSE 8080

# Run the web application
#ENTRYPOINT ["nohup /app/web-analyzer &"]
ENTRYPOINT ["nohup /root/web-analyzer &"]