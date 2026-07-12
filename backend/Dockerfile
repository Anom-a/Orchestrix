# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go module files first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /orchestrix-backend ./cmd/api/main.go

# ---- Runtime Stage ----
FROM alpine:3.21

WORKDIR /app

# Install CA certificates for HTTPS calls
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /orchestrix-backend .

# Create storage directory for uploads
RUN mkdir -p /app/storage/uploads

EXPOSE 8080

CMD ["./orchestrix-backend"]
