const { UpdateOne, FindOneAndUpdate, FindOne } = require("../db");
const { Numerics } = require("../numerics");
const { AuthServer } = require("../auth_server_comm");
const { CRLF } = require("../constants");
const { isJson } = require("../utils");
require('dotenv').config();
const _logger = require('pino')();
const logger = _logger.child({ Command: 'AUTHENTICATE' });

const AUTHENTICATE = async (params, clients, clientSocket, clientNickname, serverName) => {
    logger.info(`AUTHENTICATE Start. params: ${params}, clients: ${clients}, clientSocket: ${clientSocket}, nickName: ${clientNickname}`)
    // client needs to have negotiated sasl cap to use this command
    if (params.length === 0) {
        return Numerics["ERR_NEEDMOREPARAMS"]("AUTHENTICATE");
    }

    /**
     * 
     * @param Object parameters capabilities 
     */
    const callback = async (parameters) => {
        logger.info(`AUTHENTICATE callback. parameters: ${Object.keys(parameters)}`);
        logger.info(`CAPS ==== ${JSON.stringify(parameters["capabilities"])}`);
        if (!(parameters?.capabilities)) {
            return {"err": "callback error for AUTHENTICATE"};
        }
        const clientIP = clientSocket.remoteAddress;
        logger.info(`AUTHENTICATE callback, clientIP: ${clientIP}`);
        /**
         * Get sasl cap
         * Set authentication state
         * Increment authentication step
         */
        const findRes = await FindOneAndUpdate(
            {ip: clientIP, "state.capabilities": ["sasl"]}, 
            // {$set: {"state.auth.isAuthenticating": true}, $inc: {"state.auth.step": 1}}, 
            {$set: {"state.auth.isAuthenticating": true}}, 
            {upsert: true, returnOriginal: false, projection: {"state.auth.step": 1, _id: 0, "state.auth.failures": 1, "state.auth.type": 1}}
            );
        logger.info(`AUTHENTICATE FIND RES = ${JSON.stringify(findRes)}`);
        if (!findRes) {
            return {"err": Numerics["ERR_SASLABORTED"]("You must negotiate SASL capability to authenticate")};
        }

        // TODO
        // what prevents client from logging out which clears the auth state and trying agian?
        // need to not clear ALL auth state. Keep failures and exponetial backoff timer?
        const failures = findRes?.value?.state?.auth?.failures;
        if (failures && failures > 3) {
            return {"err": Numerics["ERR_SASLFAIL"]()};
        }

        if (!findRes?.value?.state?.auth?.step) {
            const updateRes = await UpdateOne(
                {ip: clientIP, "state.capabilities": ["sasl"]},
                {$inc: {"state.auth.step": 1}}
            );
        }

        const clientAuthStep = !findRes?.value?.state?.auth?.step ? 1 : findRes?.value?.state?.auth?.step;
        logger.info(`clientAuthStep = ${clientAuthStep}`);

        const supportedMechanisms = parameters["capabilities"]["sasl"];
        const authParams = String(params);
        if (clientAuthStep === 1 && !supportedMechanisms.includes(params)) {
            if (authParams === "*") {
                const updateRes = await UpdateOne({ip: clientIP}, {$set: {"state.auth": {}}}, {"upsert": true});
                // TODO
                // error checking here
                return {"res": Numerics["ERR_SASLABORTED"]()};
            }
            return {"res": Numerics["RPL_SASLMECHS"](supportedMechanisms, serverName)};
        }

        logger.info(`AUTHENTICATE params = ${authParams}`);

        // MUST prepend args with # character as we are querying an EXISTING auth method on the Auth Server.
        // See auth server README.md for more info.
        const modAuthParams = clientAuthStep === 1 ? authParams.toLowerCase() : authParams;
        // TODO
        // should we decrypt the client message then encrypt it again here?
        // or should we just pass the client message to the auth server?
        // we should probably do some parsing here and then do more in the auth server.
        // regardless, the data should be encrypted when server talks to auth server and back

        // TODO
        // if not on step 1, get the auth method from db and add it to args.
        // this is needed as the auth server needs to agnostically require() the module
        const authType = findRes?.value?.state?.auth?.type;
        if (!authType) {
            // set type here
            const updateRes = await UpdateOne({ip: clientIP}, {$set: {"state.auth.type": modAuthParams}}, {"upsert": true});
            logger.info(`Auth type set: ${updateRes}`);
        }

        const _authType = !authType ? modAuthParams.toLowerCase() : authType.toLowerCase()
        logger.info(`modAuthParams: ${modAuthParams}, authType: ${authType}`);

        /**
         * msg order should be <auth_type>::<step>::<params>
         */
        const args = `#auth_${_authType}::${clientAuthStep}::${modAuthParams.toLowerCase()}`
        try {
            const authServer = AuthServer();
            const authRes = await authServer.Write(args);
            logger.info(authRes);

            // TODO
            // have auth server return the failure type based on the authentication type
            // e.g. return "ERR_SASLFAIL"
            if (isJson(authRes)) {
                logger.info(`AUTHENTICATE authRes === Json: ${authRes}`)
                const authResParsed = JSON.parse(authRes);
                if (authResParsed?.err) {
                    logger.info("AUTHENTICATE authRes contains err");
                    return {"err": Numerics[authResParsed["err"]]()};
                }
            }

            if (!authRes) {
                logger.info("AUTHENTICATE error with Auth Server");
                // TODO
                // inc failure count in db
                return {"err": Numerics["ERR_SASLFAIL"]()};
            }

            // if successful (got past the if checks above) inc the step we are on
            const updateRes = await UpdateOne(
                {ip: clientIP, "state.capabilities": ["sasl"]},
                {$inc: {"state.auth.step": 1}}
            );
            // if (updateRes !== true) {
            //     logger.info(`MONGO ERROR: UpdateOne failed`);
            //     return {"err": Numerics["ERR_UNKNOWNERROR"]()}
            // }
            try {
                const jsonAuthRes = JSON.parse(authRes);
                if (jsonAuthRes?.success === "SASL authentication successful") {
                    // TODO
                    // having multiple writes after another adds them all to the same buffer w/ TCP (remember... stream!)
                    clientSocket.write(`tokenPkg::${JSON.stringify(jsonAuthRes?.tokenPkg)}` + CRLF);
                    clientSocket.write(Numerics["RPL_LOGGEDIN"](jsonAuthRes?.nickname) + CRLF);
                    clientSocket.write(Numerics["RPL_SASLSUCCESS"]() + CRLF);
                    const updateRes = await UpdateOne(
                        {ip: clientIP, "state.capabilities": ["sasl"]}, 
                        {$set: {"state.auth.success": true}}
                    );
                    // if (updateRes !== true) {
                    //     logger.info(`MONGO ERROR: UpdateOne failed`);
                    //     return {"err": Numerics["ERR_UNKNOWNERROR"]()}
                    // }
                } else {
                    logger.info(`PLAIN Auth failed for some reason...`);
                    throw new Error("Error");
                }
            } catch (notJsonObject) {
                logger.info(`authRes not json Object`);
                clientSocket.write(`AUTHENTICATE ${authRes}` + CRLF);
            }

            return null;
        } catch (error) {
            logger.info(`Error during auth server: ${error}`);
            const updateRes = await UpdateOne(
                {ip: clientIP, "state.capabilities": ["sasl"]}, 
                {$inc: {"state.auth.failures": 1}}
            );
            return {"err": Numerics["ERR_SASLFAIL"]()};
        }
    }

    // send a series of AUTHENITCATE MESSAGES 400-byte chunks ending with "AUTHENTICATE +"

    return {"req": ["capabilities", "clientIP"], "callback": callback};
};

const BASIC_AUTHENTICATE = (clientSocket) => {

};

module.exports = {
    AUTHENTICATE
}