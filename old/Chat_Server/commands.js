const { ADMIN } = require("./Commands_old/ADMIN");
const { AUTHENTICATE } = require("./Commands_old/AUTHENTICATE");
const { AWAY } = require("./Commands_old/AWAY");
const { CAP } = require("./Commands_old/CAP");
const { CONNECT } = require("./Commands_old/CONNECT");
const { ERROR } = require("./Commands_old/ERROR");
const { HELP } = require("./Commands_old/HELP");
const { INFO } = require("./Commands_old/INFO");
const { INVITE } = require("./Commands_old/INVITE");
const { JOIN } = require("./Commands_old/JOIN");
const { KICK } = require("./Commands_old/KICK");
const { KILL } = require("./Commands_old/KILL");
const { LINKS } = require("./Commands_old/LINKS");
const { LIST } = require("./Commands_old/LIST");
const { LOGIN } = require("./Commands_old/LOGIN");
const { LUSERS } = require("./Commands_old/LUSERS");
const { MODE } = require("./Commands_old/MODE");
const { MOTD } = require("./Commands_old/MOTD");
const { NAMES } = require("./Commands_old/NAMES");
const { NICK } = require("./Commands_old/NICK");
const { NOTICE } = require("./Commands_old/NOTICE");
const { OPER } = require("./Commands_old/OPER");
const { PART } = require("./Commands_old/PART");
const { PASS } = require("./Commands_old/PASS");
const { PING } = require("./Commands_old/PING");
const { PONG } = require("./Commands_old/PONG");
const { PRIVMSG } = require("./Commands_old/PRIVMSG");
const { QUIT } = require("./Commands_old/QUIT");
const { REGISTER } = require("./Commands_old/REGISTER");
const { REHASH } = require("./Commands_old/REHASH");
const { RESTART } = require("./Commands_old/RESTART");
const { SEND } = require("./Commands_old/SEND");
const { SQUIT } = require("./Commands_old/SQUIT");
const { STATS } = require("./Commands_old/STATS");
const { TIME } = require("./Commands_old/TIME");
const { TOPIC } = require("./Commands_old/TOPIC");
const { USER } = require("./Commands_old/USER");
const { USERHOST } = require("./Commands_old/USERHOST");
const { VERSION } = require("./Commands_old/VERSION");
const { WALLOPS } = require("./Commands_old/WALLOPS");
const { WHO } = require("./Commands_old/WHO");
const { WHOIS } = require("./Commands_old/WHOIS");
const { WHOWAS } = require("./Commands_old/WHOWAS");

require('dotenv').config();


/**
 * Command Mappings
 */

// const ClientCommands = {
//     "JOIN": JOIN,
//     "NICK": NICK,
// }

// const ServerCommands = {
//     "MOTD": MOTD,
//     "VERSION": VERSION,
//     "ADMIN": ADMIN,
//     "CONNECT": CONNECT,
//     "LUSERS": LUSERS,
//     "TIME": TIME,
//     "STATS": STATS,
//     "HELP": HELP,
//     "INFO": INFO,
//     "MODE": MODE,
// }

// const OperatorCommands = {

// }

const Commands = {
    "CAP": CAP,
    "JOIN": JOIN,
    "NICK": NICK,
    "SEND": SEND,
    "LOGIN": LOGIN,
    "PASS": PASS,
    "AUTHENTICATE": AUTHENTICATE,
    "LIST": LIST,
    "PRIVMSG": PRIVMSG,
    "REGISTER": REGISTER,

    // "MOTD": MOTD,
    // "VERSION": VERSION,
    // "ADMIN": ADMIN,
    // "CONNECT": CONNECT,
    // "LUSERS": LUSERS,
    // "TIME": TIME,
    // "STATS": STATS,
    "HELP": HELP,
    // "INFO": INFO,
    // "MODE": MODE,
}

const UNAUTHENTICATED_CMDS = {
    // "AUTHENTICATE": true,
    "REGISTER": true,
    "VERIFY": true,
    // "NICK": true,
    // "USER": true,
    // "CAP": tru
    "HELP": true
    // "RPL_ISUPPORT": true
}

module.exports = { 
    Commands,
    UNAUTHENTICATED_CMDS
};