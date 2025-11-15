# --- Build Stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
# This copies 'api/', 'processor/', 'pkg/', 'db/', and 'ui/'
COPY . .

# Build the api service
# The 'go build' command will find the 'ui' dir and embed it
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/api-server ./api

# --- Final Stage ---
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/api-server .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./api-server"]