FROM golang:1.22-alpine as builder

RUN apk add --no-cache git

# Install PostgreSQL client
RUN apk add --no-cache postgresql-client


WORKDIR /app

# Download goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy migrations
COPY sql ./sql

# Copy .env file
COPY .env.docker .env

# Copy entrypoint script
COPY wait-for-postgres.sh .

CMD ["./wait-for-postgres.sh"]
