const { ADMIN } = require("./Commands/ADMIN");
const { AUTHENTICATE } = require("./Commands/AUTHENTICATE");
const { AWAY } = require("./Commands/AWAY");
const { CAP } = require("./Commands/CAP");
const { CONNECT } = require("./Commands/CONNECT");
const { ERROR } = require("./Commands/ERROR");
const { HELP } = require("./Commands/HELP");
const { INFO } = require("./Commands/INFO");
const { INVITE } = require("./Commands/INVITE");
const { JOIN } = require("./Commands/JOIN");
const { KICK } = require("./Commands/KICK");
const { KILL } = require("./Commands/KILL");
const { LINKS } = require("./Commands/LINKS");
const { LIST } = require("./Commands/LIST");
const { LOGIN } = require("./Commands/LOGIN");
const { LUSERS } = require("./Commands/LUSERS");
const { MODE } = require("./Commands/MODE");
const { MOTD } = require("./Commands/MOTD");
const { NAMES } = require("./Commands/NAMES");
const { NICK } = require("./Commands/NICK");
const { NOTICE } = require("./Commands/NOTICE");
const { OPER } = require("./Commands/OPER");
const { PART } = require("./Commands/PART");
const { PASS } = require("./Commands/PASS");
const { PING } = require("./Commands/PING");
const { PONG } = require("./Commands/PONG");
const { PRIVMSG } = require("./Commands/PRIVMSG");
const { QUIT } = require("./Commands/QUIT");
const { REGISTER } = require("./Commands/REGISTER");
const { REHASH } = require("./Commands/REHASH");
const { RESTART } = require("./Commands/RESTART");
const { SEND } = require("./Commands/SEND");
const { SQUIT } = require("./Commands/SQUIT");
const { STATS } = require("./Commands/STATS");
const { TIME } = require("./Commands/TIME");
const { TOPIC } = require("./Commands/TOPIC");
const { USER } = require("./Commands/USER");
const { USERHOST } = require("./Commands/USERHOST");
const { VERSION } = require("./Commands/VERSION");
const { WALLOPS } = require("./Commands/WALLOPS");
const { WHO } = require("./Commands/WHO");
const { WHOIS } = require("./Commands/WHOIS");
const { WHOWAS } = require("./Commands/WHOWAS");

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
    "AUTHENTICATE": true,
    "REGISTER": true,
    // "CAP": true,
    "HELP": true
    // "RPL_ISUPPORT": true
}

module.exports = { 
    Commands,
    UNAUTHENTICATED_CMDS
};