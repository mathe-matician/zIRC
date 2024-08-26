const bcrypt = require('bcrypt');
const saltRounds = 10;
const myPlaintextPassword = 'test';

// bcrypt.genSalt(saltRounds, (err, salt) => {
//     bcrypt.hash(myPlaintextPassword, salt, (err, hash) => {
//         // console.log(`Salt: ${salt}`);
//         console.log(`Hash: ${hash}`);
//     });
// });

// Load hash from your password DB.
// bcrypt.compare(myPlaintextPassword, "$2b$10$FXIxXEjRiaoyjvvTYI9VTuvB9K/JkbT0SDmk27ZHQRRH/apiD/Ury", (err, result) => {
//     if (err) {
//         console.log("ERROR")
//         console.log(err)
//     } else {
//         console.log("RESULT")
//         console.log(result)
//     }
// });

const generate_hash = (password) => {
    // const salt = await bcrypt.genSalt(saltRounds);
    // console.log(`Salt: ${salt}`)
    // const hash = await bcrypt.hash(password, salt);
    // console.log(`Hash: ${hash}`)
    bcrypt.genSalt(saltRounds, (err, salt) => {
      bcrypt.hash(password, salt, (err, hash) => {
          // console.log(`Salt: ${salt}`);
          console.log(`Hash: ${hash}`);
      });
    });
  };

// const doer = async () => {
//     await generate_hash("testpassword");
// }

// doer()

generate_hash()