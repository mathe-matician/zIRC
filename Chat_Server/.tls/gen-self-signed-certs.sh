#!/bin/bash

if [[ ! -e "server.key" ]]; then
    echo "Generating server.key..."
    openssl genrsa 2048 > server.key
    chmod 400 server.key
fi

if [[ ! -e "server.crt" ]]; then
    echo "Generating server.crt..."
    openssl req -new -x509 -nodes -sha256 -days 365 -key server.key -out server.crt
fi

if [[ ! -e "s2s.key" ]]; then
    echo "Generating s2s.key..."
    openssl genrsa 2048 > s2s.key
    chmod 400 s2s.key
fi

if [[ ! -e "s2s.crt" ]]; then
    echo "Generating s2s.crt..."
    openssl req -new -x509 -nodes -sha256 -days 365 -key s2s.key -out s2s.crt
fi

if [[ ! -e "db.key" ]]; then
    echo "Generating db.key..."
    openssl genrsa 2048 > db.key
    chmod 400 db.key
fi

if [[ ! -e "db.crt" ]]; then
    echo "Generating db.key..."
    openssl req -new -x509 -nodes -sha256 -days 365 -key db.key -out db.crt
fi