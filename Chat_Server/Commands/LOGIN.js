const { NICK } = require("./NICK");
const { PASS } = require("./PASS");
const { Numerics } = require("../numerics");
const { CRLF } = require("../constants");

const _logger = require('pino')();
const logger = _logger.child({ Service: 'Chat Server', Command: "LOGIN" });

// Use AUTHENTICATE as a backend and pass the auth type?

const LOGIN = async (params, clients, clientSocket, clientNickname, serverName) => {
    if (params.length !== 2) {
      return Numerics["ERR_NEEDMOREPARAMS"]("LOGIN");
    } else {
      try {
        const nickRes = await NICK(clients, params[0], "", clientSocket.remoteAddress);
        logger.info(`NICKRES == ${nickRes}`);
        if (nickRes?.err) {
          logger.info(`NICKRES ERROR: ${nickRes["err"]}`);
          return nickRes["err"];
        }
        const passRes = await PASS(params[1], params[0], client, clientSocket.remoteAddress);
        if (passRes?.err) {
          return passRes["err"];
        }
        // TODO
        // parse nick from full source:
        const insertRes = await InsertClient(nickRes["nick"], params[1], "*");
        if (insertRes?.err) {
            return insertRes;
        }

        clients[nickRes["nick"]] = Client(nickRes["nick"], clientSocket);
        // TODO
        // temporary host name here... reformat
        // resObj.res = `${nickRes["nick"]} :${nickRes["nick"]}!${nickRes["nick"]}@localhost`;
        clientSocket.write(nickRes.res + CRLF);
        return Numerics["RPL_WELCOME"](nickRes["nick"]);
      } catch (error) {
        logger.info(`LOGIN ERROR: ${error}`)
        return {"err": Numerics["ERR_UNKNOWNERROR"]("LOGIN")};
      }
    }
};

module.exports = { LOGIN };