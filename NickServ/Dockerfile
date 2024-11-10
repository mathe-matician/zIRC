FROM golang:1.22.5-bookworm AS builder

WORKDIR /nick_serv

# caches dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY src/ src/

RUN go build nickserv

CMD ["./nickserv"]

FROM debian:bookworm-slim AS prod

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /nick_serv/nickserv /nick_serv/nickserv

CMD ["/nick_serv/nickserv"]