// const { ClientCommands } = require('./commands');
const net = require('node:net');
const fs = require('node:fs');
const readline = require("readline");
const { v4 } = require('uuid');

let nickname;
const port = "6667";
const host = "127.0.0.1";
let clientUUID; 
let client;
let tokenExists = false;
let tokenpkg = "";

console.log(`/////////////////////////////////////////////////////////////////////////////////////`);
console.log(`//// TEST CLIENT`);
console.log(`/////////////////////////////////////////////////////////////////////////////////////`);
console.log(`//// Available Commands:`);
console.log(`//// \tuse <user>:`);
console.log(`//// \t\tSelect the user to use. E.g. 0, 1, 2, 3, admin (whatever the file is named)`);
console.log(`//// \t\tIf user doesn't exist, create uuid and save it.`);
console.log(`//// \tconnect: Connects to IRC server`);
console.log(`/////////////////////////////////////////////////////////////////////////////////////\n\n`);
const reader = readline.createInterface({ input: process.stdin });
  reader.on("line", (line) => {
    if (line === "connect") {
      Connect();
    } else {
      if (line.includes("use")) {
        const res = line.split(" ");
        GetUserUUID(res[1]);
      }
  
      if (client) {
        // console.log(`C: :${nickname}@${host} ${line} \r\n`);
        if (tokenExists) {
          client.write(`tokenPkg::${tokenpkg} :${nickname}@${host} ${line} \r\n`);
        } else {
          client.write(`:${nickname}@${host} ${line} \r\n`)
        }
      }
    }
  }); 
  reader.on("close", () => {
    client.end()
  });

const GetUserUUID = (user) => {
  try {
    if (fs.existsSync(user)) {
      const allFileContents = fs.readFileSync(user, 'utf-8');
      allFileContents.split(/\r?\n/).forEach(line =>  {
        if (line.length !== 0) {
          clientUUID = line;
        }
      });
    } else {
      // generate UUID
      const uuid = v4();
      // write file with uuid for later use
      fs.writeFileSync(user, uuid);
      // set global uuid here.
      clientUUID = uuid;
    }
  } catch(err) {
    console.log(err);
  }
}

const Connect = () => {
  client = net.createConnection({ port: port, host: host }, () => {
    client.setEncoding('utf8');
  });
  client.on('connect', () => {
      console.log("Connected to IRC server");
  });
  client.on('data', (data) => {
    if (data.includes("NICK")) {
      const res = data.split(' ');
      nickname = res[res.length];
    }
    // console.log(`\nDATA:\n${data.toString()}\n`)
    const _data = data.toString();
    if (_data.includes("::")) {
      const regex = /\r\n/g;
      const arr = [..._data.matchAll(regex)];
      let crlfsplit = [];
      let splitData = [];
      if (arr.length > 1) {
        // multiple lines were sent in this tcp stream
        crlfsplit = _data.split("\r\n");
        splitData = crlfsplit[0].split("::");
        crlfsplit.shift(); // remove token package
      } else {
        splitData = _data.split("::");
      }
      // console.log(`\nSplit data:\n${splitData}\n`);
      if (splitData[0] === "tokenPkg") {
        // console.log(`\nsplitData[1]:\n${splitData[1]}\n`);
        const tokenPkg = JSON.parse(splitData[1]);
        console.log(`Token: ${tokenPkg["token"]}`);
        console.log(`Token Expr: ${tokenPkg["expr"]}`);
        tokenpkg = tokenPkg["token"] + "::" + tokenPkg["expr"];
        // token = tokenPkg[0];
        // token_expr = tokenPkg[1];
        tokenExists = true;
      }

      if (crlfsplit.length > 0) {
        for (let i = 0; i < crlfsplit.length; i++) {
          console.log(`${crlfsplit[i]}`);
        }
      }
    } else {
      console.log(`${_data}`);
    }
    // client.end();
  });
}