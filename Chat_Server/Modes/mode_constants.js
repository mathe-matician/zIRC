const MODE_PREFIX = "+";

/**
 * User Modes
 */
const MODE_USER_INVISIBLE = "i"; // If a user is set to ‘invisible’, they will not show up in commands such as WHO or NAMES unless they share a channel with the user that submitted the command. In addition, some servers hide all channels from the WHOIS reply of an invisible user they do not share with the user that submitted the command.
const MODE_USER_OPER = "o"; // If a user has this mode, this indicates that they are a network operator.
const MODE_USER_LOCAL_OPER = "O"; // If a user has this mode, this indicates that they are a server operator. A local operator has operator privileges for their server, and not for the rest of the network.
const MODE_USER_REGISTERED = "r"; // If a user has this mode, this indicates that they have logged into a user account. IRCv3 extensions such as account-notify, account-tag, and extended-join provide the account name of logged-in users, and are more accurate than trying to detect this user mode due to the capability name remaining consistent.
const MODE_USER_WALLOPS = "w"; // If a user has this mode, this indicates that they will receive WALLOPS messages from the server.

/**
 * Channel Modes
 */

/**
 * MODE_CHAN_BAN
 * This channel mode controls a list of client masks that are ‘banned’ from joining or speaking in the channel. If this mode has values, each of these values should be a client mask.
 * If this mode is set on a channel, and a client sends a JOIN request for this channel, 
 * their nickmask (the combination of nick!user@host) is compared with each banned client mask set with this mode. 
 * If they match one of these banned masks, they will receive an ERR_BANNEDFROMCHAN (474) reply and the JOIN command will fail. 
 * See the ban exception mode for more details.
 */
const MODE_CHAN_BAN = "b"; 

/**
 * MODE_CHAN_EXCEPTION
 * The standard mode letter used for it is "+e", but it SHOULD be defined in the EXCEPTS RPL_ISUPPORT parameter on connection.
 * This channel mode controls a list of client masks that are exempt from the ‘ban’ channel mode. 
 * If this mode has values, each of these values should be a client mask. 
 * If this mode is set on a channel, and a client sends a JOIN request for this channel, 
 * their nickmask is compared with each ‘exempted’ client mask. If their nickmask matches any one of the masks set by this mode, 
 * and their nickmask also matches any one of the masks set by the ban channel mode, they will not be blocked from joining due to the ban mode.
 */
const MODE_CHAN_EXCEPTION = MODE_PREFIX + "e";

/**
 * This channel mode controls whether new users may join based on the number of users who already exist in the channel. 
 * If this mode is set, its value is an integer and defines the limit of how many clients may be joined to the channel.
 * If this mode is set on a channel, and the number of users joined to that channel matches or exceeds the value of this mode, 
 * new users cannot join that channel. If a client sends a JOIN request for this channel, they will receive an ERR_CHANNELISFULL (471) reply and the command will fail.
 */
const MODE_CHAN_CLIENT_LIMIT = MODE_PREFIX + "l";

/**
 * MODE_CHAN_INVITE_ONLY
 * This channel mode controls whether new users need to be invited to the channel before being able to join.
 * If this mode is set on a channel, a user must have received an INVITE for this channel before being allowed to join it. 
 * If they have not received an invite, they will receive an ERR_INVITEONLYCHAN (473) reply and the command will fail.
 */
const MODE_CHAN_INVITE_ONLY = MODE_PREFIX + "i";

/**
 * MODE_CHAN_INVITE_EXCEPTION
 * The standard mode letter used for it is "+I", but it SHOULD be defined in the INVEX RPL_ISUPPORT parameter on connection.
 * This channel mode controls a list of channel masks that are exempt from the invite-only channel mode. 
 * If this mode has values, each of these values should be a client mask.
 * If this mode is set on a channel, and a client sends a JOIN request for that channel, 
 * their nickmask is compared with each ‘exempted’ client mask. 
 * If their nickmask matches any one of the masks set by this mode, and the channel is in invite-only mode, 
 * they do not need to require an INVITE in order to join the channel.
 */
const MODE_CHAN_INVITE_EXCEPTION = MODE_PREFIX + "I";

/**
 * MODE_CHAN_KEY
 * This mode letter sets a ‘key’ that must be supplied in order to join this channel. If this mode is set, its’ value is the key that is required. Servers may validate the value (eg. to forbid spaces, as they make it harder to use the key in JOIN messages). If the value is invalid, they SHOULD return ERR_INVALIDMODEPARAM. However, clients MUST be able to handle any of the following:

ERR_INVALIDMODEPARAM
ERR_INVALIDKEY
MODE echoed with a different key (eg. truncated or stripped of invalid characters)
the key changed ignored, and no MODE echoed if no other mode change was valid.
If this mode is set on a channel, and a client sends a JOIN request for that channel, they must supply <key> in order for the command to succeed. If they do not supply a <key>, or the key they supply does not match the value of this mode, they will receive an ERR_BADCHANNELKEY (475) reply and the command will fail.
 */
const MODE_CHAN_KEY = MODE_PREFIX + "k";

/**
 * MODE_CHAN_MODERATED
 * This channel mode controls whether users may freely talk on the channel, and does not have any value.

If this mode is set on a channel, only users who have channel privileges may send messages to that channel. The voice channel mode is designed to let a user talk in a moderated channel without giving them other channel moderation abilities, and users of higher privileges (such as halfops or chanops) may also speak in moderated channels.
 */
const MODE_CHAN_MODERATED = "m";

/**
 * MODE_CHAN_SECRET
 * This channel mode controls whether the channel is ‘secret’, and does not have any value.

A channel that is set to secret will not show up in responses to the LIST or NAMES command unless the client sending the command is joined to the channel. Likewise, secret channels will not show up in the RPL_WHOISCHANNELS (319) numeric unless the user the numeric is being sent to is joined to that channel.
 */
const MODE_CHAN_SECRET = "s";

/**
 * MODE_CHAN_PROTECTED_TOPIC
 * This channel mode controls whether channel privileges are required to set the topic, and does not have any value.

If this mode is enabled, users must have channel privileges such as halfop or operator status in order to change the topic of a channel. In a channel that does not have this mode enabled, anyone may set the topic of the channel using the TOPIC command.
 */
const MODE_CHAN_PROTECTED_TOPIC = "t";

/**
 * MODE_CHAN_NO_EXTERNAL_MESSAGES
 * This channel mode controls whether users who are not joined to the channel can send messages to it, and does not have any value.

If this mode is enabled, users MUST be joined to the channel in order to send private messages and notices to the channel. If this mode is enabled and they try to send one of these to a channel they are not joined to, they will receive an ERR_CANNOTSENDTOCHAN (404) numeric and the message will not be sent to that channel.
 */
const MODE_CHAN_NO_EXTERNAL_MESSAGES = "n";

const MODE_CHAN_FOUNDER = "q"; // e.g. ~+q // This prefix shows that the given user is the ‘founder’ of the current channel and has full moderation control over it – ie, they are considered to ‘own’ that channel by the network. This prefix is typically only used on networks that have the concept of client accounts, and ownership of channels by those accounts.
const MODE_CHAN_PROTECTED = "a"; // e.g. &+a // Users with this mode cannot be kicked and cannot have this mode removed by other protected users. In some software, they may perform actions that operators can, but at a higher privilege level than operators. This prefix is typically only used on networks that have the concept of client accounts, and ownership of channels by those accounts.
const MODE_CHAN_OPERATOR = "o"; // e.g. @+o // Users with this mode may perform channel moderation tasks such as kicking users, applying channel modes, and set other users to operator (or lower) status.
const MODE_CHAN_HALFOP = "h"; // e.g. %+h // Users with this mode may perform channel moderation tasks, but at a lower privilege level than operators. Which channel moderation tasks they can and cannot perform varies with server software and configuration.
const MODE_CHAN_VOICE = "v"; // e.g. ++v // Users with this mode may send messages to a channel that is moderated.

/**
 * Membership Prefixes
 * 
 * Users joined to a channel may get certain privileges or status in that channel based on channel modes given to them. 
 * These users are given prefixes before their nickname whenever it is associated with a channel (ie, in NAMES, WHO and WHOIS messages). 
 * The standard and common prefixes are listed here, and MUST be advertised by the server in the PREFIX RPL_ISUPPORT parameter on connection.
 */
const MEM_PREFIX_CHAN_FOUNDER = "~"; // This prefix shows that the given user is the ‘founder’ of the current channel and has full moderation control over it – ie, they are considered to ‘own’ that channel by the network. This prefix is typically only used on networks that have the concept of client accounts, and ownership of channels by those accounts.
const MEM_PREFIX_CHAN_PROTECTED = "&"; // Users with this mode cannot be kicked and cannot have this mode removed by other protected users. In some software, they may perform actions that operators can, but at a higher privilege level than operators. This prefix is typically only used on networks that have the concept of client accounts, and ownership of channels by those accounts.
const MEM_PREFIX_CHAN_OPERATOR = "@"; // Users with this mode may perform channel moderation tasks such as kicking users, applying channel modes, and set other users to operator (or lower) status.
const MEM_PREFIX_CHAN_HALFOP = "%"; // Users with this mode may perform channel moderation tasks, but at a lower privilege level than operators. Which channel moderation tasks they can and cannot perform varies with server software and configuration.
const MEM_PREFIX_CHAN_VOICE = "+"; // Users with this mode may send messages to a channel that is moderated.