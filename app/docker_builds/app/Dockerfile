# Use a multi-stage build for efficiency
# (https://docs.docker.com/build/building/multi-stage/)

###############################################################################
# Stage 1: Build Environment (Small & Temporary)
###############################################################################

FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Download and verify dependencies
RUN set -Eeux && \
    go mod download && \
    go mod verify

# Copy the rest of the project files
COPY . .

# Build the application with version injection
RUN GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s -X 'main.Version=$(grep PROJECT_VERSION .env | cut -d= -f2)'" \
    -o bin/pagination-app ./main.go

###############################################################################
# Stage 2: Final Image (Slim & Optimized)
###############################################################################

FROM alpine:latest AS final

# Install libc6 compatibility layer and curl
RUN apk add --no-cache libc6-compat curl

WORKDIR /app

# Copy the binary from stage 1
COPY --from=builder /app/bin/pagination-app /usr/local/bin/pagination-app

# Copy the .env and schema file
COPY .env /app/.env
COPY pkg/schema.sql /app/pkg/schema.sql

# Set executable permissions
RUN chmod +x /usr/local/bin/pagination-app

# Expose server port (placeholder, set via .env)
EXPOSE ${SERVER_PORT}

# Run the application
ENTRYPOINT ["/usr/local/bin/pagination-app"]
