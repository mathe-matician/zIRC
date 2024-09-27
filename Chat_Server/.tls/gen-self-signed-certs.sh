#!/bin/bash

openssl genrsa 2048 > server.key
chmod 400 server.key
openssl req -new -x509 -nodes -sha256 -days 365 -key server.key -out server.crt
openssl genrsa 2048 > s2s.key
chmod 400 s2s.key
openssl req -new -x509 -nodes -sha256 -days 365 -key s2s.key -out s2s.crt