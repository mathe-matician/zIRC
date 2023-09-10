# FROM node:16.17.0
FROM node:16.20.2

COPY package.json package-lock.json \
    numerics.js operators.js server.js \
    db.js constants.js commands.js client.js \
    chat_schema.js channel.js capabilities.js \
    auth.js auth_server_comm.js .env /chat_server/
    
COPY Modes/ /chat_server/Modes/
COPY db/ /chat_server/db/
COPY Crypto/ /chat_server/Crypto/
COPY Commands/ /chat_server/Commands/
COPY Chilkat/ /chat_server/Chilkat/

WORKDIR /chat_server/
RUN npm i --loglevel verbose

EXPOSE 6667

ENTRYPOINT [ "node", "server" ]