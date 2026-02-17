# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mms main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates and ssh client if needed
RUN apk --no-cache add ca-certificates tzdata

# Copy binary
COPY --from=builder /app/mms .

# Copy templates and static assets
COPY --from=builder /app/app/templates ./app/templates
# Assuming static assets are in app/static
COPY --from=builder /app/app/static ./app/static 

# Set environment variables
ENV PORT=3000

# Expose port
EXPOSE 3000

# Run the application
CMD ["./mms"]
