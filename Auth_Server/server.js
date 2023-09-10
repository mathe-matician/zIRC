const net = require('node:net');
const tls = require('tls');
const fs = require('node:fs');
const util = require('util');
// const chilkatManager = require('./Chilkat/chilkat_manager');
const cryptoManager = require('./Crypto/crypto_manager');
const { DB } = require('./db');
require('dotenv').config();

const host = process.env.AUTH_SERVER_HOST;
const port = process.env.AUTH_SERVER_PORT;
const statAsync = util.promisify(fs.stat);

const AuthServer = (tls=false) => {
  // chilkatManager.unlock_bundle();
  const db = DB();

  const check_admin_token = async (token) => {
    const prp_stmt_AuthCheck = db.Prepare(
      "AuthCheck", "SELECT (salt, token, token_expiration) FROM admin_tokens"
      );
    const authCheckRes = await db.Exec(prp_stmt_AuthCheck);
    console.log(`authCheckRes: ${JSON.stringify(authCheckRes)}`);
    if (!authCheckRes) {
      console.log("Failed auth check!");
      throw new Error("Error");
    }

    const checkResult = authCheckRes["row"];
    const rowValues = checkResult.substring(
      checkResult.indexOf("(") + 1,
      checkResult.lastIndexOf(")")
    ).split(",");
    const _expiration = rowValues[2];
    console.log(rowValues);
    
    const now = new Date();
    const expiration = new Date(_expiration);
    if (now > expiration) {
      // TODO
      // should these even expire...?? probably...
      console.log(`Admin auth token is expired. now: ${now}, expiration: ${expiration}`);
      throw new Error("Error");
    }

    const salt = rowValues[0];
    const hash = rowValues[1];
    const tokenHash = cryptoManager.generate_hash(token, salt);
    // const tokenHash = chilkatManager.hash_string(token, salt);
    if (tokenHash !== hash) {
      console.log("Auth token does not match")
      throw new Error("Error");
    }
  };

  const Start = async () => {
    try {
      if (tls) {
        const options = {
          key: fs.readFileSync('server-key.pem'),
          cert: fs.readFileSync('server-cert.pem'),
        
          // This is necessary only if using client certificate authentication.
          requestCert: true,
        
          // This is necessary only if the client uses a self-signed certificate.
          ca: [ fs.readFileSync('client-cert.pem') ],
        };
        
        const server = tls.createServer(options, (socket) => {
          console.log('server connected: ', socket.authorized ? 'authorized' : 'unauthorized');
          socket.setEncoding('utf8');
          socket.on('data', (data) => {
          });
          server.on("error", (err) => {
            console.log(`Server error: ${err}`);
            throw err;
          });
          c.on("end", (client) => {
            console.log(`${client} disconnected`);
          });
        });
        // });
        server.listen(port, host, () => {
          console.log('server started:', server.address());
        });
      } else {
        const server = net.createServer(async (socket) => {
          socket.setEncoding('utf8');
          socket.on('data', async (data) => {
            console.log(`{"remoteAddress": "${socket.remoteAddress}", "data": "${data}"}`);
            // 1. Decrypt data
            const decryptedData = data;
            // const decryptedData = chilkatManager.encrypt_decrypt_AES(data, false);
            // if (!decryptedData) {
            if (!decryptedData) {
              // console.log("Could not decrypt the data...");
              console.log("Could not get data");
              // TODO
              // Should we send response to IRC server?
              socket.end();
            }
            // 2. Parse and Split data
            try {
              console.log(`Trying to parse plain data == ${decryptedData}`);
              const jsonData = JSON.parse(decryptedData);
              console.log(`Data received: ${JSON.stringify(jsonData)}`);
              // 3. Get admin auth token (IRC server should only be able to talk to me!)
              if (!(jsonData?.token)) {
                console.log(`No token section included in request`);
                throw new Error("Error");
              }
              
              // 4. Check admin auth token
              console.log(`Checking ChatServer admin token: ${jsonData["token"]}`);
              await check_admin_token(jsonData["token"]);

              // TODO
              // check token expiration here??
              
              // 5. Parse message
              /**
               * # query existing thing
               *   #+ 
               *   #- 
               */

              const data = jsonData["data"];
              console.log(`Before data validation: ${data}`);
              if (data[0] !== "#") {
                console.log("Not a valid data section");
                throw new Error("Error");
              }

              if (data.includes("..")) {
                // prevents directory traversals, e.g. "./modules/../../somefile.js"
                console.log("data cannot contain ANY '..' characters");
                throw new Error("Error");
              }

              const args = data.split("::");
              console.log(`\nAuth Args:\n${args}\n`);
              const mod = `./modules/${args[0].substring(1)}.js`;
              console.log(`Searching for module: ${mod}`);
              const stats = await statAsync(mod); // if file doesn't exist this will throw error
              console.log(`Module Found!!`)
              args.shift(); // remove filename from args
              console.log(`ARGS: ${JSON.stringify(args)}`);
              const module = require(mod);
              console.log(`After mod require`);
              if (args[0] === "authcheck") {
                console.log(`args[0] === "authcheck"`)
                console.log(`Passing in args:\n${args}`);
                const authCheckRes = await module.AuthCheck(args[1]);
                socket.write(authCheckRes);
              } else {
                const res = await module.Exec(args); // pass rest of args into Exec method
                console.log(`After mod exec`);
                if (!res) {
                  console.log(`Error running module ${mod}`);
                  throw new Error("Error");
                }
              
                socket.write(res);
              }
            } catch (error) {
              console.log(error);
              socket.write("ya fucked up kid");
              socket.end();
            }
          });
          server.on("error", (err) => {
            console.log(`Server error: ${err}`);
            throw err;
          });
          socket.on("end", async () => {
            console.log(`client disconnected`);
          });
          socket.on("close", () => {
            console.log(`client closed connection`);
          });
          socket.on("drain", () => {
            // can be used to throttle uploads
            // possible can be used for rate limiting requests
            console.log(`drain event`);
          });
          socket.on("timeout", () => {
            // TODO
            console.log(`connection timeout`);
          });
        });
        server.listen(port, host, () => {
          console.log('server started:', server.address());
        })
      }
    } catch (e) {
      console.log(e.message);
    }
  }

  return { Start }
};

try {
  s = AuthServer();
  s.Start();
} catch (error) {
  console.log(`AuthServer error: ${error}`);
}
