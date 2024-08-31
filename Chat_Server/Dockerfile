FROM golang:1.22.5-bookworm as builder

WORKDIR /chat_server

# caches dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY chat/ chat/
COPY helpers/ helpers/
COPY server/ server/
COPY remote_conn/ remote_conn/
COPY client/ client/

RUN go build zirc

FROM debian:bookworm-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /chat_server/zirc /chat_server/zirc

CMD ["/chat_server/zirc"]