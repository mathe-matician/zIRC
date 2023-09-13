const { FindOne } = require("../db");
require('dotenv').config();

const NOTICE = async (params, clients, clientSocket, clientNickname) => {
    // <target>{,<target>} <text to be sent>
    
    // if (resObj?.command) {
    //     // Commands that land here:
    //     // JOIN
    //     // NICK
    //     const channelType = resObj.channelType;
    //     const channelName = resObj.channelName;
    //     if (channelType in channels && channelName in channels[channelType]) {
    //       // Send a notice to all clients on a specific channel that the client is on
    //       // TODO
    //       // if we support clients being connected to multiple channels at once, this would need to notify all clients on all channels.
    //       // additionally this should notify all privmessages and should update the FE
    //       const clients = channels[resObj.channelType][resObj.channelName]["clients"];
    //       for (const clientName of Object.keys(clients)) {
    //         if (clientName.length !== 0 && typeof(clientName) === String) {
    //           const clientObj = clients[clientName]["clientObj"];
    //           if ("socket" in Object.keys(clientObj) && clientObj["socket"] !== undefined)
    //             clientObj["socket"].write(`${resObj["client"]} ${resObj["command"]} ${resObj["channel"]}`);
    //         }
    //       }
    //     } else {
    //       resObj.res = Numerics["ERR_UNKNOWNERROR"]("NOTICE", [resObj.command]);
    //     }
    //   } else {
    //     resObj.res = Numerics["ERR_UNKNOWNERROR"]("NOTICE", [resObj.command]);
    //   }
};

module.exports = NOTICE;