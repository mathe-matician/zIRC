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

console.log(`/////////////////////////////////////////////////////////////////////////////////////`);
console.log(`//// TEST CLIENT`);
console.log(`/////////////////////////////////////////////////////////////////////////////////////`);
console.log(`//// Available Commands:`);
console.log(`//// \tuse <user>:`);
console.log(`//// \t\tSelect the user to use. E.g. 0, 1, 2, 3, admin (whatever the file is named)`);
console.log(`//// \t\tIf user doesn't exist, create uid and save it.`);
console.log(`//// \trole <role>:`);
console.log(`//// \t\tclient: Sends messages without a source prefix`);
console.log(`//// \t\tserver: Sends messages with source prefix`);
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
        // console.log(`CLIENT UID: '${clientUID}', clientUIDExists: ${clientUIDExists}`);
        if (tokenExists) {
          client.write(`tokenPkg::${tokenpkg} :${nickname}@${host} ${line} \r\n`);
        } else if (clientUIDExists) {
          client.write(`clientUID::${clientUID} :${nickname}@${host} ${line} \r\n`)
        } else {
          prefix = ""
          if (line.includes("use")) {
            const res = line.split(" ");
            if (res === "server") {
              prefix += `:${nickname}@${host} `
            }
          }
          client.write(`${prefix}${line} \r\n`)
        }
      }
    }
  }); 
  reader.on("close", () => {
    client.end()
  });

const GetUserUID = (uid) => {
  try {
    if (fs.existsSync(uid)) {
      const allFileContents = fs.readFileSync(uid, 'utf-8');
      allFileContents.split(/\r?\n/).forEach(line =>  {
        if (line.length !== 0) {
          clientUID = line;
        }
      });
    } else {
      // console.log(`UID == ${uid}`);
      clientUIDExists = true;
      // generate UID
      // const uid = v4();
      // write file with uid for later use
      fs.writeFileSync(uid, uid);
      // set global uid here.
      clientUID = uid;
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
      } else if (splitData[0] === "clientUID") {
        console.log('Recieve client uid from server')
        GetUserUID(splitData[1].replace(/\s+/g, ' ').trim());
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