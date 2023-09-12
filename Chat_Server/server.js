const net = require('node:net');
const tls = require('tls');
const Channel = require('./channel');
const fs = require('node:fs');
const {
  Commands,
  UNAUTHENTICATED_CMDS
} = require('./commands');
const { 
  ERROR,
  RPL_WELCOME,
  ERR_INPUTTOOLONG,
  Numerics,
} = require('./numerics');
const { CRLF } = require("./constants");
const { Client } = require("./client");
// const chilkatManager = require('./Chilkat/chilkat_manager');
const { UpdateOne, FindOne, InsertOne } = require('./db');
const { AuthServer } = require("./auth_server_comm");
require('dotenv').config();
// const { v4 } = require('uuid');

const _logger = require('pino')();
const logger = _logger.child({ Service: 'Chat Server' });

// Servers SHOULD pick a name which contains a dot character (".", 0x2E). 
//This can help clients disambiguate between server names and nicknames in a message source.

const Server = (
  serverName=process.env.IRC_DEFAULT_SERVER_NAME, // max 63 chars
  type="",
  host=process.env.IRC_HOST,
  tls=process.env.IRC_ENABLE_TLS === "false" ? false : true,
  timestamp=Date.now(), // datetime created
  version=process.env.IRC_VERSION, // IRCv3.2
  options={}, // e.g. supported channel types
  ) => {


    const name = serverName;

    // prepopulate the server with a General channel
    // TEST ADMIN
    const testUsers = {
      "GymleteAdmin": {
        "clientObj": Client("GymleteAdmin"),
        "prefixes": ["~"]
      }, // setup default admin
    };

    // TEST CHANNEL
    const chan1 = Channel(
      "General",
      { "+l": "1000" }, // chan limit
      "=",
      {},
      "#",
      { 
        topic: "Welcome to zIRC!",
        timestamp: new Date().toISOString(),
        nick: "zIRCAdmin"
      },
      testUsers
    );

    logger.info("Done creating channel");

    const channels = {
      "Regular": {
        "General": chan1
      },
      "Local": {
        "General": chan1
      }
    }; // mapping of channel name to Channel object
    const clients = {}; // map of NICK to USER?
    const servers = {}; // this will be object of adjacency lists
    // see https://ircv3.net/specs/extensions/account-registration#architecture for draft/account-registration capability spec
    let capabilities = null; // init in start method below
    const _parameters = {}; 
    const callbackReqs = {
      "clients": clients,
      "capabilities": null,
      "serverVersion": version,
      "isClient": true,
      "clientSocket": null,
    };

    const verifyServerName = () => {
      // const invalidServerName = //;
      // max 63 chars
      // Servers SHOULD pick a name which contains a dot character (".", 0x2E). 
      //This can help clients disambiguate between server names and nicknames in a message source.
    }

    const Configure = (options = {}, connListener = null) => {
        logger.info('configure');
    }

    const TYPE_REGULAR = "#"; // known to all servers that are connected to the network
    const TYPE_LOCAL = "&"; // server-specific. clients connected can only see and talk to other clients on the same server

    const inferChannelType = () => {
      // 

      return TYPE_REGULAR
    }

    // const CreateChannel = (client = null, name = "", modes = []) => {
    //   // e.g. common modes by defaulton channels could be "+t" = protected topic, and "+n" = no external messages
    //   // but the modes are completely up to the server

    //   if (channels.has(name)) {
    //     logger.info(`Channel ${name} already exists on server.`);
    //     return;
    //   }

    //   _Channel = await Channel(
    //     name, 
    //     modes,
    //     [client]
    //     );
    // }

    const registerClient = (message) => {
      // The recommended order of commands during registration is as follows:
      // CAP LS 302
      // PASS
      // NICK <- from client // could automatically be chosen
      // USER <- pull from database
      // Capability Negotiation
      // SASL (if negotiated)
      // CAP END
      
      logger.info(`Register Client on server ${name}`)

      // if NICK exists return ERR_NICKNAMEINUSE else return RPL_WELCOME

      // SEND RPL_ISUPPORT numeric to advertise server features to client
      // THe server MUST send at least one to the client and MAY send more. If more are sent they SHOULD be sent adjacent to each other
      return [RPL_WELCOME, "Welcome!"]
    }

    const processMessage = async (server, message, clientSocket, clients) => {
      // if command is sent from client to server with less parameters than the cmd requires
      // then the server replys with ERR_NEEDMOREPARAMS (461) e.g. return errors.ERR_NEEDMOREPARAMS;

      /**
       * Parameters section:
       * 
       * Curly braces around a part of parameter indicate that it may be repeated zero or more times.
       * E.g. <key>{, <key>}
       * this indicates that there must be at least one <key>, and that there may be additional keys separated
       * by a comma (",", 0x2C) character.
       */

  //     message         ::= ['@' <tags> SPACE] [':' <source> SPACE] <command> <parameters> <crlf>
  // SPACE           ::=  %x20 *( %x20 )   ; space character(s)
  // crlf            ::=  %x0D %x0A        ; "carriage return" "linefeed"

      // ORDER OF MESSAGE????
      // tags: Optional metadata on a message, starting with ('@', 0x40).
      // source: Optional note of where the message came from, starting with (':', 0x3A).
      // command: The specific command this message represents.
      // parameters: If it exists, data relevant to this specific command.

      // TODO auth for gymlete irc
      // TODO auth for user to server
      // TODO auth for user to channel
      // logger.info(`processMessage got data ${message}`);
      // logger.info(`TYPE OF = ${typeof(message)}`);
      // if (!message.isEncoding('utf8')) {
      //   logger.info("Client message not UTF8...");
      //   throw new Error("Client message not UTF8...");
      // }

      const tags = {};
      let source = "";
      message = message.trim();
      logger.info(`Raw message: '${message}'`);
      const captureUntilSpace = /^[^\s]*/g;

      // TODO
      // offload this to auth server?
      let _tokenPkg = "";
      let splitTokenPkg = [];
      try {
        // TODO
        // need to not check the client's expiration as this could be forged
        // need to get token, then check in db.
        // Actually I think we check the actual one after this - need to figure that out!
        if (message.startsWith("tokenPkg")) {
          logger.info(`message.startsWith("tokenPkg")`)
          const tokenpkg = message.slice(0).match(captureUntilSpace);
          message = message.slice(tokenpkg[0].length+1);
          logger.info(`\nMessage after tokenPkg parsing:\n${message}\n`)
          _tokenPkg = tokenpkg[0];
          splitTokenPkg = tokenpkg[0].split("::");
          splitTokenPkg.shift(); // rm the string 'tokenPkg' so the array only contains the token and expr
          logger.info(`\nsplitTokenPkg: ${splitTokenPkg}\n`);
          // check whether token is expired
          const now = new Date();
          const expr = new Date(splitTokenPkg[1]);
          if (now > expr) {
            logger.info(`--- Client token expired!! ---`);
            // client must re-authenticate
            clientSocket.write(Numerics["RPL_LOGGEDOUT"]());
            return null;
          }
        }
      } catch (error) {
        logger.info(`Client message does not contain a token: ${error}`);
      }
      logger.info(`tokenPkg == ${_tokenPkg}`);



      /**
       * Parse tag data
       */
      if (message[0] === "@") {
        const MAX_TAG_DATA = 4094; // actually: 4096 - '@' - ' ' == 4094
        // tags parse them
        // capture word til '=' sign, then capture until ';'
        // TODO handle out of bounds array check here
        const rawTags = message.slice(1).match(captureUntilSpace);
        message = message.slice(rawTags[0].length+2);
        if (rawTags.length > MAX_TAG_DATA) {  
          // TOOOOOO BIG TAG DATA
          logger.info(`Tag data too big yo: ${rawTags.length} > ${MAX_TAG_DATA}`);
          return Numerics["ERR_INPUTTOOLONG"];
        }
        logger.info(`RawTags = ${rawTags}`);
        const splitTags = rawTags.split(";");
        logger.info(`Tags = ${splitTags}`);
        for (const tag of splitTags) {
        	const t = tag.split("=");
          if (t.length != 2) {
          	logger.info(`ERROR PARSING TAG ${tag} DROP THIS SUCKA?`);
          } else {
            // add tags to tag object
            tags[t[0]] = t[1];
          }
        }
      } 

      /**
       * Parse message source (prefix)
       * 
       * source ::=  servername / ( nickname [ "!" user ] [ "@" host ] )
       * e.g. ":Gymlete.IRC@irc.gymlete.com"
       *      ":Gymlete.IRC"
       *      ":zach_mathe"
       *      ":zach_mathe!Zach"
       * 
       * (Optional!)
       * If there are no tags it MUST be the first character of the message itself.
       * Clients MUST NOT include a source when sending a message.
       * Servers MAY include a source on any message
       * 
       */

      logger.info(`Checking for message source ${message}`);
      if (message[0] === ":") {
        logger.info(`Source exists`)
        // from a server
        // TODO check for out of bounds (if i+1 is out of bounds of array)
        const res = message.slice(1).match(captureUntilSpace);
        message = message.slice(res[0].length+2);
        logger.info(`\nMessage after source parsing:\n${message}\n`);
        source = res;
      }
      
      /**
       * Parse command and parameters
       */
      let returnCode = "";
      let response = "";
      let returnParams = "";
      let nickname = "";
      // let clientIdentifier = "";
      let clientIP;
      let clientNick;
      let clientName;

      if (String(source).includes("@")) {
        const tmpSplit = String(source).split("@");
        clientIP = tmpSplit[1]
        if (tmpSplit.includes("!")) {
          const tmpSplit2 = tmpSplit.split("!");
          clientNick = tmpSplit2[0];
          clientName = tmpSplit2[1];
        }
      }

      // TODO
      // if this is a server sending a message, set
      // callbackReqs["isClient"] = false
      // THEN before we exit this function, check if it is false, and flip it back to true

      const formattedSource = (clientName ? "!" + clientName : "") + (clientIP ? "@" + clientIP : "");

      const clientIdentifiers = {
        nick: clientNick,
        name: clientName,
        ip: clientIP,
        source: (clientNick ? clientNick : "*") + formattedSource // wasteful, but saves doing this stupid logic checking later
      };

      logger.info(`clientIdentifiers source: ${clientIdentifiers.source}`);

      let parameters = [];
      // special parsing for PRIVMSG as the message can contain special characters like ":" and " "
      if (message.startsWith("PRIVMSG")) {
        logger.info(`Command == PRIVMSG`);
        const messageIndex = message.indexOf(":");
        const msgStart = message.substring(0, messageIndex);
        const msg = message.substring(messageIndex, message.length);
        const msgStartSplit = msgStart.trim().split(' ');
        msgStartSplit.push(msg);
        parameters = msgStartSplit;
      } else {
        parameters = message.split(" ");
      }
      logger.info(`SERVER parameters = ${parameters}`);
    
      const cmd = parameters.shift(); // get command
      // REQUIRES CLIENT MESSAGE TO ALWAYS END WITH A SPACE TO BE ABLE TO POP \r\n
      // parameters.pop(); // get rid of the \r\n
      logger.info(`Message = '${message}', Command = '${cmd}', Parameters = '${parameters}', Token = ${splitTokenPkg[0]}`);
      // const resObj = {res: "", parameters: parameters, client: clientIdentifiers};
      if (cmd in Numerics) {
        // this is a message from a server
        logger.info("Message is numeric from server");
      } else if (cmd in Commands) {
        // this is a command from a user
        logger.info("Message is command from user");

        // TODO
        // Need to pass _tokenPkg in to all commands here.
        // then in all commands that a client SHOULDN"T access without being logged in
        // need to add auth check.
        if (!(cmd in UNAUTHENTICATED_CMDS)) {
          // if it is not a command that can be run without authentication, authenticate.
          // TODO
          // get auth type form mongo. if there is no auth type in db then the user hasn't authenticated yet.
          
          // TODO
          // This check should be different.
          // We shouldn't lookup by clientIP
          // we probably should be checking token here?
          // or using client uuid as the client IP can change or be spoofed.
          // const findRes = await FindOne(
          //   {ip: clientIP}, 
          //   process.env.MONGODB_CHAT_USERS_COLLECTION_NAME, 
          //   process.env.MONGODB_NAME, {$project: "state"});
          // logger.info(`FIND RES === ${findRes}`);
          // logger.info(`FIND RES === ${JSON.stringify(findRes)}`);
          // const authType = findRes?.state?.auth?.type;
          // if (!authType) {
          //   logger.info(`No auth type set for client`);
          //   // no auth type set, so client has not authenticated.
          //   // let them know they MUST auth??
          //   // or fail silently?
          //   clientSocket.write(Numerics["ERR_NOTREGISTERED"]());
          //   return null;
          // }
          // pass only the token
          // const args = `#auth_${authType}::authcheck::${splitTokenPkg[0]}`;
          const args = `#auth_plain::authcheck::${splitTokenPkg[0]}`;
          logger.info(`CHAT SERVER BEFORE AUTH CHECK: ${args}`);
          let authServer = AuthServer();
          const authServerRes = await authServer.Write(args);
          logger.info(`Auth check res === ${authServerRes}`);
          authServer = null; // mark for garbage collection.
          
          const authRes = JSON.parse(authServerRes);
          if (authRes?.err) {
            logger.info(`ERROR: AUTHENTICATE error with Auth Server: ${authRes["err"]}`);
            return {"err": Numerics[authRes["err"]]()};
          }
        }

        const cmdRes = await Commands[cmd](parameters, clients, clientSocket);
        if (!cmdRes) {
          return null;
        }
        if (typeof(cmdRes) === "object" && cmdRes?.err) {
          return cmdRes["err"];
        }
        const req = {};
        logger.info(`callbackReqs[capabilities] === ${Object.keys(callbackReqs["capabilities"])}`)
        if (typeof(cmdRes) === "object" && "callback" in cmdRes && "req" in cmdRes) {
          for (const request of cmdRes["req"]) {
            // get requested params for callback and add them to req
            if (callbackReqs[request]) {
              logger.info(`Adding capability ${request}:${callbackReqs[request]}`);
              req[request] = callbackReqs[request];
            }
          }
          logger.info(`Before Callback: ${Object.keys(req)} value === ${req["capabilities"]}`);
          const callbackRes = await cmdRes["callback"](req);
          logger.info(`callbackRes = ${callbackRes}`);
          if (callbackRes?.err) {
            logger.info("CallbackRes ERROR");
            return callbackRes["err"];
          }
          if (callbackRes?.res) {
            return callbackRes["res"];
          } 
          if (!callbackRes) {
            logger.info(`NO CALLBACKRES`);
            return null;
          }
        } 
        // else {
        //   logger.info("Error when trying to execute callback");
        //   return Numerics["ERR_UNKNOWNCOMMAND"](cmd);
        // }

        if (cmdRes?.res) {
          logger.info(`CMDRES == ${cmdRes.res}`);
        }
        return clientIdentifiers.source ? clientIdentifiers.source + " " + cmdRes + CRLF : cmdRes + CRLF;

      } else {
        // else not a valid command or numeric
        logger.info("Not a valid command or numeric");
        return Numerics["ERR_UNKNOWNCOMMAND"](cmd);
      }

      // TODO
      // if source is from another server, return that response to that server? or does that happen automatically? I think automatically via the socket request.
      
      // clientIdentifier:
      // Format: <nick>!<user>@<host>
      // Examples:
      //    zachmathe
      //    zachmathe!zachmathe
      //    zachmathe!zachmathe@263.83.482.38

      // TODO
      // how to get the clientIdentifierf
      //const clientIdentifier = `${clients[]}`;

      // const clientIdentifier = `clientIdentifier`;

      // const res = `:${host} ${returnCode} ${response}${clientIdentifier}\r\n`;

      // return res;

      // TODO
      // if client, give client here
      // return resObj.client.source ? resObj.client.source + " " + resObj.res + CRLF : resObj.res + CRLF;
    }

    const Start = async () => {

      // chilkatManager.unlock_bundle();

      const findRes = await FindOne(
        {capabilities: {$exists: true}}, 
        process.env.MONGODB_CHAT_SERVER_COLLECTION_NAME
        );
      if (!findRes) {
        logger.info(`COULD NOT FIND CAPABILITIES! FAILING TO START`);
        return;
      }
      logger.info(`findRes == ${findRes}`);
      logger.info(`findRes == ${JSON.stringify(findRes)}`);
      capabilities = findRes["capabilities"];
      callbackReqs["capabilities"] = capabilities;

      logger.info(`CAPABILITIES == ${Object.keys(capabilities)}`);

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
            logger.info('server connected',
                        socket.authorized ? 'authorized' : 'unauthorized');
            socket.setEncoding('utf8');
            // socket.pipe(socket);
            socket.on('data', (data) => {

              // TODO: Process Message

              const client = data.toString();
              msg = client + " connected"
              logger.info(msg);
              // socket.end(client);
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
          server.listen(process.env.IRC_TLS_PORT, host, () => {
            logger.info('server started:', server.address());
          });
        } else {
          const server = net.createServer(async (socket) => {
            socket.setEncoding('utf8');
            const insertRes = await UpdateOne({"ip": socket.remoteAddress}, { $set: {"ip": socket.remoteAddress, "state": {}}}, {"upsert": true});
            if (insertRes?.err) {
              logger.info(`ERROR ${insertRes["err"]}`);
              server.emit("end");
            }
            // socket.write(helloMsg);
            // c.pipe(c);
            socket.on('data', async (data) => {
              // TODO process message
              // const msgs = data.split("\r\n");
              // let response = null;
              // for (msg of msgs) {
              //   response = processMessage(msgs);
              // }

              // TEST
              // const response = parseMessage(server, data);
              callbackReqs["clientSocket"] = socket;
              const response = await processMessage(server, data, socket, clients);

              // const _data = data.toString();
              // msg = _data + " connected"
              // msg = "{" + socket.remoteAddress + "}: " + _data
              // logger.info(msg);
              logger.info(`processMessage response = '${response}'`);
              if (response)
                socket.write(response);
              // socket.end(client); // closes the connection prematurely and throws and error?
            });
            server.on("error", (err) => {
              logger.info(`Server error: ${err}`);
              throw err;
            });
            socket.on("end", async () => {
              logger.info(`client disconnected`);
              // clear all user state before disconnecting
              callbackReqs["clientSocket"] = null;
              const updateRes = await UpdateOne({ip: socket.remoteAddress}, {$set: {state: {}}});
              // todo
              // check for errors from this res
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

          server.listen(process.env.IRC_PORT, host, async () => {
            logger.info('server started:', server.address());
            logger.info(`Max connections = ${server.maxConnections}`);
            // const updateRes = await UpdateOne({ip: socket.remoteAddress}, {$set: {state: {}}});
          })
        }
        // if (e.code === 'EADDRINUSE') {
        //   console.error('Address in use, retrying...');
        //   setTimeout(() => {`
        //     server.close();
        //     server.listen(PORT, HOST);
        //   }, 1000);
        // }
      } catch (e) {
        logger.info(e.message);
      }
    }

    return { Start }
};

try {
  s = Server();
  s.Start();
} catch (error) {
  logger.info(`Server error: ${error}`);
}


// module.exports = Server;
