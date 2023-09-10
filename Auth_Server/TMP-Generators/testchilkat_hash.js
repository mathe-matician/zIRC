const chilkatManager = require('./Chilkat/chilkat_manager');
require('dotenv').config();

chilkatManager.unlock_bundle();

// const testObj = JSON.stringify({"data": "hello world"});
const value = chilkatManager.hash_string("test");

// const value = chilkatManager.encrypt_decrypt_AES("test", true);

console.log(`Encrypted value = ${value}`);