# Build stage
FROM golang:1.23.1 AS builder

WORKDIR /app

# Copy only the go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/rent-watcher

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the config file
COPY --from=builder /app/config.json .

# Use a non-root user for better security
RUN adduser -D appuser
USER appuser

# Run the binary
CMD ["./main"]