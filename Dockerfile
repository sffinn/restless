# Multi-stage build for a small production image.
# DigitalOcean App Platform will use this Dockerfile when present.

# ---- Build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install ca-certificates (needed for HTTPS) and git if modules need it.
RUN apk add --no-cache ca-certificates git

# Copy go.mod first for better layer caching.
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# Build a statically linked binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# ---- Runtime stage ----
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage.
COPY --from=builder /app/server .

# App Platform sets PORT; we listen on 0.0.0.0:$PORT inside the container.
EXPOSE 8080

CMD ["./server"]

