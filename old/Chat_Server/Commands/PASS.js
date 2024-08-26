const { CheckDBForPassword } = require("../db");
const { Numerics } = require("../numerics");
require('dotenv').config();

const PASS = async (password, nickname, client, clientIP) => {
    const findRes = await FindOne({ip: clientIP, "registered": true}, process.env.MONGODB_CHAT_USERS_COLLECTION_NAME);
    if (!findRes) {
        return {"err": Numerics["ERR_NOTREGISTERED"]()}
    }
    
    // TODO
    // decrypt password here
    
    // only param is the password
    console.log(`PASS start parameters '${password}'`);
    if (!password || password === "" || password.length === 0) {
        return {"err": Numerics["ERR_NEEDMOREPARAMS"]("PASS")};
    }

    const res = await CheckDBForPassword(password, nickname);
    if (res?.err) {
        return res;
    }
    console.log("PASS ALL GOOD");
    return res;
};

module.exports = { PASS };