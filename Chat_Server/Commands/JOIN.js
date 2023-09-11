const { Numerics } = require("../numerics");
const { CRLF } = require("../constants");
const _logger = require('pino')();
const logger = _logger.child({ Service: 'Chat Server', Command: "JOIN" });

/**
 * JOIN
 * 
 * Parameters: <channel>{,<channel>} [<key>{,<key>}]
 * Alt Params: 0
 * 
 * @param {*} channel 
 * 
 * Description:
 *  The JOIN command indicates that the client wants to join the given channel(s), each channel using the given key for it.
 *  The server receiving the command checks whether or not the client can join the given channel, and processes the request. 
 *  Servers MUST process the parameters of this command as lists on incoming commands from clients, with the first <key> being used for the first <channel>,
 *  the second <key> being used for the second <channel>, etc. 
 * 
 *  While a client is joined to a channel, they receive all relevant information about that channel including the JOIN, PART, KICK, and MODE messages affecting the channel. 
 *  They receive all PRIVMSG and NOTICE messages sent to the channel, and they also receive QUIT messages from other clients joined to the same channel (to let them know those users have left the channel and the network). 
 *  This allows them to keep track of other channel members and channel modes.
 */
const JOIN = (socket, client, channels, parameters) => {
    // const JOIN = (channels = [], keys = []) => {
        // send server a command to join
        // whether or not the client successfully joins that channel depends on the modes currently set on the channel
        // For example, if the channel is set to invite-only mode (+i), the client only joins the channel if they have been invited by another user or they have been exempted from requiring an invite by the channel operators.
    
        //     Command Examples:
    
    //   JOIN #foobar                    ; join channel #foobar.
    
    //   JOIN &foo fubar                 ; join channel &foo using key "fubar".
    
    //   JOIN #foo,&bar fubar            ; join channel #foo using key "fubar"
    //                                   and &bar using no key.
    
    //   JOIN #foo,#bar fubar,foobar     ; join channel #foo using key "fubar".
    //                                   and channel #bar using key "foobar".
    
    //   JOIN #foo,#bar                  ; join channels #foo and #bar.
    // Message Examples:
    
    //   :WiZ JOIN #Twilight_zone        ; WiZ is joining the channel
    //                                   #Twilight_zone
    
    //   :dan-!d@localhost JOIN #test    ; dan- is joining the channel #test
        
        logger.info(`JOIN client: '${client}', parameters: '${parameters}'`);
        if (!parameters || parameters === "" || parameters.length === 0) {
            // socket.write(`${client} JOIN ${ERR_NEEDMOREPARAMS} :Not enough parameters\r\n`);
            // return {"res": false};
            return {"err": Numerics["ERR_NEEDMOREPARAMS"]("JOIN")};
        }
    
        let chan;
        if (parameters.includes(" ")) {
            // split params to get channel and keys
            const paramsSplit = parameters.split(" ");
            chan = paramsSplit[0];
        } else {
            chan = parameters;
        }
        chan = String(chan);
        logger.info(`CHAN = ${chan}, TYPE OF CHAN = ${typeof(chan)}`);
        if (chan === "0") {
            // special argument, as in, JOIN 0, or join nothing
            // server makes client leave all channels they are connected to.
            // Server will process as though the client had sent a PART command for each channel they are a member of.
        }
        const chanName = chan.substring(1, chan.length);
        const chanMask = chan[0];
    
            // If true, it is a valid channel
    
            // Regular Channel #
            // Clients can join this channel, and the first client who joins a normal channel is made a channel operator, 
            // along with the appropriate channel membership prefix. 
            // On most servers, newly-created channels have then protected topic "+t" 
            // and no external messages "+n" modes enabled, but exactly what modes new channels are given is up to the server.
            // Regular channels are persisted across the network. 
            // If two clients on different servers join the same regular channel, 
            // they’ll be able to see that each other are joined, and will see messages sent to the channel by the other client.
            // On servers that support the concept of ‘channel ownership’ (a client being able to own a channel and retain control of it with their account), 
            // clients may not receive channel operator priveledges on joining an otherwise empty channel.
            
            // Local Channel &
            // Clients can join this channel as normal, and the first client who joins a normal channel is made a channel operator, 
            // but the channel is not persisted across the network. 
            // In other words, each server has its own set of local channels that the other servers on the network don’t see.
            // If a client on server A and a client on server B join the channel &info, they will not be able to see each other or the messages each posts to their server’s local channel &info. 
            // However, if a client on server A and another client on server A join the channel &info, they will be able to see each other and the messages the other posts to that local channel.
            // Generally, the concept of channel ownership is not supported for local channels. 
            // Local channels also aren’t as widely available as regular channels. 
            // As well, some networks disable or disallow local channels as opers across the network can’t see nor administrate them.
        logger.info(`ChanMask = '${chanMask}', Type = ${typeof(chanMask)}`);
        if (chanMask !== "#" && chanMask !== "&") {
            // not a valid channel type, fail
            logger.info(`Channel Mask = ${chanMask}`);
            // socket.write(`${chan} ${ERR_BADCHANMASK} :Bad Channel Mask\r\n`);
            return {"err": Numerics["ERR_BADCHANMASK"](chan)};
        }
    
        const chanType = chanMask === "#" ? "Regular" : "Local";
    
        if (!(chanName in channels[chanType])) {
            // socket.write(`${client} ${chanName} ${ERR_NOSUCHCHANNEL} :No such channel\r\n`);
            // return {"res": false};
            return {"err": Numerics["ERR_NOSUCHCHANNEL"](chanName)};
        }
    
        const channel = channels[chanType][chanName];
        if (MODE_CHAN_CLIENT_LIMIT in channel.modes && Object.keys(channel.clients).length == channel.modes[MODE_CHAN_CLIENT_LIMIT]) {
            // Adding this user will make the channel go over its limit.
            // socket.write(Numerics["ERR_CHANNELISFULL"](client, chanName) + CRLF);
            // return {"res": false};
            return {"err": Numerics["ERR_CHANNELISFULL"](chanName)};
        }
    
        channel.AddClient(
            client.nick, 
            {
                "clientObj": Client(client.nick),
                "prefixes": ["~"] // temporary mode add
            });
    
        // If a client’s JOIN command to the server is successful, the server MUST send, in this order:
        // 1. A JOIN message with the client as the message <source> and the channel they have joined as the first parameter of the message.
        // 2. The channel’s topic (with RPL_TOPIC (332) and optionally RPL_TOPICWHOTIME (333)), and no message if the channel does not have a topic.
        // 3. A list of users currently joined to the channel (with one or more RPL_NAMREPLY (353) numerics followed by a single RPL_ENDOFNAMES (366) numeric). 
        // These RPL_NAMREPLY messages sent by the server MUST include the requesting client that has just joined the channel.
        
        // 1
        const resp1 = `${client.source} JOIN ${chan}` + CRLF;
        // 2
        // "<client> <channel> :<topic>"
        socket.write(resp1);
        const resp2 = client.source + " " + Numerics["RPL_TOPIC"](chan, channel.topic["topic"]) + CRLF;
        const respt2_optional = client.source + " " + Numerics["RPL_TOPICWHOTIME"](chan, channel.topic["nick"], channel.topic["timestamp"]) + CRLF;
        socket.write(resp2);
        socket.write(respt2_optional);
        // 3
        let membershipPrefixMapping = "";
        for (const client of Object.keys(channel.clients)) {
            if (client.length === 0 || typeof(client) !== String)
                continue;
            logger.info(`Client = '${client}'`);
            logger.info(`Props = '${Object.keys(channel["clients"][client])}'`);
            const prefixes = channel["clients"][client]["prefixes"];
            for (const prefix of prefixes) {
                membershipPrefixMapping += `${prefix}${client} `
            } 
        }
        // TODO
        // If server write request is TOO big, it MUST be broken up over multiple RPL_NAMREPLY commands, then ending with RPL_ENDOFNAMES;
        // Need to chunk data to write back to client
        // see limits.js
        socket.write(client.source + " " + Numerics["RPL_NAMREPLY"](channel.status, chanName, membershipPrefixMapping) + CRLF);
        socket.write(client.source + " " + Numerics["RPL_ENDOFNAMES"](chanName) + CRLF);
    
        return {
            "command": "JOIN",
            "channelType": chanType,
            "channelName": chanName,
            "client": client
        };
};

module.exports = { JOIN };