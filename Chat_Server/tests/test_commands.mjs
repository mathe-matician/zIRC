const name = "default user";
const host = "127.0.0.1";
const port = "6667";

const TestNICK = () => {
    const client = net.createConnection({ port: port, host: host }, () => {
        // 'connect' listener.
      //   console.log('connected to server!');
        client.setEncoding('utf8');
    
        const test1 = "NICK pickle";
        const test2 = ":pickle NICK pickle_man";
        const clientMsg = Buffer.from(name, 'utf8');
    
        client.write(clientMsg, 'utf8', console.log("Message sent sucessfully"))
    });
    client.on('data', (data) => {
      console.log(data.toString());
      // client.end();
    });
};