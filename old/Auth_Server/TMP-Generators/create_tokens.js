const cryptoManager = require('../Crypto/crypto_manager');

const salt = cryptoManager.generate_salt();
console.log(`SALT: ${salt}`);
// const password = cryptoManager.generate_salt();
// console.log(`PASSWORD: ${password}`);
const hash = cryptoManager.generate_hash("test", salt);
console.log(`HASH: ${hash}`);