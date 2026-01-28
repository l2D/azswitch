# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o azswitch ./cmd/azswitch

# Runtime stage with Azure CLI
FROM mcr.microsoft.com/azure-cli:2.82.0

LABEL org.opencontainers.image.title="azswitch"
LABEL org.opencontainers.image.description="TUI for switching Azure tenants and subscriptions"
LABEL org.opencontainers.image.source="https://github.com/l2D/azswitch"
LABEL org.opencontainers.image.licenses="MIT"

# Create non-root user
RUN tdnf install -y shadow-utils && \
    groupadd -g 1000 azswitch && \
    useradd -u 1000 -g azswitch -d /home/azswitch -s /bin/sh -m azswitch && \
    tdnf clean all

# Copy binary from builder
COPY --from=builder /app/azswitch /usr/local/bin/azswitch

# Set up Azure CLI cache directory with proper ownership
RUN mkdir -p /home/azswitch/.azure && \
    chown -R azswitch:azswitch /home/azswitch/.azure

USER azswitch
WORKDIR /home/azswitch

ENTRYPOINT ["/usr/local/bin/azswitch"]
