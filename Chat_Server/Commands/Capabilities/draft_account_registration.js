// const { RegisterClient } = require("../../db");
const _logger = require('pino')();
const logger = _logger.child({ Capability: 'draft/account-registration' });
const { Numerics } = require('../../numerics');
const { FindOneAndUpdate, UpdateOne } = require("../../db");
const { AuthServer } = require("../../auth_server_comm");
const { isJson } = require('../../utils');

/**
 * "draft/account-registration" capability allows the server to support REGISTER command
 * REGISTER <nickname> {<email> | "*"} <password>
 */
const draft_account_registration = async (accountRegistrationAttributes, register_params, clientIP) => {
    
    // TODO
    // password length
    // valid email
    // valid nickname

    // TODO
    // check all attributes:
    // before-connect
    // email-required
    // custom-account-name

    logger.info(`Start: accountRegistrationAttributes: ${accountRegistrationAttributes} register_params: ${register_params}`)

    for (attr of accountRegistrationAttributes) {
        if (attr === "email-required") {
            logger.info("email-required set");
            if (register_params.length !== 3) {
                // we currently need nick name (account name), email and password
                // params will always be 3 if email is required.
                return Numerics["ERR_REGISTERFAIL"]();
            }
        } 
        
        if (attr === "custom-account-name") {
            // then this desired account name can be different from the user’s current nickname.
            logger.info("custom-account-name set");
        } 
        
        if (attr === "before-connect") {
            logger.info("before-connect set");
        }
    }

    const updateRes = await UpdateOne(
        {nickname: register_params[0]}, 
        {
            $setOnInsert: 
            {
                "state.auth": {
                    "isAuthenticating": true,
                    "isRegistering": true
                },
                "state.capabilities": [
                    "sasl",
                    "draft/account-registration"
                ],
                ip: clientIP
            }
        }, 
        {
            upsert: true, 
            // returnOriginal: false, 
            // projection: {
            //     "state.auth.step": 1, 
            //     _id: 0, 
            //     "state.auth.failures": 1, 
            //     "state.auth.type": 1
            // }
        }
    );

    if (isJson(updateRes)) {
        logger.info(`updateRes = ${JSON.stringify(updateRes)}`)
    } else {
        logger.info(`updateRes = ${updateRes}`)
    }

    const registerPkg = {
        "nickname": register_params[0],
        "email": register_params[1],
        "password": register_params[2]
    }
    const args = `#auth_plain::register::${JSON.stringify(registerPkg)}`;

    logger.info(`CHAT SERVER BEFORE AUTH CHECK: ${args}`);
    let authServer = AuthServer();
    const authServerRes = await authServer.Write(args);
    logger.info(`Auth check res === ${authServerRes}`);
    authServer = null; // mark for garbage collection.

    authRes = JSON.parse(authServerRes);
    if (authRes?.err) {
      logger.info(`ERROR: AUTHENTICATE error with Auth Server: ${authRes["err"]}`);
      return Numerics[authRes["err"]]();
    }

    logger.info("Successfully inserted client into db for registration");
    logger.info("About to send verification email");

    // TODO
    // email verification

    return Numerics["RPL_VERIFICATIONREQUIRED"](register_params[0], register_params[1]);
};

module.exports = {
    draft_account_registration
};