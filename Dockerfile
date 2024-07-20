FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /morphic

# Use a smaller base image for the final image
FROM alpine:latest

WORKDIR /app

# Install any required runtime dependencies
RUN apk add --no-cache ca-certificates

# Expose the application port
EXPOSE 8060

# Copy the binary from the builder stage
COPY --from=builder /morphic /morphic

COPY config.yaml /app/config.yaml

CMD ["sh", "-c", "sleep 5 && /morphic server start -c /app/config.yaml"]
