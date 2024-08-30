// const { ClientCommands } = require('./commands');
const net = require('node:net');
const fs = require('node:fs');
const readline = require("readline");

let nickname;
const port = "6667";
const host = "127.0.0.1";
let clientUID = "";
let clientUIDExists = false;
let client;
let tokenExists = false;
let tokenpkg = "";
let debug = false;

console.log(`/////////////////////////////////////////////////////////////////////////////////////`);
console.log(`//// TEST CLIENT`);
console.log(`/////////////////////////////////////////////////////////////////////////////////////`);
console.log(`//// Available Commands:`);
console.log(`//// \trole <role>:`);
console.log(`//// \t\tclient (default): Sends messages without a source prefix`);
console.log(`//// \t\tserver: Sends messages with source prefix`);
console.log(`//// \tdebug: toggles debug flag. echos the message the client is about to send to the server`);
console.log(`//// \tconnect: Connects to IRC server`);
console.log(`/////////////////////////////////////////////////////////////////////////////////////\n\n`);
const reader = readline.createInterface({ input: process.stdin });
  reader.on("line", (line) => {
    if (line === "connect") {
      Connect();
    } else {
      if (line.includes("use")) {
        const res = line.split(" ");
        GetUserUID(res[1]);
      }
      // GetUserUID();
  
      if (client) {
        prefix = ""
          if (line.includes("role")) {
            const res = line.split(" ");
            if (res === "server") {
              prefix += `:${nickname}@${host} `
            }
          }
          if (line.includes("debug")) {
            debug = !debug
          }
          if (debug) {
            console.log(`Client about to write: ${prefix}${line} \\r\\n`)
            client.write(`${prefix}${line} \r\n`)
          } else {
            client.write(`${prefix}${line} \r\n`)
          }
      }
    }
  }); 
  reader.on("close", () => {
    client.end()
  });

const Connect = () => {
  client = net.createConnection({ port: port, host: host }, () => {
    client.setEncoding('utf8');
  });
  client.on('connect', () => {
      console.log("Connected to IRC server");
  });
  client.on('data', (data) => {
    const _data = data.toString();
    console.log(`${_data}`);
  });
}