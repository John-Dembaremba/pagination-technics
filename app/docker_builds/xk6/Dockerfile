# Stage 1: Build k6 with xk6-output-influxdb
FROM golang:1.23-alpine3.21 AS builder

WORKDIR /go/src

# Install git for cloning
RUN apk --no-cache add git

# Install xk6
RUN go install go.k6.io/xk6/cmd/xk6@latest

# Clone xk6-output-influxdb
RUN git clone https://github.com/grafana/xk6-output-influxdb.git

WORKDIR /go/src/xk6-output-influxdb

# Build k6 with the InfluxDB extension
RUN xk6 build --with github.com/grafana/xk6-output-influxdb=. --output /tmp/k6

# Stage 2: Final Image
FROM alpine:3.21

# Install ca-certificates and create k6 user
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 12345 -g 12345 k6

# Copy k6 binary
COPY --from=builder /tmp/k6 /usr/bin/k6

USER 12345
WORKDIR /home/k6
ENTRYPOINT ["k6"]
