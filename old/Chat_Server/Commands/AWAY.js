const { FindOne } = require("../db");
const { Numerics } = require("../numerics");
require('dotenv').config();

const AWAY = async () => {
    const findRes = await FindOne({ip: clientIP, "registered": true}, process.env.MONGODB_CHAT_USERS_COLLECTION_NAME);
    if (!findRes) {
        return {"err": Numerics["ERR_NOTREGISTERED"]()}
    }
};

module.exports = AWAY;