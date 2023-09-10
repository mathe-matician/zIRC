const { NICK } = require("./NICK");
const { PASS } = require("./PASS");
const { Numerics } = require("../numerics");
const { CRLF } = require("../constants");

const LOGIN = async (parameters, clients, clientSocket, client) => {
    if (parameters.length !== 2) {
      return Numerics["ERR_NEEDMOREPARAMS"]("LOGIN");
    } else {
      try {
        const nickRes = await NICK(clients, parameters[0], "", clientSocket.remoteAddress);
        console.log(`NICKRES == ${nickRes}`);
        if (nickRes?.err) {
          console.log(`NICKRES ERROR: ${nickRes["err"]}`);
          return nickRes["err"];
        }
        const passRes = await PASS(parameters[1], parameters[0], client, clientSocket.remoteAddress);
        if (passRes?.err) {
          return passRes["err"];
        }
        // TODO
        // parse nick from full source:
        const insertRes = await InsertClient(nickRes["nick"], parameters[1], "zach");
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
        console.log(`LOGIN ERROR: ${error}`)
        return {"err": Numerics["ERR_UNKNOWNERROR"]("LOGIN")};
      }
    }
};

module.exports = { LOGIN };