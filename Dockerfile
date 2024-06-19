# Stage 1: Build the Go application
FROM golang:1.22-alpine as builder

# Install dependencies
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -o blog-api .

FROM debian:stable-slim as final
LABEL authors="tanmay-e-patil"

# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/blog-api .
COPY --from=builder /app/.env.docker /app/.env

# Expose port 8080 to the outside world
EXPOSE 8080

CMD ["./blog-api"]

