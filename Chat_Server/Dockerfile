FROM golang:1.24.5 AS builder

WORKDIR /chat_server

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY chat/ chat/
COPY helpers/ helpers/
COPY remote_conn/ remote_conn/
COPY s2s_config.yaml s2s_config.yaml
COPY config.yaml config.yaml

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go build -gcflags="all=-N -l" -o zirc .

FROM golang:1.24.5 AS dev

WORKDIR /chat_server

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY --from=builder /chat_server /chat_server

# ENTRYPOINT [ "/go/bin/dlv", "--continue", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "./zirc" ]
CMD [ "/chat_server/zirc" ]

FROM scratch AS prod

COPY --from=builder /chat_server/zirc /chat_server/zirc
COPY --from=builder /chat_server/s2s_config.yaml /chat_server/s2s_config.yaml
COPY --from=builder /chat_server/config.yaml /chat_server/config.yaml

CMD ["/chat_server/zirc"]