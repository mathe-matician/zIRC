const SCHEMA_ChatMsg = (sender, server, message, tags) => {
    return {
        "sender": sender,
        "server": server,
        "message": message,
        "timestamp": Date.now(),
        "tags": tags
    };
};

module.exports = {
    SCHEMA_ChatMsg
};