const SCHEMA_ChatUser = (user, nickname, password, ip) => {
    return {user, nickname, password, email, ip, registered: true};
}

const SCHEMA_CAPState = (deviceUID="UUID-TEST", ip, clientVersion, capabilities=[]) => {
    return {
        deviceUID: deviceUID,
        ip: ip,
        state: {
            clientVersion: clientVersion,
            capStarted: true,
            capabilities: capabilities
        },
        registered: false
    };
}

module.exports = {
    SCHEMA_CAPState,
}
