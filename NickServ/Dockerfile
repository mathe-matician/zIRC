FROM golang:1.22.5-bookworm AS builder

WORKDIR /nickserv

# caches dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./

RUN go build nickserv

FROM debian:bookworm-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /nickserv /nickserv
WORKDIR /nickserv

CMD ["tail", "-f", "/dev/null"]