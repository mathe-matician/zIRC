/**
 * As mentioned in the numeric replies section, the first parameter of most numerics is the target of that numeric (the nickname of the client that is receiving it). Underneath the name and numeric of each reply, we list the parameters sent by this message.

Clients MUST NOT fail because the number of parameters on a given incoming numeric is larger than the number of parameters we list for that numeric here. Most IRC servers extends some of these numerics with their own special additions. For example, if a message is listed here as having 2 parameters, and your client receives it with 5 parameters, your client should not fail to parse or handle that message correctly because of the extra parameters.

Optional parameters are surrounded with the standard square brackets ([<optional>]) – this means clients MUST NOT assume they will receive this parameter from all servers, and that servers SHOULD send this parameter unless otherwise specified in the numeric description. Parameters and parts of parameters surrounded with curly brackets ({ <repeating>}) may be repeated zero or more times.

Server authors that wish to extend one of the numerics listed here SHOULD make their extension into a client capability. If your extension would be useful to other client and server software, you should consider submitting it to the IRCv3 Working Group for standardisation.

Note that for numerics with “human-readable” informational strings for the last parameter which are not designed to be parsed, such as in RPL_WELCOME, servers commonly change this last-param text. Clients SHOULD NOT rely on these sort of parameters to have exactly the same human-readable string as described in this document. Clients that rely on the format of these human-readable final informational strings may fail. We do try to note numerics where this is the case with a message like “The text used in the last param of this message varies wildly”.
 */
const { CRLF } = require("./constants");
require('dotenv').config();

const RPL_WELCOME = (nick, server="zIRC") => {
  return `Welcome to ${server} ${nick}!`;
}

const RPL_YOURHOST = "002";
const RPL_CREATED = "003";
const RPL_MYINFO = "004";
const RPL_ISUPPORT = "005";
const RPL_BOUNCE = "010";
const RPL_UMODEIS = "221";
const RPL_LUSERCLIENT = "251";
const RPL_LUSEROP = "252";
const RPL_LUSERUNKNOWN = "253";
const RPL_LUSERCHANNELS = "254";
const RPL_LUSERME = "255";
const RPL_ADMINME = "256";
const RPL_ADMINLOC1 = "257";
const RPL_ADMINLOC2 = "258";
const RPL_ADMINEMAIL = "259";
const RPL_TRYAGAIN = "263";
const RPL_LOCALUSERS = "265";
const RPL_GLOBALUSERS = "266";
const RPL_WHOISCERTFP = "276";
const RPL_NONE = "300";

const RPL_AWAY = (client, nick, message) => {
  return `${client} ${nick} 301 :${message}`;
};

const RPL_USERHOST = "302";
const RPL_UNAWAY = "305";
const RPL_NOWAWAY = "306";
const RPL_WHOREPLY = "352";
const RPL_ENDOFWHO = "315";
const RPL_WHOISREGNICK = "307";
const RPL_WHOISUSER = "311";
const RPL_WHOISSERVER = "312";
const RPL_WHOISOPERATOR = "313";
const RPL_WHOWASUSER = "314";
const RPL_WHOISIDLE = "317";
const RPL_ENDOFWHOIS = "318";
const RPL_WHOISCHANNELS = "319";
const RPL_WHOISSPECIAL = "320";

const RPL_LISTSTART = () => {
  return "Channel 321 :Users  Name";
};

const RPL_LIST = (channel, clientCount, channelTopic) => {
  return `${channel} ${clientCount} 322 :${channelTopic}`;
};

const RPL_LISTEND = () => {
  return "323 :End of /LIST";
};

const RPL_CHANNELMODEIS = "324";
const RPL_CREATIONTIME = "329";
const RPL_WHOISACCOUNT = "330";
const RPL_NOTOPIC = "331";

const RPL_TOPIC = (channelName, topic) => {
  return `${channelName} 332 :${topic}`;
}

const RPL_TOPICWHOTIME = (channelName, nick, timestamp) => {
  return `${channelName} 333 ${nick} ${timestamp}`;
}

const RPL_INVITELIST = "336";
const RPL_ENDOFINVITELIST = "337";
const RPL_WHOISACTUALLY = "338";
const RPL_INVITING = "341";
const RPL_INVEXLIST = "346";
const RPL_ENDOFINVEXLIST = "347";
const RPL_EXCEPTLIST = "348";
const RPL_ENDOFEXCEPTLIST = "349";
const RPL_VERSION = "351";

const RPL_NAMREPLY = (status, channelName, membershipPrefixMapping) => {
  return `${status} ${channelName} "353" :${membershipPrefixMapping}`;
}

const RPL_ENDOFNAMES = (channelName) => {
  return `${channelName} 366 :End of /NAMES list`;
}

const RPL_LINKS = "364";
const RPL_ENDOFLINKS = "365";
const RPL_BANLIST = "367";
const RPL_ENDOFBANLIST = "368";
const RPL_ENDOFWHOWAS = "369";
const RPL_INFO = "371";
const RPL_ENDOFINFO = "374";
const RPL_MOTDSTART = "375";
const RPL_MOTD = "372";
const RPL_ENDOFMOTD = "376";
const RPL_WHOISHOST = "378";
const RPL_WHOISMODES = "379";
const RPL_YOUREOPER = "381";
const RPL_REHASHING = "382";
const RPL_TIME = "391";

const ERR_UNKNOWNERROR = (command, subcommand=[""], info="Unknown error occurred") => {
  let parsedSubCommands;
  for (subcmd in subcommand) {
    parsedSubCommands += subcmd;
  }

  if (parsedSubCommands.length !== 0)
    parsedSubCommands += " ";

  return `${command} ${parsedSubCommands}400 :${info}`;
}

const ERR_NOSUCHNICK = (nickname="*") => {
  return `<client> ${nickname} 401 :No such nick/channel`;
}

const ERR_NOSUCHSERVER = "402";

const ERR_NOSUCHCHANNEL = (channelName) => {
  return `${channelName} 403 :No such channel`;
}

const ERR_CANNOTSENDTOCHAN = (client="*", channel="*") => {
  `${client} ${channel} 404 :Cannot send to channel`;
};

const ERR_TOOMANYCHANNELS = "405";
const ERR_WASNOSUCHNICK = "406";
const ERR_NOORIGIN = "409";

const ERR_INVALIDCAPCMD = (client="*", command="...") => {
  return `410 ${client} ${command} :Invalid CAP command`;
}

const ERR_INPUTTOOLONG = () => {
  return "417";
};

const ERR_UNKNOWNCOMMAND = (command) => {
  return `${command} 421 :Unknown command`;
}

const ERR_NOMOTD = "422";
const ERR_NONICKNAMEGIVEN = () => {
  return "431 :No nickname given";
}

const ERR_ERRONEUSNICKNAME = (nick, message="Erroneus nickname") => {
  return `${nick} 432 :${message}`;
}

const ERR_NICKNAMEINUSE = (nickname) => {
  return `${nickname} 433 :Nickname already in use`;
}

const ERR_USERNOTINCHANNEL = "441";
const ERR_NOTONCHANNEL = "442";
const ERR_USERONCHANNEL = "443";

const ERR_NOTREGISTERED = () => {
  return `451 :You have not registered`;
}

const ERR_NEEDMOREPARAMS = (command) => {
  return `${command} 461 :Not enough parameters`;
}

const ERR_ALREADYREGISTERED = (client) => {
  return `${client} 462 :You may not reregister`;
}

const ERR_CREDSMISMATCH = (client) => {
  return `${client} 464 :Credentials incorrect`;
}
const ERR_YOUREBANNEDCREEP = "465";

const ERR_CHANNELISFULL = (channelName) => {
  return `${channelName} 471 :Cannot join channel limit reached (+l)`;
}

const ERR_UNKNOWNMODE = "472";
const ERR_INVITEONLYCHAN = "473";
const ERR_BANNEDFROMCHAN = "474";
const ERR_BADCHANNELKEY = "475";

const ERR_BADCHANMASK = (channelName) => {
  return `${channelName} 476 :Bad Channel Mask`;
}

const ERR_NOPRIVILEGES = "481";
const ERR_CHANOPRIVSNEEDED = "482";
const ERR_CANTKILLSERVER = "483";
const ERR_NOOPERHOST = "491";
const ERR_UMODEUNKNOWNFLAG = "501";
const ERR_USERSDONTMATCH = "502";

const ERR_HELPNOTFOUND = (msg="I do not know anything about this", val=" ") => {
  return `:${process.env.IRC_DEFAULT_SERVER_NAME} 524${val}* :${msg} ${CRLF}`
}

const ERR_INVALIDKEY = "525";
const RPL_STARTTLS = "670";
const RPL_WHOISSECURE = "671";
const ERR_STARTTLS = "691";
const ERR_INVALIDMODEPARAM = "696";

const RPL_HELPSTART = (msg="** Help System **", val=" ") => {
  return `:${process.env.IRC_DEFAULT_SERVER_NAME} 704${val}* :${msg}`
}

const RPL_HELPTXT = (msg="", val=" ") => {
  return `:${process.env.IRC_DEFAULT_SERVER_NAME} 705${val}* :${msg}`
} 

const RPL_ENDOFHELP = (msg="", val=" ") => {
  return `:${process.env.IRC_DEFAULT_SERVER_NAME} 706${val}* :${msg} ${CRLF}`
}

const ERR_NOPRIVS = "723";

const RPL_LOGGEDIN = (nickname) => {
  return `900 :You are now logged in as ${nickname}`;
}

const RPL_LOGGEDOUT = () => {
  return `901 :You are now logged out`;
}

const ERR_NICKLOCKED = "902";

const RPL_SASLSUCCESS = () => {
  return `903 :SASL authentication successful`;
}

const ERR_SASLFAIL = (message="SASL authentication failed") => {
  return `904 :${message}`;
}
const ERR_SASLTOOLONG = "905";

const ERR_SASLABORTED = (message="SASL authentication aborted") => {
  return `906 :${message}`;
}

const ERR_SASLALREADY = "907";

const RPL_SASLMECHS = (mechanisms) => {
  return `:server 908 <nick> ${mechanisms} :are available SASL mechanisms`;
}

const Numerics = {
  "RPL_WELCOME": RPL_WELCOME,
  "ERR_UNKNOWNERROR": ERR_UNKNOWNERROR,
  "RPL_NAMREPLY": RPL_NAMREPLY,
  "RPL_ENDOFNAMES": RPL_ENDOFNAMES,
  "ERR_CHANNELISFULL": ERR_CHANNELISFULL,
  "ERR_NEEDMOREPARAMS": ERR_NEEDMOREPARAMS,
  "ERR_BADCHANMASK": ERR_BADCHANMASK,
  "ERR_NOSUCHCHANNEL": ERR_NOSUCHCHANNEL,
  "ERR_ERRONEUSNICKNAME": ERR_ERRONEUSNICKNAME,
  "ERR_UNKNOWNCOMMAND": ERR_UNKNOWNCOMMAND,
  "ERR_NICKNAMEINUSE": ERR_NICKNAMEINUSE,
  "ERR_NONICKNAMEGIVEN": ERR_NONICKNAMEGIVEN,
  "RPL_TOPIC": RPL_TOPIC,
  "RPL_TOPICWHOTIME": RPL_TOPICWHOTIME,
  "ERR_CREDSMISMATCH": ERR_CREDSMISMATCH,
  "ERR_ALREADYREGISTERED": ERR_ALREADYREGISTERED,
  "ERR_INVALIDCAPCMD": ERR_INVALIDCAPCMD,
  "ERR_NOTREGISTERED": ERR_NOTREGISTERED,
  "ERR_SASLABORTED": ERR_SASLABORTED,
  "RPL_SASLSUCCESS": RPL_SASLSUCCESS,
  "RPL_SASLMECHS": RPL_SASLMECHS,
  "ERR_SASLFAIL": ERR_SASLFAIL,
  "RPL_LOGGEDIN": RPL_LOGGEDIN,
  "RPL_LOGGEDOUT": RPL_LOGGEDOUT,
  "RPL_LISTSTART": RPL_LISTSTART,
  "RPL_LIST": RPL_LIST,
  "RPL_LISTEND": RPL_LISTEND,
  "ERR_INPUTTOOLONG": ERR_INPUTTOOLONG,
  "RPL_HELPSTART": RPL_HELPSTART,
  "RPL_HELPTXT": RPL_HELPTXT,
  "RPL_ENDOFHELP": RPL_ENDOFHELP,
  "ERR_HELPNOTFOUND": ERR_HELPNOTFOUND,
  "ERR_NOSUCHNICK": ERR_NOSUCHNICK,
  "ERR_CANNOTSENDTOCHAN": ERR_CANNOTSENDTOCHAN,
  "RPL_AWAY": RPL_AWAY,
}

module.exports = { 
  Numerics
};

// const Numerics = {
//   RPL_WELCOME: "Welcome!",
// RPL_YOURHOST = 002;
// RPL_CREATED = 003;
// RPL_MYINFO = 004;
// RPL_ISUPPORT = 005;
// RPL_BOUNCE = 010;
// RPL_UMODEIS = 221;
// RPL_LUSERCLIENT = 251;
// RPL_LUSEROP = 252;
// RPL_LUSERUNKNOWN = 253;
// RPL_LUSERCHANNELS = 254;
// RPL_LUSERME = 255;
// RPL_ADMINME = 256;
// RPL_ADMINLOC1 = 257;
// RPL_ADMINLOC2 = 258;
// RPL_ADMINEMAIL = 259;
// RPL_TRYAGAIN = 263;
// RPL_LOCALUSERS = 265;
// RPL_GLOBALUSERS = 266;
// RPL_WHOISCERTFP = 276;
// RPL_NONE = 300;
// RPL_AWAY = 301;
// RPL_USERHOST = 302;
// RPL_UNAWAY = 305;
// RPL_NOWAWAY = 306;
// RPL_WHOREPLY = 352;
// RPL_ENDOFWHO = 315;
// RPL_WHOISREGNICK = 307;
// RPL_WHOISUSER = 311;
// RPL_WHOISSERVER = 312;
// RPL_WHOISOPERATOR = 313;
// RPL_WHOWASUSER = 314;
// RPL_WHOISIDLE = 317;
// RPL_ENDOFWHOIS = 318;
// RPL_WHOISCHANNELS = 319;
// RPL_WHOISSPECIAL = 320;
// RPL_LISTSTART = 321;
// RPL_LIST = 322;
// RPL_LISTEND = 323;
// RPL_CHANNELMODEIS = 324;
// RPL_CREATIONTIME = 329;
// RPL_WHOISACCOUNT = 330;
// RPL_NOTOPIC = 331;
// RPL_TOPIC = 332;
// RPL_TOPICWHOTIME = 333;
// RPL_INVITELIST = 336;
// RPL_ENDOFINVITELIST = 337;
// RPL_WHOISACTUALLY = 338;
// RPL_INVITING = 341;
// RPL_INVEXLIST = 346;
// RPL_ENDOFINVEXLIST = 347;
// RPL_EXCEPTLIST = 348;
// RPL_ENDOFEXCEPTLIST = 349;
// RPL_VERSION = 351;
// RPL_NAMREPLY = 353;
// RPL_ENDOFNAMES = 366;
// RPL_LINKS = 364;
// RPL_ENDOFLINKS = 365;
// RPL_BANLIST = 367;
// RPL_ENDOFBANLIST = 368;
// RPL_ENDOFWHOWAS = 369;
// RPL_INFO = 371;
// RPL_ENDOFINFO = 374;
// RPL_MOTDSTART = 375;
// RPL_MOTD = 372;
// RPL_ENDOFMOTD = 376;
// RPL_WHOISHOST = 378;
// RPL_WHOISMODES = 379;
// RPL_YOUREOPER = 381;
// RPL_REHASHING = 382;
// RPL_TIME = 391;
// ERR_UNKNOWNERROR = 400;
// ERR_NOSUCHNICK = 401;
// ERR_NOSUCHSERVER = 402;
// ERR_NOSUCHCHANNEL = 403;
// ERR_CANNOTSENDTOCHAN = 404;
// ERR_TOOMANYCHANNELS = 405;
// ERR_WASNOSUCHNICK = 406;
// ERR_NOORIGIN = 409;
// ERR_INPUTTOOLONG = 417;
// ERR_UNKNOWNCOMMAND = 421;
// ERR_NOMOTD = 422;
// ERR_ERRONEUSNICKNAME = 432;
// ERR_NICKNAMEINUSE = 433;
// ERR_USERNOTINCHANNEL = 441;
// ERR_NOTONCHANNEL = 442;
// ERR_USERONCHANNEL = 443;
// ERR_NOTREGISTERED = 451;
// ERR_NEEDMOREPARAMS = 461;
// ERR_ALREADYREGISTERED = 462;
// ERR_PASSWDMISMATCH = 464;
// ERR_YOUREBANNEDCREEP = 465;
// ERR_CHANNELISFULL = 471;
// ERR_UNKNOWNMODE = 472;
// ERR_INVITEONLYCHAN = 473;
// ERR_BANNEDFROMCHAN = 474;
// ERR_BADCHANNELKEY = 475;
// ERR_BADCHANMASK = 476;
// ERR_NOPRIVILEGES = 481;
// ERR_CHANOPRIVSNEEDED = 482;
// ERR_CANTKILLSERVER = 483;
// ERR_NOOPERHOST = 491;
// ERR_UMODEUNKNOWNFLAG = 501;
// ERR_USERSDONTMATCH = 502;
// ERR_HELPNOTFOUND = 524;
// ERR_INVALIDKEY = 525;
// RPL_STARTTLS = 670;
// RPL_WHOISSECURE = 671;
// ERR_STARTTLS = 691;
// ERR_INVALIDMODEPARAM = 696;
// RPL_HELPSTART = 704;
// RPL_HELPTXT = 705;
// RPL_ENDOFHELP = 706;
// ERR_NOPRIVS = 723;
// RPL_LOGGEDIN = 900;
// RPL_LOGGEDOUT = 901;
// ERR_NICKLOCKED = 902;
// RPL_SASLSUCCESS = 903;
// ERR_SASLFAIL = 904;
// ERR_SASLTOOLONG = 905;
// ERR_SASLABORTED = 906;
// ERR_SASLALREADY = 907;
// RPL_SASLMECHS = 908;
// }


/*
RPL_WELCOME (001) 
  "<client> :Welcome to the <networkname> Network, <nick>[!<user>@<host>]"
The first message sent after client registration, this message introduces the client to the network. The text used in the last param of this message varies wildly.

Servers that implement spoofed hostmasks in any capacity SHOULD NOT include the extended (complete) hostmask in the last parameter of this reply, either for all clients or for those whose hostnames have been spoofed. This is because some clients try to extract the hostname from this final parameter of this message and resolve this hostname, in order to discover their ‘local IP address’.

Clients MUST NOT try to extract the hostname from the final parameter of this message and then attempt to resolve this hostname. This method of operation WILL BREAK and will cause issues when the server returns a spoofed hostname.

RPL_YOURHOST (002) 
  "<client> :Your host is <servername>, running version <version>"
Part of the post-registration greeting, this numeric returns the name and software/version of the server the client is currently connected to. The text used in the last param of this message varies wildly.

RPL_CREATED (003) 
  "<client> :This server was created <datetime>"
Part of the post-registration greeting, this numeric returns a human-readable date/time that the server was started or created. The text used in the last param of this message varies wildly.

RPL_MYINFO (004) 
  "<client> <servername> <version> <available user modes>
  <available channel modes> [<channel modes with a parameter>]"
Part of the post-registration greeting. Clients SHOULD discover available features using RPL_ISUPPORT tokens rather than the mode letters listed in this reply.

RPL_ISUPPORT (005) 
  "<client> <1-13 tokens> :are supported by this server"
The ABNF representation for an RPL_ISUPPORT token is:

  token      =  *1"-" parameter / parameter *1( "=" value )
  parameter  =  1*20 letter
  value      =  * letpun
  letter     =  ALPHA / DIGIT
  punct      =  %d33-47 / %d58-64 / %d91-96 / %d123-126
  letpun     =  letter / punct
As the maximum number of message parameters to any reply is 15, the maximum number of RPL_ISUPPORT tokens that can be advertised is 13. To counter this, a server MAY issue multiple RPL_ISUPPORT numerics. A server MUST issue at least one RPL_ISUPPORT numeric after client registration has completed. It MUST be issued before further commands from the client are processed.

When clients send a VERSION command to an external server (i.e. not the one they’re currently connected to), they receive the appropriate information from that server. That external server’s ISUPPORT tokens are sent to the client using the 105 (RPL_REMOTEISUPPORT) numeric instead of 005, to ensure that clients don’t process and start using these tokens sent by an external server. The format of the 105 message is exactly the same as RPL_ISUPPORT – the numeric itself is the only difference.

A token is of the form PARAMETER, PARAMETER=VALUE or -PARAMETER. Servers MUST send the parameter as upper-case text.

Tokens of the form PARAMETER or PARAMETER=VALUE are used to advertise features or information to clients. A parameter MAY have a default value and value MAY be empty when sent by servers. Unless otherwise stated, when a parameter contains a value, the value MUST be treated as being case sensitive. The value MAY contain multiple fields, if this is the case the fields SHOULD be delimited with a comma character (",", 0x2C). The value MAY contain escape sequences: \x20 for the space character (" ", 0x20), \x5C for the backslash character ("\", 0x5C) and \x3D for the equal character ("=", 0x3D).

If the value of a parameter changes, the server SHOULD re-advertise the parameter with the new value in an RPL_ISUPPORT reply. An example of this is a client becoming an IRC operator and their CHANLIMIT changing.

Tokens of the form -PARAMETER are used to negate a previously specified parameter. If the client receives a token like this, the client MUST consider that parameter to be removed and revert to the behaviour that would occur if the parameter was not specified. The client MUST act as though the paramater is no longer advertised to it. These tokens are intended to allow servers to change their features without disconnecting clients. Tokens of this form MUST NOT contain a value field.

The server MAY negate parameters which have not been previously advertised; in this case, the client MUST ignore the token.

A single RPL_ISUPPORT reply MUST NOT contain the same parameter multiple times nor advertise and negate the same parameter. However, the server is free to advertise or negate the same parameter in separate replies.

See the Feature Advertisement section for more details on this numeric. A list of parameters is available in the RPL_ISUPPORT Parameters section.

RPL_BOUNCE (010) 
  "<client> <hostname> <port> :<info>"
Sent to the client to redirect it to another server. The <info> text varies between server software and reasons for the redirection.

Because this numeric does not specify whether to enable SSL and is not interpreted correctly by all clients, it is recommended that this not be used.

This numeric is also known as RPL_REDIR by some software.

RPL_UMODEIS (221) 
  "<client> <user modes>"
Sent to a client to inform that client of their currently-set user modes.

RPL_LUSERCLIENT (251) 
  "<client> :There are <u> users and <i> invisible on <s> servers"
Sent as a reply to the LUSERS command. <u>, <i>, and <s> are non-negative integers, and represent the number of total users, invisible users, and other servers connected to this server.

RPL_LUSEROP (252) 
  "<client> <ops> :operator(s) online"
Sent as a reply to the LUSERS command. <ops> is a positive integer and represents the number of IRC operators connected to this server. The text used in the last param of this message may vary.

RPL_LUSERUNKNOWN (253) 
  "<client> <connections> :unknown connection(s)"
Sent as a reply to the LUSERS command. <connections> is a positive integer and represents the number of connections to this server that are currently in an unknown state. The text used in the last param of this message may vary.

RPL_LUSERCHANNELS (254) 
  "<client> <channels> :channels formed"
Sent as a reply to the LUSERS command. <channels> is a positive integer and represents the number of channels that currently exist on this server. The text used in the last param of this message may vary.

RPL_LUSERME (255) 
  "<client> :I have <c> clients and <s> servers"
Sent as a reply to the LUSERS command. <c> and <s> are non-negative integers and represent the number of clients and other servers connected to this server, respectively.

RPL_ADMINME (256) 
  "<client> [<server>] :Administrative info"
Sent as a reply to an ADMIN command, this numeric establishes the name of the server whose administrative info is being provided. The text used in the last param of this message may vary.

<server> is optional and MAY be included in responses, the server can also be gained from the <source> of this message.

RPL_ADMINLOC1 (257) 
  "<client> :<info>"
Sent as a reply to an ADMIN command, <info> is a string intended to provide information about the location of the server (i.e. city, state and country). The text used in the last param of this message varies wildly.

RPL_ADMINLOC2 (258) 
  "<client> :<info>"
Sent as a reply to an ADMIN command, <info> is a string intended to provide information about whoever runs the server (i.e. details of the institution hosting it). The text used in the last param of this message varies wildly.

RPL_ADMINEMAIL (259) 
  "<client> :<info>"
Sent as a reply to an ADMIN command, <info> MUST contain the email address to contact the administrator(s) of the server. The text used in the last param of this message varies wildly.

RPL_TRYAGAIN (263) 
  "<client> <command> :Please wait a while and try again."
When a server drops a command without processing it, this numeric MUST be sent to inform the client. The text used in the last param of this message varies wildly, and commonly provides the client with more information about why the command could not be processed (i.e., due to rate-limiting).

RPL_LOCALUSERS (265) 
  "<client> [<u> <m>] :Current local users <u>, max <m>"
Sent as a reply to the LUSERS command. <u> and <m> are non-negative integers and represent the number of clients currently and the maximum number of clients that have been connected directly to this server at one time, respectively.

The two optional parameters SHOULD be supplied to allow clients to better extract these numbers.

RPL_GLOBALUSERS (266) 
  "<client> [<u> <m>] :Current global users <u>, max <m>"
Sent as a reply to the LUSERS command. <u> and <m> are non-negative integers. <u> represents the number of clients currently connected to this server, globally (directly and through other server links). <m> represents the maximum number of clients that have been connected to this server at one time, globally.

The two optional parameters SHOULD be supplied to allow clients to better extract these numbers.

RPL_WHOISCERTFP (276) 
  "<client> <nick> :has client certificate fingerprint <fingerprint>"
Sent as a reply to the WHOIS command, this numeric shows the SSL/TLS certificate fingerprint used by the client with the nickname <nick>. Clients MUST only be sent this numeric if they are either using the WHOIS command on themselves or they are an operator.

RPL_NONE (300) 
  Undefined format
RPL_NONE is a dummy numeric. It does not have a defined use nor format.

RPL_AWAY (301) 
  "<client> <nick> :<message>"
Indicates that the user with the nickname <nick> is currently away and sends the away message that they set.

RPL_USERHOST (302) 
  "<client> :[<reply>{ <reply>}]"
Sent as a reply to the USERHOST command, this numeric lists nicknames and the information associated with them. The last parameter of this numeric (if there are any results) is a list of <reply> values, delimited by a SPACE character (' ', 0x20).

The ABNF representation for <reply> is:

  reply   =  nickname [ isop ] "=" isaway hostname
  isop    =  "*"
  isaway  =  ( "+" / "-" )
<isop> is included if the user with the nickname of <nickname> has registered as an operator. <isaway> represents whether that user has set an [away] message. "+" represents that the user is not away, and "-" represents that the user is away.

RPL_UNAWAY (305) 
  "<client> :You are no longer marked as being away"
Sent as a reply to the AWAY command, this lets the client know that they are no longer set as being away. The text used in the last param of this message may vary.

RPL_NOWAWAY (306) 
  "<client> :You have been marked as being away"
Sent as a reply to the AWAY command, this lets the client know that they are set as being away. The text used in the last param of this message may vary.

RPL_WHOREPLY (352) 
  "<client> <channel> <username> <host> <server> <nick> <flags> :<hopcount> <realname>"
Sent as a reply to the WHO command, this numeric gives information about the client with the nickname <nick>. Refer to RPL_WHOISUSER (311) for the meaning of the fields <username>, <host> and <realname>. <server> is the name of the server the client is connected to. <channel> is an arbitrary channel the client is joined to or a literal asterisk character ('*', 0x2A) if no channel is returned. <hopcount> is the number of intermediate servers between the client issuing the WHO command and the client <nick>, it might be unreliable so clients SHOULD ignore it.

<flags> contains the following characters, in this order:

Away status: the letter H ('H', 0x48) to indicate that the user is here, or the letter G ('G', 0x47) to indicate that the user is gone.
Optionally, a literal asterisk character ('*', 0x2A) to indicate that the user is a server operator.
Optionally, the highest channel membership prefix that the client has in <channel>, if the client has one.
Optionally, one or more user mode characters and other arbitrary server-specific flags.
RPL_ENDOFWHO (315) 
  "<client> <mask> :End of WHO list"
Sent as a reply to the WHO command, this numeric indicates the end of a WHO response for the mask <mask>.

<mask> MUST be the same <mask> parameter sent by the client in its WHO message, but MAY be casefolded.

This numeric is sent after all other WHO response numerics have been sent to the client.

RPL_WHOISREGNICK (307) 
  "<client> <nick> :has identified for this nick"
Sent as a reply to the WHOIS command, this numeric indicates that the client with the nickname <nick> was authenticated as the owner of this nick on the network.

See also RPL_WHOISACCOUNT, for information on the account name of the user.

RPL_WHOISUSER (311) 
  "<client> <nick> <username> <host> * :<realname>"
Sent as a reply to the WHOIS command, this numeric shows details about the client with the nickname <nick>. <username> and <realname> represent the names set by the USER command (though <username> may be set by the server in other ways). <host> represents the host used for the client in nickmasks (which may or may not be a real hostname or IP address). <host> CANNOT start with a colon (':', 0x3A) as this would get parsed as a trailing parameter – IPv6 addresses such as "::1" are prefixed with a zero ('0', 0x30) to ensure this. The second-last parameter is a literal asterisk character ('*', 0x2A) and does not mean anything.

RPL_WHOISSERVER (312) 
  "<client> <nick> <server> :<server info>"
Sent as a reply to the WHOIS (or WHOWAS) command, this numeric shows which server the client with the nickname <nick> is (or was) connected to. <server> is the name of the server (as used in message prefixes). <server info> is a string containing a description of that server.

RPL_WHOISOPERATOR (313) 
  "<client> <nick> :is an IRC operator"
Sent as a reply to the WHOIS command, this numeric indicates that the client with the nickname <nick> is an operator. This command MAY also indicate what type or level of operator the client is by changing the text in the last parameter of this numeric. The text used in the last param of this message varies wildly, and SHOULD be displayed as-is by IRC clients to their users.

RPL_WHOWASUSER (314) 
  "<client> <nick> <username> <host> * :<realname>"
Sent as a reply to the WHOWAS command, this numeric shows details about one of the last clients that used the nickname <nick>. The purpose of each argument is the same as with the RPL_WHOISUSER (311) numeric.

RPL_WHOISIDLE (317) 
  "<client> <nick> <secs> <signon> :seconds idle, signon time"
Sent as a reply to the WHOIS command, this numeric indicates how long the client with the nickname <nick> has been idle. <secs> is the number of seconds since the client has been active. Servers generally denote specific commands (for instance, perhaps JOIN, PRIVMSG, NOTICE, etc) as updating the ‘idle time’, and calculate this off when the idle time was last updated. <signon> is a unix timestamp representing when the user joined the network. The text used in the last param of this message may vary.

RPL_ENDOFWHOIS (318) 
  "<client> <nick> :End of /WHOIS list"
Sent as a reply to the WHOIS command, this numeric indicates the end of a WHOIS response for the client with the nickname <nick>.

<nick> MUST be exactly the <nick> parameter sent by the client in its WHOIS message. This means the case MUST be preserved, and if the client sent multiple nicks, this MUST be the comma-separated list of nicks, even if some of them were dropped.

This numeric is sent after all other WHOIS response numerics have been sent to the client.

RPL_WHOISCHANNELS (319) 
  "<client> <nick> :[prefix]<channel>{ [prefix]<channel>}
Sent as a reply to the WHOIS command, this numeric lists the channels that the client with the nickname <nick> is joined to and their status in these channels. <prefix> is the highest channel membership prefix that the client has in that channel, if the client has one. <channel> is the name of a channel that the client is joined to. The last parameter of this numeric is a list of [prefix]<channel> pairs, delimited by a SPACE character (' ', 0x20).

RPL_WHOISCHANNELS can be sent multiple times in the same whois reply, if the target is on too many channels to fit in a single message.

The channels in this response are affected by the secret channel mode and the invisible user mode, and may be affected by other modes depending on server software and configuration.

RPL_WHOISSPECIAL (320) 
  "<client> <nick> :blah blah blah"
Sent as a reply to the WHOIS command, this numeric is used for extra human-readable information on the client with nickname <nick>. This should only be used for non-essential information that does not need to be machine-readable or understood by client software.

RPL_LISTSTART (321) 
  "<client> Channel :Users  Name"
Sent as a reply to the LIST command, this numeric marks the start of a channel list. As noted in the command description, this numeric MAY be skipped by the server so clients MUST NOT depend on receiving it.

RPL_LIST (322) 
  "<client> <channel> <client count> :<topic>"
Sent as a reply to the LIST command, this numeric sends information about a channel to the client. <channel> is the name of the channel. <client count> is an integer indicating how many clients are joined to that channel. <topic> is the channel’s topic (as set by the TOPIC command).

RPL_LISTEND (323) 
  "<client> :End of /LIST"
Sent as a reply to the LIST command, this numeric indicates the end of a LIST response.

RPL_CHANNELMODEIS (324) 
  "<client> <channel> <modestring> <mode arguments>..."
Sent to a client to inform them of the currently-set modes of a channel. <channel> is the name of the channel. <modestring> and <mode arguments> are a mode string and the mode arguments (delimited as separate parameters) as defined in the MODE message description.

RPL_CREATIONTIME (329) 
  "<client> <channel> <creationtime>"
Sent to a client to inform them of the creation time of a channel. <channel> is the name of the channel. <creationtime> is a unix timestamp representing when the channel was created on the network.

RPL_WHOISACCOUNT (330) 
  "<client> <nick> <account> :is logged in as"
Sent as a reply to the WHOIS command, this numeric indicates that the client with the nickname <nick> was authenticated as the owner of <account>.

This does not necessarily mean the user owns their current nickname, which is covered byRPL_WHOISREGNICK.

RPL_NOTOPIC (331) 
  "<client> <channel> :No topic is set"
Sent to a client when joining a channel to inform them that the channel with the name <channel> does not have any topic set.

RPL_TOPIC (332) 
  "<client> <channel> :<topic>"
Sent to a client when joining the <channel> to inform them of the current topic of the channel.

RPL_TOPICWHOTIME (333) 
  "<client> <channel> <nick> <setat>"
Sent to a client to let them know who set the topic (<nick>) and when they set it (<setat> is a unix timestamp). Sent after RPL_TOPIC (332).

RPL_INVITELIST (336) 
  "<client> <channel>"
Sent to a client as a reply to the INVITE command when used with no parameter, to indicate a channel the client was invited to.

This numeric should not be confused with RPL_INVEXLIST (346), which is used as a reply to MODE.

Some rare implementations use 346 instead of 336 for this reply.

RPL_ENDOFINVITELIST (337) 
  "<client> :End of /INVITE list"
Sent as a reply to the INVITE command when used with no parameter, this numeric indicates the end of invitations a client received.

This numeric should not be confused with RPL_ENDOFINVEXLIST (347), which is used as a reply to MODE.

Some rare implementations use 347 instead of 337 for this reply.

RPL_WHOISACTUALLY (338) 
  "<client> <nick> :is actually ..."
  "<client> <nick> <host|ip> :Is actually using host"
  "<client> <nick> <username>@<hostname> <ip> :Is actually using host"
Sent as a reply to the WHOIS and WHOWAS commands, this numeric shows details about the client with the nickname <nick>.

<username> represents the name set by the USER command (though <username> may be set by the server in other ways).

<host> and <ip> represent the real host and IP address the client is connecting from. <host> CANNOT start with a colon (':', 0x3A) as this would get parsed as a trailing parameter – IPv6 addresses such as "::1" are prefixed with a zero ('0', 0x30) to ensure this. The resulting IPv6 is equivalent, as this is a partial expansion of the :: shorthand.

See also: RPL_WHOISHOST (378), for similar semantics on other servers.

RPL_INVITING (341) 
  "<client> <nick> <channel>"
Sent as a reply to the INVITE command to indicate that the attempt was successful and the client with the nickname <nick> has been invited to <channel>.

RPL_INVEXLIST (346) 
  "<client> <channel> <mask>"
Sent as a reply to the MODE command, when clients are viewing the current entries on a channel’s invite-exception list. <mask> is the given mask on the invite-exception list.

This numeric should not be confused with RPL_INVITELIST (336), which is used as a reply to INVITE.

This numeric is sometimes erroneously called RPL_INVITELIST, as this was the name used in RFC2812.

RPL_ENDOFINVEXLIST (347) 
  "<client> <channel> :End of Channel Invite Exception List"
Sent as a reply to the MODE command, this numeric indicates the end of a channel’s invite-exception list.

This numeric should not be confused with RPL_ENDOFINVITELIST (337), which is used as a reply to INVITE.

This numeric is sometimes erroneously called RPL_ENDOFINVITELIST, as this was the name used in RFC2812.

RPL_EXCEPTLIST (348) 
  "<client> <channel> <mask>"
Sent as a reply to the MODE command, when clients are viewing the current entries on a channel’s exception list. <mask> is the given mask on the exception list.

RPL_ENDOFEXCEPTLIST (349) 
  "<client> <channel> :End of channel exception list"
Sent as a reply to the MODE command, this numeric indicates the end of a channel’s exception list.

RPL_VERSION (351) 
  "<client> <version> <server> :<comments>"
Sent as a reply to the VERSION command, this numeric indicates information about the desired server. <version> is the name and version of the software being used (including any revision information). <server> is the name of the server. <comments> may contain any further comments or details about the specific version of the server.

RPL_NAMREPLY (353) 
  "<client> <symbol> <channel> :[prefix]<nick>{ [prefix]<nick>}"
Sent as a reply to the NAMES command, this numeric lists the clients that are joined to <channel> and their status in that channel.

<symbol> notes the status of the channel. It can be one of the following:

("=", 0x3D) - Public channel.
("@", 0x40) - Secret channel (secret channel mode "+s").
("*", 0x2A) - Private channel (was "+p", no longer widely used today).
<nick> is the nickname of a client joined to that channel, and <prefix> is the highest channel membership prefix that client has in the channel, if they have one. The last parameter of this numeric is a list of [prefix]<nick> pairs, delimited by a SPACE character (' ', 0x20).

RPL_ENDOFNAMES (366) 
  "<client> <channel> :End of /NAMES list"
Sent as a reply to the NAMES command, this numeric specifies the end of a list of channel member names.

RPL_LINKS (364) 
  "<client> * <server> :<hopcount> <server info>"
Sent as a reply to the LINKS command, this numeric specifies one of the known servers on the network.

<server info> is a string containing a description of that server.

RPL_ENDOFLINKS (365) 
  "<client> * :End of /LINKS list"
Sent as a reply to the LINKS command, this numeric specifies the end of a list of channel member names.

RPL_BANLIST (367) 
  "<client> <channel> <mask> [<who> <set-ts>]"
Sent as a reply to the MODE command, when clients are viewing the current entries on a channel’s ban list. <mask> is the given mask on the ban list.

<who> and <set-ts> are optional and MAY be included in responses. <who> is either the nickname or nickmask of the client that set the ban, or a server name, and <set-ts> is the UNIX timestamp of when the ban was set.

RPL_ENDOFBANLIST (368) 
  "<client> <channel> :End of channel ban list"
Sent as a reply to the MODE command, this numeric indicates the end of a channel’s ban list.

RPL_ENDOFWHOWAS (369) 
  "<client> <nick> :End of WHOWAS"
Sent as a reply to the WHOWAS command, this numeric indicates the end of a WHOWAS reponse for the nickname <nick>. This numeric is sent after all other WHOWAS response numerics have been sent to the client.

RPL_INFO (371) 
  "<client> :<string>"
Sent as a reply to the INFO command, this numeric returns human-readable information describing the server: e.g. its version, list of authors and contributors, and any other miscellaneous information which may be considered to be relevant.

RPL_ENDOFINFO (374) 
  "<client> :End of INFO list"
Indicates the end of an INFO response.

RPL_MOTDSTART (375) 
  "<client> :- <server> Message of the day - "
Indicates the start of the Message of the Day to the client. The text used in the last param of this message may vary, and SHOULD be displayed as-is by IRC clients to their users.

RPL_MOTD (372) 
  "<client> :<line of the motd>"
When sending the Message of the Day to the client, servers reply with each line of the MOTD as this numeric. MOTD lines MAY be wrapped to 80 characters by the server.

RPL_ENDOFMOTD (376) 
  "<client> :End of /MOTD command."
Indicates the end of the Message of the Day to the client. The text used in the last param of this message may vary.

RPL_WHOISHOST (378) 
  "<client> <nick> :is connecting from *@localhost 127.0.0.1"
Sent as a reply to the WHOIS command, this numeric shows details about where the client with nickname <nick> is connecting from.

See also: RPL_WHOISACTUALLY (338), for similar semantics on other servers.

RPL_WHOISMODES (379) 
  "<client> <nick> :is using modes +ailosw"
Sent as a reply to the WHOIS command, this numeric shows the client what user modes the target users has.

RPL_YOUREOPER (381) 
  "<client> :You are now an IRC operator"
Sent to a client which has just successfully issued an OPER command and gained operator status. The text used in the last param of this message varies wildly.

RPL_REHASHING (382) 
  "<client> <config file> :Rehashing"
Sent to an operator which has just successfully issued a REHASH command. The text used in the last param of this message may vary.

RPL_TIME (391) 
  "<client> <server> [<timestamp> [<TS offset>]] :<human-readable time>"
Reply to the TIME command. Typically only contains the human-readable time, but it may include a UNIX timestamp.

Clients SHOULD NOT parse the human-readable time.

<TS offset> is used by some servers using a TS-based server-to-server protocol (eg. TS6), and represents the offset between the server’s system time, and the TS of the network. A positive value means the server is lagging behind the TS of the network. Clients SHOULD ignore its value.

ERR_UNKNOWNERROR (400) 
  "<client> <command>{ <subcommand>} :<info>"
Indicates that the given command/subcommand could not be processed. <subcommand> may repeat for more specific subcommands.

For example, for an issue with a hypothetical command PACK, this may be returned:

  :example.com 400 dan!~d@n PACK :Could not process multiple invalid parameters
For an issue with a hypothetical command PACK with the subcommand BOX, this may be returned:

  :example.com 400 dan!~d@n PACK BOX :Could not find box to pack
This numeric indicates a very generalised error (which <info> should further explain). If there is another more specific numeric which represents the error occuring, that should be used instead.

ERR_NOSUCHNICK (401) 
  "<client> <nickname> :No such nick/channel"
Indicates that no client can be found for the supplied nickname. The text used in the last param of this message may vary.

ERR_NOSUCHSERVER (402) 
  "<client> <server name> :No such server"
Indicates that the given server name does not exist. The text used in the last param of this message may vary.

ERR_NOSUCHCHANNEL (403) 
  "<client> <channel> :No such channel"
Indicates that no channel can be found for the supplied channel name. The text used in the last param of this message may vary.

ERR_CANNOTSENDTOCHAN (404) 
  "<client> <channel> :Cannot send to channel"
Indicates that the PRIVMSG / NOTICE could not be delivered to <channel>. The text used in the last param of this message may vary.

This is generally sent in response to channel modes, such as a channel being moderated and the client not having permission to speak on the channel, or not being joined to a channel with the no external messages mode set.

ERR_TOOMANYCHANNELS (405) 
  "<client> <channel> :You have joined too many channels"
Indicates that the JOIN command failed because the client has joined their maximum number of channels. The text used in the last param of this message may vary.

ERR_WASNOSUCHNICK (406) 
  "<client> :There was no such nickname"
Returned as a reply to WHOWAS to indicate there is no history information for that nickname.

ERR_NOORIGIN (409) 
  "<client> :No origin specified"
Indicates a PING or PONG message missing the originator parameter which is required by old IRC servers. Nowadays, this may be used by some servers when the PING <token> is empty.

ERR_INPUTTOOLONG (417) 
  "<client> :Input line was too long"
Indicates a given line does not follow the specified size limits (512 bytes for the main section, 4094 or 8191 bytes for the tag section).

ERR_UNKNOWNCOMMAND (421) 
  "<client> <command> :Unknown command"
Sent to a registered client to indicate that the command they sent isn’t known by the server. The text used in the last param of this message may vary.

ERR_NOMOTD (422) 
  "<client> :MOTD File is missing"
Indicates that the Message of the Day file does not exist or could not be found. The text used in the last param of this message may vary.

ERR_ERRONEUSNICKNAME (432) 
  "<client> <nick> :Erroneus nickname"
Returned when a NICK command cannot be successfully completed as the desired nickname contains characters that are disallowed by the server. See the wire format section for more information on characters which are allowed in various IRC servers. The text used in the last param of this message may vary.

ERR_NICKNAMEINUSE (433) 
  "<client> <nick> :Nickname is already in use"
Returned when a NICK command cannot be successfully completed as the desired nickname is already in use on the network. The text used in the last param of this message may vary.

ERR_USERNOTINCHANNEL (441) 
  "<client> <nick> <channel> :They aren't on that channel"
Returned when a client tries to perform a channel+nick affecting command, when the nick isn’t joined to the channel (for example, MODE #channel +o nick).

ERR_NOTONCHANNEL (442) 
  "<client> <channel> :You're not on that channel"
Returned when a client tries to perform a channel-affecting command on a channel which the client isn’t a part of.

ERR_USERONCHANNEL (443) 
  "<client> <nick> <channel> :is already on channel"
Returned when a client tries to invite <nick> to a channel they’re already joined to.

ERR_NOTREGISTERED (451) 
  "<client> :You have not registered"
Returned when a client command cannot be parsed as they are not yet registered. Servers offer only a limited subset of commands until clients are properly registered to the server. The text used in the last param of this message may vary.

ERR_NEEDMOREPARAMS (461) 
  "<client> <command> :Not enough parameters"
Returned when a client command cannot be parsed because not enough parameters were supplied. The text used in the last param of this message may vary.

ERR_ALREADYREGISTERED (462) 
  "<client> :You may not reregister"
Returned when a client tries to change a detail that can only be set during registration (such as resending the PASS or USER after registration). The text used in the last param of this message varies.

ERR_PASSWDMISMATCH (464) 
  "<client> :Password incorrect"
Returned to indicate that the connection could not be registered as the password was either incorrect or not supplied. The text used in the last param of this message may vary.

ERR_YOUREBANNEDCREEP (465) 
  "<client> :You are banned from this server."
Returned to indicate that the server has been configured to explicitly deny connections from this client. The text used in the last param of this message varies wildly and typically also contains the reason for the ban and/or ban details, and SHOULD be displayed as-is by IRC clients to their users.

ERR_CHANNELISFULL (471) 
  "<client> <channel> :Cannot join channel (+l)"
Returned to indicate that a JOIN command failed because the client limit mode has been set and the maximum number of users are already joined to the channel. The text used in the last param of this message may vary.

ERR_UNKNOWNMODE (472) 
  "<client> <modechar> :is unknown mode char to me"
Indicates that a mode character used by a client is not recognized by the server. The text used in the last param of this message may vary.

ERR_INVITEONLYCHAN (473) 
  "<client> <channel> :Cannot join channel (+i)"
Returned to indicate that a JOIN command failed because the channel is set to [invite-only] mode and the client has not been invited to the channel or had an invite exception set for them. The text used in the last param of this message may vary.

ERR_BANNEDFROMCHAN (474) 
  "<client> <channel> :Cannot join channel (+b)"
Returned to indicate that a JOIN command failed because the client has been banned from the channel and has not had a ban exception set for them. The text used in the last param of this message may vary.

ERR_BADCHANNELKEY (475) 
  "<client> <channel> :Cannot join channel (+k)"
Returned to indicate that a JOIN command failed because the channel requires a key and the key was either incorrect or not supplied. The text used in the last param of this message may vary.

Not to be confused with ERR_INVALIDKEY, which may be returned when setting a key.

ERR_BADCHANMASK (476) 
  "<channel> :Bad Channel Mask"
Indicates the supplied channel name is not a valid.

This is similar to, but stronger than, ERR_NOSUCHCHANNEL (403), which indicates that the channel does not exist, but that it may be a valid name.

The text used in the last param of this message may vary.

ERR_NOPRIVILEGES (481) 
  "<client> :Permission Denied- You're not an IRC operator"
Indicates that the command failed because the user is not an IRC operator. The text used in the last param of this message may vary.

ERR_CHANOPRIVSNEEDED (482) 
  "<client> <channel> :You're not channel operator"
Indicates that a command failed because the client does not have the appropriate channel privileges. This numeric can apply for different prefixes such as halfop, operator, etc. The text used in the last param of this message may vary.

ERR_CANTKILLSERVER (483) 
  "<client> :You cant kill a server!"
Indicates that a KILL command failed because the user tried to kill a server. The text used in the last param of this message may vary.

ERR_NOOPERHOST (491) 
  "<client> :No O-lines for your host"
Indicates that an OPER command failed because the server has not been configured to allow connections from this client’s host to become an operator. The text used in the last param of this message may vary.

ERR_UMODEUNKNOWNFLAG (501) 
  "<client> :Unknown MODE flag"
Indicates that a MODE command affecting a user contained a MODE letter that was not recognized. The text used in the last param of this message may vary.

ERR_USERSDONTMATCH (502) 
  "<client> :Cant change mode for other users"
Indicates that a MODE command affecting a user failed because they were trying to set or view modes for other users. The text used in the last param of this message varies, for instance when trying to view modes for another user, a server may send: "Can't view modes for other users".

ERR_HELPNOTFOUND (524) 
  "<client> <subject> :No help available on this topic"
Indicates that a HELP command requested help on a subject the server does not know about.

The <subject> MUST be the one requested by the client, but may be casefolded; unless it would be an invalid parameter, in which case it MUST be *.

ERR_INVALIDKEY (525) 
"<client> <target chan> :Key is not well-formed"
Indicates the value of a key channel mode change (+k) was rejected.

Not to be confused with ERR_BADCHANNELKEY, which is returned when someone tries to join a channel.

RPL_STARTTLS (670) 
  "<client> :STARTTLS successful, proceed with TLS handshake"
This numeric is used by the IRCv3 tls extension and indicates that the client may begin a TLS handshake. For more information on this numeric, see the linked IRCv3 specification.

The text used in the last param of this message varies wildly.

RPL_WHOISSECURE (671) 
  "<client> <nick> :is using a secure connection"
Sent as a reply to the WHOIS command, this numeric shows the client is connecting to the server in a way the server considers reasonably safe from eavesdropping (e.g. connecting from localhost, using TLS, using Tor).

ERR_STARTTLS (691) 
  "<client> :STARTTLS failed (Wrong moon phase)"
This numeric is used by the IRCv3 tls extension and indicates that a server-side error occured and the STARTTLS command failed. For more information on this numeric, see the linked IRCv3 specification.

The text used in the last param of this message varies wildly.

ERR_INVALIDMODEPARAM (696) 
"<client> <target chan/user> <mode char> <parameter> :<description>"
Indicates that there was a problem with a mode parameter. Replaces various implementation-specific mode-specific numerics.

RPL_HELPSTART (704) 
"<client> <subject> :<first line of help section>"
Indicates the start of a reply to a HELP command. The text used in the last parameter of this message may vary, and SHOULD be displayed as-is by IRC clients to their users; possibly emphasized as the title of the help section.

The <subject> MUST be the one requested by the client, but may be casefolded; unless it would be an invalid parameter, in which case it MUST be *.

RPL_HELPTXT (705) 
"<client> <subject> :<line of help text>"
Returns a line of HELP text to the client. Lines MAY be wrapped to a certain line length by the server. Note that the final line MUST be a RPL_ENDOFHELP (706) numeric.

The <subject> MUST be the one requested by the client, but may be casefolded; unless it would be an invalid parameter, in which case it MUST be *.

RPL_ENDOFHELP (706) 
"<client> <subject> :<last line of help text>"
Returns the final HELP line to the client.

The <subject> MUST be the one requested by the client, but may be casefolded; unless it would be an invalid parameter, in which case it MUST be *.

ERR_NOPRIVS (723) 
  "<client> <priv> :Insufficient oper privileges."
Sent by a server to alert an IRC operator that they they do not have the specific operator privilege required by this server/network to perform the command or action they requested. The text used in the last param of this message may vary.

<priv> is a string that has meaning in the server software, and allows an operator the privileges to perform certain commands or actions. These strings are server-defined and may refer to one or multiple commands or actions that may be performed by IRC operators.

Examples of the sorts of privilege strings used by server software today include: kline, dline, unkline, kill, kill:remote, die, remoteban, connect, connect:remote, rehash.

RPL_LOGGEDIN (900) 
  "<client> <nick>!<user>@<host> <account> :You are now logged in as <username>"
This numeric indicates that the client was logged into the specified account (whether by SASL authentication or otherwise). For more information on this numeric, see the IRCv3 sasl-3.1 extension.

The text used in the last param of this message varies wildly.

RPL_LOGGEDOUT (901) 
  "<client> <nick>!<user>@<host> :You are now logged out"
This numeric indicates that the client was logged out of their account. For more information on this numeric, see the IRCv3 sasl-3.1 extension.

The text used in the last param of this message varies wildly.

ERR_NICKLOCKED (902) 
  "<client> :You must use a nick assigned to you"
This numeric indicates that SASL authentication failed because the account is currently locked out, held, or otherwise administratively made unavailable. For more information on this numeric, see the IRCv3 sasl-3.1 extension.

The text used in the last param of this message varies wildly.

RPL_SASLSUCCESS (903) 
  "<client> :SASL authentication successful"
This numeric indicates that SASL authentication was completed successfully, and is normally sent along with RPL_LOGGEDIN (900). For more information on this numeric, see the IRCv3 sasl-3.1 extension.

The text used in the last param of this message varies wildly.

ERR_SASLFAIL (904) 
  "<client> :SASL authentication failed"
This numeric indicates that SASL authentication failed because of invalid credentials or other errors not explicitly mentioned by other numerics. For more information on this numeric, see the IRCv3 sasl-3.1 extension.

The text used in the last param of this message varies wildly.

ERR_SASLTOOLONG (905) 
  "<client> :SASL message too long"
This numeric indicates that SASL authentication failed because the AUTHENTICATE command sent by the client was too long (i.e. the parameter was longer than 400 bytes). For more information on this numeric, see the IRCv3 sasl-3.1 extension.

The text used in the last param of this message varies wildly.

ERR_SASLABORTED (906) 
  "<client> :SASL authentication aborted"
This numeric indicates that SASL authentication failed because the client sent an AUTHENTICATE command with the parameter ('*', 0x2A). For more information on this numeric, see the IRCv3 sasl-3.1 extension.

The text used in the last param of this message varies wildly.

ERR_SASLALREADY (907) 
  "<client> :You have already authenticated using SASL"
This numeric indicates that SASL authentication failed because the client has already authenticated using SASL and reauthentication is not available or has been administratively disabled. For more information on this numeric, see the IRCv3 sasl-3.1 and sasl-3.2 extensions.

The text used in the last param of this message varies wildly.

RPL_SASLMECHS (908) 
  "<client> <mechanisms> :are available SASL mechanisms"
This numeric specifies the mechanisms supported for SASL authentication. <mechanisms> is a list of SASL mechanisms, delimited by a comma (',', 0x2C). For more information on this numeric, see the IRCv3 sasl-3.1 extension.

IRCv3.2 also specifies this information in the sasl client capability value. For more information on this, see the IRCv3 sasl-3.2 extension.

The text used in the last param of this message varies wildly.

RPL_ISUPPORT Parameters
Used to advertise features to clients, the RPL_ISUPPORT (005) numeric lists parameters that let the client know which features are active and their value, if any.

The parameters listed here are standardised and/or widely-advertised by IRC servers today and do not include deprecated parameters. Servers SHOULD support at least the following parameters where appropriate, and may advertise any others. For a more extensive list of parameters advertised by this numeric, see the irc-defs RPL_ISUPPORT list.

Certain parameters described here may not be standardised nor widely-advertised. These parameters are noted with the descriptor "Status: Proposed". However, we try to be conservative with the parameters we’re proposing, both in terms of having a small number of them and them being fairly understandable extensions to the current widely-used parameters.

If a ‘default value’ is listed for a parameter, this is the assumed value of the parameter until and unless it is advertised by the server. This is primarily to interoperate with servers that don’t advertise particular well-known and well-used parameters. If an ‘empty value’ is listed for a parameter, this is the assumed value of the parameter if it is advertised without a value.

AWAYLEN Parameter
  Format: AWAYLEN=<number>
The AWAYLEN parameter indicates the maximum length for the <reason> of an AWAY command. If an AWAY <reason> has more characters than this parameter, it may be silently truncated by the server before being passed on to other clients. Clients MAY receive an AWAY <reason> that has more characters than this parameter.

The value MUST be specified and MUST be a positive integer.

Examples:

  AWAYLEN=200

  AWAYLEN=307
CASEMAPPING Parameter
  Format: CASEMAPPING=<casemap>
The CASEMAPPING parameter indicates what method the server uses to compare equality of case-insensitive strings (such as channel names and nicks).

The value MUST be specified and MUST be a string representing the method that the server uses.

The specified casemappings are as follows:

ascii: Defines the characters a to z to be considered the lower-case equivalents of the characters A to Z only.
rfc1459: Same as 'ascii', with the addition of the characters '{', '}', '|', and '^' being considered the lower-case equivalents of the characters '[', ']', '\', and '~' respectively.
rfc1459-strict: Same casemapping as 'ascii', with the characters '{', '}', and '|' being the lower-case equivalents of '[', ']', and '\', respectively. Note that the difference between this and rfc1459 above is that in rfc1459-strict, '^' and '~' are not casefolded.
rfc7613: Proposed casemapping which defines a method based on PRECIS, allowing additional Unicode characters to be correctly casemapped [link].
The value MUST be specified and is a string. Servers MAY advertise alternate casemappings to those above, but clients MAY NOT be able to understand or perform them.

Servers SHOULD AVOID using the rfc1459 casemapping unless explicitly required for compatibility reasons or for linking with servers using it. The equivalency of the extra characters is not necessary nor useful today, and issues such as incorrect implementations and a conflict between matching masks exists.

Examples:

  CASEMAPPING=ascii

  CASEMAPPING=rfc1459
CHANLIMIT Parameter
  Format: CHANLIMIT=<prefixes>:[limit]{,<prefixes>:[limit]}
The CHANLIMIT parameter indicates the number of channels a client may join.

The value MUST be specified and is a list of "<prefixes>:<limit>" pairs, delimited by a comma (',', 0x2C). <prefixes> is a list of channel prefix characters as defined in the CHANTYPES parameter. <limit> is OPTIONAL and if specified is a positive integer indicating the maximum number of these types of channels a client may join. If there is no limit to the number of these channels a client may join, <limit> will not be specified.

Clients should not assume other clients are limited to what is specified in the CHANLIMIT parameter.

Examples:

  CHANLIMIT=#:25           ; indicates that clients may join 25 '#' channels

  CHANLIMIT=#&:50          ; indicates that clients may join 50 '#' and 50 '&' channels

  CHANLIMIT=#:70,&:        ; indicates that clients may join 70 '#' channels and any
                           number of '&' channels
CHANMODES Parameter
  Format: CHANMODES=A,B,C,D[,X,Y...]
The CHANMODES parameter specifies the channel modes available and which types of arguments they do or do not take when using them with the MODE command.

The value lists the channel mode letters of Type A, B, C, and D, respectively, delimited by a comma (',', 0x2C). The channel mode types are defined in the the MODE message description.

To allow for future extensions, a server MAY send additional types, delimited by a comma (',', 0x2C). However, server authors SHOULD NOT extend this parameter without good reason, and SHOULD CONSIDER whether their mode would work as one of the existing types instead. The behaviour of any additional types is undefined.

Server MUST NOT list modes in this parameter that are also advertised in the PREFIX parameter. However, modes within the PREFIX parameter may be treated as type B modes.

Examples:

  CHANMODES=b,k,l,imnpst

  CHANMODES=beI,k,l,BCMNORScimnpstz

  CHANMODES=beI,kfL,lj,psmntirRcOAQKVCuzNSMTGZ
CHANNELLEN Parameter
  Format: CHANNELLEN=<string>
The CHANNELLEN parameter specifies the maximum length of a channel name that a client may join. A client elsewhere on the network MAY join a channel with a larger name, but network administrators should take care to ensure this value stays consistent across the network.

The value MUST be specified and MUST be a positive integer.

Examples:

  CHANNELLEN=32

  CHANNELLEN=50

  CHANNELLEN=64
CHANTYPES Parameter
   Format: CHANTYPES=[string]
  Default: CHANTYPES=#
The CHANTYPES parameter indicates the channel prefix characters that are available on the current server. Common channel types are listed in the Channel Types section.

The value is OPTIONAL and if not specified indicates that no channel types are supported.

Examples:

  CHANTYPES=#

  CHANTYPES=&#

  CHANTYPES=#&
ELIST Parameter
  Format: ELIST=<string>
The ELIST parameter indicates that the server supports search extensions to the LIST command.

The value MUST be specified, and is a non-delimited list of letters, each of which denote an extension. The letters MUST be treated as being case-insensitive.

The following search extensions are defined:

C: Searching based on channel creation time, via the "C<val" and "C>val" modifiers to search for a channel that was created either less than val minutes ago, or more than val minutes ago, respectively
M: Searching based on a mask.
N: Searching based on a non-matching !mask. i.e., the opposite of M.
T: Searching based on topic set time, via the "T<val" and "T>val" modifiers to search for a topic time that was set less than val minutes ago, or more than val minutes ago, respectively.
U: Searching based on user count within the channel, via the "<val" and ">val" modifiers to search for a channel that has less or more than val users, respectively.
Examples:

  ELIST=MNUCT

  ELIST=MU

  ELIST=CMNTU
A widespread bug in existing implementations is to swap the semantics of "C<val" with "C>val", and/or "T<val" with "T>val", due to ambiguous legacy specifications. You should check the server you are using implements them as expected.

EXCEPTS Parameter
  Format: EXCEPTS=[character]
   Empty: e
The EXCEPTS parameter indicates that the server supports ban exceptions, as specified in the ban exception channel mode section.

The value is OPTIONAL and when not specified indicates that the letter "e" is used as the channel mode for ban exceptions. If the value is specified, the character indicates the letter which is used for ban exceptions.

Examples:

  EXCEPTS

  EXCEPTS=e
EXTBAN Parameter
  Format: EXTBAN=[<prefix>],<types>
The EXTBAN parameter indicates the types of “extended ban masks” that the server supports.

<prefix> denotes the character that indicates an extban to the server and <types> is a list of characters indicating the types of extended bans the server supports. If <prefix> does not exist then the server does not require a prefix for extbans, and they should be sent with no prefix.

Extbans may allow clients to issue bans based on account name, SSL certificate fingerprints and other attributes, based on what the server supports.

Extban masks SHOULD also be supported for the ban exception and invite exception modes.

Ensure that extban masks are actually typically supported in ban exception and invite exception modes.

We should include a list of 'typical' extban characters and their associated meaning, but make sure we specify that these are not standardised and may change based on server software. See also the irc-defs EXTBAN list.

Examples:

  EXTBAN=~,cqnr

  EXTBAN=~,qjncrRa

  EXTBAN=,ABCNOQRSTUcjmprsz
HOSTLEN Parameter
  Format: HOSTLEN=<number>
  Status: Proposed
The HOSTLEN parameter indicates the maximum length that a hostname may be on the server (whether cloaked, spoofed, or a looked-up domain name). Networks SHOULD be consistent with this value across different servers.

If a looked-up domain name is longer than this length, the server SHOULD opt to use the IP address instead, so that the hostname is underneath this length.

The value MUST be specified and MUST be a positive integer.

Examples:

  HOSTLEN=63
  HOSTLEN=64
INVEX Parameter
  Format: INVEX=[character]
   Empty: I
The INVEX parameter indicates that the server supports invite exceptions, as specified in the invite exception channel mode section.

The value is OPTIONAL and when not specified indicates that the letter "I" is used as the channel mode for invite exceptions. If the value is specified, the character indicates the letter which is used for invite exceptions.

Examples:

  INVEX

  INVEX=I
KICKLEN Parameter
  Format: KICKLEN=<length>
The KICKLEN parameter indicates the maximum length for the <reason> of a KICK command. If a KICK <reason> has more characters than this parameter, it may be silently truncated by the server before being passed on to other clients. Clients MAY receive a KICK <reason> that has more characters than this parameter.

The value MUST be specified and MUST be a positive integer.

Examples:

  KICKLEN=255

  KICKLEN=307
MAXLIST Parameter
  Format: MAXLIST=<modes>:<limit>{,<modes>:<limit>}
The MAXLIST parameter specifies how many “variable” modes of type A that have been defined in the CHANMODES parameter that a client may set in total on a channel.

The value MUST be specified and is a list of <modes>:<limit> pairs, delimited by a comma (',', 0x2C). <modes> is a list of type A modes defined in CHANMODES. <limit> is a positive integer specifying the maximum number of entries that all of the modes in <modes>, combined, may set on a channel.

A client MUST NOT make any assumptions on how many mode entries may actually exist on any given channel. This limit only applies to the client setting new modes of the given types, and other clients may have different limits.

Examples:

  MAXLIST=beI:25           ; indicates that a client may set up to a total of 25 of a
                           combination of "b", "e", and "I" modes.

  MAXLIST=b:60,e:60,I:60   ; indicates that a client may set up to 60 "b" modes,
                           "e" modes, and 60 "I" modes.

  MAXLIST=beI:100,q:50     ; indicates that a client may set up to a total of 100 of
                           a combination of "b", "e", and "I" modes, and that they
                           may set up to 50 "q" modes.
MAXTARGETS Parameter
  Format: MAXTARGETS=[number]
The MAXTARGETS parameter specifies the maximum number of targets a PRIVMSG or NOTICE command may have, and may apply to other commands based on server software.

The value is OPTIONAL and if specified, [number] is a positive integer representing the maximum number of targets those commands may have. If there is no limit, then [number] MAY not be specified.

The TARGMAX parameter SHOULD be advertised instead of or in addition to this parameter. TARGMAX is intended to replace MAXTARGETS as that parameter is more clear about which commands limits apply to.

Examples:

  MAXTARGETS=4

  MAXTARGETS=20
MODES Parameter
  Format: MODES=[number]
The MODES parameter specifies how many ‘variable’ modes may be set on a channel by a single MODE command from a client. A ‘variable’ mode is defined as being a type A, B or C mode as defined in the CHANMODES parameter, or in the channel modes specified in the PREFIX parameter.

A client SHOULD NOT issue more ‘variable’ modes than this in a single MODE command. A server MAY however issue more ‘variable’ modes than this in a single MODE message. The value is OPTIONAL and when not specified indicates that there is no limit to the number of ‘variable’ modes that may be set in a single client MODE command.

If the value is specified, it MUST be a positive integer.

Examples:

  MODES=4

  MODES=12

  MODES=20
NETWORK Parameter
  Format: NETWORK=<string>
The NETWORK parameter indicates the name of the IRC network that the client is connected to. This parameter is advertised for INFORMATIONAL PURPOSES ONLY. Clients SHOULD NOT use this value to make assumptions about supported features on the server as networks may change server software and configuration at any time.

Examples:

  NETWORK=EFNet

  NETWORK=Rizon

  NETWORK=Example\x20Network
NICKLEN Parameter
   Format: NICKLEN=<number>
The NICKLEN parameter indicates the maximum length of a nickname that a client may set. Clients on the network MAY have longer nicks than this.

The value MUST be specified and MUST be a positive integer. 30 or 31 are typical values for this parameter advertised by servers today.

Examples:

  NICKLEN=9

  NICKLEN=30

  NICKLEN=31
PREFIX Parameter
   Format: PREFIX=[(modes)prefixes]
  Default: PREFIX=(ov)@+
Within channels, clients can have different statuses, denoted by single-character prefixes. The PREFIX parameter specifies these prefixes and the channel mode characters that they are mapped to. There is a one-to-one mapping between prefixes and channel modes. The prefixes in this parameter are in descending order, from the prefix that gives the most privileges to the prefix that gives the least.

The typical prefixes advertised in this parameter are listed in the Channel Membership Prefixes section.

The value is OPTIONAL and when it is not specified indicates that no prefixes are supported.

Examples:

  PREFIX=(ov)@+

  PREFIX=(ohv)@%+

  PREFIX=(qaohv)~&@%+
SAFELIST Parameter
  Format: SAFELIST
If SAFELIST parameter is advertised, the server ensures that a client may perform the LIST command without being disconnected due to the large volume of data the LIST command generates.

The SAFELIST parameter MUST NOT be specified with a value.

Examples:

  SAFELIST
SILENCE Parameter
  Format: SILENCE[=<limit>]
The SILENCE parameter indicates the maximum number of entries a client can have in their silence list.

The value is OPTIONAL and if specified is a positive integer. If the value is not specified, the server does not support the SILENCE command.

Most IRC clients also include client-side filter/ignore lists as an alternative to this command.

Examples:

  SILENCE

  SILENCE=15

  SILENCE=32
STATUSMSG Parameter
  Format: STATUSMSG=<string>
The STATUSMSG parameter indicates that the server supports a method for clients to send a message via the PRIVMSG / NOTICE commands to those people on a channel with (one of) the specified channel membership prefixes.

The value MUST be specified and MUST be a list of prefixes as specified in the PREFIX parameter. Most servers today advertise every prefix in their PREFIX parameter in STATUSMSG.

Examples:

  STATUSMSG=@+

  STATUSMSG=@%+

  STATUSMSG=~&@%+
TARGMAX Parameter
  Format: TARGMAX=[<command>:[limit]{,<command>:[limit]}]
Certain client commands MAY contain multiple targets, delimited by a comma (',', 0x2C). The TARGMAX parameter defines the maximum number of targets allowed for commands which accept multiple targets. If this parameter is not advertised or a value is not sent then a client SHOULD assume that no commands except the JOIN and PART commands accept multiple parameters.

The value is OPTIONAL and is a set of <command>:<limit> pairs, delimited by a comma (',', 0x2C). <command> is the name of a client command. <limit> is the maximum number of targets which that command accepts. If <limit> is specified, it is a positive integer. If <limit> is not specified, then there is no maximum number of targets for that command. Clients MUST treat <command> as case-insensitive.

Examples:

  TARGMAX=PRIVMSG:3,WHOIS:1,JOIN:

  TARGMAX=NAMES:1,LIST:1,KICK:1,WHOIS:1,PRIVMSG:4,NOTICE:4,ACCEPT:,MONITOR:

  TARGMAX=ACCEPT:,KICK:1,LIST:1,NAMES:1,NOTICE:4,PRIVMSG:4,WHOIS:1
TOPICLEN Parameter
  Format: TOPICLEN=<number>
The TOPICLEN parameter indicates the maximum length of a topic that a client may set on a channel. Channels on the network MAY have topics with longer lengths than this.

The value MUST be specified and MUST be a positive integer. 307 is the typical value for this parameter advertised by servers today.

Examples:

  TOPICLEN=307

  TOPICLEN=390
USERLEN Parameter
  Format: USERLEN=<number>
  Status: Proposed
The USERLEN parameter indicates the maximum length that a username may be on the server. Networks SHOULD be consistent with this value across different servers. As noted in the USER message, the tilde prefix ("~"), if it exists, contributes to the length of the username and would be included in this parameter.

The value MUST be specified and MUST be a positive integer.

Examples:

  USERLEN=12
  USERLEN=18
Current Architectural Problems
There are a number of recognized problems with the IRC protocol. This section only addresses the problems related to the architecture of the protocol.

Scalability
It is widely recognized that this protocol may not scale sufficiently well when used in a very large arena. The main problem comes from the requirement that all servers know about all other servers, clients, and channels, and that information regarding them be updated as soon as it changes.

Server-to-server protocols can attempt to alleviate this by, for example, only sending ‘necessary’ state information to leaf servers. These sort of optimisations are implementation-specific and are not covered in this document. However, server authors should take great care in their protocols to ensure race conditions and other network instability does not result from these attempts to improve the scalability of their protocol.

Reliability
As the only network configuration used for IRC servers is that of a spanning tree, each link between two servers is an obvious and serious point of failure.

Software authors are and have been experimenting with alternative topologies such as mesh networks. However, there is not yet a production implementation or specification of any topology other than spanning-tree.

Implementation Notes
The IRC protocol is reasonably complex. When writing software that interacts with it, there are certain choices that are implementation-defined, as well as certain areas that are commonly incorrectly implemented.

This section raises discussion, questions, and recommendations intended to help implementors. In particular, the advice/discussion here may be sloppy compared to the above, and the questions may be less well-defined or without strict answers, but regardless should help you when writing software that interacts with the IRC protocol.

Character Encodings
Character encodings in IRC are hard. UTF-8 is recommended, the mess of Latin-1/ISO-8859-1(5)/CP1252 also seems common, but all sorts of other encodings are also used in practice. Particularly on networks that support other languages, and were created before UTF-8 became as widespread as it has.

When sending, we always recommend UTF-8. When decoding, we generally recommend trying UTF-8 and falling back to Latin-1 (what has been called the Hybrid encoding).

For clients, this is fine. Even if they incorrectly decode a private message, the user should see that the message has been decoded incorrectly and be able to resolve the issue (hopefully by telling the sending user to use UTF-8).

However, servers are in a trickier position (especially for PRIVMSG/NOTICE or any other command that takes arbitrary user input such as USER, TOPIC, etc). Servers should simply treat this input from the user as a character array they accept and then spit out again, no trouble.

Servers implemented in languages with first-class Unicode strings may wish to treat IRC lines and messages as Unicode text internally. For servers to treat messages in this way, they need to decode lines as they’re received and later encode the lines before they’re sent out.

This presents an issue. What if the line from the user is decoded incorrectly, modified (eg. by casefolding), and then sent out? (see also: Mojibake). What these servers may instead do is either:

follow the lead of the majority of existing servers and treat these parameters as byte arrays not to be parsed or decoded in any way.
attempt to decode all incoming lines as UTF-8 (possibly using Hybrid encoding like clients do) and if the line cannot be decoded it is ignored or returns an error. The IRCv3 UTF8ONLY specification allows them to signal this to clients.
The former ensures all messages are sent correctly, and the latter simplifies server implementations and allows clients to disable decoding heuristics.

Message Parsing and Assembly
Message parsing/assembly is one area where implementations can differ wildly, and is a common vector for both security issues and general runtime problems.

Message Parsing is turning raw IRC messages into the various message parts (tags, prefix, command, parameters). Message Assembly is the opposite – taking the various message parts and creating an IRC line to be sent over the wire.

Implementors should ensure that their message parsing and assembly responds in expected ways, by running their software through test cases. I recommend these public-domain irc-parser-tests, which are reasonably extensive.

Trailing
Trailing is a completely normal parameter, except for the fact that it can contain spaces. When parsing messages, the ‘normal params’ and trailing should be appended and returned as a single list containing all the message params.

This is an example of an incorrect parser, that specifically separates normal params and trailing. When returning messages after parsing, don’t return a struct/object containing these variables:

  Message
      .Tags
      .Source
      .Verb
      .Params (containing all but the trailing param)
      .Trailing (containing just the trailing param)
Trailing is a normal parameter. Separating the parameter types in this way will cause many breakages and weird issues, as logic code will depend on the final param being in either .Params or .Trailing, when the simple fact is that it can be in either. Make sure that your message parser instead outputs parsed messages more like this:

  Message
      .Tags
      .Source
      .Verb
      .Params (including all normal params, and the trailing param if it exists)
This will make sure that you don’t run into silly trailing parameter errors.

Direct String Comparisons on IRC Lines
Some software decides that the best way to process incoming lines is with something along the lines of this:

  Line = NewIRCLineFromSocket()
  If Line.StartsWith("PART") {
      Part(...etc...)
  } Else If Line.StartsWith("QUIT") {
      Quit(...etc...)
  }
This is bad. This will break. Here’s why: Any IRC message can choose to include or not include the source.

If you directly compare the beginning of lines like this, then you will break when servers decide to start including sources on messages (for example, some newer IRCds decide to include the source on all messages that they output). This results in clients that don’t correctly parse incoming messages and break as a result.

Instead, you should make sure that you send incoming lines through a message parser, and then do things based on what’s output by that parser. For instance:

  Message = IRCMessageParser(Line)
  If Message.Verb == "PART" {
        Part(...etc...)
  } Else If Message.Verb == "QUIT" {
        Quit(...etc...)
  }
This will ensure that your software doesn’t break when clients or servers send extra, or omit unnecessary, message elements.

Something to keep in mind is that the message verb is always case insensitive, so you should casemap it appropriately before doing comparisons similar to the above. In my own IRC libraries, I convert the verb to uppercase before returning the message.
*/