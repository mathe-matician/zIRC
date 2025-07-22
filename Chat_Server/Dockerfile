FROM golang:1.24.2 AS builder

ARG DEBUG_PORT=4000

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
COPY s2s_config.yaml s2s_config.yaml
COPY config.yaml config.yaml

# TODO
# add args to specify build tags
# any additional flags for prod stage?
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o zirc .

FROM golang:1.24.2 AS dev

ARG DEBUG_PORT=4000

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY --from=builder /chat_server/zirc /chat_server/zirc
COPY --from=builder /chat_server/s2s_config.yaml /chat_server/s2s_config.yaml
COPY --from=builder /chat_server/config.yaml /chat_server/config.yaml

EXPOSE 4000

ENV DLV_DEBUG_PORT=${DEBUG_PORT}
# TODO
# dlv starts and must be manually resumed instead of the process just starting and then the user
# specifying debug commands
# helpful if you want to immediately start the debugger, but shouldn't be the default behavior
# ENTRYPOINT [ "/go/bin/dlv", "--listen=:4000", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/chat_server/zirc" ]
# CMD [ "/go/bin/dlv", , "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/chat_server/zirc" ]
CMD ["/chat_server/zirc"]

# FROM scratch AS prod

# COPY --from=builder /chat_server/zirc /chat_server/zirc
# COPY --from=builder /chat_server/s2s_config.yaml /chat_server/s2s_config.yaml
# COPY --from=builder /chat_server/config.yaml /chat_server/config.yaml

# CMD ["/chat_server/zirc"]