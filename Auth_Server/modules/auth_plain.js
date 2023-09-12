const { DB } = require("../db");
const cryptoManager = require('../Crypto/crypto_manager');
require('dotenv').config();

const db = DB();

const AuthCheck = async (token) => {
  console.log(`\nAUTH_PLAIN AuthCheck start\nWith token: ${token}`)
  const prp_stmt_AuthCheck = db.Prepare(
    "Auth_PLAIN_AuthCheck", "SELECT (token_expr) FROM user_tokens WHERE token = $1", [token]
    );
  const authExpr = await db.Exec(prp_stmt_AuthCheck);
  console.log(`Auth_PLAIN_AuthCheck Expiration: ${authExpr}`);
  if (!authExpr) {
    console.log("Failed auth check when checking DB");
    return '{"err": "ERR_SASLFAIL"}';
    // return "failure";
  }
  console.log(`Successfully got token_expr from DB!`);

  const now = new Date();
  const expiration = new Date(authExpr);
  if (now > expiration) {
    console.log(`Auth_PLAIN_AuthCheck token is expired. now: ${now}, expiration: ${expiration}`);
    // throw new Error("Error");
    return '{"err": "ERR_SASLFAIL"}';
  }
  return '{"success": true}';
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
    console.log(`AUTH_PLAIN Module start. Args: ${JSON.stringify(args)}`);
    // if (args.length !== 2) {
    //     console.log("Plain authentication is missing arguments.");
    //     return null;
    // }

    switch (args[0]) {
        case "1":
            console.log("PLAIN AUTH STEP 1");
            // TODO
            // check if email, realname, and nickname state exists,
            // if not tell the user to set those specific pieces.
            // this data must be sent and parsed here!
            // OR
            // we do this in the chat app itself before even reaching the auth server..
            // probably that.
            return "+" // initial ack from server, user would then send `AUTHENTICATE <base64 encoded username:pw>`

        case "2":
            // const db = DB();
            const _split = args[1].split(':');
            const email = _split[0];
            const password = _split[1];
            const prp_stmt_AuthCheck = db.Prepare(
              "Auth_PLAIN", "SELECT (id, salt, password, nickname) FROM users WHERE email = $1", [email]
              );
            const authCheckRes = await db.Exec(prp_stmt_AuthCheck);
            console.log(`Auth: ${JSON.stringify(authCheckRes)}`);
            if (!authCheckRes) {
              console.log("Failed auth check when checking DB");
              // return "failure";
              return '{"err": "ERR_SASLFAIL"}';
            }
            console.log(`Successfully got creds from DB!`);
            const checkResult = authCheckRes["row"];
            const rowValues = checkResult.substring(
              checkResult.indexOf("(") + 1,
              checkResult.lastIndexOf(")")
            ).split(",");
            console.log(rowValues);
            
            const userId = rowValues[0];
            const dbSalt = rowValues[1];
            const pwHash = rowValues[2];
            const nickname = rowValues[3];
            const hash = cryptoManager.generate_hash(password, dbSalt);
            if (hash !== pwHash) {
              console.log("Password doesn't match")
              // return "Credentials incorrect";
              return '{"err": "ERR_SASLFAIL"}';
            }

            const tokenPkg = cryptoManager.session_token();
            console.log(`Before prep statement. userId: ${userId}, tokenPkg["token"]: ${tokenPkg["token"]}, tokenPkg["expr"]: ${tokenPkg["expr"]}`);
            const prp_stmt_SessionId = db.Prepare(
              "Auth_PLAIN_SessionToken", 
              "INSERT INTO user_tokens(user_id, token, token_expr) VALUES ($1, $2, $3)", 
              [userId, tokenPkg["token"], tokenPkg["expr"]]
              );
            const sessionIdRes = await db.Exec(prp_stmt_SessionId);
            console.log(`sessionIdRes: ${sessionIdRes}`);
            // if (!sessionIdRes) {
            //   console.log("Failed to insert sessionId and sessionExpr");
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
    AuthCheck
};