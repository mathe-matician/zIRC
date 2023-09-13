const { FindOne, UpdateOne } = require("../db");
const { Numerics } = require("../numerics");

const _logger = require('pino')();
const logger = _logger.child({ Service: 'Chat Server', Command: "NICK" });

//parameters, clients, clientSocket
const NICK = async (nickname, clients, clientSocket, clientNickname) => {
    logger.info(`NICK cmd start. nickname = '${nickname}', clients: '${clients}'`);
    const clientIP = clientSocket.remoteAddress;
    const findRes = await FindOne(
        { $or: [{ip: clientIP, "state.capStarted": true}, {ip: clientIP, "registered": true}]}, 
        process.env.MONGODB_CHAT_USERS_COLLECTION_NAME
        );
    logger.info(`NICK findRes = ${Object.keys(findRes)}`);
    if (!findRes) {
        return {"err": Numerics["ERR_NOTREGISTERED"]()}
    }

    if (!nickname) {
        return {"err": Numerics["ERR_NONICKNAMEGIVEN"]()};
    }

    const invalidNickName = /^([$:#&])|([\\0\s,*?!@.])/g;
    if (nickname === undefined || nickname.length === 0) {
        logger.info(`Nickname must be at least 1 character in length: '${nickname}'`);
        return {"err": Numerics["ERR_ERRONEUSNICKNAME"](nickname, "Nickname must be at least 1 character in length")};
    } else if (nickname.length > 64) {
        logger.info("Nickname cannot be longer than 64 characters.");
        return {"err": Numerics["ERR_ERRONEUSNICKNAME"](nickname, "Nickname cannot be longer than 64 characters")};
    }

    let err = "";
    while ((res = invalidNickName.exec(nickname)) !== null) {
        if (res.index === 0) {
            err += `Nickname cannot start with "${res[0]}" character. `;
        } else {
            err += `Nickname cannot contain "${res[0]}" characters. `;
        }
    }

    if (err) {
        return {"err": Numerics["ERR_ERRONEUSNICKNAME"](nickname, err)};
    }

    try {
        const findRes = await FindOne({ "nickname": nickname }, process.env.MONGODB_CHAT_USERS_COLLECTION_NAME);
        if (findRes) {
            logger.info(`Nickname ${nickname} already in use`);
            return {"err": Numerics["ERR_NICKNAMEINUSE"](nickname)};
        }
        const insertRes = await UpdateOne({"ip": clientIP}, {$set: {"nickname": nickname}}, {"upsert": true});
        if (!insertRes || insertRes?.err) {
            logger.info(`Error inserting/updating NICK`);
            return {"err": Numerics["ERR_UNKNOWNERROR"]("CAP", "LS")};
        }
        return {"command": "NICK", "nick": nickname, "res": `:${clientIP} NICK ${nickname}`};
    } catch (error) {
        console.error(`DOIT ERROR == ${error}`);
        return error;
    }
    

    // if (nickname in clients) {
    //     // check whether nickname exists in clients object
    //     // "<client> <nick> :Nickname is already in use"
    //     logger.info(`nickname ${nickname} already in use`);
    //     return {"err": Numerics["ERR_NICKNAMEINUSE"](nickname)};
    // }
    // return {"command": "NICK", "nick": nickname, "res": `:${srcNick} NICK ${nickname}`}
};

module.exports = { NICK };