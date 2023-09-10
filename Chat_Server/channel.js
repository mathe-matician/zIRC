const { UpdateOne, FindOneAndUpdate, FindOne } = require("./db");
const { Numerics } = require("./numerics");
const { CRLF } = require("./constants");
require('dotenv').config();

/*
A channel is a named group of one or more clients. All clients in the channel will receive all messages addressed to that channel. The channel is created implicitly when the first client joins it, and the channel ceases to exist when the last client leaves it. While the channel exists, any client can reference the channel using the name of the channel. Networks that support the concept of ‘channel ownership’ may persist specific channels in some way while no clients are connected to them.

Channel names are strings (beginning with specified prefix characters). Apart from the requirement of the first character being a valid channel type prefix character; the only restriction on a channel name is that it may not contain any spaces (' ', 0x20), a control G / BELL ('^G', 0x07), or a comma (',', 0x2C) (which is used as a list item separator by the protocol).

There are several types of channels used in the IRC protocol. The first standard type of channel is a regular channel, which is known to all servers that are connected to the network. The prefix character for this type of channel is ('#', 0x23). The second type are server-specific or local channels, where the clients connected can only see and talk to other clients on the same server. The prefix character for this type of channel is ('&', 0x26). Other types of channels are described in the Channel Types section.

Along with various channel types, there are also channel modes that can alter the characteristics and behaviour of individual channels. See the Channel Modes section for more information on these.

To create a new channel or become part of an existing channel, a user is required to join the channel using the JOIN command. If the channel doesn’t exist prior to joining, the channel is created and the creating user becomes a channel operator. If the channel already exists, whether or not the client successfully joins that channel depends on the modes currently set on the channel. For example, if the channel is set to invite-only mode (+i), the client only joins the channel if they have been invited by another user or they have been exempted from requiring an invite by the channel operators.

Channels also contain a topic. The topic is a line shown to all users when they join the channel, and all users in the channel are notified when the topic of a channel is changed. Channel topics commonly state channel rules, links, quotes from channel members, a general description of the channel, or whatever the channel operators want to share with the clients in their channel.

A user may be joined to several channels at once, but a limit may be imposed by the server as to how many channels a client can be in at one time. This limit is specified by the CHANLIMIT RPL_ISUPPORT parameter. See the Feature Advertisement section for more details on RPL_ISUPPORT.

If the IRC network becomes disjoint because of a split between servers, the channel on either side is composed of only those clients which are connected to servers on the respective sides of the split, possibly ceasing to exist on one side. When the split is healed, the connecting servers ensure the network state is consistent between them.
*/

const Channel = (
    channelName = "", 
    _modes = {
        "+l": true, // Client Limit Channel Mode. Int
        "+s": true, // ("@", 0x40) - Secret channel (secret channel mode "+s").
        "+p": true, // ("*", 0x2A) - Private channel (was "+p", no longer widely used today).
    }, // e.g. "+t" = protected topic, and "+n" = no external messages
    _status = "=", // ("=") - default public channel
    _operators = {}, // first client that joins becomes a channel operator
    type = "#", // by default create a regular channel type
    _topic = {
        name: "Welcome to zIRC!",
        timestamp: new Date().toISOString(),
        nick: "default"
    },
    _clients = null,
    ) => {

    // TODO
    // clean user input
    // probably just clear whitespace from ends
    // and null characters
    const modes = _modes;
    const operators = _operators;
    const topic = _topic;
    const status = _status;
    const clients = _clients;
    
    const verifyChannelName = () => {
        // channel names cannot contain spaces, ',', or BELL character ^G.
        const invalidChannelName = /[\s,]|\\a|\\0/g;
        if (channelName === undefined || channelName.length === 0) {
            console.log("Channel name must be at least 1 character in length.");
            throw new Error("Channel name must be at least 1 character in length.");
        } else if (channelName.length > 64) {
            // Idk if this is true or not
            console.log("Channel name cannot be longer than 64 characters.");
            throw new Error("Channel name cannot be longer than 64 characters.");
        }

        let err = "";
        while ((res = invalidChannelName.exec(channelName)) !== null) {
            if (res.index === 0) {
                throw new Error("Invalid channel type");
            } else {
                err += `Channel name cannot contain "${res[0]}" characters.\n`;
            }
        }

        if (err) {
            throw new Error(err);
        }

        console.log(`${channelName} is a valid Channel name.`);
        const _name = type + channelName;
        console.log(`Creating Channel ${_name}`);
        return _name;
    }

    const verifyModes = () => {
        // if (modes) {
        //     for (let i = 0; i < modes.length; i++) {
        //         if (!_modes.has(modes[i]))
        //             _modes[modes[i]] =  "";
        //     }
        // }
    }

    const verifyOperators = () => {
        // TODO
        // loop through modes and check if _modes has it already
        // _modes.set(modes)
        if (operators) {

        }
    }

    const AddOperator = (client = null) => {
        if (!_operators.has(client.name)) {
            _operators[client.name] = "";
        }
        
    }

    const AddMode = (mode = "") => {
        if (!modes.has(mode)) {
            modes[mode] = true;
        }
    }

    const ChangeTopic = (newTopic = null) => {
        // TODO
        // Auth check
        if (!newTopic) {
            return;
        }
        topic = newTopic;
    }

    const AddClient = (clientName, props) => {
        clients[clientName] = props;
    }

    try {
        const name = verifyChannelName();
        // TODO
        // insert into db
        
        // verifyModes();
        // verifyOperators();
        console.log("Done verifying channle name");
        return { 
            name,
            topic,
            status,
            modes,
            operators,
            clients,
            AddClient,
        };
    } catch (error) {
        console.log(error);
    }
}

module.exports = Channel;