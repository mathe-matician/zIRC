const MongoClient = require('mongodb').MongoClient;
const url = 'mongodb://localhost:27017';
const dbName = 'Gymlete_Coach';
const ChatMessage = require(__dirname + "../../schemas/javascript/chat_schema.js");

// Major Images update
// 
//  Adding any previously taken pics to an 'Images' collection
//  Changing 'athlete.attributes.profile_pic_path' -> 'athlete.attributes.profile_pic_id'

let main = async () => {

    MongoClient.connect(url, {useUnifiedTopology: true}, async function(err, client) {
        if (err) {
            console.log('Error Connecting to Mongo... \n' + err);
            return;
        }
        const db = client.db(dbName);

        try {
            const collectionName = 'Chat';
            const exists = (await (await db.listCollections().toArray()).findIndex((item) => item.name === collectionName) !== -1)
            console.log(exists);
            if (!exists) {
                console.log("Creating Chat collection");
                db.createCollection('Chat');
            }
            
            let meFind = { $and: [{'attributes.first_name' : 'Zach'}, {'attributes.last_name' : 'Mathe'}] };
            var athDoc = await db.collection('Athletes').findOne(meFind);
            if (!athDoc) {
                throw 'Couldn\'t find me...';
            }
            // insert test message
            console.log("Inserting test message");
            const msg = ChatMessage(athDoc._id, "127.0.0.1", "General", "Hello World", Date.now(), {"testing": true}, 302);
            const res = await db.collection("Chat").insertOne(msg);
            const msg2 = ChatMessage(athDoc._id, "127.0.0.1", "General", "Hows it going", Date.now(), {"testing": true}, 302);
            const res2 = await db.collection("Chat").insertOne(msg2);
            const msg3 = ChatMessage(athDoc._id, "127.0.0.1", "General", "Good how are you?", Date.now(), {"testing": true}, 302);
            const res3 = await db.collection("Chat").insertOne(msg3);
            const msg4 = ChatMessage(athDoc._id, "127.0.0.1", "General", "fine and dandy", Date.now(), {"testing": true}, 302);
            const res4 = await db.collection("Chat").insertOne(msg4);
            const msg5 = ChatMessage(athDoc._id, "127.0.0.1", "General", "ohhhh boy", Date.now(), {"testing": true}, 302);
            const res5 = await db.collection("Chat").insertOne(msg5);
            console.log('Finished creating Chat collection and adding test message');
        } catch (err) {
            console.log('Error occurred while updating...\n' + err);
        } finally {
            await client.close();
        }
    });
}

main();