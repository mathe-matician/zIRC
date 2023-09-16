// const { RegisterClient } = require("../../db");
const _logger = require('pino')();
const logger = _logger.child({ Capability: 'draft/account-registration' });

/**
 * "draft/account-registration" capability allows the server to support REGISTER command
 * REGISTER <nickname> {<email> | "*"} <password>
 */
const draft_account_registration = async (accountRegistrationAttributes, register_params) => {
    
    // TODO
    // password length
    // valid email
    // valid nickname
    
    // TODO
    // pass this to AUTHENTICATE?

    // TODO
    // email verification

    // TODO
    // check all attributes:
    // before-connect
    // email-required
    // custom-account-name
    // const res = await RegisterClient(params[0], params[1], params[2]);
    // if (res?.err) {
    //     return "ERRROR REGISTERING";
    // }
    // return res;

    logger.info(`Start: accountRegistrationAttributes: ${accountRegistrationAttributes} register_params: ${register_params}`)

    const email_required = () => {

    }

    const custom_account_name = () => {
        
    }

    const before_connect = () => {
        
    }

    for (attr of accountRegistrationAttributes) {
        if (attr === "email-required") {
            // REGISTER needs an email with it
            logger.info("email-required set");
        } 
        
        if (attr === "custom-account-name") {
            // then this desired account name can be different from the user’s current nickname.
            logger.info("custom-account-name set");
        } 
        
        if (attr === "before-connect") {
            logger.info("before-connect set");
        }
    }

    return 
};

module.exports = {
    draft_account_registration
};