# --- Build Stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the processor service
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/processor-server ./processor

# --- Final Stage ---
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/processor-server .

# Command to run the executable
CMD ["./processor-server"]