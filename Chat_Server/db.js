const { 
    Numerics
  } = require("./numerics");
const { MongoClient, ServerApiVersion } = require('mongodb');
// const { chilkat } = require('./Chilkat/chilkat_manager');
const cryptoManager = require('./Crypto/crypto_manager');
const { SCHEMA_CAPState } = require('./db/schema/chatuser.schema');

require('dotenv').config();

const DEFAULT_COLLECTION = process.env.MONGODB_CHAT_USERS_COLLECTION_NAME;
const DEFAULT_DB = process.env.MONGODB_NAME;

const mongoConfig = {
    serverApi: {
        version: ServerApiVersion.v1,
        strict: true,
        deprecationErrors: true,
        useUnifiedTopology: true,
    }
}

const GetMongoConn = () => {
    return new MongoClient(process.env.MONGODB_CONN_STRING, mongoConfig);
}

const Find = async (query, options, collectionName=DEFAULT_COLLECTION, database=DEFAULT_DB) => {
  console.log(`Find start: ${query}`);
    const mongoClient = GetMongoConn();
    try {
        await mongoClient.connect();
        const db = await mongoClient.db(database);
        const findRes = await db.collection(collectionName).find(query, options);
        console.log(`Find res = ${findRes}`);
        if (!findRes) {
          console.log(`No channels`);
        }
        return findRes
    } catch (error) {
        console.log(`Find ERROR == ${error}`);
        return {"err": Numerics["ERR_UNKNOWNERROR"]()};
    } finally {
      await mongoClient.close();
    }
};

const FindOne = async (
  query, 
  collectionName=process.env.MONGODB_CHAT_USERS_COLLECTION_NAME, 
  database=process.env.MONGODB_NAME, 
  options={}
  ) => {
    const mongoClient = GetMongoConn();
    try {
        await mongoClient.connect();
        const db = await mongoClient.db(database);
        const findRes = await db.collection(collectionName).findOne(query, options);
        // if (insertRes.hasWriteError() || insertRes.hasWriteConcernError()) {
        //     console.log(`InsertCAPState error: ${insertRes}`);
        //     return {"err": Numerics["ERR_UNKNOWNERROR"]("CAP")};
        // }
        return findRes
    } catch (error) {
        console.log(`FindOne ERROR: ${error}, Query: ${JSON.stringify(query)}, Database: ${database}, Collection: ${collectionName}`);
        return {"err": Numerics["ERR_UNKNOWNERROR"]()};
    } finally {
      await mongoClient.close();
    }
};

const InsertCAPState = async (deviceUID=null, clientIP=null, requestedServerVersion=null) => {
    console.log(`InsertCAPState start. deviceUID: ${deviceUID}, clientIP: ${clientIP}, requestedServerVersion: ${requestedServerVersion}`)
    if (!deviceUID || !clientIP || !requestedServerVersion) {
        return null;
    }

    const mongoClient = GetMongoConn();
    try {
        await mongoClient.connect();
        const db = await mongoClient.db(process.env.MONGODB_NAME);
        const insertRes = await db.collection(
          process.env.MONGODB_CHAT_USERS_COLLECTION_NAME)
          .insertOne(SCHEMA_CAPState(deviceUID, clientIP, requestedServerVersion), {"upsert": true});
        console.log(`Insert res = ${insertRes}`);
        // if (insertRes.hasWriteError() || insertRes.hasWriteConcernError()) {
        //     console.log(`InsertCAPState error: ${insertRes}`);
        //     return {"err": Numerics["ERR_UNKNOWNERROR"]("CAP")};
        // }
        return true; // just to make sure error isn't thrown when returned
    } catch (error) {
        console.log(`RegisterClient ERROR == ${error}`);
        return {"err": Numerics["ERR_UNKNOWNERROR"]()};
    } finally {
      await mongoClient.close();
    }
};

const MongoDB = async () => {

    // Set the default write concern for the db
    const Init = async () => {
        const mongoClient = GetMongoConn();
        try {
            await mongoClient.connect();
            const db = await mongoClient.db(process.env.MONGODB_NAME);
            db.adminCommand({
                "setDefaultRWConcern" : 1,
                "defaultWriteConcern" : {
                  "w" : 2
                },
                "defaultReadConcern" : { "level" : "majority" }
              });
            db.runCommand({
              "collMod": "tickets",
              "index": {
                "keyPattern": { "lastModifiedDate": 1 },
                "expireAfterSeconds": 10
              }
            });              
        } catch (error) {
            console.log(`RegisterClient ERROR == ${error}`);
            return {"err": Numerics["ERR_UNKNOWNERROR"]()};
        } finally {
          await mongoClient.close();
        }
    }

    return { Init }
}

const InsertOne = async (
  doc={}, 
  collection=process.env.MONGODB_CHAT_USERS_COLLECTION_NAME,
  database=process.env.MONGODB_NAME
  ) => {
  console.log(`InsertOne start. filter: ${Object.keys(doc)}}`);

  const mongoClient = GetMongoConn();
  try {
      await mongoClient.connect();
      const db = await mongoClient.db(database);
      const insertRes = await db.collection(collection).insertOne(doc);
      console.log(`Insert res = ${insertRes}`);
      // if (insertRes.hasWriteError() || insertRes.hasWriteConcernError()) {
      //     console.log(`InsertCAPState error: ${insertRes}`);
      //     return {"err": Numerics["ERR_UNKNOWNERROR"](command)};
      // }
      return true; // just so error isn't thrown
  } catch (error) {
      console.log(`RegisterClient ERROR == ${error}`);
      return {"err": Numerics["ERR_UNKNOWNERROR"]()};
  } finally {
    await mongoClient.close();
  }
};

const UpdateOne = async (
  filter={}, 
  keyValues, 
  options={"upsert": false}, 
  collection=process.env.MONGODB_CHAT_USERS_COLLECTION_NAME,
  database=process.env.MONGODB_NAME
  ) => {
    console.log(`UpdateOne start. filter: ${Object.keys(filter)}, keyValues=${Object.keys(keyValues)}, options=${Object.keys(options)}`);

    const mongoClient = GetMongoConn();
    try {
        await mongoClient.connect();
        const db = await mongoClient.db(database);
        const insertRes = await db.collection(collection).updateOne(filter, keyValues, options);
        console.log(`Insert res = ${insertRes}`);
        // if (insertRes.hasWriteError() || insertRes.hasWriteConcernError()) {
        //     console.log(`InsertCAPState error: ${insertRes}`);
        //     return {"err": Numerics["ERR_UNKNOWNERROR"](command)};
        // }
        return true; // just so error isn't thrown
    } catch (error) {
        console.log(`RegisterClient ERROR == ${error}`);
        return {"err": Numerics["ERR_UNKNOWNERROR"]()};
    } finally {
      await mongoClient.close();
    }
};

const RegisterClient = async (client, nickname, password) => {
    console.log(`RegisterClient: nickname: ${nickname}, password: ${password}`);
    const mongoClient = GetMongoConn();
    try {
        await mongoClient.connect();
        const db = await mongoClient.db(process.env.MONGODB_NAME);
        // var crypt = new chilkat.Crypt2();
        // crypt.HashAlgorithm = process.env.CHILKAT_PASSWORD_HASH_ALGORITHM;
        // The byte representation of the string matters when hashing. 
        // Tell Chilkat that we want to use the utf-8 byte representation.
        crypt.Charset = "utf-8";
        hashBytes = cryptoManager.generate_hash(password);
        // hashBytes = crypt.HashString(password);
        // Let's examine the hash as a hex string.
        // var sb = new chilkat.StringBuilder();
        // sb.AppendEncoded(hashBytes, process.env.CHILKAT_ENCODING_TYPE);
        // console.log("SHA256 hash = " + sb.GetAsString());
        console.log("SHA256 hash = " + hashBytes);
        console.log(`INSERTING TO ${process.env.MONGODB_CHAT_USERS_COLLECTION_NAME}`);

        //TODO
        // edit this insert to create a client db schema object to insert
        // this will include state.
        const insertRes = await db.collection(
          process.env.MONGODB_CHAT_USERS_COLLECTION_NAME)
          .insertOne({ "user": client, "nickname": nickname, "password": hashBytes, "registered": true });
        if (!insertRes) {
            console.log(`Error REGISTERING client`);
            return {"err": Numerics["ERR_UNKNOWNERROR"]("REGISTER")};
        } 
        return {"res": "registered all good"};
    } catch (error) {
        console.log(`RegisterClient ERROR == ${error}`);
        return {"err": Numerics["ERR_UNKNOWNERROR"]()};
    } finally {
      await mongoClient.close();
    }
};

const CheckDBForPassword = async (password, nickname, client="defaultClient") => {
    const mongoClient = GetMongoConn();
    try {
        await mongoClient.connect();
        const db = await mongoClient.db(process.env.GYMLETE_DB_NAME);
        const searchNickRes = await db.collection(
          process.env.MONGODB_CHAT_USERS_COLLECTION_NAME)
          .findOne({ "nickname": nickname, "password": password });
        if (searchNickRes) {
            return {"err": Numerics["ERR_PASSWDMISMATCH"](client)};
        }
        return {"res": "all good"};
    } catch (error) {
        console.log(`CheckDBForPassword ERROR == ${error}`);
        return {"err": Numerics["ERR_UNKNOWNERROR"]()};
    } finally {
      await mongoClient.close();
    }
};

// Updates AND returns doc
const FindOneAndUpdate = async (filter, query, options) => {
  console.log(`FindOneAndUpdate start. filter: ${Object.keys(filter)}, keyValues=${Object.keys(query)}, options=${Object.keys(options)}`);

  const mongoClient = GetMongoConn();
  try {
      await mongoClient.connect();
      const db = await mongoClient.db(process.env.MONGODB_NAME);
      const insertRes = await db.collection(
        process.env.MONGODB_CHAT_USERS_COLLECTION_NAME)
        .findOneAndUpdate(filter, query, options);
      console.log(`FindOneAndUpdate res = ${insertRes}`);
      return insertRes;
  } catch (error) {
      console.log(`FindOneAndUpdate ERROR == ${error}`);
      return {"err": Numerics["ERR_UNKNOWNERROR"]()};
  } finally {
    await mongoClient.close();
  }
};

const CheckDBForNickName = async (nickname="zach", srcNick="srcNick") => {
    const mongoClient = GetMongoConn();
    try {
        await mongoClient.connect();
        const db = await mongoClient.db(process.env.GYMLETE_DB_NAME);
        const searchNickRes = await db.collection(
          process.env.MONGODB_CHAT_USERS_COLLECTION_NAME)
          .findOne({ "nickname": nickname });
        if (searchNickRes) {
            console.log(`INSIDE: nickname ${nickname} already in use`);
            return {"err": Numerics["ERR_NICKNAMEINUSE"](nickname)};
        }
        return {"command": "NICK", "nick": nickname, "res": `:${srcNick} NICK ${nickname}`};
    } catch (error) {
        console.log(`CheckDBForNickName ERROR == ${error}`);
        return {"err": Numerics["ERR_UNKNOWNERROR"]()};
    } finally {
      await mongoClient.close();
    }
};

const InsertClient = async (nickname, password, clientName) => {
    const mongoClient = GetMongoConn();
    try {
        await mongoClient.connect();
        const db = await mongoClient.db(process.env.GYMLETE_DB_NAME);
        // var crypt = new chilkat.Crypt2();
        // crypt.HashAlgorithm = process.env.CHILKAT_PASSWORD_HASH_ALGORITHM;
        // The byte representation of the string matters when hashing. 
        // Tell Chilkat that we want to use the utf-8 byte representation.
        // crypt.Charset = "utf-8";
        // hashBytes = crypt.HashString(password);
        hashBytes = cryptoManager.generate_hash(password);
        // Let's examine the hash as a hex string.
        // var sb = new chilkat.StringBuilder();
        // sb.AppendEncoded(hashBytes, process.env.CHILKAT_ENCODING_TYPE);
        // console.log("SHA256 hash = " + sb.GetAsString());
        console.log("SHA256 hash = " + hashBytes);
        console.log(`Inserting to : ${process.env.MONGODB_CHAT_USERS_COLLECTION_NAME}`);
        const insertRes = await db.collection(
          process.env.MONGODB_CHAT_USERS_COLLECTION_NAME)
          .insertOne({ "nickname": nickname, "password": hashBytes });
        if (!insertRes) {
            console.log(`Error inserting client into database`);
            return {"err": Numerics["ERR_UNKNOWNERROR"](command, subcommand)};
        }
        return {"res": "Insert successful!"};
    } catch (error) {
        console.log(`InsertClient == ${error}`);
        return {"err": Numerics["ERR_UNKNOWNERROR"](command, subcommand)};
    } finally {
      await mongoClient.close();
    }
};

const CollectionExists = async (collectionName, database=process.env.MONGODB_CHAT_MESSAGE_DB_NAME) => {
  const mongoClient = GetMongoConn();
    try {
        await mongoClient.connect();
        const db = await mongoClient.db(database);
        console.log(`collectionName: ${collectionName}, collectionName type: ${typeof(collectionName)}`)
        const collections = await db.listCollections().toArray();
        for (const col of collections) {
          console.log(`Searching for '${collectionName}' ?= '${col?.name}'`)
          if (col?.name === collectionName) {
            return true;
          }
        }

        return false;
    } catch (error) {
        console.log(`CollectionExists == ${error}`);
        return false;
    } finally {
      await mongoClient.close();
    }
};

module.exports = {
    MongoDB,
    InsertCAPState,
    RegisterClient,
    CheckDBForPassword,
    CheckDBForNickName,
    InsertClient,
    UpdateOne,
    FindOne,
    Find,
    InsertOne,
    FindOneAndUpdate,
    CollectionExists
};