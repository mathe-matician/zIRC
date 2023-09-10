const net = require('node:net');
const tls = require('tls');
// const chilkatManager = require('./Chilkat/chilkat_manager');
require('dotenv').config();

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
            console.log(`AuthServer.Send() before sending: ${args}`)
            // console.log(`AuthServer.Send() before encrypt: ${args}`)
            // const encryptedArgs = chilkatManager.encrypt_decrypt_AES(args, true);
            // client.write(encryptedArgs);
            client.write(args);
        });
        client.on('connect', async () => {
            console.log("Connected to Auth Server");
        });
        client.on('data', async (data) => {
            console.log(`Auth Server response: ${data.toString()}`);
            resolve(data.toString())
            client.end();
        });
        client.on('error', async (error) => {
            console.log("AuthServer error");
            reject(error);
        });
    });
    }

    return { Write };
}

module.exports = {
    AuthServer,
}