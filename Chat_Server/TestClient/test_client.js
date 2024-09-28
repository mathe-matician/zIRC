// const { ClientCommands } = require('./commands');
const net = require('node:net');
const fs = require('node:fs');
const readline = require("readline");

let nickname;
let port = "6667";
let host = "127.0.0.1";
let clientUID = "";
let clientUIDExists = false;
let client;
let tokenExists = false;
let tokenpkg = "";
let debug = false;
let mode = "client";

console.log(`/////////////////////////////////////////////////////////////////////////////////////`);
console.log(`//// TEST CLIENT`);
console.log(`/////////////////////////////////////////////////////////////////////////////////////`);
console.log(`//// Available Commands:`);
console.log(`//// \trole <role>:`);
console.log(`//// \t\tclient (default): Sends messages without a source prefix`);
console.log(`//// \t\tserver: Sends messages with source prefix`);
console.log(`//// \ttest_mode <mode>:`);
console.log(`//// \t\tclient: sends messages as client (normal)`);
console.log(`//// \t\tserver: Sends messages as "servermanager"`);
console.log(`//// \tdebug: toggles debug flag. echos the message the client is about to send to the server`);
console.log(`//// \tconnect: Connects to IRC server`);
console.log(`/////////////////////////////////////////////////////////////////////////////////////\n\n`);
const reader = readline.createInterface({ input: process.stdin });
  reader.on("line", (line) => {
    if (line === "connect") {
      Connect();
    } else if (line.includes("test_mode")) {
      // console.log(`Line includes test_mode: ${line}`)
      const res = line.split(" ");
      if (res[1] === "server") {
        mode = "server"
        port = "7000"
      } else {
        mode = "client"
      }
    } else {
      // if (line.includes("use")) {
      //   const res = line.split(" ");
      //   GetUserUID(res[1]);
      // }
      // GetUserUID();
  
      if (client) {
        prefix = ""
        if (line.includes("role")) {
          const res = line.split(" ");
          if (res === "server") {
            prefix += `:${nickname}@${host} `
          }
        } else if (line.includes("debug")) {
          debug = !debug
        } else {
          if (debug) {
            console.log(`Client about to write: ${prefix}${line} \\r\\n`)
            if (mode === "client") {
              client.write(`${prefix}${line} \r\n`)
            } else {
              client.write(`${line}`)
            }
          } else {
            if (mode === "client") {
              client.write(`${prefix}${line} \r\n`)
            } else {
              client.write(`${line}`)
            }
          }
        }
      }
    }
  }); 
  reader.on("close", () => {
    client.end()
  });

const Connect = () => {
  if (mode === "server") {
    port = "7000"
  }
  client = net.createConnection({ port: port, host: host }, () => {
    client.setEncoding('utf8');
  });
  client.on('connect', () => {
      console.log(`Connected to IRC server: ${host}:${port}`);
  });
  client.on('data', (data) => {
    const _data = data.toString();
    console.log(`${_data}`);
  });
  client.on('end', (data) => {
    console.log(`Server terminated connection`);
  });
}