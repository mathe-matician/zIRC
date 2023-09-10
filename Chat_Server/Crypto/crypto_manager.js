const crypto = require('crypto');
require('dotenv').config();

const generate_salt = (randomness=64) => {
  return crypto.randomBytes(randomness).toString("hex");
}

const generate_hash = (value, salt=null) => {
  let _salt = ""; 
  if (!salt) {
      _salt = generate_salt();
      console.log(`SALT: ${_salt}`);
  } else {
      _salt = salt;
  }

  const hash = crypto.pbkdf2Sync(value, _salt, 1000, 64, "sha256").toString("hex");
  console.log(`HASH: ${hash}`);
  return hash;
};

// const _iv = process.env.CRYPT_AES_IV;
// const _iv = "dd20d485701a465289b7a981b0de0e1d";
const iv = Buffer.from("dd20d485701a465289b7a981b0de0e1d");

const decrypt_AES = (value, salt) => {
  const algorithm = 'aes-192-cbc';
  const password = 'Password used to generate key';

  const key = crypto.scryptSync(password, salt, 24);
    // The IV is usually passed along with the ciphertext.
    // const iv = Buffer.alloc(16, 0); // Initialization vector.

    const decipher = crypto.createDecipheriv(algorithm, key, iv);

    let decrypted = '';
    decipher.on('readable', () => {
      let chunk;
      while (null !== (chunk = decipher.read())) {
        decrypted += chunk.toString('utf8');
      }
    });
    decipher.on('end', () => {
      console.log(decrypted);
      // Prints: some clear text data
    });

    // Encrypted with same algorithm, key and iv.
    // const encrypted =
    //   'bdfcd1d9a5d4fb1191248f2ca1a87a81';
    decipher.write(value, 'hex');
    decipher.end();
}

const encrypt_AES = (value, salt) => {
  const algorithm = 'aes-256-cbc';
  const password = 'Password used to generate key';

  // First, we'll generate the key. The key length is dependent on the algorithm.
  // In this case for aes192, it is 24 bytes (192 bits).
  crypto.scrypt(password, salt, 24, (err, key) => {
    if (err) throw err;
    // Then, we'll generate a random initialization vector
    const cipher = crypto.createCipheriv(algorithm, key, iv);
    let encrypted = '';
    cipher.setEncoding('hex');
  
    cipher.on('data', (chunk) => encrypted += chunk);
    cipher.on('end', () => {
      console.log(encrypted)
      decrypt_AES(encrypted, salt, iv)
    });
  
    cipher.write(value);
    cipher.end();
    // crypto.randomFill(new Uint8Array(16), (err, iv) => {
    //   if (err) throw err;
    //   console.log(`Encrypt iv: ${iv}`);
    //   // Once we have the key and iv, we can create and use the cipher...
    //   const cipher = crypto.createCipheriv(algorithm, key, iv);
    
    //   let encrypted = '';
    //   cipher.setEncoding('hex');
    
    //   cipher.on('data', (chunk) => encrypted += chunk);
    //   cipher.on('end', () => {
    //     console.log(encrypted)
    //     decrypt_AES(encrypted, salt, iv)
    //   });
    
    //   cipher.write(value);
    //   cipher.end();
    // });
  });
};

// const salt = generate_salt();
// encrypt_AES("hello world", salt);

// decrypt_AES("bdfcd1d9a5d4fb1191248f2ca1a87a81", salt);

module.exports = {
  generate_hash,
  generate_salt,
  encrypt_AES,
  decrypt_AES
};