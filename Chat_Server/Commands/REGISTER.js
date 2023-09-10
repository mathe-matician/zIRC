const { RegisterClient } = require("../db");

/**
 * "draft/account-registration" capability allows the server to support REGISTER command
 * REGISTER <account> {<email> | "*"} <password>
 */
const REGISTER = async (parameters, client="Zach") => {
    const res = await RegisterClient(client, parameters[0], parameters[1]);
    if (res?.err) {
        return "ERRROR REGISTERING";
    }
    return res;
    // /**
    //  * 
    //  * @param  {...any} parameters 
    //  * @param parameters[0]: clients object
    //  * @param parameters[1]: capabilities
    //  * @returns 
    //  */
    // const callback = (...parameters) => {
    //     console.log(`parameters= ${parameters}`);
    //     const capabilities = parameters[1];
    //     if ("draft/account-registration" in capabilities) {
    //         const capabilityList = capabilities.split(",");
    //         for (cap of capabilityList) {
    //             if (cap === "custom-account-name") {
    //                 // then this desired account name can be different from the user’s current nickname.
    //             }
    //         }
    //     }

    //     return {"response": `${nickname} :${nickname}!${nickname}@${parameters[1]}`};
    // };

    // // request capabilities to check whether draft/account-registration exists on server as it is needed
    // // for this REGISTER command to work
    // return {"req": ["clients", "capabilities"], "callback": callback};
};

module.exports = REGISTER;