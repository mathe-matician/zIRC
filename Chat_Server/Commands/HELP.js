const { Numerics } = require("../numerics");
const { FindOne } = require("../db");
const { CRLF } = require("../constants");
require('dotenv').config();
const _logger = require('pino')();
const logger = _logger.child({ Service: 'Chat Server', Command: "HELP" });

// TODO
// implement Trie here for more free text searching.
// e.g. HELP C -> all commands/keywords that start with C.

const Commands = {
    "CAP": true,
    "JOIN": true,
    "NICK": true,
    "SEND": true,
    "PASS": true,
    "AUTHENTICATE": true,
    "USER": true,
    "LIST": true,
    "PRIVMSG": true,
    "HELP": true
}

const HELP_USERCMDS = (clientSocket) => {
    logger.info(`Help listing user commands`)
    clientSocket.write(Numerics["RPL_HELPSTART"]("** Help User Commands **") + "\n");
    clientSocket.write(Numerics["RPL_HELPTXT"]("") + "\n");
    const cmds = Object.keys(Commands);
    for (let i = 0; i < cmds.length; i++) {
        if (i !== cmds.length-1) {
            clientSocket.write(Numerics["RPL_HELPTXT"](cmds[i]) + "\n");
        } else {
            clientSocket.write(Numerics["RPL_ENDOFHELP"](cmds[i]));
        }
    }
}

const HELP_DEFAULT = (clientSocket, addlMsgs) => {
    clientSocket.write(Numerics["RPL_HELPSTART"]() + "\n");
    if (!addlMsgs, addlMsgs.length !== 0) {
        logger.info(`Adding additional messages`)
        for (let i = 0; i < addlMsgs.length; i++) {
            clientSocket.write(Numerics["RPL_HELPTXT"](addlMsgs[i]) + "\n");
        }
    }
    clientSocket.write(Numerics["RPL_HELPTXT"]("Try HELP <commmand> for specific help") + "\n");
    // clientSocket.write(Numerics["RPL_HELPTXT"]("For example:") + "\n");
    // clientSocket.write(Numerics["RPL_HELPTXT"]("\t'HELP AUTHENTICATE' to see documentation on the AUTHENTICATE command") + "\n");
    // clientSocket.write(Numerics["RPL_HELPTXT"]("\t'HELP C' will list all available 'C' keywords that are available.") + "\n");
    clientSocket.write(Numerics["RPL_HELPTXT"]("HELP USERCMDS to list available commands,") + "\n");
    clientSocket.write(Numerics["RPL_ENDOFHELP"]("or join the #help channel"));
}

const HELP_SPECIFIC = async (cmd, clientSocket) => {
    logger.info(`Help specific command`)
    if (!(cmd in Commands)) {
        clientSocket.write(Numerics["ERR_HELPNOTFOUND"]() + "\n");
        HELP_DEFAULT(clientSocket, ["", "I do not know anything about this", ""]);
        return;
    }
    logger.info(`HELP_SPECIFIC searching for '${cmd}' typeof: ${typeof(cmd)}`)
    const cmdRes = await FindOne(
        {"cmd": String(cmd)},
        process.env.MONGODB_CHAT_COMMAND_DESCRIPTIONS_COLLECTION_NAME,
        process.env.MONGODB_CHAT_SERVER_CONFIG_DB_NAME,
        {$projection: {"description": 1}}
    );
    if (!cmdRes) {
        clientSocket.write(Numerics["ERR_HELPNOTFOUND"]() + "\n");
        HELP_DEFAULT(clientSocket, ["", "I do not know anything about this", ""]);
        return;
    }
    clientSocket.write(Numerics["RPL_HELPSTART"](`** The ${cmd} command **`) + "\n");
    const lines = cmdRes?.description.split("\n");
    for (let i = 0; i < lines.length; i++) {
        if (i !== lines.length-1) {
            clientSocket.write(Numerics["RPL_HELPTXT"](lines[i]) + "\n");
        } else {
            clientSocket.write(Numerics["RPL_ENDOFHELP"](lines[i]));
        }
    }
}

const HELP = async (params, clients, clientSocket, clientNickname, serverName) => {
    logger.info(`HELP start. params: ${params}`);
    if (params.length > 0) {
        logger.info(`params.length > 0: '${params}'`)
        if (params[0] === "USERCMDS") {
            HELP_USERCMDS(clientSocket);
        } else {
            await HELP_SPECIFIC(params, clientSocket);
        }
    } else {
        HELP_DEFAULT(clientSocket, [""])
    }
};

module.exports = {
    HELP
};