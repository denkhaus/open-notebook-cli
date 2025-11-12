# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o onb ./cmd/onb

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S opennotebook && \
    adduser -u 1001 -S opennotebook -G opennotebook

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/onb .

# Change ownership to non-root user
RUN chown opennotebook:opennotebook /app/onb

# Switch to non-root user
USER opennotebook

# Expose configuration directory as volume
VOLUME ["/app/config"]

# Set environment variables
ENV OPEN_NOTEBOOK_API_URL=http://localhost:5055
ENV OPEN_NOTEBOOK_OUTPUT=table
ENV OPEN_NOTEBOOK_TIMEOUT=300

# Add binary to PATH
ENV PATH="/app:${PATH}"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD onb auth check || exit 1

# Default command
ENTRYPOINT ["onb"]
CMD ["--help"]