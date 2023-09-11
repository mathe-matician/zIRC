const net = require('node:net');
const tls = require('tls');
// const chilkatManager = require('./Chilkat/chilkat_manager');
require('dotenv').config();
const _logger = require('pino')();
const logger = _logger.child({ Service: 'Chat Server', Module: "auth_server_comm.js" });

const AuthServer = () => {

    const host = process.env.AUTH_SERVER_HOST;
    const port = process.env.AUTH_SERVER_PORT;

    const Write = async (_args) => {
        return new Promise((resolve, reject) => {
        client = net.createConnection({ port: port, host: host }, async () => {
            client.setEncoding('utf8');
            const argsObj = {
                "token": process.env.AUTH_SERVER_ADMIN_KEY,
                "data": _args
            }
            const args = JSON.stringify(argsObj);
            logger.info(`AuthServer.Send() before sending: ${args}`)
            // logger.info(`AuthServer.Send() before encrypt: ${args}`)
            // const encryptedArgs = chilkatManager.encrypt_decrypt_AES(args, true);
            // client.write(encryptedArgs);
            client.write(args);
        });
        client.on('connect', async () => {
            logger.info("Connected to Auth Server");
        });
        client.on('data', async (data) => {
            logger.info(`Auth Server response: ${data.toString()}`);
            resolve(data.toString())
            client.end();
        });
        client.on('error', async (error) => {
            logger.info("AuthServer error");
            reject(error);
        });
    });
    }

    return { Write };
}

module.exports = {
    AuthServer,
}