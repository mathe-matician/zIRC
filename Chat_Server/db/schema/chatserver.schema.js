const SCHEMA_Server = (
    version="302",
    serverName="",
    host="",
    tls="",
    capabilities={"sasl": "PLAIN"},
    channels={},
    ) => {
    return {
        "version": version,
        "serverName": serverName,
        "host": host,
        "tls": tls,
        "capabilities": capabilities,
        "channels": channels
    }
};

module.exports = {
    SCHEMA_Server
};