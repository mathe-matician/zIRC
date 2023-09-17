const net = require('node:net');
const tls = require('tls');
const fs = require('node:fs');
const util = require('util');
const cryptoManager = require('./Crypto/crypto_manager');
const { DB } = require('./db');
require('dotenv').config();
const bcrypt = require('bcrypt');
const saltRounds = 10;
const _logger = require('pino')();
const logger = _logger.child({ Service: 'Auth Server' });

const host = process.env.AUTH_SERVER_HOST;
const port = process.env.AUTH_SERVER_PORT;
const statAsync = util.promisify(fs.stat);

const AuthServer = (tls=false) => {
  const db = DB();

  const check_admin_token = async (token) => {
    const prp_stmt_AuthCheck = db.Prepare(
      "AuthCheck", 
      "SELECT (token, token_expr) FROM admin_tokens WHERE id = 1"
      );
    const authCheckRes = await db.Exec(prp_stmt_AuthCheck);
    logger.info(`authCheckRes: ${JSON.stringify(authCheckRes)}`);
    if (!authCheckRes) {
      logger.error("Failed auth check!");
      throw new Error("Error");
    }

    const checkResult = authCheckRes["row"];
    const rowValues = checkResult.substring(
      checkResult.indexOf("(") + 1,
      checkResult.lastIndexOf(")")
    ).split(",");
    const _expiration = rowValues[1];
    logger.info(rowValues);
    
    const now = new Date();
    const expiration = new Date(_expiration);
    if (now > expiration) {
      // TODO
      // should these even expire...?? probably...
      logger.error(`Admin auth token is expired.\nNow: ${now}\nExpiration: ${expiration}`);
      throw new Error("Error");
    }

    const compareRes = await cryptoManager.compare_password(token, rowValues[0]);
    if (!compareRes) {
      logger.error("Server Admin key does not match Auth Server Token");
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
          logger.info('server connected: ', socket.authorized ? 'authorized' : 'unauthorized');
          socket.setEncoding('utf8');
          socket.on('data', (data) => {
          });
          server.on("error", (err) => {
            logger.info(`Server error: ${err}`);
            throw err;
          });
          c.on("end", (client) => {
            logger.info(`${client} disconnected`);
          });
        });
        // });
        server.listen(port, host, () => {
          logger.info('server started:', server.address());
        });
      } else {
        const server = net.createServer(async (socket) => {
          socket.setEncoding('utf8');
          socket.on('data', async (data) => {
            logger.info(`{"remoteAddress": "${socket.remoteAddress}", "data": "${data}"}`);
            // 1. Decrypt data
            const decryptedData = data;
            if (!decryptedData) {
              // logger.info("Could not decrypt the data...");
              logger.info("Could not get data");
              // TODO
              // Should we send response to IRC server?
              socket.end();
            }
            // 2. Parse and Split data
            try {
              logger.info(`Trying to parse plain data == ${decryptedData}`);
              const jsonData = JSON.parse(decryptedData);
              logger.info(`Data received: ${JSON.stringify(jsonData)}`);
              // 3. Get admin auth token (IRC server should only be able to talk to me!)
              if (!(jsonData?.token)) {
                logger.info(`No token section included in request`);
                throw new Error("Error");
              }
              
              // 4. Check admin auth token
              logger.info(`Checking ChatServer admin token: ${jsonData["token"]}`);
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
              logger.info(`Before data validation: ${data}`);
              if (data[0] !== "#") {
                logger.info("Not a valid data section");
                throw new Error("Error");
              }

              if (data.includes("..")) {
                // prevents directory traversals, e.g. "./modules/../../somefile.js"
                logger.info("data cannot contain ANY '..' characters");
                throw new Error("Error");
              }

              const args = data.split("::");
              logger.info(`\nAuth Args:\n${args}\n`);
              if (args.length > 3) {
                // args will always be 3...
                // whether it is auth_type::authcheck::splittokenpkg[0]
                // OR auth_${_authType}::${clientAuthStep}::${modAuthParams.toLowerCase()}
                // currently only these two types are sent.
                // this prevents the user from sending a string containing extra :: characters
                // and causing unexpected behavior
                logger.info("Auth Server: args are greater than 3... This is unexpected!");
                throw new Error("Error");
              }

              if (args[0].includes(".")) {
                // avoid directory traversal by checking whether args[0] contains any periods.
                logger.info("ERROR: File name contains periods... not allowed.");
                throw new Error("Error");
              }

              const mod = `./modules/${args[0].substring(1)}.js`;
              logger.info(`Searching for module: ${mod}`);
              const stats = await statAsync(mod); // if file doesn't exist this will throw error
              logger.info(`Module Found!!`)
              args.shift(); // remove filename from args
              logger.info(`ARGS: ${JSON.stringify(args)}`);
              const module = require(mod);
              logger.info(`After mod require`);

              if (args[0] === "authcheck") {
                logger.info(`args[0] === "authcheck"`);
                logger.info(`Passing in args:\n${args}`);
                const authCheckRes = await module.AuthCheck(args[1]);
                logger.info(`AuthCheck returned: ${authCheckRes}`);
                socket.write(authCheckRes);
              } else if (args[0] === "register") {
                logger.info("args[0] === register");
                const registerRes = await module.Register(args[1]);
                logger.info(`Register returned: ${registerRes}`);
                socket.write(registerRes);
              } else {
                const res = await module.Exec(args); // pass rest of args into Exec method
                logger.info(`After mod exec`);
                if (!res) {
                  logger.info(`Error running module ${mod}`);
                  throw new Error("Error");
                  // socket.write('{"err": "ERR_UNKNOWNERROR"}');
                }
              
                socket.write(res);
              }
            } catch (error) {
              logger.info(error);
              socket.write('{"err": "ERR_UNKNOWNERROR"}');
              socket.end();
            }
          });
          server.on("error", (err) => {
            logger.info(`Server error: ${err}`);
            throw err;
          });
          socket.on("end", async () => {
            logger.info(`client disconnected`);
          });
          socket.on("close", () => {
            logger.info(`client closed connection`);
          });
          socket.on("drain", () => {
            // can be used to throttle uploads
            // possible can be used for rate limiting requests
            logger.info(`drain event`);
          });
          socket.on("timeout", () => {
            // TODO
            logger.info(`connection timeout`);
          });
        });
        server.listen(port, host, () => {
          logger.info('server started:', server.address());
        })
      }
    } catch (e) {
      logger.info(e.message);
      return {"err": ""}
    }
  }

  return { Start }
};

try {
  s = AuthServer();
  s.Start();
} catch (error) {
  logger.info(`AuthServer error: ${error}`);
}
