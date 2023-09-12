const { InsertOne, CollectionExists, FindOne } = require("../db");
const { Numerics } = require("../numerics");
const { SCHEMA_ChatMsg } = require("../db/schema/chatmsg.schema");
require('dotenv').config();

const _logger = require('pino')();
const logger = _logger.child({ Service: 'Chat Server', Command: "PRIVMSG" });

/**
 * PRIVMSG
 * 
 * Parameters: <target>{,<target>} <text to be sent>
 * 
 * Description: The PRIVMSG command is used to send private messages between users, as well as to send messages to channels. 
 *              <target> is the nickname of a client or the name of a channel.
 */

/**
 * 
 * @param {array} params array containing target and message
 * @param {*} clients 
 * @param {*} clientSocket 
 * @returns 
 */
const PRIVMSG = async (params, clients, clientSocket) => {
    logger.info(`PRIVMSG start:\nparams: ${params}`);
    if (params.length !== 2) {
        return {"err": Numerics["ERR_NEEDMOREPARAMS"]("PRIVMSG")};
    }
    // TODO
    // if we need to forward this message onto another server we need to account for that.
    // return a special value to this server to then search for client on other servers?
    // TODO - the above would be if the client sent the server which they wanted to talk to to this server... i think
    try {
        // this command will always contain a ':' which is the prefix to the message
        // PRIVMSG <user/chan/etc> :<msg>
        // e.g. PRIVMSG Marley :Hey Mar, what's up?
        const message = params[1].substring(1, params[1].length);
        const target = params[0].trim();
        const chanTypeIndex = target.indexOf("#");
        if (chanTypeIndex !== -1) {
            // then it is a channel
            // parse any ops before the channel
            const ops = target.substring(0, chanTypeIndex);
            const targetName = target.substring(chanTypeIndex, target.length);
            logger.info(`PRIVMSG: searching for collection ${targetName}`);
            const colExists = await CollectionExists(targetName);
            if (!colExists) {
                logger.info(`Collection '${targetName}' does not exist`);
                return {"err": Numerics["ERR_NOSUCHCHANNEL"](targetName)};
            }
            // insert message
            const insertRes = await InsertOne(
                SCHEMA_ChatMsg(
                "userid",
                "server",
                message,
                "tags"
            ), 
            targetName, 
            process.env.MONGODB_CHAT_MESSAGE_DB_NAME);
            // TODO
            // error checking for insertRes
            return null;
        }

        // check if destination nickname exists
        logger.info(`Searching users for nickname '${target}'`);
        const nicknameRes = await FindOne({nickname: target});
        if (!nicknameRes) {
            logger.error(`Nickname does not exist!`);
            return {"err": Numerics["ERR_NOSUCHNICK"](target)};
        }
        logger.info(`NicknameRes == ${JSON.stringify(nicknameRes)}`);

        // get source nickname from clientIP
        logger.info(`Searching for ${clientSocket.remoteAddress}'s nickname...`);
        const senderNicknameRes = await FindOne(
            {ip: clientSocket.remoteAddress},
            process.env.MONGODB_CHAT_USERS_COLLECTION_NAME, 
            process.env.MONGODB_NAME, 
            {nickname: 1}
        );

        const collectionName = `${target}_usernick`
        const insertRes = await InsertOne(
            SCHEMA_ChatMsg(
                "userid",
                "server",
                message,
                "tags"
            ), 
            collectionName,
            process.env.MONGODB_CHAT_MESSAGE_DB_NAME
            );

        // how does it work? just anyone can private message anyone?
        // yes
        // allow for blocking though
        return null;

        // return {"req": ["clients", "capabilities", "serverVersion", "isClient", "clientIP"], "callback": callback};

    } catch (error) {
        logger.info(`PRIVMSG error: ${error}`);
        return {"err": Numerics["ERR_CANNOTSENDTOCHAN"]()};
    }
};

module.exports = {
    PRIVMSG
};