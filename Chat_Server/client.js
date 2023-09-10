// const { ClientCommands } = require('./commands');
const net = require('node:net');
const { Buffer } = require('node:buffer');
const fs = require('node:fs');

/**
 * ClientServer is an abstraction of a Server to connect to.
 * 
 * @param {*} name 
 * @param {*} host 
 * @param {*} port // probably don't even need the port included here as we don't necessarily want to allow users to choose which port to connect to. We just default to use TLS.
 */
const ClientServer = (
    name = "default user",
    host = "127.0.0.1",
    port = "6667", // probably don't even need the port included here as we don't necessarily want to allow users to choose which port to connect to. We just default to use TLS.
    ) => {
    // const Connect = (d = "default user") => {
    //     // net.createConnection(port, host);
        
    // }
    const MAX_CLIENT_MESSAGE_LEN = 512;

    console.log(`${name} connecting to server ${host}:${port}`);
    const client = net.createConnection({ port: port, host: host }, () => {
        // 'connect' listener.
      //   console.log('connected to server!');
        client.setEncoding('utf8');

        const clientMsg = Buffer.from(name, 'utf8');

        // TODO
        // chunk client message if greater than 512 bytes
        // +2 here represents the \r\n characters at the end of every message
        if ((clientMsg.length + 2) > MAX_CLIENT_MESSAGE_LEN) {
            console.log("Message greater than 512, will chunk that bad boy");
        }

        // write defaults to UTF8
        if (clientMsg.length != 0 && !client.write(clientMsg, 'utf8', console.log("Message sent sucessfully"))) {
            // write was not successfull
            console.log("WRITE NOT FINISHED");
        } 
      });
      client.on('data', (data) => {
        console.log(data.toString());
        // client.end();
      });
      client.on('end', () => {
      //   console.log('disconnected from server');
      });
}

const Client = (
    nickname = "",
    _socket = null,
    ) => {
    // TODO
    // clean user input

    // const commands = ClientCommands;
    const socket = _socket;

    const verifyNickName = () => {
        const invalidNickName = /^([$:#&])|([\\0\s,*?!@.])/g;
        
        if (nickname === undefined || nickname.length === 0) {
            console.log("Nickname must be at least 1 character in length.");
            throw new Error("Nickname must be at least 1 character in length.");
        } else if (nickname.length > 64) {
            console.log("Nickname cannot be longer than 64 characters.");
            throw new Error("Nickname cannot be longer than 64 characters.");
        }

        let err = "";
        while ((res = invalidNickName.exec(nickname)) !== null) {
            if (res.index === 0) {
                err += `Nickname cannot start with "${res[0]}" character.\n`;
            } else {
                err += `Nickname cannot contain "${res[0]}" characters.\n`;
            }
        }

        if (err) {
            throw new Error(err);
        }

        console.log(`${nickname} is a valid nickname.`);

    }

    verifyNickName();

    const getNickName = () => console.log(nickname);

    const sendMessage = (message = "") => {
        console.log(`${nickname} sent message ${message}`)
    }

    const JoinServer = (server = "") => {

    }

    const leaveServer = () => {

    }

    const joinChannel = (name = "") => {

    }

    const composeMessage = (message = "") => {
        console.log("Composing message")

        /**
         * Message format:
         * 
         * [] == optional
         * SPACE == ' '
         * 
         * message ::= ['@' <tags> SPACE] [':' <source> SPACE] <command> <parameters> <crlf>
         */

        // tags: optional metadata on a message starting with ('@', 0x40)
        // source: Optional note of where the message came from, starting with (':', 0x3A)
        // command: The specific command this message represents
        // parameters: if it exists, data relevant to this specific command

        // Servers limit message size to 512 bytes in length including the CR-LF characters.
        // Implementations which include message tags need to allow additional bytes for the tags section
        // of a message; clients must allow 8191 additional bytes and servers must allow 4096 additional bytes.

        // Examples:
        // :irc.example.com CAP LS * :multi-prefix extended-join sasl
        // @id=234AB :dan!d@localhost PRIVMSG #chan :Hey what's up! CAP REQ :sasl

        message += "\r\n"
    }

    // TEST
    // ClientServer(nickname);

    return { 
        getNickName, 
        sendMessage, 
        nickname,
        socket
    };
}

// try {
//     // c2 = Client("zach")
//     // c1 = Client("jeff")
//     // c3 = Client("peter")
//     // c1.sendMessage("hello")
// } catch (error) {
//     console.log(error)
// }

module.exports = {
    Client
};