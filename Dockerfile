# Start from the official Golang image as a build stage
FROM golang:1.24.1 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./

# this is to download dependencies form golang
RUN go env -w GOPROXY=http://proxy.golang.org,direct

RUN apt-get update && apt-get install -y ca-certificates

RUN go mod download

# Copy the application source code
COPY . .

RUN chmod +x ./scripts/start.sh

# Build the Go application
RUN go build -o web-analyzer .

# Used a alpine minimal image for deployment. But there was issue while connecting to https due to certificates issues
FROM ubuntu:latest

# Install necessary certificates for HTTPS requests
# which are required to make secure HTTPS requests from within the container.
RUN  apt update &&  apt install -y ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built application from the builder stage
COPY --from=builder /app/web-analyzer .

# Expose the application port
EXPOSE 8080 6060 9090

RUN chmod +x web-analyzer

# Run the web application
ENTRYPOINT ["/root/web-analyzer"]