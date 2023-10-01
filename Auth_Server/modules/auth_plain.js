const { DB } = require("../db");
const cryptoManager = require('../Crypto/crypto_manager');
const { isJson } = require("../utils");
const { Mailer } = require("../Mailer/mailer");
require('dotenv').config();
const _logger = require('pino')();
const logger = _logger.child({ Service: 'Auth Server', Auth_Module: "auth_plain" });

const db = DB();

const Register = async (registerPkg) => {
  logger.info(`Register auth_plain start. registerPkg = ${registerPkg}`);
  if (!isJson(registerPkg)) {
    logger.error("registerPkg is not json??");
    return '{"err": "ERR_REGISTERFAIL"}';
  }

  const parsedRegisterPkg = JSON.parse(registerPkg);
  const password_hash = await cryptoManager.gen_hash(parsedRegisterPkg["password"])

  // generate email verification token here
  const token = cryptoManager.gen_random_bytes();
  const token_expr = new Date()
  token_expr.setHours(token_expr.getHours() + parseInt(process.env.EMAIL_TOKEN_EXPIRATION_HOURS));

  const prp_stmt_Register = db.Prepare(
    "Auth_PLAIN_Register",
    `
    WITH user_insrt AS (
      INSERT INTO users VALUES (default, $1, $2, $3)
      ON CONFLICT DO NOTHING
    RETURNING email
    )
    INSERT INTO email_tokens VALUES (default, (select email from user_insrt), $4, $5) ON CONFLICT DO NOTHING;
    `,
    [
      parsedRegisterPkg?.email, 
      password_hash, 
      parsedRegisterPkg?.nickname,
      token,
      token_expr.toISOString()
    ]
  );

  const registerRes = await db.Exec(prp_stmt_Register);
  if (isJson(registerRes)) {
    logger.info(`registerRes: ${JSON.stringify(registerRes)}`);
  } else {
    logger.info(`registerRes: ${registerRes}`);
  }
  
  if (!registerRes) {
    logger.error("Failed registration when trying to insert into DB");
    return '{"err": "ERR_REGISTERFAIL"}';
  }

  logger.info(`Successfully inserted new user into db`);

  // TODO
  // send email with token
  // const MailerInstance = Mailer();


  return `{"res": ""}`;
};

const AuthCheck = async (token) => {
  logger.info(`AUTH_PLAIN AuthCheck start - With token: ${token}`)

  // Select both the token_expr and the nickname of the client making the request
  // the current client nickname is used in many different commands, so we just return it for all of them.
  const prp_stmt_AuthCheck = db.Prepare(
    "Auth_PLAIN_AuthCheck",
    `WITH pkg AS (
      SELECT token_expr, user_id
      FROM user_tokens
      WHERE token = $1
      )
      SELECT p.token_expr, u.nickname
      FROM pkg p
      JOIN users u ON p.user_id = u.id`,
      [token]
  );
  const authExpr = await db.Exec(prp_stmt_AuthCheck);
  logger.info(`Auth_PLAIN_AuthCheck Expiration: ${JSON.stringify(authExpr)}`);
  if (!authExpr) {
    logger.info("Failed auth check when checking DB");
    return '{"err": "ERR_SASLFAIL"}';
  }
  logger.info(`Successfully got token_expr and nickname from DB!`);

  const now = new Date();
  const expiration = new Date(authExpr?.token_expr);
  if (now > expiration) {
    logger.info(`Auth_PLAIN_AuthCheck token is expired. now: ${now}, expiration: ${expiration}`);
    return '{"err": "ERR_SASLFAIL"}';
  }
  logger.info("AuthCheck Successful!");
  return `{"nickname": "${authExpr?.nickname}"}`;
}

/**
 * 
 * @param {array} args
 * Iteration 1:
 *      @param {string} args[0] iteration 
 *      @param {string} args[1] user
 *      @param {string} args[2] nickname?
 * Iteration 2:
 *      @param {string} args[0] iteration 
 *      @param {string} args[1] email:password
 * @returns 
 */
const Exec = async (args) => {
    logger.info(`AUTH_PLAIN Module start. Args: ${JSON.stringify(args)}`);
    // if (args.length !== 2) {
    //     logger.info("Plain authentication is missing arguments.");
    //     return null;
    // }

    switch (args[0]) {
        case "1":
            logger.info("PLAIN AUTH STEP 1");
            // TODO
            // check if email, realname, and nickname state exists,
            // if not tell the user to set those specific pieces.
            // this data must be sent and parsed here!
            // OR
            // we do this in the chat app itself before even reaching the auth server..
            // probably that.

            // TODO
            // check if this is a valid sasl mechanism?
            // do we already do this?
            return "+" // initial ack from server, user would then send `AUTHENTICATE <base64 encoded username:pw>`

        case "2":
            // const db = DB();
            const _split = args[1].split(':');
            const email = _split[0].trim();
            const findRes = await FindOneAndUpdate(
                // TODO search for email && ip as I can't think of a situation wehre the ip woudl change MID authentication
                {email: email},
                {$set: 
                    {
                        "state.auth.isAuthenticating": true,
                        "state.capabilities": ["sasl"],
                        ip: "*", // TODO need to send client ip in auth message?
                        "state.auth.type": "plain",
                        "state.auth.step": 2,
                        "state.auth.failures": 0
                    }
                }, 
                {
                    upsert: true, 
                    returnOriginal: false, 
                    projection: {
                        _id: 0, 
                        "state.auth.failures": 1
                    }
                }
            );
            if (!findRes) {
              return '{"err": "ERR_SASLFAIL"}';
            }

            logger.info(`auth plain step 2 findRes: ${JSON.stringify(findRes)}`);
            const password = _split[1].trim();
            const prp_stmt_AuthCheck = db.Prepare(
              "Auth_PLAIN", "SELECT (id, salt, password, nickname) FROM users WHERE email = $1", [email]
              );
            const authCheckRes = await db.Exec(prp_stmt_AuthCheck);
            logger.info(`Auth: ${JSON.stringify(authCheckRes)}`);
            if (!authCheckRes) {
              logger.info("Failed auth check when checking DB");
              // return "failure";
              return '{"err": "ERR_SASLFAIL"}';
            }

            logger.info(`Successfully got creds from DB!`);
            const checkResult = authCheckRes["row"];
            const rowValues = checkResult.substring(
              checkResult.indexOf("(") + 1,
              checkResult.lastIndexOf(")")
            ).split(",");
            logger.info(rowValues);
            
            const userId = rowValues[0];
            const dbSalt = rowValues[1];
            const pwHash = rowValues[2];
            const nickname = rowValues[3];
            const hash = cryptoManager.generate_hash(password, dbSalt);
            if (hash !== pwHash) {
              logger.info("Password doesn't match")
              // return "Credentials incorrect";
              return '{"err": "ERR_SASLFAIL"}';
            }

            const tokenPkg = cryptoManager.session_token();
            logger.info(`Before prep statement. userId: ${userId}, tokenPkg["token"]: ${tokenPkg["token"]}, tokenPkg["expr"]: ${tokenPkg["expr"]}`);
            const prp_stmt_SessionId = db.Prepare(
              "Auth_PLAIN_SessionToken", 
              "INSERT INTO user_tokens(user_id, token, token_expr) VALUES ($1, $2, $3)", 
              [userId, tokenPkg["token"], tokenPkg["expr"]]
              );
            const sessionIdRes = await db.Exec(prp_stmt_SessionId);
            logger.info(`sessionIdRes: ${sessionIdRes}`);
            // if (!sessionIdRes) {
            //   logger.info("Failed to insert sessionId and sessionExpr");
            //   return "failure";
            // }

            return JSON.stringify({
              "success": "SASL authentication successful", 
              "nickname": nickname,
              "tokenPkg": tokenPkg
            });
            
        default:
            return '{"err": "ERR_SASLFAIL"}';
    }
};

module.exports = {
    Exec,
    AuthCheck,
    Register
};