const { Numerics } = require("../numerics");
const { AuthServer } = require("../auth_server_comm");
const _logger = require('pino')();
const logger = _logger.child({ Command: 'VERIFY' });

const VERIFY = async (params, clients, clientSocket, clientNickname, serverName) => {
    // ERR_INVALIDCODE
    // ERR_TEMPORARILYUNAVAILABLE
    // ERR_ACCOUNTREQUIRED
    // ERR_ALREADYAUTHENTICATED
    // RPL_SUCCESS Account successfully registered
    
    // VERIFY <account|nickname> <code>
    if (params.length < 2) {
        return {"err": Numerics["ERR_NEEDMOREPARAMS"]("VERIFY")}
    }

    if (params.length > 3) {
        return {"err": Numerics["ERR_UNKNOWNERROR"]("VERIFY")};
    }

    if (params[0] === "RESEND") {
        // resend a new verification code

        // at this point if the client is trying to resend the verification code
        // they should have a nickname in mongodb and a nickname in the auth server user table.
        // so we check whether there exists a user with the email <email>
        // if this is true, then we send another verification email.

        // TODO
        // check if params[1] is a valid email.

        const pkg = `{"account": "${params[1]}", "email": "${params[2]}"}`;

        const args = `#auth_plain::emailcheck::${pkg}`;
        logger.info(`CHAT SERVER BEFORE AUTH CHECK: ${args}`);
        let authServer = AuthServer();
        const authServerRes = await authServer.Write(args);
        logger.info(`EmailCheck res === ${authServerRes}`);
        authServer = null; // mark for garbage collection.

        emailCheckRes = JSON.parse(authServerRes);
        if (emailCheckRes?.err) {
            logger.error(`ERROR during email check ${emailCheckRes["err"]}`);
            return Numerics[emailCheckRes["err"]]();
        }
        
        // get params[1] which should be email
        // pending verification; verification code has been sent to
        return Numerics["RPL_VERIFICATIONREQUIRED"](
            clientNickname,
            `Pending verification; verification code has been sent to`,
            params[1]
        ); 
    }

    const pkg = `{"account": "${params[1]}", "email": "${params[2]}"}`;

    const args = `#auth_plain::emailcheck::${pkg}`;
    logger.info(`CHAT SERVER BEFORE AUTH CHECK: ${args}`);
    let authServer = AuthServer();
    const authServerRes = await authServer.Write(args);
    logger.info(`EmailCheck res === ${authServerRes}`);
    authServer = null; // mark for garbage collection.

    emailCheckRes = JSON.parse(authServerRes);
    if (emailCheckRes?.err) {
        logger.error(`ERROR during email check ${emailCheckRes["err"]}`);
        return Numerics[emailCheckRes["err"]]();
    }

    const parsedRes = JSON.parse(emailCheckRes?.res);

    if (parsedRes !== "RPL_SUCCESS") {
        return {"err": Numerics("ERR_UNKNOWNERROR")("VERIFY")};
    }

    // if we got here without an error, this should be RPL_SUCCESS
    return {"res": Numerics(parsedRes)("VERIFY")};
}

module.exports = {
    VERIFY
};