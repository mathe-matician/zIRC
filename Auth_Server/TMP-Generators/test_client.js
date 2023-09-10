// const { ClientCommands } = require('./commands');
const net = require('node:net');
const fs = require('node:fs');
const readline = require("readline");
const chilkatManager = require('../Chilkat/chilkat_manager');
require('dotenv').config();

chilkatManager.unlock_bundle();

const host = process.env.AUTH_SERVER_HOST;
const port = process.env.AUTH_SERVER_PORT;

let client;

console.log(`/////////////////////////////////////////////////////////////////////////////////////`);
console.log(`//// AUTH SERVER TEST CLIENT`);
console.log(`/////////////////////////////////////////////////////////////////////////////////////`);
console.log(`//// Available Commands:`);
console.log(`//// \tconnect: Connects to Auth Server`);
console.log(`/////////////////////////////////////////////////////////////////////////////////////\n\n`);
const reader = readline.createInterface({ input: process.stdin });
  reader.on("line", (line) => {
    if (line === "connect") {
      Connect();
    } else {
      if (typeof(line) === "object") {
        client.write(JSON.stringify(line));
      } else {
        const data = `{"token": "test", "data": "${line}"}`
        client.write(chilkatManager.encrypt_decrypt_AES(data, true));
      }
    }
  }); 
  reader.on("close", () => {
    client.end();
  });

const Connect = () => {
  client = net.createConnection({ port: port, host: host }, () => {
    client.setEncoding('utf8');
  });
  client.on('connect', () => {
      console.log("Connected to Auth Server");
  });
  client.on('data', (data) => {
    console.log(data.toString());
  });
}