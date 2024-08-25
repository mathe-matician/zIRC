const { Numerics } = require("../numerics");
const _logger = require('pino')();
const logger = _logger.child({ Command: 'REGISTER' });

/**
 * "draft/account-registration" capability allows the server to support REGISTER command
 * REGISTER <nickname> {<email> | "*"} <password>
 */
const REGISTER = async (params, clients, clientSocket, clientNickname, serverName) => {
    
    // TODO
    // password length
    // valid email
    // valid nickname
    
    // TODO
    // pass this to AUTHENTICATE?

    // TODO
    // email verification
    // const res = await RegisterClient(params[0], params[1], params[2]);
    // if (res?.err) {
    //     return "ERRROR REGISTERING";
    // }
    // return res;

    logger.info(`Start: params: ${JSON.stringify(params)}`);

    const callback = async (parameters) => {
        console.log(`REGISTER parameters= ${JSON.stringify(parameters)}`);
        try {
            const capabilities = parameters?.capabilities;
            if (capabilities.length !== 0 && "draft/account-registration" in capabilities) {
                const { draft_account_registration } = require("./Capabilities/draft_account_registration");
                let accountRegistrationAttributes = [];
                if ("draft/account-registration" in capabilities) {
                    // When there are multiple attributes to the capability
                    accountRegistrationAttributes = capabilities["draft/account-registration"].split(",");
                }
                    const draftAccountRegRes = await draft_account_registration(accountRegistrationAttributes, params, clientSocket.remoteAddress);
                    return {"res": draftAccountRegRes};
            } else {
                return {"err": Numerics["ERR_INVALIDCAPCMD"](clientNickname, "CAP REQ draft/account-registration")}
            }
        } catch (error) {
            logger.error(error)
            return {"err": Numerics["ERR_REGISTERFAIL"]()};
        }

        // return {"res": `${nickname} :${nickname}!${nickname}@${parameters[1]}`};
    };

    // request capabilities to check whether draft/account-registration exists on server as it is needed
    // for this REGISTER command to work
    return {"req": ["capabilities"], "callback": callback};
};

module.exports = {
    REGISTER
};