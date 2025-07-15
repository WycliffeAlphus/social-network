# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# CGO_ENABLED=0 is important for creating a static binary
# -ldflags="-s -w" reduces the binary size
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Create a minimal production image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS communication if needed (e.g., if Go app makes external calls)
# For pure Go SQLite drivers (like go-sqlite3), this might not be strictly necessary,
# but it's good practice for general compatibility.
RUN apk --no-cache add ca-certificates sqlite-libs

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port Go application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]