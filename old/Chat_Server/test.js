const net = require('node:net');
const { MongoClient, ServerApiVersion } = require('mongodb');
require('dotenv').config();

const client = new MongoClient(process.env.GYMLETE_MONGODB_CONN_STRING,  {
  serverApi: {
      version: ServerApiVersion.v1,
      strict: true,
      deprecationErrors: true,
  }
});

const DoIt = async (nickname="zach") => {
  try {
      await client.connect();
      const db = await client.db(process.env.GYMLETE_DB_NAME);
      const searchNickRes = await db.collection(
        process.env.GYMLETE_DB_USER_COLLECTION_NAME)
        .findOne({ "chat.nickname": nickname });
      if (searchNickRes) {
          console.log(`INSIDE: nickname ${nickname} already in use`);
          // client.close();
          return `nickname ${nickname} already in use`;
      }
      return "Success!!";
  } catch (error) {
      console.log(`ERROR == ${error}`);
  } finally {
    await client.close();
  }
}

const commands = {
  "DoIt": DoIt
}

const processmsg = async () => {
  // const resObj = {res: ""};
  // server.emit("LOGIN", resObj);
  const res = await commands["DoIt"]();
  console.log(`RESULT === ${res}`);
  console.log("Some other process stuff1")
  console.log("Some other process stuff2")
  console.log("Some other process stuff3")
  return res;
}

const server = net.createServer(async (socket) => {
    socket.setEncoding('utf8');
    socket.on('data', async (data) => {
      // const res = await DoIt();
      // console.log(`res == ${res}`);
      const res = await processmsg();
      
      console.log("MAIN SOME OTHER STUFF");
      console.log("MAIN SOME OTHER STUFF");
      console.log(`MAIN RES = ${res}`)
    });

    // server.on("LOGIN", async (resObj) => {
    //   const res = await DoIt("zach");
    //   console.log(`LOGIN res == ${res}`);
    //   console.log("LOGIN SOME OTHER STUFFFF");
    // });
  });

server.listen("6667", "127.0.0.1", () => {
    console.log('server started:', server.address());
    console.log(`Max connections = ${server.maxConnections}`);
});