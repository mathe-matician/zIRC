const SCHEMA_ChatMsg = (sender, server, message, prefixes, tags) => {
    return {
        "sender": sender,
        "server": server,
        "message": message,
        "timestamp": Date.now(),
        "prefixes": prefixes,
        "tags": tags
    };
};

module.exports = {
    SCHEMA_ChatMsg
};