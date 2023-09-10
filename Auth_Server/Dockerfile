# FROM node:v16.17.0
FROM node:16.20.2

COPY package.json package-lock.json server.js db.js .env /auth_server/

COPY modules/ /auth_server/modules/
COPY Crypto/ /auth_server/Crypto/
COPY Chilkat/ /auth_server/Chilkat/

WORKDIR /auth_server/
RUN npm i --loglevel verbose

ENTRYPOINT [ "node", "server" ]