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
const PRIVMSG = async (params, clients, clientSocket, clientNickname, serverName) => {
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
        if (message === "") {
            return {"err": Numerics["ERR_NOTEXTTOSEND"](clientNickname)}
        }

        const target = params[0].trim();
        const chanTypeIndex = target.indexOf("#");
        if (chanTypeIndex !== -1) {
            // then it is a channel
            // parse any ops before the channel
            logger.info("PRIVMSG: message is to channel")
            const ops = target.substring(0, chanTypeIndex);
            logger.info(`PRIVMSG: ops == ${ops}`);
            const targetName = target.substring(chanTypeIndex, target.length);
            if (targetName === "") {
                return {"err": Numerics["ERR_NORECIPIENT"](clientNickname, "PRIVMSG")};
            }

            logger.info(`PRIVMSG: searching for collection ${targetName}`);
            const colExists = await CollectionExists(targetName);
            if (!colExists) {
                logger.error(`Collection '${targetName}' does not exist`);
                return {"err": Numerics["ERR_NOSUCHCHANNEL"](targetName)};
            }
            // insert message
            const insertRes = await InsertOne(
                    SCHEMA_ChatMsg(
                        clientNickname,
                        serverName,
                        message,
                        ops,
                        "tags"
                    ), 
                targetName, 
                process.env.MONGODB_CHAT_MESSAGE_DB_NAME
            );
            logger.info(`insertRes == ${JSON.stringify(insertRes)}`);
            if (!insertRes) {
                return {"err": Numerics["ERR_CANNOTSENDTOCHAN"](clientNickname, targetName)};
            }
            return null;
        }
        let targetTrimmed = "";
        if (target[0] === "$") {
            // When a message is sent from a server to a client, the target is prefixed with a $
            logger.info(`Message is from Server to client`);
            targetTrimmed = str.substring(1);
        } else {
            logger.info(`Message is from client to client`);
            targetTrimmed = target;
            logger.info(`Sending message from ${clientNickname} to ${targetTrimmed}`);
        }

        // check if destination nickname exists
        logger.info(`Searching users for nickname '${targetTrimmed}'`);
        const nicknameRes = await FindOne({nickname: targetTrimmed});
        if (!nicknameRes) {
            logger.error(`Nickname does not exist!`);
            return {"err": Numerics["ERR_NOSUCHNICK"](targetTrimmed, clientNickname)};
        }
        logger.info(`NicknameRes == ${JSON.stringify(nicknameRes)}`);

        logger.info(`Sending message from ${clientNickname} to ${targetTrimmed}`);

        const collectionName = `${targetTrimmed}_${clientNickname}`
        const insertRes = await InsertOne(
            SCHEMA_ChatMsg(
                clientNickname,
                serverName,
                message,
                "",
                "tags"
            ), 
            collectionName,
            process.env.MONGODB_CHAT_MESSAGE_DB_NAME
            );
        if (!insertRes) {
            logger.error("Insert message failure");
            throw new Error("Error");
        }

        logger.info(`PRIVMSG insertRes = ${JSON.stringify(insertRes)}`)
        // allow for blocking though
        return null;

        // return {"req": ["clients", "capabilities", "serverVersion", "isClient", "clientIP"], "callback": callback};

    } catch (error) {
        logger.error(`PRIVMSG error: ${error}`);
        return {"err": Numerics["ERR_NOSUCHNICK"](clientNickname)};
    }
};

module.exports = {
    PRIVMSG
};