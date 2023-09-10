const { 
    InsertCAPState,
    UpdateOne,
    FindOne,
} = require('../db');

const { Numerics } = require("../numerics");
const { SCHEMA_CAPState } = require("../db/schema/chatuser.schema");

/**
 * Connection Messages
 */
const CAP_LS = async (clientIP, requestedServerVersion, capabilities, serverVersion, deviceUID) => {
    console.log("CAP_LS start");
    console.log(`clientIP: ${clientIP}, requestedServerVersion: ${requestedServerVersion}, capabilities: ${capabilities}, serverVersion: ${serverVersion}, `);
    if (!requestedServerVersion) {
        const findRes = await FindOne({ip: clientIP, "state.capStarted": true}, process.env.MONGODB_CHAT_USERS_COLLECTION_NAME);
        if (!findRes) {
            return {"err": Numerics["ERR_NEEDMOREPARAMS"]("CAP LS")}
        }
        // if here, the client is already in capability negotiation
    } else {
        if (parseInt(requestedServerVersion) < serverVersion) {
            console.log(`requestedServerVersion: ${requestedServerVersion} < serverVersion: ${serverVersion}`)
            return {"err": `You are using an unsupported client version ${requestedServerVersion}`}
        }
        const version = requestedServerVersion > serverVersion ? serverVersion : requestedServerVersion;
        const insertRes = await UpdateOne({"ip": clientIP}, {$set: SCHEMA_CAPState("UUID-TEST", clientIP, version)}, {"upsert": true});
        // const insertCAPStateRes = await InsertCAPState("1234", clientIP, version);
        if (!insertRes || insertRes?.err) {
            console.log(`Error inserting CAP state`);
            return {"err": Numerics["ERR_UNKNOWNERROR"]("CAP", "LS")};
        }
    }

    if (capabilities.length === 0)
        return ":";

    let caps = "";
    Object.entries(capabilities).forEach(([k,v]) => {
        const vals = v === null ? "" : `=${v}`;
        console.log(`CAPABILITIES: ${k}=${vals}`)
        caps += `${k}${vals} `;
    });
    console.log(`All CAPS LS: '${caps}'`);
    return `LS :${caps}`;
};

const CAP_LIST = async (capabilities) => {
    console.log("CAP_LIST start");

};

const CAP_REQ = async (capabilities, requestedCaps, ip) => {
    console.log(`CAP_REQ start: capabilities: ${Object.keys(capabilities)}, requestedCaps: ${requestedCaps}`);
    let caps;
    let toInsert = { "state.capabilities": [], "state.capStarted": true};
    if (requestedCaps.includes(",")) {
        const reqCaps = requestedCaps.split(",");
        for (const cap in reqCaps) {
            console.log(`CAP = ${cap.trim()}`);
            const trimmedCap = cap.trim();
            if (trimmedCap in Object.keys(capabilities)) {
                toInsert["capabilities"].push(trimmedCap);
                caps += `${cap},`;
            }
        }
    } else if (requestedCaps in capabilities) {
        toInsert["state.capabilities"].push(requestedCaps);
        caps = requestedCaps;
    } else {
        caps = Numerics["ERR_INVALIDCAPCMD"]("*", requestedCaps);
    }

    const insertRes = await UpdateOne({"ip": ip}, { $set: toInsert }, {"upsert": true});
    console.log(`CAP REQ insertRes = ${insertRes}`);
    if (insertRes?.err) {
        return insertRes["err"];
    }
    return {"immediateWrite": `ACK ${caps}`};
};

const CAP_ACK = async (capabilities) => {
    console.log("CAP_ACK start");

};

const CAP_NAK = async (capabilities) => {
    console.log("CAP_NAK start");

};

const CAP_END = async (clientIP) => {
    console.log("CAP_END start");
    
    // TODO
    // check success state of auth
    // if successful and we are here
    // then display Numeric 001 welcome message
    const findRes = await FindOne({ip: clientIP, "state.auth.success": true});
    if (findRes) {
        // unset unneccessary auth state for when the client is actually authenticated
        const updateRes = await UpdateOne(
            {ip: clientIP}, 
            {$unset: {"state.auth.isAuthenticating": "", "state.auth.step": ""}}
        );
        // if (updateRes !== true) {
        //     console.log(`MONGO ERROR: UpdateOne failed`);
        //     return {"err": Numerics["ERR_UNKNOWNERROR"]()};
        // }
        return Numerics["RPL_WELCOME"]();
    }

    // clear all state associated with the current capability negotiation
    const updateRes = await UpdateOne({ip: clientIP}, {$set: {state: {}}});
    if (updateRes?.err) {
        return updateRes["err"];
    }

    // TODO
    // what is the return value here to client?
};

const CAP_NEW = async (capabilities) => {
    console.log("CAP_NEW start");

};

const CAP_DEL = async (capabilities) => {
    console.log("CAP_DEL start");

};

const CAP_SHARED_CMDS = {
    "LS": true,
    "LIST": true,
};

const CAP_SERVER_CMDS = {
    "ACK": CAP_ACK,
    "NAK": CAP_NAK,
    "NEW": CAP_NEW,
    "DEL": CAP_DEL,
};

const CAP_CLIENT_CMDS = {
    "REQ": CAP_REQ,
    "END": CAP_END,
};

const CAP = async (subcommand, clients, clientSocket, deviceUID) => {
    // TODO 
    // parse subcommand
    // TODO check if can split as may not be able to.
    console.log(`CAP start, subcommand: ${subcommand}`);
    // if (subcommand[0])
    // const nickname = nickname ? nickname : "*";

    // if (!(subcommand[0] in CAP_CLIENT_CMDS)) {
    //     return Numerics["ERR_INVALIDCAPCMD"]("*", clientSocket.remoteAddress)
    // }

    // if (!subcommand[1])
    // const returnVal = {"res": }

    /**
     * 
     * @param  {...any} parameters 
     * @param parameters[0]: clients object
     * @param parameters[1]: capabilities
     * @param parameters[2]: server IRC version
     * @param parameters[3]: isClient
     * @returns 
     */
    const callback = async (parameters) => {
        console.log(Object.keys(parameters));
        if (Object.keys(parameters).length !== 4) {
            console.log("Not enough params!!!");
            console.log(Object.keys(parameters));
            return {"err": Numerics["ERR_NEEDMOREPARAMS"]("CAP " + subcommand[0])}
        }
        const clients = parameters["clients"];
        const capabilities = parameters["capabilities"];
        const serverVersion = parseInt(parameters["serverVersion"]);
        const isClient = parameters["isClient"];
        // First determine whether user is client or sever

        const clientIP = clientSocket.remoteAddress;

        console.log(`CAP callback, clientHost: ${clientIP}`);
        console.log(`Before switches = subcommand[0] '${subcommand[0]}'`);
        let res;
        if (subcommand[0] in CAP_SHARED_CMDS) {
            switch (subcommand[0]) {
                case "LS":
                    console.log(`Switch LS`);
                    const requestedVersion = subcommand.length < 2 ? null : subcommand[1];
                    res = await CAP_LS(clientIP, requestedVersion, capabilities, serverVersion, deviceUID);
                    break;
                case "LIST":
                    console.log(`Switch LIST`);
                    res = await CAP_LIST();
                    break;
            }
        } else if (isClient && subcommand[0] in CAP_CLIENT_CMDS) {
            console.log(`isClient`);
            switch (subcommand[0]) {
                case "REQ":
                    console.log(`REQ subcommand1 = ${subcommand[1]}`);
                    res = await CAP_REQ(capabilities, subcommand[1], clientIP);
                    break;
                case "END":
                    res = await CAP_END(clientIP);
                    break;
                default:
                    res = Numerics["ERR_INVALIDCAPCMD"]("*", clientIP);
                    break;
            }
            // res = await CAP_CLIENT_CMDS[subcommand[0]](clientSocket.remoteAddress, subcommand[1], capabilities, serverVersion);
        } else if (!(isClient) && subcommand[0] in CAP_SERVER_CMDS) {
            console.log(`is not client`);
            let res;
            switch (subcommand[0]) {
                case "ACK":
                    res = await CAP_ACK(clientSocket);
                    break;
                case "NAK":
                    res = await CAP_NAK();
                    break;
                case "NEW":
                    res = await CAP_NEW();
                    break;
                case "DEL":
                    res = await CAP_DEL();
                    break;
                default:
                    Numerics["ERR_INVALIDCAPCMD"]("*", clientIP)
                    break;
            }
            // res = await CAP_SERVER_CMDS[subcommand[0]](clientSocket.remoteAddress, subcommand[1], capabilities, serverVersion); 
        } else {
            console.log("CAP Subcommand error");
            return {"err": Numerics["ERR_INVALIDCAPCMD"]("*", "CAP " + subcommand[0])}
        }

        if (res?.err) {
            return res;
        }

        if (res?.immediateWrite) {
            console.log(`CAP immediate write: ${res["immediateWrite"]}`);
            clientSocket.write(`CAP ${clientIP} ${res["immediateWrite"]}`);
            return null;
        }

        return {"res": `CAP ${clientIP} ${res}`};
    };

    return {"req": ["clients", "capabilities", "serverVersion", "isClient", "clientIP"], "callback": callback};
};

module.exports = { CAP };