FROM golang:1.22.5-bookworm AS builder

WORKDIR /chat_server

# caches dependencies
COPY go.mod go.sum ./
# RUN --mount=type=cache,target=/root/.cache,id=chat_server_cache \
#     --mount=type=cache,target=/go/pkg/mod,id=go_pkg_cache \
#     go mod download
RUN go mod download

COPY main.go ./
COPY chat/ chat/
COPY helpers/ helpers/
COPY remote_conn/ remote_conn/

RUN go build zirc

FROM debian:bookworm-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /chat_server/zirc /chat_server/zirc

CMD ["/chat_server/zirc"]