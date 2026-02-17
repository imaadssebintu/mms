FROM alpine:latest

WORKDIR /app

# Install ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Copy the pre-built binary (built in CI)
COPY mms .

# Copy templates and static assets
COPY app/templates ./app/templates

# Set environment variables
ENV PORT=8000

# Expose port
EXPOSE 8000

# Run the application
CMD ["./mms"]
