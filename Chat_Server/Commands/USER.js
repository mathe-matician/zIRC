const { UpdateOne } = require("../db");
const { Numerics } = require("../numerics");
require('dotenv').config();

// cmd: <username> 0 * <realname>
// The USER command is used at the beginning of a connection to specify the username and realname of a new user.
// It must be noted that <realname> must be the last parameter because it may contain SPACE (' ', 0x20) characters, 
// and should be prefixed with a colon (:) if required.

/**
 * 
 * @param {array} params
 *     params[0]: email
 *     params[1]: realname
 * @param {*} clients 
 * @param {*} clientSocket 
 */
const USER = async (params, clients, clientSocket) => {
    // TODO
    // Check if already registered during CAP
    // If a client tries to send the USER command after they have already completed registration with the server, the ERR_ALREADYREGISTERED reply should be sent and the attempt should fail.
    if (params.length !== 2 || params[0].length === 0 || params[1].length === 0) {
        // if not correct number of params OR
        // if no email, request more params
        return {'err': Numerics["ERR_NEEDMOREPARAMS"]("USER")}
    }

    // TODO
    // contact mail server to verify email

    // update mongo auth state email
    // all state is sent when the user tries to use the AUTHENTICATE command
    const updateRes = await UpdateOne(
        {ip: clientSocket.remoteAddress},
        {$set: { 
            "state.auth.email": params[0],
            "state.auth.realname": params[1]
        }},
        { upsert: true }
    );
};

module.exports = USER;