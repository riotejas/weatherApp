# Build stage
FROM golang:1.22-alpine AS builder
LABEL authors="russell"

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app ./cmd/api/...

# Run stage
FROM alpine:latest

WORKDIR /weatherApp/

# Copy the binary from the builder stage
COPY --from=builder /app /weatherApp

EXPOSE 8080

CMD ["/weatherApp/app"]
