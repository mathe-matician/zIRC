const { Find } = require("../db");
const { Numerics } = require("../numerics");
const { CRLF } = require("../constants");
require('dotenv').config();

const _logger = require('pino')();
const logger = _logger.child({ Service: 'Chat Server', Command: "LIST" });

/**
 * ELIST object used for parsing
 */
const ELIST = async () => {

};

/**
 * 
 * @param {*} params 
 * @param {*} clients 
 * @param {*} clientSocket 
 * @returns 
 */
const LIST = async (params, clients, clientSocket, clientNickname) => {
    logger.info(`LIST Start.`);
    // Parameters: [<channel>{,<channel>}] [<elistcond>{,<elistcond>}]
    // both parameters are OPTIONAL

    if (params.length !== 0) {
        if (params[0].includes(",")) {
            const splitParams = params[0].split(",");
            let i = 0;
            for (i; i < params.length; i++) {
                if (params[i][0] !== "#") {
                    break;
                }
            }
        }

        // TODO
        // if param[i][0] in ELIST
        // then
        // parse expression
        // for (i; i < params.length; i++) {
        //     if (param[i][0]) {
        //         break;
        //     }
        // }
    }

    // The first possible parameter to this command is a list of channel names, delimited by a comma (",", 0x2C) character.
    // If this parameter is given, the information for only the given channels is returned. 
    // If this parameter is not given, the information about all visible channels (those not hidden by the secret channel mode rules) is returned.

    // TODO
    // The second possible parameter to this command is a list of conditions as defined in the ELIST RPL_ISUPPORT parameter, 
    // delimited by a comma (",", 0x2C) character. Clients MUST NOT submit an ELIST condition unless the server has explicitly defined support for that condition with the ELIST token. 
    // If this parameter is supplied, the server filters the returned list of channels with the given conditions as specified in the ELIST documentation.

    // In response to a successful LIST command, the server MAY send one RPL_LISTSTART numeric,
    // MUST send back zero or more RPL_LIST numerics, and MUST send back one RPL_LISTEND numeric.
    
    clientSocket.write(Numerics["RPL_LISTSTART"]());

    const channels = await Find({});

    for (const chan of channels) {
        clientSocket.write(Numerics["RPL_LIST"](
            chan["name"], 
            chan["clients"].length,
            chan["topic"])
        );
    }

    clientSocket.write(Numerics["RPL_LISTEND"]());

    const callback = () => {

    }

    return {"req": ["clients", "capabilities", "serverVersion", "isClient", "clientIP"], "callback": callback};
};

module.exports = {
    LIST
};