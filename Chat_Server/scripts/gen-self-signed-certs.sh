#!/bin/bash

TLS_DIR="./Chat_server/.tls"
mkdir -p $TLS_DIR

SERVER_KEY_PATH="$TLS_DIR/server.key"
if [[ ! -e "$SERVER_KEY_PATH" ]]; then
    echo "Generating $SERVER_KEY_PATH..."
    openssl genrsa 2048 > $SERVER_KEY_PATH
    chmod 400 $SERVER_KEY_PATH
fi

SERVER_CRT_PATH="$TLS_DIR/server.crt"
if [[ ! -e "$SERVER_CRT_PATH" ]]; then
    echo "Generating $SERVER_CRT_PATH..."
    openssl req -new -x509 -nodes -sha256 -days 365 -key $SERVER_KEY_PATH -out $SERVER_CRT_PATH
fi

S2S_KEY_PATH="$TLS_DIR/s2s.key"
if [[ ! -e "$S2S_KEY_PATH" ]]; then
    echo "Generating $S2S_KEY_PATH..."
    openssl genrsa 2048 > $S2S_KEY_PATH
    chmod 400 $S2S_KEY_PATH
fi

S2S_CRT_PATH="$TLS_DIR/s2s.crt"
if [[ ! -e "$S2S_CRT_PATH" ]]; then
    echo "Generating $S2S_CRT_PATH..."
    openssl req -new -x509 -nodes -sha256 -days 365 -key $S2S_KEY_PATH -out $S2S_CRT_PATH
fi

DB_KEY_PATH="$TLS_DIR/db.key"
if [[ ! -e "$DB_KEY_PATH" ]]; then
    echo "Generating $DB_KEY_PATH..."
    openssl genrsa 2048 > $DB_KEY_PATH
    chmod 400 $DB_KEY_PATH
fi

DB_CRT_PATH="$TLS_DIR/db.crt"
if [[ ! -e $DB_CRT_PATH ]]; then
    echo "Generating db.key..."
    openssl req -new -x509 -nodes -sha256 -days 365 -key $DB_KEY_PATH -out $DB_CRT_PATH
fi