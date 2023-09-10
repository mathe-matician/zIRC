const SCHEMA_ChatMsg = (userId, server, message, tags) => {
    return {
        "userId": userId,
        "server": server,
        "message": message,
        "timestamp": Date.now(),
        "tags": tags
    };
};

module.exports = {
    SCHEMA_ChatMsg
};